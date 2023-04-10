package global

import (
	auth "be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/pkg/database"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repositorys interface {
	Authentication(module *string) gin.HandlerFunc
	CheckPermission(module string) gin.HandlerFunc
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
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		err := auth.ValidateToken(tokenString)
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}

		var jwtTokenCheck string

		query := `SELECT "token" FROM public.users where token=$1 and is_login = true`

		if err := m.DB.QueryRow(query, &tokenString).Scan(
			&jwtTokenCheck,
		); err != nil && err != sql.ErrNoRows {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "token is expired 1"})
			context.Abort()
			return
		}
		if jwtTokenCheck == "" {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "token is expired 2"})
			context.Abort()
			return
		}

		jwtPayload, err := auth.ParseJwtToken(tokenString)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"codes": http.StatusUnauthorized, "messages": http.StatusUnauthorized})
			context.Abort()
			return
		}
		log.Println(jwtPayload)
		context.Set("user_id", jwtPayload.ID)
		context.Set("email", jwtPayload.Email)
		context.Set("user_role", jwtPayload.UserRole)
		context.Set("user_role_id", jwtPayload.UserRoleID)
		context.Set("name_user", jwtPayload.Name)
		context.Set("token", tokenString)

		log.Println("module ", module )
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

func (m *repositorys) CheckPermission(module string) gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		jwtPayload, err := auth.ParseJwtToken(tokenString)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"codes": http.StatusUnauthorized, "messages": http.StatusUnauthorized})
			context.Abort()
			return
		}
		permission := 0 
		log.Println("permit ", permission )
		query := `SELECT count(rp.id) FROM public.route_permissions rp 
		join modules m on rp.module_id = m.id::text 
		where m."key"  = $1 and (rp.relation_id = $2 or rp.relation_id = $3)  and rp.deleted_at IS NULL and rp.can_view = true`

		if err := m.DB.QueryRow(query, module, jwtPayload.ID, jwtPayload.UserRoleID).Scan(
			&permission,
		); 
		err != nil && err != sql.ErrNoRows {
			context.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "User Don't have Permission"})
			context.Abort()
			return
		}
		log.Println("permit ", permission)
		if(permission <= 0 ){
			context.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "User Don't have Permission"})
			context.Abort()
			return
		}
	
		context.Next()
	}
}
