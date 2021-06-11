package utils

import "strings"

var magicTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

func IsJPEG(incipit []byte) bool {
	iniciptStr := string(incipit)
	jpgMime := magicTable["jpeg"]
	if strings.HasPrefix(iniciptStr, jpgMime) {
		return true
	} else {
		return false
	}
}
