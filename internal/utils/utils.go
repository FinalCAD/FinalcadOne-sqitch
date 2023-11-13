package utils

import (
	"os"
	"path/filepath"
	"strconv"
)

func CreateFile(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func GetenvBool(key string, fallback bool) (bool, error) {
	s := Getenv(key, "")
	if len(s) == 0 {
		return fallback, nil
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}
