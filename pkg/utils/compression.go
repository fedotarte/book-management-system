package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"os"

	"github.com/nfnt/resize"
)

// CompressAndSaveImage сжимает изображение перед сохранением
func CompressAndSaveImage(inputFile image.Image, outputPath string, quality int) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer out.Close()

	err = jpeg.Encode(out, inputFile, &jpeg.Options{Quality: quality})
	if err != nil {
		return fmt.Errorf("ошибка кодирования изображения: %w", err)
	}

	return nil
}

// DecodeImage декодирует изображение из multipart.File
func DecodeImage(file multipart.File, format string) (image.Image, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования изображения: %w", err)
	}

	// Если PNG, конвертируем в JPEG
	if format == "png" {
		img = resize.Resize(600, 900, img, resize.Lanczos3)
	}

	return img, nil
}
