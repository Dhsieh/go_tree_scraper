package utils

import (
	"net/http"
	"strings"
)

var magicTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

// Check the []byte to see if the contents are JPEG or not
// Does 2 different types of checks of the header
func IsJPEG(incipit []byte) bool {
	iniciptStr := string(incipit)
	jpgMime := magicTable["jpeg"]
	if strings.HasPrefix(iniciptStr, jpgMime) && http.DetectContentType(incipit[:512]) == "image/jpeg" {
		return true
	} else {
		return false
	}
}
