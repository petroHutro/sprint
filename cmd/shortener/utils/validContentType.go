package utils

import (
	"errors"
	"strings"
)

func ValidContentType(str string) error {
	contentType := strings.Split(str, "; ")
	if len(contentType) > 0 && contentType[0] == "text/plain" {
		if len(contentType) == 2 {
			validCharsets := map[string]bool{
				"charset=utf-8":        true,
				"charset=iso-8859-1":   true,
				"charset=us-ascii":     true,
				"charset=windows-1251": true,
			}
			if validCharsets[contentType[1]] {
				return nil
			}
		} else if len(contentType) == 1 {
			return nil
		}
	}
	return errors.New("wrong fornat Content Type")
}
