package core

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

func init() {
	initValidTokenList()
}

// 初始化有效的github token列表
func initValidTokenList() {
	// 为了安全起见，应该等待请求完成并处理其响应。
	var wg sync.WaitGroup
	for _, token := range config.Copilot.Tokens {
		wg.Add(1)
		go func(token string) {
			defer wg.Done()
			if getGithubApi(token) {
				validGithubTokens[token] = true
			}
		}(token)
	}
	wg.Wait()
}

// 请求github api
func getGithubApi(token string) bool {
	githubApiCount++
	// 设置请求头
	headers := map[string]string{
		"Authorization": "token " + token,
		/*"editor-version":        c.GetHeader("editor-version"),
		"editor-plugin-version": c.GetHeader("editor-plugin-version"),
		"user-agent":            c.GetHeader("user-agent"),
		"accept":                c.GetHeader("accept"),
		"accept-encoding":       c.GetHeader("accept-encoding"),*/
	}
	// 发起GET请求
	response, err := resty.New().R().SetHeaders(headers).Get(config.Copilot.GithubApiUrl)
	if err != nil {
		// 处理请求错误
		return false
	}
	// 判断响应状态码
	if response.StatusCode() == http.StatusOK {
		// 响应状态码为200 OK
		copilotTokenData := map[string]interface{}{}
		err = json.Unmarshal(response.Body(), &copilotTokenData)
		if err != nil {
			// 处理JSON解析错误
			return false
		}
		// cache token
		copilotTokens[token] = copilotTokenData
		return true
	} else {
		// 处理其他状态码
		delete(validGithubTokens, token)
		return false
	}
}

// 获取copilot token
func getGithubCopilotToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestCount++
		if err := authentication(c); err != nil {
			badRequest(c)
			return
		}
		token := getRandomGithubToken(validGithubTokens)
		// 从map中获取github token对应的copilot token
		if copilotToken, exists := copilotTokens[token]; exists {
			if !isTokenExpired(copilotToken) {
				proxyResponse(c, copilotToken)
				return
			}
		}
		if getGithubApi(token) {
			proxyResponse(c, copilotTokens[token])
		} else {
			badRequest(c)
		}
	}
}

func getState() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":      version,
			"valid_tokens": len(validGithubTokens),
			"succeed":      successCount,
			"failed":       requestCount - successCount,
			"github_api":   githubApiCount,
		})
	}
}

// 验证请求代理请求token
func authentication(c *gin.Context) error {
	if config.Secret != "" {
		token := c.GetHeader("Authorization")
		tokenStr := strings.ReplaceAll(token, " ", "")
		configSecret := strings.ReplaceAll(config.Secret, " ", "")
		if tokenStr != "token"+configSecret {
			return errors.New("unauthorized")
		}
	}
	return nil
}

// 检测copilot token是否过期
func isTokenExpired(copilotToken map[string]interface{}) bool {
	if expiresAt, ok := copilotToken["expires_at"].(float64); ok {
		currentTime := time.Now().Unix()
		expiresAtInt64 := int64(expiresAt)
		return expiresAtInt64 <= currentTime+60
	}
	return true
}

// 重置请求计数
func resetRequestCount() {
	requestCountMutex.Lock()
	defer requestCountMutex.Unlock()
	requestCount = 0
	successCount = 0
}

// 从map中随机获取一个github token
func getRandomGithubToken(m map[string]bool) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return "" // 返回空字符串或处理其他错误情况
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(keys))
	return keys[randomIndex]
}
