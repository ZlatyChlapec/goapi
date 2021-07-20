package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Tokens []authToken
}

func NewConfig() (config *Config) {
	config = &Config{}
	err := json.Unmarshal([]byte(os.Getenv("TOKENS")), &config.Tokens)
	if err != nil {
		panic(fmt.Sprintf("NewConfig, %v", err))
	}
	// We are associating images to usernames so we can't have two with same name.
	checkDuplicits(&config.Tokens)
	return
}

func checkDuplicits(tokens *[]authToken) {
	unique := make(map[string]interface{})
	for _, token := range *tokens {
		if _, ok := unique[token.User]; ok {
			panic(fmt.Sprintf("checkDuplicits.duplicitUser, %v", token.User))
		} else {
			unique[token.User] = nil
		}
	}
}