package main

import (
    "github.com/asynkron/protoactor-go/actor"
    //"github.com/asynkron/protoactor-go/remote"
    "math/rand"
    "reddit-clone/engine/messages"
	"reddit-clone/metrics"
    "time"
	"fmt"
	"log"
)

type UserStatus struct {
    Username string
    Online   bool
	LastActive time.Time
}

type Simulator struct {
    system *actor.ActorSystem
    root   *actor.PID
    users  []UserStatus 
}

func GenerateZipfDistribution(numSubreddits int, s float64) []int {
    counts := make([]int, numSubreddits)
    total := 0

    // Calculate the sum of the first numSubreddits terms of the Zipf distribution
    for i := 1; i <= numSubreddits; i++ {
        total += int(1 / (float64(i) * float64(i)))
    }

    // Generate member counts based on Zipf distribution
    for i := 1; i <= numSubreddits; i++ {
        counts[i-1] = int(float64(total) / (float64(i) * float64(i)) * float64(numSubreddits))
    }

    return counts
}

func NewSimulator(system *actor.ActorSystem, rootAddress string) *Simulator {
    root := actor.NewPID(rootAddress, "root")
    return &Simulator{
        system: system,
        root:   root,
        users:  make([]UserStatus, 0),
    }
}

func (s *Simulator) SimulateUsers(count int) {
    for i := 0; i < count; i++ {
        username := fmt.Sprintf("user%d", i)
        s.system.Root.Send(s.root, &messages.RegisterAccount{Username: username})
        
        // Add each user with an initial online status
        s.users = append(s.users, UserStatus{Username: username, Online: true}) // Start all users as online
        fmt.Printf("User registered: %s\n", username)
		metrics.LogEvent("registration")  // Log user registration
    }
}

func (s *Simulator) SimulateSubredditPosts(numSubreddits int) {
    memberCounts := GenerateZipfDistribution(numSubreddits, 1.0)

    for _, count := range memberCounts {
        numPosts := count / 10 // Adjust this divisor as needed for realistic post counts

        for j := 0; j < numPosts; j++ {
            subredditName := fmt.Sprintf("subreddit%d", rand.Intn(numSubreddits))
            postContent := fmt.Sprintf("Post content from auto-generated post in %s", subredditName)

            if rand.Float32() < 0.3 { // 30% chance to be a re-post
                postContent += " (Re-post)"
            }

            fmt.Printf("Auto-generating post in %s: %s\n", subredditName, postContent)
            s.system.Root.Send(s.root, &messages.Post{
                SubredditName: subredditName,
                Username:      "auto_user",
                Content:       postContent,
            })
		}
    }
}

func (s *Simulator) ToggleUserConnection() {
    for i := range s.users {
        if rand.Float32() < 0.1 { // 10% chance to toggle each user's connection state
            s.users[i].Online = !s.users[i].Online
            if s.users[i].Online {
                fmt.Printf("User %s is now ONLINE\n", s.users[i].Username)
            } else {
                fmt.Printf("User %s is now OFFLINE\n", s.users[i].Username)
            }
        }
    }
}

func (s *Simulator) CheckInactivity() {
    for {
        time.Sleep(5 * time.Second) // Check every 5 seconds

        for i := range s.users {
            if !s.users[i].Online && time.Since(s.users[i].LastActive) > 10*time.Second { 
                s.users[i].Online = false
                fmt.Printf("User %s is now OFFLINE due to inactivity.\n", s.users[i].Username)
            }
        }
    }
}

func (s *Simulator) SimulateConnectionFluctuations() {
    for {
        time.Sleep(time.Second * 60) // Wait for 10 seconds before toggling
        s.ToggleUserConnection()           // Toggle connection state
    }
}

