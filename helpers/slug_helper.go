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


func Slugify(s string) string {
	s = strings.ToLower(s) // 1. Jadi huruf kecil
	s = nonAlphanumericRegex.ReplaceAllString(s, "") // 2. Hapus karakter aneh
	s = strings.ReplaceAll(s, " ", "-") // 3. Ganti spasi dengan strip
	return s
}

func FindUniqueFilename(slug string, ext string) string {
	i := 1
	for {
		
		filename := fmt.Sprintf("%s-%d%s", slug, i, ext)
		
		filePath := filepath.Join("assets", filename)

		
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			
			return filename
		}
		i++
	}
}