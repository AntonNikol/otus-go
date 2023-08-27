package main

import (
	"github.com/cheggaaa/pb"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
)

var (
	ErrNotFound              = errors.New("file not found")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// Copy копирует данные из исходного файла в целевой файл с учетом смещения и ограничения.
func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := openFile(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("file %v close error: %v", fromPath, err)
		}
	}()

	fileSize, err := getFileSize(file)
	if err != nil {
		return err
	}

	// Если неизвестна длина (например, файл /dev/urandom)
	if fileSize == 0 {
		return ErrUnsupportedFile
	}

	err = validateOffset(fileSize, offset)
	if err != nil {
		return err
	}

	bytesToCopy := fileSize - offset
	if limit > 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	_, err = file.Seek(offset, 0)
	if err != nil {
		return errors.Wrap(err, "unable to set offset in file")
	}

	tmpFile, err := createTempFile()
	if err != nil {
		return err
	}

	err = copyDataWithProgress(file, tmpFile, bytesToCopy)
	if err != nil {
		return err
	}

	err = os.Rename(tmpFile.Name(), toPath)
	if err != nil {
		return errors.Wrap(err, "rename error")
	}

	return nil
}

// openFile открывает файл по указанному пути и возвращает его указатель
func openFile(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 444)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "file open error")
	}
	return file, nil
}

// getFileSize получает размер файла
func getFileSize(file *os.File) (int64, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, errors.Wrap(err, "unable to load file info")
	}
	return fileInfo.Size(), nil
}

// validateOffset проверяет, что смещение (offset) не превышает размер файла.
func validateOffset(fileSize, offset int64) error {
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}
	return nil
}

// createTempFile создает временный файл и возвращает его указатель.
func createTempFile() (*os.File, error) {
	tmpFile, err := os.CreateTemp(".", "tmp_file_")
	if err != nil {
		return nil, errors.Wrap(err, "unable to create tmp file")
	}
	return tmpFile, nil
}

// copyDataWithProgress копирует данные из reader в writer с отображением прогресса.
func copyDataWithProgress(reader io.Reader, writer io.Writer, bytesToCopy int64) error {
	bar := pb.New64(bytesToCopy)
	bar.Start()
	barReader := bar.NewProxyReader(reader)

	_, err := io.CopyN(writer, barReader, bytesToCopy)
	if err != nil && err != io.EOF {
		return errors.Wrap(err, "unable to copy data")
	}

	bar.Finish()
	return nil
}
