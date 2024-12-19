package actors

import (
    "github.com/asynkron/protoactor-go/actor"
    "reddit-clone/engine/messages"
	"math/rand"
	"time"
	"fmt"
)

var subredditNames = []string{
    "Technology",
    "Science",
    "Gaming",
    "Movies",
    "Books",
    "USA",
    "Florida",
    "DOSP",
    "CISE",
    "UF",
    "Cooking",
    "Travel",
    "Music",
}


// Initialize random seed in init function
func init() {
    rand.Seed(time.Now().UnixNano())
}

type SubredditActor struct {
    name    string
    members map[string]bool
    posts   []*actor.PID
}


func NewSubredditActor(name string) *SubredditActor {
    return &SubredditActor{
        name:    name,
        members: make(map[string]bool),
    }
}

func (s *SubredditActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
    case *messages.CreateSubreddit:
        randomIndex := rand.Intn(len(subredditNames))
        s.name = subredditNames[randomIndex]
        fmt.Printf("%s Created subreddit: %s\n", 
            time.Now().Format("2006/01/02 15:04:05"), 
            s.name)
    case *messages.JoinSubreddit:
        s.members[msg.Username] = true
    case *messages.LeaveSubreddit:
        delete(s.members, msg.Username)
        //fmt.Printf("User %s left subreddit %s\n", msg.Username, s.name)
	case *messages.Post:
        props := actor.PropsFromProducer(func() actor.Actor { return NewPostActor(msg.Content, msg.Username) })
        pid, _ := context.SpawnNamed(props, msg.Content)
        s.posts = append(s.posts, pid)
    }
}