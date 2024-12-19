package main

import (
    "github.com/asynkron/protoactor-go/actor"
    "github.com/asynkron/protoactor-go/remote"
    "reddit-clone/engine/actors"
    "reddit-clone/engine/messages"
	"log"
    "net/http"
    "reddit-clone/api"
    "reddit-clone/api/handlers"
)

func main() {
    system := actor.NewActorSystem()

    // Configure the remote system
    remoteConfig := remote.Configure("localhost", 8080)
    remoting := remote.NewRemote(system, remoteConfig)
    remoting.Start()

    // Create a root actor to manage users and subreddits
    rootProps := actor.PropsFromProducer(func() actor.Actor {
        return &RootActor{
            users:      make(map[string]*actor.PID),
            subreddits: make(map[string]*actor.PID),
        }
    })

    

    // Register the root actor
    remoting.Register("root", rootProps)
	
	// Spawn the root actor
    root, err := system.Root.SpawnNamed(rootProps, "root")
    if err != nil {
        log.Fatal("Failed to spawn root actor:", err)
    }

    // Initialize handlers
    userHandler := handlers.NewUserHandler(system, root)
    postHandler := handlers.NewPostHandler(system, root)
	subredditHandler := handlers.NewSubredditHandler(system, root)
    commentHandler := handlers.NewCommentHandler(system, root)
	// Setup router
    router := api.SetupRouter(userHandler, postHandler, subredditHandler, commentHandler)

    // Start HTTP server in a separate goroutine
    go func() {
        log.Println("Starting REST API server on :8081") // Note: Using 8081 since 8080 is used by Proto.Actor
        if err := http.ListenAndServe(":8081", router); err != nil {
            log.Fatal("HTTP server error:", err)
        }
    }()
	
	// Keep the program running
    select {}
}

type RootActor struct {
    users      map[string]*actor.PID
    subreddits map[string]*actor.PID
}

func (r *RootActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
    case *messages.RegisterAccount:
        props := actor.PropsFromProducer(func() actor.Actor { return actors.NewUserActor(msg.Username) })
        pid, _ := context.SpawnNamed(props, msg.Username)
        r.users[msg.Username] = pid
    case *messages.CreateSubreddit:
        props := actor.PropsFromProducer(func() actor.Actor { return actors.NewSubredditActor(msg.Name) })
        pid, _ := context.SpawnNamed(props, msg.Name)
        r.subreddits[msg.Name] = pid
    }
}