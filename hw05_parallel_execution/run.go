package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded     = errors.New("errors limit exceeded")
	ErrNumberGoroutinesInvalid = errors.New("number goroutines should be positive")
	ErrNoTasks                 = errors.New("no tasks provided")
)

type Task func() error

func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return ErrNoTasks
	}

	if n <= 0 {
		return ErrNumberGoroutinesInvalid
	}

	// Если m меньше или равно 0 - игнорируем все ошибки
	if m <= 0 {
		m = len(tasks) + 1
	}

	// Создаем синк группу
	wg := sync.WaitGroup{}
	wg.Add(n)

	// Подсчет количества ошибок
	var errCount int32

	// Канал для задач
	taskChan := make(chan Task)

	// Запускаем горутины
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range taskChan {
				if err := task(); err != nil {
					// Инкрементируем счетчик ошибок
					atomic.AddInt32(&errCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		// Пока количество ошибок не превысило допустимое складываем задачи в канал
		if atomic.LoadInt32(&errCount) >= int32(m) {
			break
		}
		taskChan <- task
	}
	close(taskChan)
	wg.Wait()

	if errCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
