package helpers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"time"
)

type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func VerifyRecaptcha(token string) error {
	secretKey := os.Getenv("RECAPTCHA_SECRET_KEY")
	if secretKey == "" {
		return errors.New("recaptcha secret key not found")
	}

	
	if secretKey == "dummy-secret" {
		return nil
	}

	// URL API Google ReCaptcha
	apiURL := "https://www.google.com/recaptcha/api/siteverify"

	// Data yang dikirim
	data := url.Values{}
	data.Set("secret", secretKey)
	data.Set("response", token)

	// Kirim POST Request ke Google
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response JSON dari Google
	var result RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Cek status sukses
	if !result.Success {
		return errors.New("verifikasi captcha gagal, silakan coba lagi")
	}

	return nil
}