package req

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
)

type MultipartFormData struct {
	FormValues map[string]string
	Files      map[string]*multipart.FileHeader
}

func ParseMultipartForm(r *http.Request, maxMemory int64) (*MultipartFormData, error) {
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		return nil, err
	}

	formData := &MultipartFormData{
		FormValues: make(map[string]string),
		Files:      make(map[string]*multipart.FileHeader),
	}

	// Получаем текстовые поля
	for key, values := range r.MultipartForm.Value {
		if len(values) > 0 {
			formData.FormValues[key] = values[0]
		}
	}

	// Получаем файлы
	for key, files := range r.MultipartForm.File {
		if len(files) > 0 {
			formData.Files[key] = files[0]
		}
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
