package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CORS struct {
	Mode      string          `mapstructure:"mode" json:"mode" yaml:"mode"`
	Whitelist []CORSWhitelist `mapstructure:"whitelist" json:"whitelist" yaml:"whitelist"`
}
type CORSWhitelist struct {
	AllowOrigin      string `mapstructure:"allow-origin" json:"allow-origin" yaml:"allow-origin"`
	AllowMethods     string `mapstructure:"allow-methods" json:"allow-methods" yaml:"allow-methods"`
	AllowHeaders     string `mapstructure:"allow-headers" json:"allow-headers" yaml:"allow-headers"`
	ExposeHeaders    string `mapstructure:"expose-headers" json:"expose-headers" yaml:"expose-headers"`
	AllowCredentials bool   `mapstructure:"allow-credentials" json:"allow-credentials" yaml:"allow-credentials"`
}

var coreCfg = CORS{}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, New-Token, New-Expires-At")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func CorsByRules() gin.HandlerFunc {
	if coreCfg.Mode == "allow-all" {
		return Cors()
	}
	return func(c *gin.Context) {
		whitelist := checkCors(c.GetHeader("origin"))

		if whitelist != nil {
			c.Header("Access-Control-Allow-Origin", whitelist.AllowOrigin)
			c.Header("Access-Control-Allow-Headers", whitelist.AllowHeaders)
			c.Header("Access-Control-Allow-Methods", whitelist.AllowMethods)
			c.Header("Access-Control-Expose-Headers", whitelist.ExposeHeaders)
			if whitelist.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		if whitelist == nil && coreCfg.Mode == "strict-whitelist" && !(c.Request.Method == "GET" && c.Request.URL.Path == "/health") {
			c.AbortWithStatus(http.StatusForbidden)
		} else {
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusNoContent)
			}
		}

		c.Next()
	}
}

func checkCors(currentOrigin string) *CORSWhitelist {
	for i := range coreCfg.Whitelist {
		whitelistCopy := coreCfg.Whitelist[i]
		if subtle.ConstantTimeCompare([]byte(currentOrigin), []byte(whitelistCopy.AllowOrigin)) == 1 {
			return &whitelistCopy
		}
	}
	return nil
}
