package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/Jeffail/tunny"
)

var pool *tunny.Pool

type Job struct {
	Data string
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func main() {
	pool = tunny.NewFunc(2, worker)
	defer pool.Close()

	for i := 0; i < 100; i++ {
		postJob(time.Now().String())
	}
	select {}
}

func worker(work interface{}) interface{} {
	switch w := work.(type) {
	case *Job:
		return w.build()
	}

	return "Couldn't find work type"
}

func (j *Job) build() string {
	log.Printf("Starting to process job %s ", j.Data)
	time.Sleep(time.Duration(random(1, 3)) * time.Second)
	log.Printf("Done %v", j.Data)
	return "OK"
}

func postJob(value string) {
	go processJob(pool, value)
}

func processJob(pool *tunny.Pool, data string) {
	j := &Job{Data: data}
	_, err := pool.ProcessTimed(j, time.Minute*30)
	if err == tunny.ErrJobTimedOut {
		log.Printf("problem to process job %v", err)
	}
}
