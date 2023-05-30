package model

import "database/sql"

type UploadedFilesMenuResponse struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
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

type UploadedFilesMenuResultSet struct {
	Id          string         `json:"id"`
	Type        string         `json:"type"`
	Report_Code string         `json:"report_code"`
	Report_Name string         `json:"report_name"`
	Is_Uploaded bool           `json:"is_uploaded"`
	File_Name   sql.NullString `json:"file_name"`
	File_Path   sql.NullString `json:"file_path"`
	File_Size   sql.NullInt64  `json:"file_size"`
	Created_By  string         `json:"created_by"`
	Created_At  int64          `json:"created_at"`
	Updated_By  sql.NullString `json:"updated_by"`
	Updated_At  sql.NullInt64  `json:"updated-at"`
}
