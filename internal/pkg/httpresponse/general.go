package httpresponse

import (
	"log"
	"net/http"
)

type Body struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Error      string      `json:"error,omitempty"`
	Data       interface{} `json:"data"`
	Pagination *pagination `json:"pagination,omitempty"`
}

type pagination struct {
	TotalData        interface{} `json:"total_data"`
	TotalDataPerPage interface{} `json:"total_data_perpage"`
	Limit            interface{} `json:"limit"`
	TotalPage        interface{} `json:"total_page"`
	Page             interface{} `json:"page"`
}

func Format(message string, err error, data ...interface{}) (statusCode int, b *Body) {
	var (
		msg string
		d   interface{}

		pg = pagination{}

		code int
	)

	switch message {
	case DATANOTFOUND_400:
		code = http.StatusBadRequest
		msg = DATANOTFOUND_400
		break
	case CREATESUCCESS_200:
		code = http.StatusOK
		msg = CREATESUCCESS_200
		break
	case CREATEFAILED_400:
		code = http.StatusBadRequest
		msg = CREATEFAILED_400
		break
	case CREATEDUPLICATE_400:
		code = http.StatusBadRequest
		msg = CREATEDUPLICATE_400
		break
	case READSUCCESS_200:
		code = http.StatusOK
		msg = READSUCCESS_200
		break
	case CONTENTNOTFOUND_404:
		code = http.StatusNotFound
		msg = CONTENTNOTFOUND_404
		break
	case READFAILED_400:
		code = http.StatusBadRequest
		msg = READFAILED_400
		break
	case UPDATESUCCESS_200:
		code = http.StatusOK
		msg = UPDATESUCCESS_200
		break
	case UPDATEFAILED_400:
		code = http.StatusBadRequest
		msg = UPDATEFAILED_400
		break
	case DELETESUCCESS_200:
		code = http.StatusOK
		msg = DELETESUCCESS_200
		break
	case DELETEFAILED_400:
		code = http.StatusBadRequest
		msg = DELETEFAILED_400
		break
	case UPLOADSUCCESS_200:
		code = http.StatusOK
		msg = UPLOADSUCCESS_200
		break
	case UPLOADFAILED_400:
		code = http.StatusBadRequest
		msg = UPLOADFAILED_400
		break
	case DOWNLOADSUCCESS_200:
		code = http.StatusOK
		msg = DOWNLOADSUCCESS_200
		break
	case DOWNLOADFAILED_400:
		code = http.StatusBadRequest
		msg = DOWNLOADFAILED_400
		break
	default:
		code = http.StatusInternalServerError
		msg = ERR_GENERAL_400
		break
	}

	if err != nil {
		log.Println(err)
	}

	if len(data) >= 1 {
		d = data[0]
	}

	b = &Body{
		Code:    code,
		Data:    d,
		Message: msg,
	}

	if err != nil {
		b.Error = err.Error()
	}

	if len(data) > 1 {
		pg.TotalData = data[1]
		pg.Page = data[2]
		pg.Limit = data[3]
		pg.TotalDataPerPage = data[4]
		pg.TotalPage = data[5]

		b.Pagination = &pg
	}

	statusCode = code

	return
}
