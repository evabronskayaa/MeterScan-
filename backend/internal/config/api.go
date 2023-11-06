package config

import (
	"encoding/json"
	"flag"
	"os"
)

type ApiConfig struct {
	Port int `json:"port"`

	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Schema   string `json:"schema"`
	} `json:"database"`

	Mail struct {
		Server   string `json:"server"`
		Port     int    `json:"port"`
		Login    string `json:"login"`
		Password string `json:"password"`
	} `json:"mail"`

	JWTSecret string `json:"jwt_secret"`

	ReCaptchaSecret string `json:"re_captcha_secret"`

	GRPCServer string `json:"grpc_server"`
}

func (c *ApiConfig) Load() error {
	configPath := flag.String("config", "api-config.json", "Web server config file location")
	flag.Parse()

	if file, err := os.ReadFile(*configPath); err != nil {
		return err
	} else if err := json.Unmarshal(file, c); err != nil {
		return err
	} else {
		return nil
	}
}
