package middleware

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-remix/appo"
	"go-remix/utils"
	"log"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			c.JSON(http.StatusUnauthorized, appo.NewAuthorization("权限不足"))
			c.Abort()
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := utils.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, appo.NewAuthorization("权限不足"))
			c.Abort()
			return
		}

		// 验证通过后获取claim 中的userId
		userId := claims.UserId

		session := sessions.Default(c)
		sToken := session.Get("access_token")

		if sToken == nil {
			c.JSON(http.StatusUnauthorized, appo.NewAuthorization("权限不足"))
			c.Abort()
			return
		}

		fmt.Println(token, userId, sToken.(string))

		c.Set("userId", userId)

		session.Set("access_token", token)
		if err := session.Save(); err != nil {
			log.Printf("Failed recreate the session: %v\n", err.Error())
		}

		c.Next()
	}
}
