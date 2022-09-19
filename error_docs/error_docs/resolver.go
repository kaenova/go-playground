package errordocs

func ErrorDocs(num int, err ...error) error {
	idx := searchCodeInDocs(num)
	if idx != -1 {
		return Docs[idx]
	}
	return nil
}

func searchCodeInDocs(num int) int {
	for i := 0; i < len(Docs); i++ {
		if num == Docs[i].Code {
			return num
		}
	}
	return -1
}
