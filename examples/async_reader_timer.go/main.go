package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/equalsgibson/concur/concur"
)

// https://pkg.go.dev/runtime/pprof#hdr-Profiling_a_Go_program
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// Create a context that we can pass into the async loop function
	ctx := context.Background()

	iterations := uint(0)

	// Create a new ASyncReader that will print the current iteration every 5 seconds.
	// This could also fetch data from an API or database at specific intervals, or
	// set up an asynchronous connection to a datasource (i.e. a Websocket)
	reader := concur.NewAsyncReader(
		func(ctx context.Context) (uint, error) {
			timer := time.NewTimer(time.Second * 1)
			defer timer.Stop()

			for {
				select {
				case <-timer.C:
					if iterations >= 300 {
						return 0, errors.New("end of the example - thanks for using the concur package")
					}
					iterations++
					return iterations, nil
				case <-ctx.Done():
					return 0, ctx.Err()
				}
			}
		},
	)

	// Create a goroutine for the Loop function to run in, so that the main program is not
	// stopped from continuing execution while we wait for an asynchronous update.
	go reader.Loop(ctx)
	// Defer reader.Close, so that the reader.Loop() goroutine can return
	defer reader.Close()

	// Listen for updates on the reader.Updates() channel and check the context has not been cancelled.
	for {
		select {
		case update := <-reader.Updates():
			if update.Err != nil {
				log.Printf("Got an error response from Loop fetch function: %s", update.Err.Error())

				if *memprofile != "" {
					f, err := os.Create(*memprofile)
					if err != nil {
						log.Fatal("could not create memory profile: ", err)
					}
					defer f.Close() // error handling omitted for example
					runtime.GC()    // get up-to-date statistics
					if err := pprof.WriteHeapProfile(f); err != nil {
						log.Fatal("could not write memory profile: ", err)
					}
				}

				return
			}

			log.Printf("Current Iteration: %d - Current goroutines: %d", iterations, runtime.NumGoroutine())
		case <-ctx.Done():
			log.Printf("Context was cancelled or hit deadline: %s", ctx.Err())

			if *memprofile != "" {
				f, err := os.Create(*memprofile)
				if err != nil {
					log.Fatal("could not create memory profile: ", err)
				}
				defer f.Close() // error handling omitted for example
				runtime.GC()    // get up-to-date statistics
				if err := pprof.WriteHeapProfile(f); err != nil {
					log.Fatal("could not write memory profile: ", err)
				}
			}

			return
		}
	}
}
