package req

import (
	"fmt"
	"github.com/shopspring/decimal"
	"mime/multipart"
	"net/http"
	"strconv"
)

type MultipartFormData struct {
	FormValues map[string]string
	Files      map[string][]*multipart.FileHeader
}

func ParseMultipartForm(r *http.Request, maxMemory int64) (*MultipartFormData, error) {
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		return nil, err
	}

	formData := &MultipartFormData{
		FormValues: make(map[string]string),
		Files:      make(map[string][]*multipart.FileHeader),
	}

	// Получаем текстовые поля
	for key, values := range r.MultipartForm.Value {
		if len(values) > 0 {
			formData.FormValues[key] = values[0]
		}
	}

	// Получаем файлы (все файлы с одним ключом)
	for key, files := range r.MultipartForm.File {
		formData.Files[key] = files
	}

	return formData, nil
}

func GetIntFromForm(formData *MultipartFormData, key string) (int, error) {
	value, exists := formData.FormValues[key]
	if !exists {
		return 0, fmt.Errorf("field %s is required", key)
	}
	return strconv.Atoi(value)
}

func GetDecimalFromForm(formData *MultipartFormData, key string) (decimal.Decimal, error) {
	value, exists := formData.FormValues[key]
	if !exists {
		return decimal.Zero, fmt.Errorf("field %s is required", key)
	}

	d, err := decimal.NewFromString(value)
	if err != nil {
		return decimal.Zero, fmt.Errorf("invalid decimal in field %s: %w", key, err)
	}

	return d, nil
}

// Для получения одного файла
func GetSingleFileFromForm(formData *MultipartFormData, key string) (*multipart.FileHeader, error) {
	files, exists := formData.Files[key]
	if !exists || len(files) == 0 {
		return nil, fmt.Errorf("file %s is required", key)
	}
	return files[0], nil
}
