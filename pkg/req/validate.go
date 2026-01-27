package req

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"mime/multipart"
	"strings"
	"sync"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// Инициализация валидатора с кастомными правилами (вызывается один раз)
func initValidator() {
	once.Do(func() {
		validate = validator.New()

		// Регистрируем кастомные валидаторы
		err := validate.RegisterValidation("image_type", validateImageType)
		if err != nil {
			return
		}
		// Можно добавить другие кастомные валидаторы здесь
	})
}

// Валидатор для проверки типа изображения
func validateImageType(fl validator.FieldLevel) bool {
	file, ok := fl.Field().Interface().(*multipart.FileHeader)
	if !ok || file == nil {
		return false
	}

	// Проверка Content-Type
	contentType := file.Header.Get("Content-Type")
	allowedTypes := []string{"image/png", "image/jpeg", "image/jpg"}

	for _, t := range allowedTypes {
		if contentType == t {
			return true
		}
	}

	// Дополнительная проверка расширения файла
	filename := strings.ToLower(file.Filename)
	return strings.HasSuffix(filename, ".png") ||
		strings.HasSuffix(filename, ".jpg") ||
		strings.HasSuffix(filename, ".jpeg")
}

func IsValid[T any](payload T) error {
	// Инициализируем валидатор (произойдет только один раз)
	initValidator()

	err := validate.Struct(payload)

	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			firstError := validationErrors[0]
			fieldName := strings.ToLower(firstError.Field())

			var message string
			switch firstError.Tag() {
			case "required":
				message = fmt.Sprintf("Field %s is required", fieldName)
			case "email":
				message = fmt.Sprintf("Field %s must be a valid email", fieldName)
			case "min":
				message = fmt.Sprintf("Field %s must be at least %s", fieldName, firstError.Param())
			case "max":
				message = fmt.Sprintf("Field %s must be at most %s", fieldName, firstError.Param())
			case "len":
				message = fmt.Sprintf("Field %s must be %s characters long", fieldName, firstError.Param())
			case "gte":
				message = fmt.Sprintf("Field %s must be greater than or equal to %s", fieldName, firstError.Param())
			case "lte":
				message = fmt.Sprintf("Field %s must be less than or equal to %s", fieldName, firstError.Param())
			case "gt":
				message = fmt.Sprintf("Field %s must be greater than %s", fieldName, firstError.Param())
			case "lt":
				message = fmt.Sprintf("Field %s must be less than %s", fieldName, firstError.Param())
			case "oneof":
				message = fmt.Sprintf("Field %s must be one of [%s]", fieldName, firstError.Param())
			case "image_type":
				message = fmt.Sprintf("Field %s must be a PNG or JPEG image", fieldName)
			default:
				message = fmt.Sprintf("Field %s validation failed on %s", fieldName, firstError.Tag())
			}

			return fmt.Errorf(message)
		}
	}

	return err
}
