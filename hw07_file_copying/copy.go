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

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Открытие файла
	file, err := os.OpenFile(fromPath, os.O_RDONLY, 444)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return errors.Wrap(err, "file open error")
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("file %v close error: %v", fromPath, err)
		}
	}()

	// Данные о файле
	fileInfo, err := file.Stat()
	if err != nil {
		return errors.Wrap(err, "unable load file info")
	}

	// Определяем длину файла в байтах
	fileSize := fileInfo.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	log.Printf("fileseze %v", fileSize)

	// Если неизвестна длина (например, /dev/urandom)
	if fileSize == 0 {
		return ErrUnsupportedFile
	}

	// Определение количества копируемых байт
	bytesToCopy := fileSize - offset
	if limit > 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	_, err = file.Seek(offset, 0)
	if err != nil {
		return errors.Wrap(err, "unable set offset in file")
	}

	log.Printf("bytesToCopy %v", bytesToCopy)

	// Если limit равен 0, копируем все доступные данные
	if limit == 0 {
		bytesToCopy = fileSize - offset
	}

	// Создание и настройка прогресс-бара
	bar := pb.New64(bytesToCopy)
	// Создание функции обратного вызова для прогресс-бара
	barReader := bar.NewProxyReader(file)

	// Создаем временный файл и копируем данные в него
	tmpFile, err := os.CreateTemp(".", "tmp_file_")
	if err != nil {
		return errors.Wrap(err, "unable create tmp file")
	}

	log.Printf("Временный файл создан %v", nil)

	// Копируем нужную часть файла
	_, err = io.CopyN(tmpFile, barReader, bytesToCopy)
	if err != nil && err != io.EOF {
		return errors.Wrap(err, "unable copy file")
	}

	log.Printf("Копируется байт %v", bytesToCopy)

	// Копирование завершено, переименование временного файла
	err = os.Rename(tmpFile.Name(), toPath)
	if err != nil {
		return errors.Wrap(err, "rename error")
	}
	bar.Finish()

	return nil
}
