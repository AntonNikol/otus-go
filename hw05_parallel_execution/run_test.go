package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.Eventually(t, func() bool {
			return atomic.LoadInt32(&runTasksCount) == int32(tasksCount)
		}, time.Second*5, time.Millisecond*500, "not all tasks were completed")
	})

	t.Run("task with errors unlimited", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Equalf(t, nil, err, "actual err - %v", err)
	})

	t.Run("the number of goroutines is 0", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		tasks = append(tasks, func() error { return nil })
		workersCount := 0
		maxErrorsCount := 0

		err := Run(tasks, workersCount, maxErrorsCount)
		require.Equalf(t, err, ErrNumberGoroutinesInvalid, "actual err - %v", err)
	})

	t.Run("no tasks provided", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		workersCount := 0
		maxErrorsCount := 0

		err := Run(tasks, workersCount, maxErrorsCount)
		require.Equalf(t, err, ErrNoTasks, "actual err - %v", err)
	})

	t.Run("no tasks provided", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		workersCount := 0
		maxErrorsCount := 0

		err := Run(tasks, workersCount, maxErrorsCount)
		require.Equalf(t, err, ErrNoTasks, "actual err - %v", err)
	})

	t.Run("no tasks provided", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		workersCount := 0
		maxErrorsCount := 0

		err := Run(tasks, workersCount, maxErrorsCount)
		require.Equalf(t, err, ErrNoTasks, "actual err - %v", err)
	})

	t.Run("task with exist nil element", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount-1; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}
		tasks = append(tasks, nil)

		workersCount := 10
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount-1), "not all tasks were completed")
	})
}
