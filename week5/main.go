package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"
)

var ServiceAccessQueue map[string][]int64
var ok bool

// Check Service Access is blocked or not with a boolean type
func ServiceCircuitBreaker(serviceName string, count uint, timeWindow int64) bool {
	currentTime := time.Now().Unix()
	if ServiceAccessQueue == nil {
		ServiceAccessQueue = make(map[string][]int64)
	}
	if _, ok = ServiceAccessQueue[serviceName]; !ok {
		ServiceAccessQueue[serviceName] = make([]int64, 0)
	}
	// The access queue is not full
	if uint(len(ServiceAccessQueue[serviceName])) < count {
		ServiceAccessQueue[serviceName] = append(ServiceAccessQueue[serviceName], currentTime)
		return true
	}
	// If the queue is full, select the earliest element
	earliestTime := ServiceAccessQueue[serviceName][0]
	// To check whether the access is within time window or not
	if currentTime-earliestTime <= timeWindow {
		return false
	} else {
		// If passed the time window, current access can replace the earliest element
		ServiceAccessQueue[serviceName] = ServiceAccessQueue[serviceName][1:]
		ServiceAccessQueue[serviceName] = append(ServiceAccessQueue[serviceName], currentTime)
	}
	return true
}

func main() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100000; i++ {
		log.Println(strconv.Itoa(i) + " " + strconv.FormatBool(ServiceCircuitBreaker("testSvc", 1000, 20)))
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	}
}
