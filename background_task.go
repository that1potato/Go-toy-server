package main

import (
	"log"
	"sync"
	"time"
)

const TASK_EXECUTION_TIME = time.Second * 60

func startBackgroundTask(wg *sync.WaitGroup, num int) error {
	log.Printf("Starting %d tasks...", num)
	for range num {
		wg.Add(1)
		go func() {
			if err := doTask(wg); err != nil {
				log.Fatal("An error occurred during task execution: ", err)
			}
		}()
	}
	return nil
}

func doTask(wg *sync.WaitGroup) error {
	defer wg.Done()

	// do some work here
	time.Sleep(TASK_EXECUTION_TIME)
	return nil
}
