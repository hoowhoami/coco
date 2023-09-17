package core

import "sync"

var (
	copilotTokens     = make(map[string]map[string]interface{})
	requestCountMutex sync.Mutex
	validGithubTokens = make(map[string]bool)
	githubApiCount    = 0
	requestCount      = 0
	successCount      = 0
	config            Config
	version           = "1.0.0"
	banner            = `
                           _ _       _   
  ___ ___   ___ ___  _ __ (_) | ___ | |_ 
 / __/ _ \ / __/ _ \| '_ \| | |/ _ \| __|
| (_| (_) | (_| (_) | |_) | | | (_) | |_ 
 \___\___/ \___\___/| .__/|_|_|\___/ \__|
                    |_|
`
)
