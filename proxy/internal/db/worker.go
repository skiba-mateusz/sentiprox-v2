package db

import (
	"time"

	"go.uber.org/zap"
)

const (
	defaultBatchSize = 500
)

type Worker[T any] struct {
	logger    *zap.SugaredLogger
	batchSize int
	interval  time.Duration
	ch        chan T
	flush     func([]T) error
	done      chan struct{}
}

func NewWorker[T any](l *zap.SugaredLogger, flush func([]T) error, batchSize int, interval time.Duration) *Worker[T] {
	if batchSize <= 0 {
		batchSize = defaultBatchSize
	}

	w := &Worker[T]{
		logger:    l,
		batchSize: batchSize,
		interval:  interval,
		ch:        make(chan T, batchSize*4),
		flush:     flush,
		done:      make(chan struct{}),
	}

	go w.run()
	return w
}

func (w *Worker[T]) Enqueue(item T) {
	select {
	case w.ch <- item:
	default:
		w.logger.Warn("worker: buf full, omitting item")
	}
}

func (w *Worker[T]) Close() {
	close(w.ch)
	<-w.done
}

func (w *Worker[T]) run() {
	defer close(w.done)
	batch := make([]T, 0, w.batchSize)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		if err := w.flush(batch); err != nil {
			w.logger.Errorw("worker: flush failed", "err", err)
		}
		batch = make([]T, 0, w.batchSize)
	}

	for {
		select {
		case item, ok := <-w.ch:
			if !ok {
				ticker.Stop()
				flush()
				return
			}
			batch = append(batch, item)
			if len(batch) >= w.batchSize {
				flush()
				ticker.Reset(w.interval)
			}
		case <-ticker.C:
			flush()
		}
	}
}
