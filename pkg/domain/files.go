package domain

import "time"

// File 用户上传的附件
type File struct {
	Name       string
	URI        string
	CreateTime time.Time
	UsedBy     interface{}
}
