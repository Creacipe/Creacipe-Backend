package helpers

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"time"
)

// EmailConfig menyimpan konfigurasi email (Brevo HTTP API)
type EmailConfig struct {
	APIKey    string
	FromEmail string
	FromName  string
}

// brevoEmailRequest adalah struktur request body untuk Brevo API
type brevoEmailRequest struct {
	Sender      brevoContact   `json:"sender"`
	To          []brevoContact `json:"to"`
	Subject     string         `json:"subject"`
	HTMLContent string         `json:"htmlContent"`
}

type brevoContact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetEmailConfig membaca konfigurasi email dari environment variables
func GetEmailConfig() EmailConfig {
	return EmailConfig{
		APIKey:    os.Getenv("BREVO_API_KEY"),
		FromEmail: os.Getenv("SMTP_FROM_EMAIL"),
		FromName:  os.Getenv("SMTP_FROM_NAME"),
	}
}

// GenerateVerificationCode menghasilkan kode verifikasi 6 digit
func GenerateVerificationCode() string {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		
		return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	}
	return fmt.Sprintf("%06d", n.Int64())
}

// SendVerificationEmail mengirim email verifikasi via Brevo HTTP API
func SendVerificationEmail(toEmail, toName, verificationCode, purpose string) error {
	config := GetEmailConfig()

	if config.APIKey == "test-key" {
		return nil
	}

	var subject, htmlBody string

	switch purpose {
	case "reset_password":
		subject = "Kode Verifikasi Reset Password - Creacipe"
		htmlBody = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
        .code-box { background: white; border: 2px dashed #667eea; padding: 20px; text-align: center; margin: 20px 0; border-radius: 8px; }
        .code { font-size: 32px; font-weight: bold; color: #667eea; letter-spacing: 5px; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Reset Password</h1>
        </div>
        <div class="content">
            <p>Halo <strong>%s</strong>,</p>
            <p>Anda telah meminta untuk mereset password akun Creacipe Anda. Gunakan kode verifikasi berikut:</p>
            <div class="code-box">
                <div class="code">%s</div>
            </div>
            <p><strong>Kode ini berlaku selama 10 menit.</strong></p>
            <p>Jika Anda tidak meminta reset password, abaikan email ini.</p>
            <div class="footer">
                <p>Email ini dikirim secara otomatis, mohon tidak membalas.</p>
                <p>&copy; 2025 Creacipe. All rights reserved.</p>
            </div>
        </div>
    </div>
</body>
</html>
`, toName, verificationCode)

	case "change_email":
		subject = "Kode Verifikasi Ubah Email - Creacipe"
		htmlBody = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
        .code-box { background: white; border: 2px dashed #667eea; padding: 20px; text-align: center; margin: 20px 0; border-radius: 8px; }
        .code { font-size: 32px; font-weight: bold; color: #667eea; letter-spacing: 5px; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Ubah Email</h1>
        </div>
        <div class="content">
            <p>Halo <strong>%s</strong>,</p>
            <p>Anda telah meminta untuk mengubah email akun Creacipe Anda. Gunakan kode verifikasi berikut:</p>
            <div class="code-box">
                <div class="code">%s</div>
            </div>
            <p><strong>Kode ini berlaku selama 10 menit.</strong></p>
            <p>Jika Anda tidak meminta perubahan email, segera hubungi kami.</p>
            <div class="footer">
                <p>Email ini dikirim secara otomatis, mohon tidak membalas.</p>
                <p>&copy; 2025 Creacipe. All rights reserved.</p>
            </div>
        </div>
    </div>
</body>
</html>
`, toName, verificationCode)

	default:
		return fmt.Errorf("purpose tidak valid: %s", purpose)
	}

	// Kirim email via Brevo HTTP API
	reqBody := brevoEmailRequest{
		Sender: brevoContact{
			Name:  config.FromName,
			Email: config.FromEmail,
		},
		To: []brevoContact{
			{Name: toName, Email: toEmail},
		},
		Subject:     subject,
		HTMLContent: htmlBody,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("gagal membuat request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("gagal membuat HTTP request: %v", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("api-key", config.APIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("gagal mengirim email: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gagal mengirim email (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
