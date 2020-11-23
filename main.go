package main

import (
	"log"
	"task-pool/pool"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

func main() {
	// Start CPU Profile
	f, err := os.Create("/tmp/task-pool-cpu.prof")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatalln(err.Error())
	}
	defer pprof.StopCPUProfile()

	wg := new(sync.WaitGroup)

	var workerPool pool.WorkerPool
	workerPoolConfig := pool.Config{
		Concurrency: 10,
		OnTaskComplete: func() {
			wg.Done()
			log.Printf("task complete, %d remaining\n", workerPool.RemainingTasks())
		},
	}
	workerPool = pool.NewWorkerPool(workerPoolConfig)

	// Generate a bunch of mock tasks.
	minDuration := 1
	maxDuration := 5

	nTasks := 100
	wg.Add(nTasks)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < nTasks; i++ {
		workerPool.AddTask(pool.Task{
			// Assign a variable amount of "work" to do.
			Sleep: rand.Intn(maxDuration - minDuration + 1) + minDuration,
		})
	}

	wg.Wait()

	// Write memory profile to disk.
	memoryFile, err := 	os.Create("/tmp/task-pool-memory.prof")
	if err != nil {
		log.Fatalln("Failed to open file for memory profile: ", err)
	}
	defer memoryFile.Close()
	runtime.GC()
	if err = pprof.WriteHeapProfile(memoryFile); err != nil {
		log.Fatalln("could not write memory profile: ", err)
	}
}
