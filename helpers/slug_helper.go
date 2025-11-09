// LOKASI: helpers/slug_helper.go

package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Regex untuk menghapus karakter non-alfanumerik (kecuali spasi)
var nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9 ]+`)

// Slugify mengubah string menjadi format yang aman untuk URL/filename.
func Slugify(s string) string {
	s = strings.ToLower(s) // 1. Jadi huruf kecil
	s = nonAlphanumericRegex.ReplaceAllString(s, "") // 2. Hapus karakter aneh
	s = strings.ReplaceAll(s, " ", "-") // 3. Ganti spasi dengan strip
	return s
}

// FindUniqueFilename menemukan nama file yang unik di folder 'assets'.
// Sesuai logikamu: "slug-1.ext", "slug-2.ext", ...
func FindUniqueFilename(slug string, ext string) string {
	i := 1
	for {
		// Buat nama file (misal: "ayam-goreng-1.jpg")
		filename := fmt.Sprintf("%s-%d%s", slug, i, ext)
		// Tentukan path lengkap (misal: "assets/ayam-goreng-1.jpg")
		filePath := filepath.Join("assets", filename)

		// Cek apakah file di path itu sudah ada
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// Jika file TIDAK ada, kita temukan nama unik!
			return filename
		}
		// Jika file sudah ada, 'i' bertambah dan loop berlanjut
		i++
	}
}