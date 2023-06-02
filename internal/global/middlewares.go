package global

import (
	auth "be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"be-idx-tsg/internal/pkg/database"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type GetTokens struct {
	Token string `json:"token"`
}
type Repositorys interface {
	Authentication(module *string) gin.HandlerFunc
}

type repositorys struct {
	DB *sqlx.DB
}

func NewRepositorys() Repositorys {
	return &repositorys{
		DB: database.Init().MySql,
	}
}

func (m *repositorys) Authentication(module *string) gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			model.GenerateTokenEmptyResponse(context)
			context.Abort()
			return
		}
		err := auth.ValidateToken(tokenString)
		if err != nil {
			model.GenerateTokenErrorResponse(context, err)
			context.Abort()
			return
		}

		jwtTokenCheck, err := GetToken(tokenString)
		if err != nil {
			model.GenerateTokenErrorResponse(context, err)
			context.Abort()
			return
		}

		if jwtTokenCheck.Token == "" {
			model.GenerateTokenEmptyResponse(context)
			context.Abort()
			return
		}

		jwtPayload, err := auth.ParseJwtToken(tokenString)
		if err != nil {
			model.GenerateTokenErrorResponse(context, err)
			context.Abort()
			return
		}
		context.Set("user_id", jwtPayload.ID)
		context.Set("email", jwtPayload.Email)
		context.Set("token", tokenString)
		context.Set("name_user", jwtPayload.UserName)
		context.Set("type", jwtPayload.GroupType)
		context.Set("external_type", jwtPayload.ExternalType)
		context.Set("user_role", jwtPayload.UserRole)
		context.Set("user_role_id", jwtPayload.UserRoleID)
		context.Set("company_name", jwtPayload.CompanyName)
		context.Set("company_code", jwtPayload.CompanyCode)
		context.Set("company_id", jwtPayload.CompanyId)
		context.Set("name", jwtPayload.Name)
		context.Set("user_form_role", jwtPayload.UserFormRole)

		log.Println("module ", module)
		if module != nil {
			value, error := auth.CheckPermission(tokenString, module)
			if !value.Status || error != nil {
				context.JSON(value.Code, gin.H{"codes": value.Code, "messages": value.Message, "status": value.Status})
				context.Abort()
				return
			}
		}
		context.Next()
	}
}

func GetToken(token string) (*GetTokens, error) {

	err_host := godotenv.Load(".env")
	if err_host != nil {
		fmt.Println("err_host ", err_host)
	}
	host := os.Getenv("SERVICE_AUTH_HOST")
	url := host + "/user/get-token?token=" + token
	var payload string
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("payloadBytes ", err)
	}

	bodyReq := bytes.NewReader(payloadBytes)
	Request, err := http.NewRequest("GET", url, bodyReq)
	if err != nil {
		log.Println("error")
		log.Println(err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(Request)
	if err != nil {
		log.Println("[TAP-debug] [err] [PostRequest][Do]", err)
		log.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("[TAP-debug] [err] [PostRequest][ReadAll]", err)
		log.Println(err)
		return nil, err
	}
	datas := &GetTokens{}
	log.Printf(datas.Token)
	errorUM := json.Unmarshal([]byte(body), datas)

	return datas, errorUM
}
