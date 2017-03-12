package util

func TotalPages(pageSize int, totalRecords int)(int){
	return (totalRecords + pageSize - 1) / pageSize
}

func Range(pageNo int, pageSize int, totalRecords int)(int, int){
    if totalRecords == 0 {
        return 0, 0
    }
    totalPages := TotalPages(pageSize, totalRecords)
    if pageNo > totalPages {
        pageNo = totalPages
    } else if pageNo < 1 {
        pageNo = 1
    }

    begin := (pageNo-1)*pageSize
    end := begin + pageSize
    if end > totalRecords {
        end = totalRecords
    }
    return begin, end
}