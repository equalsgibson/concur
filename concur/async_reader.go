package concur

import (
	"context"
	"time"
)

type FetchResult[T any] struct {
	Item T
	Err  error
}

type FetchFunc[T any] func(ctx context.Context) (T, error)

func NewAsyncReader[T any](
	fetcher func(context.Context) (T, error),
) *AsyncReader[T] {
	return &AsyncReader[T]{
		fetcher: fetcher,
		closing: make(chan bool),
		updates: make(chan FetchResult[T]),
	}
}

type AsyncReader[T any] struct {
	fetcher FetchFunc[T]
	closing chan bool
	updates chan FetchResult[T]
}

func (r *AsyncReader[T]) Updates() <-chan FetchResult[T] {
	return r.updates
}

func (r *AsyncReader[T]) Close() {
	r.closing <- true
}

func (r *AsyncReader[T]) Loop(ctx context.Context) {
	var fetchCompleted chan FetchResult[T]

	var queue []FetchResult[T]
	var err error

	for {
		var first FetchResult[T]
		var updates chan FetchResult[T]

		if len(queue) > 0 {
			first = queue[0]
			updates = r.updates
		}

		var startFetch <-chan time.Time

		if fetchCompleted == nil && err == nil {
			startFetch = time.After(0)
		}

		select {
		case <-r.closing:
			close(r.updates)
			return

		case <-startFetch:
			fetchCompleted = make(chan FetchResult[T], 1)

			go func() {
				fetched, err := r.fetcher(ctx)
				fetchCompleted <- FetchResult[T]{fetched, err}
			}()

		case fetchResult := <-fetchCompleted:
			err = fetchResult.Err
			queue = append(queue, fetchResult)
			fetchCompleted = nil

		case updates <- first:
			queue = queue[1:]
		}
	}
}
