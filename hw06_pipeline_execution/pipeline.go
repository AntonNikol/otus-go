package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// ExecutePipeline Выполняет пайплайн
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		if stage != nil {
			stageCh := createStageChannel(in, done)
			in = stage(stageCh)
		}
	}
	return in
}

// createStageChannel Создает канал для передачи данных между стейджами
func createStageChannel(in In, done In) Out {
	stageCh := make(Bi)
	go func() {
		defer close(stageCh)
		for {
			select {
			// если получили сигнал о завершении
			case <-done:
				return
			// если канал закрыт
			case v, ok := <-in:
				if !ok {
					return
				}
				// передаем данные в канал
				stageCh <- v
			}
		}
	}()
	return stageCh
}
