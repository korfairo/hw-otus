package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	directory, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read directory")
	}

	envMap := make(Environment, len(directory))

	for _, entry := range directory {
		env, val, err := processDirEntry(dir, entry)
		if err != nil {
			fmt.Println("failed to process directory entry:", err)
			continue
		}

		envMap[env] = val
	}

	return envMap, nil
}

func processDirEntry(dirPath string, entry os.DirEntry) (env string, value EnvValue, err error) {
	fi, err := entry.Info()
	if err != nil {
		return "", EnvValue{}, errors.Wrap(err, "failed to get dir entry info")
	}

	if fi.IsDir() {
		return "", EnvValue{}, errors.Wrap(err, "subdirectory processing not supported")
	}

	fileName := fi.Name()
	if strings.Contains(fileName, "=") {
		return "", EnvValue{}, errors.New("file name contains '='")
	}

	fileSize := fi.Size()
	if fileSize == 0 {
		return fileName,
			EnvValue{Value: "", NeedRemove: true},
			nil
	}

	filePath := fmt.Sprintf("%s/%s", dirPath, fileName)
	firstLine, err := readFileFirstLine(filePath)
	if err != nil {
		return "", EnvValue{}, errors.Wrap(err, "failed to get first line in file")
	}

	envVal := strings.TrimRight(firstLine, " ")
	envVal = strings.Replace(envVal, "\x00", "\n", -1)

	return fileName,
		EnvValue{Value: envVal, NeedRemove: false},
		nil
}

func readFileFirstLine(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if scanner.Err() != nil {
		return "", errors.Wrap(err, "failed to read first line")
	}

	return scanner.Text(), nil
}
