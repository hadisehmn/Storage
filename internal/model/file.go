package model

import "time"

type File struct {
	ID        string
	UserID    string
	FileName  string
	FilePath  string
	CreatedAt time.Time
}
