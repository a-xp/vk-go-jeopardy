package domain

type Int64Slice []int64

func (slice Int64Slice) Search(value int64) bool {
	for i := range slice {
		if slice[i] == value {
			return true
		}
	}
	return false
}
