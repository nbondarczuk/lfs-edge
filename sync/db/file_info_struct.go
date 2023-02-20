package db

type FileInfo struct {
	ID     uint64
	Exist  bool
	Status int32
	Size   int64
	Name   string
}

