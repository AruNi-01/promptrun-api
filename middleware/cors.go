package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"regexp"
)

// Cors 跨域配置
func Cors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost"}                                         // 允许所有域名
	config.AllowCredentials = true                                                             // 允许携带 cookie
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"} // 允许的请求方法
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Cookie"}       // 允许的请求头

	// 生产环境需要配置实际的跨域域名，否则 403
	if gin.Mode() == gin.ReleaseMode {
		config.AllowOrigins = []string{"https://promptrun.0x3f4.run", "https://promptrun.shop", "http://promptrun.0x3f4.run", "http://promptrun.shop"}
	} else {
		// 测试环境下模糊匹配本地开头的请求
		config.AllowOriginFunc = func(origin string) bool {
			if regexp.MustCompile(`^http://127\.0\.0\.1:\d+$`).MatchString(origin) {
				return true
			}
			if regexp.MustCompile(`^http://localhost:\d+$`).MatchString(origin) {
				return true
			}
			return false
		}
	}

	return cors.New(config)
}
