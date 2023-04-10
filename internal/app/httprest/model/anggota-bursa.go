package model

// "time"

type GetAnggotaBursa struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
type APIResponseGetABMinModel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    GetAnggotaBursa
}

//not used
type GetAllAnggotaBursa struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
type APIResponseGetAllABMinModel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []GetAllAnggotaBursa
}
