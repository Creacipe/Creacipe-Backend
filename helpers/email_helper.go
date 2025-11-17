package helpers

import (
	"crypto/rand"
	"fmt"
	"math/big"
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

// DefaultEmailConfig adalah konfigurasi email default
// 
// SETUP BREVO (RECOMMENDED - GRATIS 300 EMAIL/HARI SELAMANYA):
// 1. Daftar: https://app.brevo.com/account/register
// 2. Verifikasi email Anda
// 3. Login ke dashboard Brevo
// 4. Menu: Settings → SMTP & API → SMTP
// 5. Klik "Generate a new SMTP key" atau "Create a new SMTP key"
// 6. Copy SMTP key yang muncul (format: xsmtpsib-xxxxxxxxxxxxx)
// 7. Paste ke SMTPPassword di bawah
// 8. SMTPUsername = email yang digunakan untuk login Brevo
//
// KEUNGGULAN BREVO:
// ✅ 300 email/hari (9000 email/bulan) - GRATIS SELAMANYA
// ✅ Dashboard analytics untuk tracking email
// ✅ Support custom domain (bisa pakai noreply@yourdomain.com)
// ✅ Reliable untuk production
// ✅ Tidak perlu kartu kredit
//
// ALTERNATIF GRATIS LAINNYA:
// - Gmail App Password (untuk development/testing)
// - SMTP2GO (1000 email/bulan gratis)
// - Resend (3000 email/bulan gratis)
// - Mailgun (1000 email/bulan gratis setelah trial)
//
// TODO: Pindahkan ke environment variables untuk production
var DefaultEmailConfig = EmailConfig{
	SMTPHost:     "smtp-relay.brevo.com",              // SMTP Brevo
	SMTPPort:     587,                                 // Port TLS
	SMTPUsername: "9bbc43001@smtp-brevo.com",       // Email login Brevo Anda
	SMTPPassword: "xsmtpsib-d4817aa3809f71457a121b456db7cd23c9c00f09ac5b3288b20240dc3362c8d6-1sBqGHRSblPJ1ZII",           // SMTP Key dari Brevo dashboard
	FromEmail:    "stoksolo01@gmail.com",           // Email pengirim (bisa custom atau pakai email Brevo)
	FromName:     "Creacipe",                         // Nama pengirim yang muncul di email
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
	config := DefaultEmailConfig

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
