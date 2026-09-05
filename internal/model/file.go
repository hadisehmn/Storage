package model

import "time"

type File struct {
	ID        string
	UserID    string
	FileName  string
	FilePath  string
	CreatedAt time.Time
}

type FileResponse struct {
	ID        string    `json:"id"`
	FileName  string    `json:"file_name"`
	CreatedAt time.Time `json:"created_at"`
}
