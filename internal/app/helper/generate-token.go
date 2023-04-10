package helper

import (
	"be-idx-tsg/internal/app/httprest/model"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/dgrijalva/jwt-go.v3"
)

var jwtKey = []byte("4Q1S3CR3TK3Y")

type JWTClaim struct {
	ID              string                `json:"id"`
	Name            *string               `json:"name"`
	Email           *string               `json:"email"`
	UserName        *string               `json:"user_name"`
	UserRole        *string               `json:"user_role"`
	UserRoleID      *string               `json:"user_role_id"`
	UserType        *string               `json:"user_type"`
	GroupType       *string               `json:"group_type"`
	CompanyId       *string               `json:"company_id"`
	CompanyName     *string               `json:"company_name"`
	CompanyCode     *string               `json:"company_code"`
	FirstLogin      bool                  `json:"first_login"`
	PasswordExpired *string               `json:"password_expired"`
	UserFormRole    []*model.UserRoleForm `json:"user_form_role"`
	jwt.StandardClaims
}

func GenerateJWT(Auth *model.AuthenticationResponse, expiredLogin int) (tokenString string, expiredTime string, err error) {
	log.Println("company id ", Auth.CompanyId)
	expirationTime := time.Now().Add(time.Duration(expiredLogin) * time.Minute)
	claims := &JWTClaim{
		ID:              Auth.ID,
		Name:            Auth.Name,
		Email:           Auth.Email,
		UserName:        Auth.UserName,
		CompanyId:       Auth.CompanyId,
		CompanyName:     Auth.CompanyName,
		CompanyCode:     Auth.CompanyCode,
		FirstLogin:      Auth.FirstLogin,
		PasswordExpired: Auth.PasswordExpired,
		GroupType:       Auth.GroupType,
		UserRoleID:      Auth.UserRoleID,
		UserRole:        Auth.UserRole,
		UserFormRole:    Auth.UserRoleForm,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	expiredTime = time.Now().Add(time.Duration(expiredLogin) * time.Minute).Format("2006-01-02 15:04:05")
	return
}
func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return
		}
		log.Print("ada error", err)
		return
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}

func ParseJwtToken(tokenStr string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return &JWTClaim{}, fmt.Errorf("failed to claim token: %v", err)
	}

	return claims, nil
}
