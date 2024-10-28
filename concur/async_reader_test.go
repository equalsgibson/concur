package concur_test

import (
	"context"
	"errors"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/equalsgibson/concur/concur"
)

func Test_AsyncReader_ValidateMessageOrder(t *testing.T) {
	mockData := []string{
		"This",
		"Is",
		"The",
		"Way",
	}

	counter := 0

	mockReader := concur.NewAsyncReader(func(context.Context) (string, error) {
		if len(mockData) == counter {
			return "", errors.New("no more data")
		}

		counter++

		return mockData[counter-1], nil
	})

	ctx := context.Background()
	go mockReader.Loop(ctx)

	actualData := []string{}

	for update := range mockReader.Updates() {
		if update.Item == "" {
			break
		}

		actualData = append(actualData, update.Item)
	}

	mockReader.Close()

	if !reflect.DeepEqual(actualData, mockData) {
		t.Fatalf("expected updates to match original data - \nactualData:\t%+v,\nmockData:\t%+v", actualData, mockData)
	}
}

func Test_AsyncReader_CancelContext(t *testing.T) {
	mockData := []string{
		"This",
		"Is",
		"The",
		"Way",
		"We",
		"Should",
		"Break",
		"Here",
		"And",
		"Not",
		"See",
		"Additional",
		"Messages",
	}

	counter := 0

	mockReader := concur.NewAsyncReader(func(ctx context.Context) (string, error) {
		if len(mockData) == counter {
			return "", errors.New("no more data")
		}

		counter++

		return mockData[counter-1], nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mockReader.Loop(ctx)

	actualData := []string{}

	for {
		select {
		case update := <-mockReader.Updates():
			if update.Item == "And" {
				cancel()

				continue
			}

			actualData = append(actualData, update.Item)
		case <-ctx.Done():
			mockReader.Close()
			if !errors.Is(ctx.Err(), context.Canceled) {
				t.Fatalf("Did not get expected error for context - got: %s", ctx.Err())
			}

			return
		}
	}
}

func Test_AsyncReader_ValidateCleanup(t *testing.T) {
	initalGoroutineCount := runtime.NumGoroutine()

	counter := 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockReader := concur.NewAsyncReader(func(context.Context) (string, error) {
		counter++
		return "", nil
	})

	go func() {
		for {
			if counter >= 5 {
				cancel()
				break
			}
		}
	}()

	go mockReader.Loop(ctx)

	for {
		select {
		case <-mockReader.Updates():
			// Do nothing - we are just testing cleanup here.
		case <-ctx.Done():
			mockReader.Close()
			time.Sleep(time.Millisecond * 50)

			if runtime.NumGoroutine() > initalGoroutineCount {
				t.Fatalf(
					"Found %d goroutines running at the end of the test but we started with %d",
					runtime.NumGoroutine(),
					initalGoroutineCount,
				)
			}

			return
		}
	}
}
