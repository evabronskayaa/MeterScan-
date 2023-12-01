package media

import (
	"backend/internal/errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
)

const (
	ErrCreateDir  errors.SimpleError = "Произошла ошибка при создании директории"
	ErrCreateFile errors.SimpleError = "Произошла ошибка при создании файла"
	ErrSaveFile   errors.SimpleError = "Произошла ошибка при сохранении файла"
)

func GetPath(path, file string) (string, error) {
	path = fmt.Sprintf("./media/%v", path)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", ErrCreateDir
	}

	return filepath.Join(path, file), nil
}

func SaveFile(path, fileName string, file multipart.File) (string, error) {
	dstPath, err := GetPath(path, fileName)
	if err != nil {
		return "", err
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return "", ErrCreateFile
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", ErrSaveFile
	}

	return dstPath, nil
}

func SaveData(path, fileName string, file []byte) (string, error) {
	dstPath, err := GetPath(path, fileName)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(dstPath, file, fs.ModePerm)
	if err != nil {
		return "", ErrSaveFile
	}

	return dstPath, nil
}
