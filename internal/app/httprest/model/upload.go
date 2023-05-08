package model

type UploadFileResponse struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	Filepath string `json:"file_path"`
}
