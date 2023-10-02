/*
 * @Author: SpenserCai
 * @Date: 2023-09-29 19:27:29
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-02 21:02:04
 * @Description: file content
 */
package middleware

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"log"

	DbotUser "github.com/SpenserCai/sd-webui-discord/user"

	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJwt(bearerHeader string) (interface{}, error) {
	bearerToken := strings.Split(bearerHeader, " ")[1]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error decoding token")
		}
		return []byte(global.Config.WebSite.Api.JwtSecret), nil
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if token.Valid {
		// 验证token是否过期
		if claims["exp"].(float64) < float64(time.Now().Unix()) {
			return nil, errors.New("token expired")
		}
		return claims["id"].(string), nil
	}
	return nil, errors.New("invalid token")
}

func BuildJwt(userInfo DbotUser.UserInfo, other map[string]string) (string, error) {
	claims := jwt.MapClaims{
		"id":   userInfo.Id,
		"name": userInfo.Name,
		"role": userInfo.Roles,
	}
	for k, v := range other {
		claims[k] = v
	}
	// 超时时间7天
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(global.Config.WebSite.Api.JwtSecret))
}

// 刷新jwt
func RefreshJwt(bearerHeader string) (string, error) {
	bearerToken := strings.Split(bearerHeader, " ")[1]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("error decoding token")
			return nil, fmt.Errorf("error decoding token")
		}
		return []byte(global.Config.WebSite.Api.JwtSecret), nil
	})
	if err != nil {
		log.Println(err)
		return "", err
	}
	if token.Valid {
		// 超时时间7天
		claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(global.Config.WebSite.Api.JwtSecret)
	}
	return "", errors.New("invalid token")
}

// 解码jwt
func DecodeJwt(bearerHeader string) (jwt.MapClaims, error) {
	bearerToken := strings.Split(bearerHeader, " ")[1]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("error decoding token")
			return nil, fmt.Errorf("error decoding token")
		}
		return []byte(global.Config.WebSite.Api.JwtSecret), nil
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
