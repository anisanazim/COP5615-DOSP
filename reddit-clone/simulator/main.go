package main

import (
    "github.com/asynkron/protoactor-go/actor"
    "github.com/asynkron/protoactor-go/remote"	
    "reddit-clone/metrics"
    "time"
    "log"
	"fmt"
)

func main() {
    metrics.InitMetrics()  // Initialize metrics at the start

    system := actor.NewActorSystem()
    
    remoteConfig := remote.Configure("localhost", 0)
    remoting := remote.NewRemote(system, remoteConfig)
    remoting.Start()
    
    simulator := NewSimulator(system, "localhost:8080")

    elapsedTime := metrics.MeasureExecutionTime(func() {
        simulator.SimulateUsers(10000)
    })
    log.Printf("Simulating users took %s", elapsedTime)
    
    elapsedTime = metrics.MeasureExecutionTime(func() {
        simulator.SimulateSubredditPosts(100)
    })
    log.Printf("Simulating subreddit posts took %s", elapsedTime)

    go simulator.SimulateConnectionFluctuations()
    
    
    go simulator.SimulateActivity()
    
    
    go func() {
        for {
            time.Sleep(5 * time.Minute)
            metrics.GenerateReport()
        }
    }()
    
   done := make(chan struct{})

    go func() {
        time.Sleep(180 * time.Second) 
        close(done)
    }()

    // Keep the program running until done is closed
    select {
    case <-done:
	    metrics.GenerateReport()
        fmt.Println("Exiting gracefully...")
        return
}
}