func (s *Simulator) SimulateActivity() {
    for {
        time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000))) 
        
        userIndex := rand.Intn(len(s.users))
        user := s.users[userIndex]
		
		user.Online = true
        user.LastActive = time.Now()

        // Check if the selected user is online before proceeding with actions
        if user.Online {
            action := rand.Intn(10)
        
            switch action {
            case 0:
                start := time.Now()
				subredditName := fmt.Sprintf("subreddit%d", rand.Intn(100))
                fmt.Printf("User %s creating subreddit: %s\n", user.Username, subredditName)
                s.system.Root.Send(s.root, &messages.CreateSubreddit{
                    Name:    subredditName,
                    Creator: user.Username,
                })
				elapsed := time.Since(start)
                log.Printf("Creating subreddit took %s", elapsed)
				metrics.LogEvent("create_subreddit") 
            case 1:
                start := time.Now()
				subredditName := fmt.Sprintf("subreddit%d", rand.Intn(100))
                fmt.Printf("User %s joining subreddit: %s\n", user.Username, subredditName)
                s.system.Root.Send(s.root, &messages.JoinSubreddit{
                    SubredditName: subredditName,
                    Username:      user.Username,
                })
				elapsed := time.Since(start)
				log.Printf("joining subreddit took %s", elapsed)
				metrics.LogEvent("join_subreddit") 
				
            case 2:
                start := time.Now()
				subredditName := fmt.Sprintf("subreddit%d", rand.Intn(100))
                postContent := fmt.Sprintf("Post content from %s", user.Username)
                fmt.Printf("User %s posting in subreddit %s: %s\n", user.Username, subredditName, postContent)
                s.system.Root.Send(s.root, &messages.Post{
                    SubredditName: subredditName,
                    Username:      user.Username,
                    Content:       postContent,
                })
				elapsed := time.Since(start)
				log.Printf("posting in subreddit took %s", elapsed)
				metrics.LogEvent("post") 
				
            case 3:
                start := time.Now()
				postID := fmt.Sprintf("post%d", rand.Intn(100))
                var parentCommentID *string 
                if rand.Float32() < 0.5 {
                    id := fmt.Sprintf("comment%d", rand.Intn(50))
                    parentCommentID = &id
                }
                commentContent := fmt.Sprintf("Comment from %s", user.Username)
                fmt.Printf("User %s commenting on %v: %s\n", user.Username, parentCommentID, commentContent)
                s.system.Root.Send(s.root, &messages.Comment{
                    PostID:          postID,
                    ParentCommentID: parentCommentID,
                    Username:        user.Username,
                    Content:         commentContent,
                })
				elapsed := time.Since(start)
				log.Printf("commenting took %s", elapsed)
				metrics.LogEvent("comment")
				
            case 4:
                start := time.Now()
				postID := fmt.Sprintf("post%d", rand.Intn(10)) 
                isUpvote := rand.Intn(2) == 0
                voteType := "upvoting"
                if !isUpvote {
                    voteType = "downvoting"
                }
	            
                if user.Online{
                    fmt.Printf("User %s %s post %s\n", user.Username, voteType, postID)
					log.Printf("VOTE: User %s %s post %s", user.Username, voteType, postID) 
                    
                   
                    s.system.Root.Send(s.root, &messages.Vote{
                        PostID:   postID,
                        Username: user.Username,
                        IsUpvote: isUpvote,
                    })
                    if isUpvote {
                        fmt.Printf("User %s will gain karma for upvoting post %s\n", user.Username, postID)
                    } else {
                        fmt.Printf("User %s will lose karma for downvoting post %s\n", user.Username, postID)
                    }
	            
                    
					future := s.system.Root.RequestFuture(s.root, &messages.GetKarma{Username: user.Username}, 0*time.Second) 
                    result, err := future.Result()
                    if err != nil {
                        //fmt.Printf("Error retrieving karma for %s: %v\n", user.Username, err)
                    } else {
                        fmt.Printf("Current karma for user %s: %d\n", user.Username, result.(int32))
                    }
                } else {
                    fmt.Printf("User %s is offline and cannot vote on post %s.\n", user.Username, postID)
                }
				elapsed := time.Since(start)
				log.Printf("Voting process for user %s took %s", user.Username, elapsed)
				metrics.LogEvent("vote")
		    case 5:
                start := time.Now()
				subredditName := fmt.Sprintf("subreddit%d", rand.Intn(100))
                fmt.Printf("User %s leaving subreddit: %s\n", user.Username, subredditName)
                s.system.Root.Send(s.root, &messages.LeaveSubreddit{
                    SubredditName: subredditName,
                    Username:      user.Username,
                })
				elapsed := time.Since(start)
				log.Printf("leaving subreddit took %s", elapsed)
				metrics.LogEvent("leave_subreddit")  
				
		    case 6:
				start := time.Now()
				toUserIndex := rand.Intn(len(s.users))
                toUser := s.users[toUserIndex].Username 
                content := fmt.Sprintf("DM from %s to %s", user.Username, toUser)
                fmt.Printf("User %s sending DM to %s: %s\n", user.Username, toUser, content)
                s.system.Root.Send(s.root, &messages.SendDirectMessage{
                    From:    user.Username,
                    To:      toUser,
                    Content: content,
                })
				elapsed := time.Since(start)
				log.Printf("sending DM took %s", elapsed)
				metrics.LogEvent("send_dm") 
			case 7:
                start := time.Now()
                fmt.Printf("User %s retrieving feed\n", user.Username)
                future := s.system.Root.RequestFuture(s.root, &messages.GetFeed{Username: user.Username}, 0*time.Second)
                result, err := future.Result()
                if err != nil {
                   // fmt.Printf("Error retrieving feed for %s: %v\n", user.Username, err)
                } else if feedResponse, ok := result.(*messages.FeedResponse); ok {
                    fmt.Printf("Feed for %s: %d posts\n", user.Username, len(feedResponse.Posts))
                } else {
                    fmt.Printf("Unexpected response type for user %s\n", user.Username)
                }
                elapsed := time.Since(start)
                log.Printf("retrieving feed took %s", elapsed)
                metrics.LogEvent("retrieve_feed")
			case 8:
                start := time.Now()
				subredditName := fmt.Sprintf("subreddit%d", rand.Intn(10))
                fmt.Printf("User %s subscribing to subreddit: %s\n", user.Username, subredditName)
                s.system.Root.Send(s.root, &messages.Subscribe{
                    Username:      user.Username,
                    SubredditName: subredditName,
                })
				elapsed := time.Since(start)
				log.Printf("Subscribing to subreddit took %s", elapsed)
				metrics.LogEvent("subscribe_subreddit") 
			case 9: 
                start := time.Now()
				postID := fmt.Sprintf("post%d", rand.Intn(10)) 
                fmt.Printf("User %s reposting post %s\n", user.Username, postID)
                s.system.Root.Send(s.root, &messages.Repost{
                    Username: user.Username,
                    PostID:   postID,
                })
				elapsed := time.Since(start)
				log.Printf("reposting post took %s", elapsed)
				metrics.LogEvent("repost") 
            }
			if time.Since(user.LastActive) > 15*time.Second { 
                user.Online = false
                fmt.Printf("User %s is now OFFLINE due to inactivity.\n", user.Username)
            }
		}else {
           fmt.Printf("User %s is offline and cannot perform actions.\n", user.Username)
       }
			
    }
}