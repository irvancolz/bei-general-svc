package model

type UploadedFilesMenuResponse struct {
	Id          string `json:"id"`
	Report_Type string `json:"report_type"`
	Report_Code string `json:"report_code"`
	Report_Name string `json:"report_name"`
	Is_Uploaded bool   `json:"is_uploaded"`
	File_Name   string `json:"file_name"`
	File_Path   string `json:"file_path"`
	File_Size   int64  `json:"file_size"`
	Created_By  string `json:"created_by"`
	Created_At  int64  `json:"created_at"`
	Updated_By  string `json:"updated_by"`
	Updated_At  int64  `json:"updated-at"`
}
