/*
 * @Author: SpenserCai
 * @Date: 2023-09-29 19:27:29
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-29 20:10:01
 * @Description: file content
 */
package middleware

import (
	"errors"
	"fmt"
	"strings"

	"log"

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
		return global.Config.WebSite.Api.JwtSecret, nil
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if token.Valid {
		return claims["id"].(string), nil
	}
	return nil, errors.New("invalid token")
}
