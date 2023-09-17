package core

import (
	"crypto/tls"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func ServerStart() {
	engine := setupGinEngine()
	setupRoutes(engine)
	initAndStartServer(engine)
	console()
}

// 创建和配置Gin引擎
func setupGinEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	engine := gin.New()
	// 设置信任的代理
	if err := engine.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatal(err)
	}
	return engine
}

// 定义路由和中间件
func setupRoutes(engine *gin.Engine) {
	routerGroup := engine.Group("/", domainMiddleware(config.Server.Domain))
	routerGroup.GET("/copilot_internal/v2/token", getGithubCopilotToken())
	routerGroup.GET("/state", getState())

}

// 域名中间件
func domainMiddleware(domain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if domain != "" {
			// 检查域名是否匹配
			requestDomain := strings.Split(c.Request.Host, ":")[0]
			if requestDomain == domain || requestDomain == "127.0.0.1" {
				c.Next()
			} else {
				c.String(403, "Forbidden")
				c.Abort()
			}
		} else {
			c.Next()
		}
	}
}

// 初始化和启动服务器
func initAndStartServer(engine *gin.Engine) {
	listenAddress := config.Server.Host + ":" + strconv.Itoa(config.Server.Port)
	server := createTLSServer(engine, listenAddress)
	go func() {
		if config.Server.Port != 443 {
			err := engine.Run(listenAddress)
			log.Fatal(err)
		} else {
			err := server.ListenAndServeTLS(config.Server.CertPath, config.Server.KeyPath)
			log.Fatal(err)
		}
	}()
}

// 创建TLS服务器配置
func createTLSServer(engine *gin.Engine, address string) *http.Server {
	return &http.Server{
		Addr: address,
		TLSConfig: &tls.Config{
			NextProtos: []string{"http/1.1", "http/1.2", "http/2"},
		},
		Handler: engine,
	}
}

// 服务器响应
func proxyResponse(c *gin.Context, respDataMap map[string]interface{}) {
	// 将map转换为JSON字符串
	responseJSON, err := json.Marshal(respDataMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON marshaling error"})
	}
	// 请求成功统计
	successCount++
	// 将JSON字符串作为响应体返回
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(responseJSON))
}

// 请求错误
func badRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message":           "Bad credentials",
		"documentation_url": "https://docs.github.com/rest"})
}
