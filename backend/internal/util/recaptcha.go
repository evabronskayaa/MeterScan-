package util

import (
	errors2 "backend/internal/errors"
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"time"
)

const ErrReCaptchaFailed errors2.SimpleError = "Ошибка рекапчти"

type ReCaptcha struct {
	Timeout time.Duration
	Secret  string
}

type Response struct {
	Success bool `json:"success"`

	Score  float64 `json:"score"`
	Action string  `json:"action"`

	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

// Verify https://developers.google.com/recaptcha/docs/verify
func (reCaptcha ReCaptcha) Verify(userResponse, ip string) error {
	if ip == "::1" || reCaptcha.Secret == "secret" {
		return nil
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/x-www-form-urlencoded")

	req.SetRequestURI("https://www.google.com/recaptcha/api/siteverify")

	args := req.PostArgs()
	args.Set("secret", reCaptcha.Secret)
	args.Set("response", userResponse)

	if ip != "" {
		args.Set("remoteip", ip)
	}

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	if err := fasthttp.DoTimeout(req, res, reCaptcha.Timeout); err != nil {
		if errors.Is(err, fasthttp.ErrTimeout) {
			return nil
		}
		return ErrReCaptchaFailed
	}

	response := &Response{}
	if err := json.Unmarshal(res.Body(), response); err != nil {
		return ErrReCaptchaFailed
	}

	if response.Success && response.Action == "" || response.Score >= 0.5 {
		return nil
	}
	return ErrReCaptchaFailed
}
