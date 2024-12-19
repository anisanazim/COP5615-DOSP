package actors

import (
    "github.com/asynkron/protoactor-go/actor"
    "reddit-clone/engine/messages"
	"fmt"
)

type UserActor struct {
    username string
    karma    int32
    subreddits map[string]bool
    directMessages []*messages.SendDirectMessage
    feed []*messages.Post
	subscriptions  []string
}

const maxKarmaChange = 100

func (u *UserActor) updateKarma(amount int32) {
    if amount > maxKarmaChange {
        amount = maxKarmaChange
    } else if amount < -maxKarmaChange {
        amount = -maxKarmaChange
    }
    u.karma += amount
}

func NewUserActor(username string) *UserActor {
    return &UserActor{
        username: username,
        subreddits: make(map[string]bool),
		karma:    0,
		subscriptions: make([]string, 0),
    }
}
// add a method to run periodically
func (u *UserActor) decayKarma() {
    u.karma = int32(float64(u.karma) * 0.99) // 1% decay
}

func (u *UserActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
    case *messages.JoinSubreddit:
        u.subreddits[msg.SubredditName] = true
    case *messages.LeaveSubreddit:
        delete(u.subreddits, msg.SubredditName)
    
	case *messages.Vote:
    if msg.IsUpvote {
        u.updateKarma(1)
    } else {
        u.updateKarma(-1)
    }
	
	case *messages.SendDirectMessage:
        u.directMessages = append(u.directMessages, msg)
    case *messages.GetDirectMessages:
        context.Respond(u.directMessages)
    case *messages.GetFeed:
        fmt.Printf("Received GetFeed request for user: %s\n", u.username)
        // Ensure feed is populated
        if len(u.feed) == 0 {
            context.Respond(&messages.FeedResponse{Posts: []*messages.Post{}})
            fmt.Printf("Responded with empty feed for user: %s\n", u.username)
        } else {
            context.Respond(&messages.FeedResponse{Posts: u.feed})
            fmt.Printf("Responded with feed for user: %s, %d posts\n", u.username, len(u.feed))
        }
    case *messages.Post:
        if u.subreddits[msg.SubredditName] {
            u.feed = append(u.feed, msg)
        }
	
	case *messages.UpdateKarma:
    u.updateKarma(msg.Amount)
    fmt.Printf("Updated karma for user %s: %d\n", u.username, u.karma)
	
	case *messages.GetKarma:
	    fmt.Printf("Received request for karma from user %s\n", u.username)
        context.Respond(u.karma)
	case *messages.Subscribe:
        fmt.Printf("User %s subscribed to subreddit %s\n", msg.Username, msg.SubredditName)
     
        for _, sub := range u.subscriptions {
            if sub == msg.SubredditName {
                fmt.Printf("User %s is already subscribed to subreddit %s\n", msg.Username, msg.SubredditName)
                return 
            }
        }
        u.subscriptions = append(u.subscriptions, msg.SubredditName)
        fmt.Printf("Updated subscriptions for user %s: %v\n", msg.Username, u.subscriptions)
	case *messages.Repost:
        fmt.Printf("User %s reposted post %s\n", msg.Username, msg.PostID)
        newPostContent := fmt.Sprintf("Repost of post %s by user %s", msg.PostID, msg.Username)

        newPost := &messages.Post{
            SubredditName: "General", 
            Username:      msg.Username,
            Content:       newPostContent,
        }
        context.Send(context.Parent(), newPost) 
        fmt.Printf("New post created from repost by user %s: %s\n", msg.Username, newPostContent)
    }
}