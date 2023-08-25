package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for i := 0; i < len(stages); i++ {
		in = execStage(stages[i], in, done)
	}

	return in
}

func execStage(stage Stage, in, done In) Out {
	inCh := make(Bi)

	go func() {
		defer close(inCh)
		for {
			select {
			case <-done:
				return
			case item, ok := <-in:
				if !ok {
					return
				}
				inCh <- item
			}
		}
	}()

	return stage(inCh)
}
