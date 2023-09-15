package global

import (
	"be-idx-tsg/internal/app/helper"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func (m *repositorys) AuthenticationTest(module *string) gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		jwtPayload := &helper.JWTClaim2{
			ID: "",
			Email: "",
			ExternalType: "participant",
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
		
		context.Next()
	}
}
