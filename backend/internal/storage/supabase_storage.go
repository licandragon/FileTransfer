package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// Se define Interfaz de Storage
type FileStorage interface {
	UploadFile(ctx context.Context, bucket string, path string, fileHeader *multipart.FileHeader) (string, error)
	DeleteFile(ctx context.Context, bucket string, path string) error
	CreateSignedURL(ctx context.Context, bucket string, path string, expiresIn int) (string, error)
}

type supabaseStorage struct {
	url        string
	apiKey     string
	httpClient *http.Client
}

func NewSupabaseStorage(url, apiKey string) FileStorage {
	return &supabaseStorage{
		url:    url,
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// UploadFile implements [services.FileStorage].
func (s *supabaseStorage) UploadFile(ctx context.Context, bucket string, path string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer file.Close()

	// Convertir a buffer para la petición HTTP
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return "", fmt.Errorf("error al leer contenido del archivo: %w", err)
	}

	fullURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.url, bucket, path)

	fmt.Println("URL:", fullURL)
	fmt.Println("API KEY LENGTH:", len(s.apiKey))

	req, err := http.NewRequestWithContext(ctx, "POST", fullURL, buf)
	if err != nil {
		return "", err
	}

	// Seteamos los headers que Supabase exige
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("apikey", s.apiKey)
	req.Header.Set("Content-Type", fileHeader.Header.Get("Content-Type"))

	fmt.Println(fileHeader.Header.Get("Content-Type"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("supabase storage error (%d): %s", resp.StatusCode, string(body))
	}

	return fullURL, nil
}

// DeleteFile elimina un archivo del bucket.
func (s *supabaseStorage) DeleteFile(ctx context.Context, bucket string, path string) error {
	fullURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.url, bucket, path)

	req, err := http.NewRequestWithContext(ctx, "DELETE", fullURL, nil)
	if err != nil {
		return fmt.Errorf("error creando petición DELETE: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("apikey", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error en petición DELETE: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase storage error al eliminar (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// CreateSignedURL genera una URL firmada temporal.
func (s *supabaseStorage) CreateSignedURL(ctx context.Context, bucket string, path string, expiresIn int) (string, error) {
	fullURL := fmt.Sprintf("%s/storage/v1/object/sign/%s/%s", s.url, bucket, path)

	bodyJSON, err := json.Marshal(map[string]int{
		"expiresIn": expiresIn,
	})
	if err != nil {
		return "", fmt.Errorf("error creando JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fullURL, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return "", fmt.Errorf("error creando petición de firma: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("apikey", s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error en petición de firma: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta de firma: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("supabase storage error al firmar (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		SignedURL string `json:"signedURL"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error decodificando URL firmada: %w", err)
	}

	baseURL := strings.TrimRight(s.url, "/")

	// Asegurar que la URL sea absoluta
	if !strings.HasPrefix(result.SignedURL, "/storage/v1") {
		result.SignedURL = "/storage/v1" + result.SignedURL
	}

	result.SignedURL = baseURL + result.SignedURL

	return result.SignedURL, nil
}
