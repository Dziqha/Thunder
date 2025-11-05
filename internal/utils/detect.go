package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// DetectMainFile mencari file main.go yang valid (berisi package main dan func main)
func DetectMainFile(root string) (string, error) {
	var mainFiles []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// skip direktori yang umum di-skip
			if info.Name() == "vendor" || info.Name() == "tmp" || info.Name() == ".git" || info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(info.Name(), ".go") {
			file, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			hasMainPackage := false
			hasMainFunc := false

			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if strings.HasPrefix(line, "package main") {
					hasMainPackage = true
				}
				if strings.HasPrefix(line, "func main()") {
					hasMainFunc = true
				}
				if hasMainPackage && hasMainFunc {
					mainFiles = append(mainFiles, path)
					break
				}
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if len(mainFiles) == 0 {
		return "", os.ErrNotExist
	}

	// kalau lebih dari satu, ambil yang pertama
	return mainFiles[0], nil
}
