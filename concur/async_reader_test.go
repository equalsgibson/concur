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

	mockReader := concur.NewAsyncReader[string](func(context.Context) (string, error) {
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

func Test_AsyncReader_ValidateCleanup(t *testing.T) {
	initalGoroutineCount := runtime.NumGoroutine()

	counter := 0

	mockReader := concur.NewAsyncReader[string](func(context.Context) (string, error) {
		counter++
		return "anything", nil
	})

	go func() {
		for {
			if counter >= 5 {
				mockReader.Close()
				break
			}
		}
	}()

	ctx := context.Background()
	go mockReader.Loop(ctx)

	for range mockReader.Updates() {
		// Do nothing - just testing goroutines
	}

	time.Sleep(time.Millisecond * 50)

	if runtime.NumGoroutine() > initalGoroutineCount {
		t.Fatalf(
			"Found %d goroutines running at the end of the test but we started with %d",
			runtime.NumGoroutine(),
			initalGoroutineCount,
		)
	}
}
