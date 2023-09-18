package core

import (
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"time"
)

// 控制台显示信息
func console() {
	fmt.Println(color.WhiteString("-----------------------------------------------------------------------"))
	fmt.Println(color.HiBlueString(banner))
	fmt.Println(color.WhiteString("[Version]") + " " + version)
	fmt.Println(color.WhiteString("-----------------------------------------------------------------------"))
	var url = "http://127.0.0.1:" + strconv.Itoa(config.Server.Port)
	var jetStr = color.WhiteString("[Jetbrains Endpoint]")
	var vscStr = color.WhiteString("[Vscode    Endpoint]")
	var valid = color.WhiteString("[Valid       Tokens]")
	fmt.Println(jetStr + ": " + color.HiBlueString(url+"/copilot_internal/v2/token"))
	fmt.Println(vscStr + ": " + color.HiBlueString(url))
	fmt.Println(valid + ": " + color.HiBlueString(strconv.Itoa(len(validGithubTokens))))
	fmt.Println(color.WhiteString("-----------------------------------------------------------------------"))
	for {
		requestCountMutex.Lock()
		sCount := successCount
		tCount := requestCount
		gCount := githubApiCount
		requestCountMutex.Unlock()
		currentDateTime := time.Now().Format("2006-01-02 15:04:05")
		currentTime := time.Now().Format("15:04:05")
		if "00:00:00" == currentTime {
			// 每天凌晨重置
			resetRequestCount()
		}
		var s2 = color.WhiteString("[Succeed]")
		var s3 = color.WhiteString("[Failed]")
		var s4 = color.WhiteString("[GithubApi]")
		// 打印文本
		fmt.Printf("\033[G%s  -  %s: %s    %s: %s    %s: %s  ",
			color.HiYellowString(currentDateTime),
			s2, color.GreenString(strconv.Itoa(sCount)),
			s3, color.RedString(strconv.Itoa(tCount-sCount)),
			s4, color.CyanString(strconv.Itoa(gCount)))
		time.Sleep(1 * time.Second) //
	}
}
