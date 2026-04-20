package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// Se define Interfaz de Storage
type FileStorage interface {
	UploadFile(ctx context.Context, bucket string, path string, fileHeader *multipart.FileHeader) (string, error)
	DeleteFile(ctx context.Context, bucket string, path string) error
	CreateSignedURL(ctx context.Context, bucket string, path string, expiresIn int) (string, error)
}

type supabaseStorage struct {
	url    string
	apiKey string
}

func NewSupabaseStorage(url, apiKey string) FileStorage {
	return &supabaseStorage{
		url:    url,
		apiKey: apiKey,
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

// DeleteFile implements [services.FileStorage].
func (s *supabaseStorage) DeleteFile(ctx context.Context, bucket string, path string) error {
	panic("unimplemented")
}

// CreateSignedURL implements [services.FileStorage].
func (s *supabaseStorage) CreateSignedURL(ctx context.Context, bucket string, path string, expiresIn int) (string, error) {
	panic("unimplemented")
}
