package config

import (
	"encoding/json"
	"os"
)

type ApiConfig struct {
	Port int `json:"port"`

	Database struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
		Schema   string `json:"schema"`
	} `json:"database"`

	JWTSecret string `json:"jwt_secret"`

	ReCaptchaSecret string `json:"re_captcha_secret"`

	GRPCServer string `json:"grpc_server"`
}

func (c *ApiConfig) Load() error {
	if file, err := os.ReadFile("api-config.json"); err != nil {
		return err
	} else if err := json.Unmarshal(file, c); err != nil {
		return err
	} else {
		return nil
	}
}
