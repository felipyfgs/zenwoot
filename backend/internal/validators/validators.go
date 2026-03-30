package validators

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) (bool, string) {
	if len(password) < 6 {
		return false, "password must be at least 6 characters"
	}
	if len(password) > 128 {
		return false, "password must be less than 128 characters"
	}
	return true, ""
}

func ValidateRequired(value, fieldName string) (bool, string) {
	if strings.TrimSpace(value) == "" {
		return false, fieldName + " is required"
	}
	return true, ""
}

func ValidatePhoneNumber(phone string) bool {
	if phone == "" {
		return false
	}
	return phoneRegex.MatchString(phone)
}

func ValidatePageParams(page, limit int) (int, int, bool) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit, true
}

func ValidateFileType(filename, contentType string) (bool, string) {
	allowedImages := []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
	allowedDocs := []string{"application/pdf", "application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document"}
	allowed := append(allowedImages, allowedDocs...)

	for _, allowedType := range allowed {
		if contentType == allowedType {
			return true, ""
		}
	}
	return false, "file type not allowed"
}

func ValidateMaxFileSize(size int64, maxMB int) (bool, string) {
	maxBytes := int64(maxMB * 1024 * 1024)
	if size > maxBytes {
		return false, "file size exceeds maximum of " + string(rune(maxMB+'0')) + "MB"
	}
	return true, ""
}
