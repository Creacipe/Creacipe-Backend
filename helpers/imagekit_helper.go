// LOKASI: helpers/imagekit_helper.go

package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ImageKitResponse adalah struktur response dari ImageKit API
type ImageKitResponse struct {
	FileID       string `json:"fileId"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
	Height       int    `json:"height"`
	Width        int    `json:"width"`
	Size         int    `json:"size"`
	FilePath     string `json:"filePath"`
	FileType     string `json:"fileType"`
}

// ImageKitErrorResponse adalah struktur error response dari ImageKit
type ImageKitErrorResponse struct {
	Message string `json:"message"`
	Help    string `json:"help"`
}

// UploadToImageKit mengunggah file ke ImageKit dan mengembalikan URL
// folder: "menus" untuk gambar resep, "profiles" untuk foto profil
func UploadToImageKit(file multipart.File, fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Ambil kredensial dari environment
	privateKey := os.Getenv("IMAGEKIT_PRIVATE_KEY")
	publicKey := os.Getenv("IMAGEKIT_PUBLIC_KEY")
	urlEndpoint := os.Getenv("IMAGEKIT_URL_ENDPOINT")

	if privateKey == "" || publicKey == "" || urlEndpoint == "" {
		return "", fmt.Errorf("imagekit credentials not configured in environment")
	}

	// Baca file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	
	base64File := base64.StdEncoding.EncodeToString(fileBytes)

	// Buat nama file unik
	ext := filepath.Ext(fileHeader.Filename)
	baseName := strings.TrimSuffix(fileHeader.Filename, ext)
	uniqueName := fmt.Sprintf("%s_%d%s", Slugify(baseName), time.Now().UnixNano(), ext)

	// Buat multipart form untuk upload
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Tambah field file 
	writer.WriteField("file", base64File)
	// Nama file
	writer.WriteField("fileName", uniqueName)
	// Folder tujuan
	writer.WriteField("folder", "/creacipe/"+folder)
	// Gunakan unique filename
	writer.WriteField("useUniqueFileName", "true")

	writer.Close()

	// Buat HTTP request
	req, err := http.NewRequest("POST", "https://upload.imagekit.io/api/v1/files/upload", &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	
	req.Header.Set("Content-Type", writer.FormDataContentType())

	
	auth := base64.StdEncoding.EncodeToString([]byte(privateKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)

	// Kirim request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload to imagekit: %v", err)
	}
	defer resp.Body.Close()

	// Baca response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Cek status code
	if resp.StatusCode != http.StatusOK {
		var errResp ImageKitErrorResponse
		json.Unmarshal(respBody, &errResp)
		return "", fmt.Errorf("imagekit upload failed: %s", errResp.Message)
	}

	// Parse response
	var ikResp ImageKitResponse
	if err := json.Unmarshal(respBody, &ikResp); err != nil {
		return "", fmt.Errorf("failed to parse imagekit response: %v", err)
	}

	return ikResp.URL, nil
}

// UploadToImageKitFromBytes mengunggah file dari bytes ke ImageKit
func UploadToImageKitFromBytes(fileBytes []byte, filename string, folder string) (string, error) {
	// Ambil kredensial dari environment
	privateKey := os.Getenv("IMAGEKIT_PRIVATE_KEY")
	urlEndpoint := os.Getenv("IMAGEKIT_URL_ENDPOINT")

	if privateKey == "" || urlEndpoint == "" {
		return "", fmt.Errorf("imagekit credentials not configured in environment")
	}

	
	base64File := base64.StdEncoding.EncodeToString(fileBytes)

	
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)
	uniqueName := fmt.Sprintf("%s_%d%s", Slugify(baseName), time.Now().UnixNano(), ext)

	
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	writer.WriteField("file", base64File)
	writer.WriteField("fileName", uniqueName)
	writer.WriteField("folder", "/creacipe/"+folder)
	writer.WriteField("useUniqueFileName", "true")

	writer.Close()

	// Buat HTTP request
	req, err := http.NewRequest("POST", "https://upload.imagekit.io/api/v1/files/upload", &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	auth := base64.StdEncoding.EncodeToString([]byte(privateKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload to imagekit: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ImageKitErrorResponse
		json.Unmarshal(respBody, &errResp)
		return "", fmt.Errorf("imagekit upload failed: %s", errResp.Message)
	}

	var ikResp ImageKitResponse
	if err := json.Unmarshal(respBody, &ikResp); err != nil {
		return "", fmt.Errorf("failed to parse imagekit response: %v", err)
	}

	return ikResp.URL, nil
}
