package helpers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os" // <--- Import ini WAJIB ada untuk membaca .env
	"time"

	"gopkg.in/gomail.v2"
)

// EmailConfig menyimpan konfigurasi SMTP
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// Konfigurasi tetap disini, KECUALI Password yang kita kosongkan/isi dummy
var DefaultEmailConfig = EmailConfig{
	SMTPHost:     "smtp-relay.brevo.com",
	SMTPPort:     587,
	SMTPUsername: "9bbc43001@smtp-brevo.com",
	SMTPPassword: "", // <--- Dikosongkan disini, nanti diisi otomatis dari .env
	FromEmail:    "stoksolo01@gmail.com",
	FromName:     "Creacipe",
}

// GenerateVerificationCode menghasilkan kode verifikasi 6 digit
func GenerateVerificationCode() string {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		// Fallback jika crypto/rand gagal
		return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	}
	return fmt.Sprintf("%06d", n.Int64())
}

// SendVerificationEmail mengirim email verifikasi
func SendVerificationEmail(toEmail, toName, verificationCode, purpose string) error {
	// Salin konfigurasi default
	config := DefaultEmailConfig

	// --- PERUBAHAN DISINI ---
	// Ambil password dari .env, bukan dari variabel global di atas
	config.SMTPPassword = os.Getenv("SMTPPassword")

	// Pengecekan keamanan sederhana (opsional)
	if config.SMTPPassword == "" {
		return fmt.Errorf("SMTPPassword tidak ditemukan di .env")
	}
	// ------------------------

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", config.FromName, config.FromEmail))
	m.SetHeader("To", toEmail)

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

	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("gagal mengirim email: %v", err)
	}

	return nil
}