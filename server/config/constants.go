package config

const MaxUploadSize = 32 << 20 // 32 MB

var AllowedImageExtensions = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	"gif":  true,
	"bmp":  true,
	"webp": true,
}
