package main

import (
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported or empty file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	inputFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err != nil {
		return errors.Wrap(err, "failed to open input file")
	}
	defer func() {
		err = inputFile.Close()
		if err != nil {
			fmt.Println("failed to close input file: ", err)
		}
	}()

	inputFileInfo, err := inputFile.Stat()
	if err != nil {
		return errors.Wrap(err, "failed to get input file stats")
	}

	if inputFileInfo.Size() == 0 {
		return ErrUnsupportedFile
	}

	if offset > inputFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	_, err = inputFile.Seek(offset, 0)
	if err != nil {
		return errors.Wrap(err, "failed to move offset in input file")
	}

	outputFile, err := os.Create(toPath)
	if err != nil {
		return errors.Wrap(err, "failed to create output file")
	}
	defer func() {
		err = outputFile.Close()
		if err != nil {
			fmt.Println("failed to close output file: ", err)
		}
	}()

	if limit <= 0 {
		limit = inputFileInfo.Size() - offset
	}

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(inputFile)
	defer bar.Finish()

	_, err = io.CopyN(outputFile, barReader, limit)
	if err != nil && err != io.EOF {
		return errors.Wrap(err, "failed to copy data")
	}

	return nil
}
