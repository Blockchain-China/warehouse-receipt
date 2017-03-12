package member

import (
    "crypto/ecdsa"
    "crypto/rand"
    "crypto/x509"
    "errors"
    "time"
    "fmt"
    "os"

    "github.com/golang/protobuf/proto"
    "github.com/golang/protobuf/ptypes/timestamp"
    "github.com/hyperledger/fabric/core/peer"
    "github.com/hyperledger/fabric/core/config"
    "github.com/hyperledger/fabric/core/crypto/primitives"
    "github.com/hyperledger/fabric/core/crypto/primitives/ecies"
    "github.com/hyperledger/fabric/membersrvc/ca"
    pb "github.com/hyperledger/fabric/membersrvc/protos"
    "golang.org/x/net/context"
    "github.com/spf13/viper"
    "google.golang.org/grpc"

    "github.com/hyperledger/fabric/work/receipt2/util"
)

type User struct {
    enrollID               string
    enrollPwd              []byte
    enrollPrivKey          *ecdsa.PrivateKey
    role                   int
    affiliation            string
    registrarRoles         []string
    registrarDelegateRoles []string
}

type Member struct {
    MemberCode string
    PubKey string
}

var (
    // NVP related objects
    peerClientConn *grpc.ClientConn
    ecapClient   pb.ECAPClient
    ecaaClient   pb.ECAAClient

    admin = User{enrollID: "admin", enrollPwd: []byte("Xurw3yU9zI0l")}
)

func Register(member Member)(msg string) {
    // Initialize a non-validating peer.
    if err := initNVP(); err != nil {
        fmt.Printf("Failed initiliazing NVP [%s]", err)
        return "{\"Msg\":\""+err.Error()+"\"}"
    }

    newUser := User{enrollID: member.MemberCode, role: 1, affiliation: "institution_a"}

    if admin.enrollPrivKey == nil {
        err := enrollUser(&admin)
        if err != nil {
            fmt.Printf("Failed to enroll admin: [%s]\n", err)
            return "{\"Msg\":\""+err.Error()+"\"}"
        }

        if err = os.RemoveAll("../mdb"); err != nil {
            fmt.Printf("Failed removing [../mdb] [%s]\n", err)
            return "{\"Msg\":\""+err.Error()+"\"}"
        }
    }
    
    err := registerUser(admin, &newUser)
    if err != nil {
        if err.Error() == "User is already registered" {
            fmt.Printf("User is already registered")
        } else {
            fmt.Printf("Failed to register user: [%s]\n", err)
        }
        return "{\"Msg\":\""+err.Error()+"\"}"  
    } else {
        util.PutMember(newUser.enrollID+"_pwd", string(newUser.enrollPwd))
        util.PutMember(newUser.enrollID, member.PubKey)
        return "{\"Msg\":\"OK\"}"
    }
}    

func initNVP() (err error) {
    if err = initPeerClient(); err != nil {
        fmt.Printf("Failed to initPeerClient [%s]\n", err)
        return
    }
    return
}

func initPeerClient() (err error) {
    config.SetupTestConfig("../")
    viper.Set("ledger.blockchain.deploy-system-chaincode", "false")
    viper.Set("peer.validator.validity-period.verification", "false")

    primitives.SetSecurityLevel("SHA3", 256)

    //peerClientConn, err = peer.NewPeerClientConnection()
    peerClientConn, err = peer.NewPeerClientConnectionWithAddress("0.0.0.0:7054")
    if err != nil {
        fmt.Printf("error connection to server at host:port = %s\n", viper.GetString("peer.address"))
        return
    }
    
    ecapClient = pb.NewECAPClient(peerClientConn)
    ecaaClient = pb.NewECAAClient(peerClientConn)

    return
}

