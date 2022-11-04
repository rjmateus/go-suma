package middleware

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/gin-gonic/gin"
	"github.com/rjmateus/go-suma/config"
	channelRepo "github.com/rjmateus/go-suma/repositories/channel"
	"log"
	"net/http"
	"strings"
	"time"
)

type jwtTokenData struct {
	Exp          int64
	Iat          int64
	Nbf          int64
	Jti          string
	Org          int
	OnlyChannels []string
}

func (token jwtTokenData) isExpired() bool {
	return token.Exp < time.Now().Unix()
}

func (token jwtTokenData) verifyChannelAccess(channel string) bool {
	if len(token.OnlyChannels) == 0 {
		return false
	}

	for _, value := range token.OnlyChannels {
		if value == channel {
			return true
		}
	}

	return false
}

func JwtAuthenticationTokenMiddleware(app *config.Application) func(c *gin.Context) {
	return func(c *gin.Context) {

		if !app.Config.CheckDownloadToken() {
			return
		}

		tokenString, error := getTokenFromRequest(c)
		if error != nil {
			c.Abort()
			return
		}

		//01 Verify token exists in the database
		// FIXME

		//02 decodes token
		sharedKey := []byte("27605eb97bc7a8ed6e80d451ed3149c9de2ee477b49075a764c6186ab6f62f20")

		dst := make([]byte, hex.DecodedLen(len(sharedKey)))
		n, err := hex.Decode(dst, sharedKey)
		if err != nil {
			log.Fatal(err)
		}

		payload, _, err := jose.Decode(tokenString, dst[:n])

		if err == nil {
			//go use token
			//fmt.Printf("\npayload = %v\n", payload)

			//and/or use headers
			//fmt.Printf("\nheaders = %v\n", headers)

			var token jwtTokenData
			err := json.Unmarshal([]byte(payload), &token)
			if err != nil {
				c.String(http.StatusForbidden, "This token Data is not valid")
				c.Abort()
			}

			//03 checks token is not expired
			if token.isExpired() {
				c.String(http.StatusForbidden, "Token expired")
				c.Abort()
			}
			//04 Verify access to channel? should it be made in here on in a different middleware
			channel := c.Param("channel")
			if !token.verifyChannelAccess(channel) {
				c.String(http.StatusForbidden, "Forbidden: Token does not provide access to channel %s", channel)
				c.Abort()
			}

			if !channelRepo.IsAccessibleBy(app.DBGorm, channel, token.Org) {
				c.String(http.StatusForbidden, "Forbidden: Token does not provide access to channel %s", channel)
				c.Abort()
			}

		} else {
			c.String(http.StatusForbidden, "This token is not valid")
			c.Abort()
		}

	}
}

func getTokenFromRequest(c *gin.Context) (string, error) {
	header := c.GetHeader("X-Mgr-Auth")
	if len(strings.TrimSpace(header)) == 0 {
		authorizationHeader := c.GetHeader("Authorization")
		if len(authorizationHeader) > 0 && strings.HasPrefix(authorizationHeader, "Basic") {
			encodedData := strings.TrimSpace(authorizationHeader[len("Basic"):])
			headerBytes, err := base64.StdEncoding.DecodeString(encodedData)
			if err != nil {
				panic(err)
			}
			header = string(headerBytes)
		}
	}

	queryStrings := c.Request.URL.Query()

	if len(queryStrings) == 0 && len(strings.TrimSpace(header)) == 0 {
		c.String(http.StatusForbidden, "You need a token to access %s", c.Request.URL.Path)
		return "", errors.New(fmt.Sprintf("You need a token to access %s", c.Request.URL.Path))
	} else {

		if (len(queryStrings) > 1 && len(header) == 0) || (len(queryStrings) > 0 && len(header) > 0) {
			c.String(http.StatusBadRequest, "Bad Request: Only one token is accepted")
			return "", errors.New("bad Request: Only one token is accepted")
		}

		if len(queryStrings) > 0 {
			result := ""
			for k, _ := range queryStrings {
				result = k
				break
			}
			return result, nil
		} else {
			return header, nil
		}
	}
}
