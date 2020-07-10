package validators

import (
	"api/util/stringmaker"
	"errors"
	"net/url"
	"strings"
)

var validationErrorMap = map[string]error{
	"InvalidLongURL":        errors.New("Error, the provided long url is not valid: %v \n The URL must have a Scheme, Host, and/or a Path."),
	"InvalidLengthShortURL": errors.New("Error, the provided short url is not 8 chars long: %v \n Current API only accepts 8 char long shortURLs."),
	"InvalidCharsShortURL":  errors.New("Error, the provided short url has an invalid character: %v \n Current API only accepts: " + stringmaker.GetCurrentCharset()),
}

// Check if the URL has a Scheme, Host, and/or a Path.
func IsUrl(str string) (bool, error) {
	u, err := url.Parse(str)
	if err == nil && u.Scheme != "" && u.Host != "" {
		return true, nil
	}
	return false, validationErrorMap["InvalidLongURL"]
}

// Check if a shortURL is 8 characters long.
func IsShort8(str string) (bool, error) {
	if len(str) == 8 {
		return true, nil
	}
	return false, validationErrorMap["InvalidLengthShortURL"]
}

// Check if shortURL has valid chars:
func IsShortValidChars(str string) (bool, error) {
	for _, char := range str {
		if strings.ContainsRune(stringmaker.GetCurrentCharset(), char) {
			continue
		} else {
			return false, validationErrorMap["InvalidCharsShortURL"]
		}
	}
	return true, nil
}
