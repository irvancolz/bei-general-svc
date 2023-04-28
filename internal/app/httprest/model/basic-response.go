package model

import (
	"be-idx-tsg/internal/pkg/httpresponse"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	error_type_general              = "general"
	error_type_asset_directory      = "asset_directory"
	error_type_invalid_json_request = "invalid_json_request"
	error_type_flow                 = "flow"
	error_type_flow_data_not_found  = "flow_data_not_found"
	error_type_flow_create          = "flow_create"
	error_type_flow_account_blocked = "flow_account_blocked"
	error_type_already_login        = "already_login"
	error_type_flow_update          = "flow_update"
	error_type_flow_delete          = "flow_delete"
	error_type_flow_read            = "flow_read"
	error_type_file_not_found       = "file_not_found"
	error_type_file_related         = "file_related"
error_type_token                 = "token"
)

type BaseErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func generateErrorResponse(c *gin.Context, flag string, message string, http_code int) {
	response := BaseErrorResponse{
		Code:    strconv.Itoa(http_code) + "-" + flag,
		Message: message,
	}

	c.JSON(http_code,
		response)
}

func GenerateInvalidJsonResponse(c *gin.Context, err error) {
	generateErrorResponse(c, error_type_invalid_json_request, httpresponse.ERR_REQUESTBODY_400+". "+err.Error(), http.StatusBadRequest)
}

func GenerateFlowErrorResponse(c *gin.Context, err error) {
	generateErrorResponse(c, error_type_flow, err.Error(), http.StatusBadRequest)
}

func GenerateFlowErrorFromMessageResponse(c *gin.Context, message string) {
	err := errors.New(message)
	generateErrorResponse(c, error_type_flow, err.Error(), http.StatusBadRequest)
}

func GenerateInsertErrorResponse(c *gin.Context, err error) {
	generateErrorResponse(c, error_type_flow_create, httpresponse.CREATEFAILED_400+". "+err.Error(), http.StatusBadRequest)
}

func GenerateUpdateErrorResponse(c *gin.Context, err error) {
	generateErrorResponse(c, error_type_flow_update, httpresponse.UPDATEFAILED_400+". "+err.Error(), http.StatusBadRequest)
}

func GenerateDeleteErrorResponse(c *gin.Context, err error) {
	generateErrorResponse(c, error_type_flow_delete, httpresponse.DELETEFAILED_400+". "+err.Error(), http.StatusBadRequest)
}

func GenerateReadErrorResponse(c *gin.Context, err error) {
	generateErrorResponse(c, error_type_flow_read, httpresponse.READFAILED_400+". "+err.Error(), http.StatusBadRequest)
}

func GenerateInternalErrorResponse(c *gin.Context, err error) {
	generateErrorResponse(c, error_type_asset_directory, err.Error(), http.StatusInternalServerError)
}

func GenerateIFileNotFoundErrorResponse(c *gin.Context, err error) {
	generateErrorResponse(c, error_type_file_not_found, "file not found. "+err.Error(), http.StatusNotFound)
}

func GenerateRemoveFileErrorResponse(c *gin.Context, err error) {
	generateErrorResponse(c, error_type_file_related, "gagal menghapus file. "+err.Error(), http.StatusInternalServerError)
}

func GenerateTokenEmptyResponse(c *gin.Context) {
	generateErrorResponse(c, error_type_token, "token is empty", http.StatusUnauthorized)
}

func GenerateTokenErrorResponse(c *gin.Context, err error) {
	generateErrorResponse(c, error_type_token, err.Error(), http.StatusUnauthorized)
}

