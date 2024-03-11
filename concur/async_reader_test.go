package concur_test

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/equalsgibson/concur/concur"
)

func Test_AsyncReader_ValidateMessageOrder(t *testing.T) {
	mockData := []string{
		"This",
		"Is",
		"The",
		"Way",
	}

	// expectedData := mockData

	counter := 0

	mockReader := concur.NewAsyncReader[string](func(context.Context) (string, error) {
		counter++
		if len(mockData) > 0 {
			result := mockData[0]
			mockData = mockData[1:]

			return result, nil
		}

		return "", errors.New("no more data")
	})

	ctx := context.Background()
	go mockReader.Loop(ctx)

	actualData := []string{}
	go func() {
		for update := range mockReader.Updates() {
			if update.Err != nil {
				t.Fatal(update.Err)
			}
			log.Println(update.Item)

			actualData = append(actualData, update.Item)
		}
	}()

	for {
		if counter != 4 {
			continue
		}
		break
	}

	mockReader.Close()

	t.Fatalf("%+v", actualData)
}

func Test_AsyncReader_ValidateCleanup(t *testing.T) {

}

func Test_AsyncReader_ValidateMultipleGoRountines(t *testing.T) {

}
