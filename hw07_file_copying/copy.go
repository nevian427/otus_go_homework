package main

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOffsetInvalid         = errors.New("non-negative offset expected on non-regular input")
	ErrLimitNegative         = errors.New("non-negative limit expected")
	ErrLimitInvalid          = errors.New("expected limit on non-regular input")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if limit < 0 {
		return ErrLimitNegative
	}

	inFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	inInfo, err := inFile.Stat()
	if err != nil {
		return err
	}

	if !inInfo.Mode().IsRegular() {
		if limit < 1 {
			return ErrLimitInvalid
		}
		if offset < 0 {
			return ErrOffsetInvalid
		}
	} else if limit == 0 {
		limit = inInfo.Size()
	}

	if offset != 0 {
		absoffset := offset
		whence := io.SeekStart
		switch { // такой вот извращённый обход линтера, ругающегося на вложенные IF - кажется было бы понятнее с ними
		case offset < 0:
			whence = io.SeekEnd
			absoffset = -offset

			if limit > absoffset {
				limit = absoffset
			}
			fallthrough
		default:
			if inInfo.Mode().IsRegular() && inInfo.Size() < absoffset {
				return ErrOffsetExceedsFileSize
			}

			if _, err := inFile.Seek(offset, whence); err != nil {
				return err
			}
		}
	}

	outFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// create bar
	bar := pb.Start64(limit).SetRefreshRate(time.Second)

	// create proxy reader
	reader := bar.NewProxyReader(inFile)

	// and copy from reader
	_, err = io.CopyN(outFile, reader, limit)
	bar.Finish()

	if errors.Is(err, io.EOF) {
		return nil
	}

	return err
}