// helper function for multiple tests
func enrollUser(user *User) error {

    // ecap := &ca.ECAP{eca}

    // Phase 1 of the protocol: Generate crypto material
    signPriv, err := primitives.NewECDSAKey()
    user.enrollPrivKey = signPriv
    if err != nil {
        return err
    }
    signPub, err := x509.MarshalPKIXPublicKey(&signPriv.PublicKey)
    if err != nil {
        return err
    }

    encPriv, err := primitives.NewECDSAKey()
    if err != nil {
        return err
    }
    encPub, err := x509.MarshalPKIXPublicKey(&encPriv.PublicKey)
    if err != nil {
        return err
    }

    req := &pb.ECertCreateReq{
        Ts:   &timestamp.Timestamp{Seconds: time.Now().Unix(), Nanos: 0},
        Id:   &pb.Identity{Id: user.enrollID},
        Tok:  &pb.Token{Tok: user.enrollPwd},
        Sign: &pb.PublicKey{Type: pb.CryptoType_ECDSA, Key: signPub},
        Enc:  &pb.PublicKey{Type: pb.CryptoType_ECDSA, Key: encPub},
        Sig:  nil}

    resp, err := ecapClient.CreateCertificatePair(context.Background(), req)
    if err != nil {
        fmt.Printf("Failed to CreateCertificatePair: [%s]\n", err)
        return err
    }

    //Phase 2 of the protocol
    spi := ecies.NewSPI()
    eciesKey, err := spi.NewPrivateKey(nil, encPriv)
    if err != nil {
        return err
    }

    ecies, err := spi.NewAsymmetricCipherFromPublicKey(eciesKey)
    if err != nil {
        return err
    }

    out, err := ecies.Process(resp.Tok.Tok)
    if err != nil {
        return err
    }

    req.Tok.Tok = out
    req.Sig = nil

    hash := primitives.NewHash()
    raw, _ := proto.Marshal(req)
    hash.Write(raw)

    r, s, err := ecdsa.Sign(rand.Reader, signPriv, hash.Sum(nil))
    if err != nil {
        return err
    }
    R, _ := r.MarshalText()
    S, _ := s.MarshalText()
    req.Sig = &pb.Signature{Type: pb.CryptoType_ECDSA, R: R, S: S}

    resp, err = ecapClient.CreateCertificatePair(context.Background(), req)
    if err != nil {
        return err
    }

    // Verify we got valid crypto material back
    x509SignCert, err := primitives.DERToX509Certificate(resp.Certs.Sign)
    if err != nil {
        return err
    }

    _, err = primitives.GetCriticalExtension(x509SignCert, ca.ECertSubjectRole)
    if err != nil {
        return err
    }

    x509EncCert, err := primitives.DERToX509Certificate(resp.Certs.Enc)
    if err != nil {
        return err
    }

    _, err = primitives.GetCriticalExtension(x509EncCert, ca.ECertSubjectRole)
    if err != nil {
        return err
    }

    return nil
}

func registerUser(registrar User, user *User) error {

    // ecaa := &ca.ECAA{eca}

    // create req
    req := &pb.RegisterUserReq{
        Id:          &pb.Identity{Id: user.enrollID},
        Role:        pb.Role(user.role),
        Affiliation: user.affiliation,
        Registrar: &pb.Registrar{
            Id:            &pb.Identity{Id: registrar.enrollID},
            Roles:         user.registrarRoles,
            DelegateRoles: user.registrarDelegateRoles,
        },
        Sig: nil}

    //sign the req
    hash := primitives.NewHash()
    raw, _ := proto.Marshal(req)
    hash.Write(raw)

    r, s, err := ecdsa.Sign(rand.Reader, registrar.enrollPrivKey, hash.Sum(nil))
    if err != nil {
        msg := "Failed to register user. Error (ECDSA) signing request: " + err.Error()
        fmt.Println(msg)
        return errors.New(msg)
    }
    R, _ := r.MarshalText()
    S, _ := s.MarshalText()
    req.Sig = &pb.Signature{Type: pb.CryptoType_ECDSA, R: R, S: S}

    token, err := ecaaClient.RegisterUser(context.Background(), req)
    if err != nil {
        fmt.Printf("Failed to RegisterUser: [%s]\n", err)
        return err
    }

    if token == nil {
        return errors.New("Failed to obtain token")
    }

    //need the token for later tests
    user.enrollPwd = token.Tok

    return nil
}
