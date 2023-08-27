package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	inputFilePath  = "./testdata/input.txt"
	outputFilePath = "output.txt"
)

func prepareFileForTest(t *testing.T) (*os.File, int64) {
	file, err := os.Open(inputFilePath)
	require.NoError(t, err)

	fileInfo, err := file.Stat()
	require.NoError(t, err)

	return file, fileInfo.Size()
}

func removeTestFile(t *testing.T, filePath string) {
	err := os.Remove(filePath)
	require.NoError(t, err)
}

func TestCopy(t *testing.T) {
	t.Run("file not found", func(t *testing.T) {
		err := Copy("test.txt", "output.txt", 0, 0)
		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		file, fileSize := prepareFileForTest(t)

		err := Copy(file.Name(), "output.txt", fileSize+10, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("offset 0 limit 10", func(t *testing.T) {
		file, _ := prepareFileForTest(t)

		err := Copy(file.Name(), outputFilePath, 0, 10)
		require.NoError(t, err)
		require.FileExists(t, outputFilePath)

		// Читаем файл в котором записаные корректные результаты копирования
		expected, err := os.ReadFile("testdata/out_offset0_limit10.txt")
		require.NoError(t, err)

		copied, err := os.ReadFile(outputFilePath)
		require.NoError(t, err)

		// Сравниваем данные
		if !bytes.Equal(expected, copied) {
			t.Errorf("copy file does not math expected")
		}
		removeTestFile(t, outputFilePath)
	})

	t.Run("offset 100 limit 1000", func(t *testing.T) {
		file, _ := prepareFileForTest(t)

		err := Copy(file.Name(), outputFilePath, 100, 1000)
		require.NoError(t, err)
		require.FileExists(t, outputFilePath)

		// Читаем файл в котором записаные корректные результаты копирования
		expected, err := os.ReadFile("testdata/out_offset100_limit1000.txt")
		require.NoError(t, err)

		copied, err := os.ReadFile(outputFilePath)
		require.NoError(t, err)

		// Сравниваем данные
		if !bytes.Equal(expected, copied) {
			t.Errorf("copy file does not math expected")
		}
		removeTestFile(t, outputFilePath)
	})
}
