package main

import (
	workerPool "github.com/Alexandr-Fox/worker-pool/worker-pool"
	"github.com/go-loremipsum/loremipsum"
)

func main() {
	loremIpsumGenerator := loremipsum.NewWithSeed(1234)
	wp := workerPool.NewWorkerPull()

	wp.AddWorkers(5)

	for i := 0; i < 15; i++ {
		wp.AddMessage(loremIpsumGenerator.Words(5))
	}

	wp.DeleteWorkers(3)
	wp.Close()
}
