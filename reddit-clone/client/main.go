package main

import (
    "fmt"
    "sync"
    "time"
    "log"
	"math/rand"
)

var subredditNames = []string{"Technology", "Science", "Gaming", "Movies", "Books", "USA", "DOSP", "CISE", "UF"}

// Metrics structure to track operations
type Metrics struct {
    sync.Mutex
    totalUsers           int
    totalPosts          int
    totalComments       int
    totalVotes          int
    subredditsCreated   int
    subredditJoins      int
    leaveSubredditCount int
    sendDMCount         int
    retrieveFeedCount   int
    subscribeCount      int
    repostCount         int
}

func main() {
    var start time.Time
	var operationStart time.Time
	rand.Seed(time.Now().UnixNano())
    start = time.Now() 
    log.Println("Starting Reddit Clone Client...")
	

    metrics := &Metrics{}
    
    clients := make([]*Client, 10)
    for i := range clients {
        clients[i] = NewClient("http://localhost:8081")
    }
	
	var wg sync.WaitGroup
    for i := 0; i < len(clients); i++ {
        wg.Add(1)
        go func(id int, client *Client) {
            defer wg.Done()
            username := fmt.Sprintf("user%d", id)
            operationStart = time.Now()
            
			// User Registration
            if err := client.Register(username); err != nil {
                log.Printf("Failed to register user %s: %v\n", username, err)
                return
            }
            log.Printf("User registered: %s\n", username)
            metrics.Lock()
            metrics.totalUsers++
            metrics.Unlock()

            // Subreddit Operations
            randomIndex := rand.Intn(len(subredditNames))
            subredditName := subredditNames[randomIndex]
            if err := client.CreateSubreddit(subredditName); err == nil {
                log.Printf("%s Created subreddit: %s\n", 
                    time.Now().Format("2006/01/02 15:04:05"),
                    subredditName)
                metrics.Lock()
                metrics.subredditsCreated++
                metrics.Unlock()
            }
			
			

            // Join and Leave Subreddit
			var subredditNames = []string{"Technology", "Science", "Gaming", "Movies", "DOSP", "CISE", "UF", "USA", "Books",}
			joinIndex := rand.Intn(len(subredditNames))
            subredditToJoin := subredditNames[joinIndex]
            if err := client.JoinSubreddit(subredditToJoin); err == nil {
                log.Printf("%s User %s joining subreddit: %s\n", 
                    time.Now().Format("2006/01/02 15:04:05"),
                    username, 
                    subredditToJoin)
                metrics.Lock()
                metrics.subredditJoins++
                metrics.Unlock()
            }
            
            // Use a different index for leaving
            leaveIndex := (joinIndex + 1) % len(subredditNames)
            subredditToLeave := subredditNames[leaveIndex]
            if err := client.LeaveSubreddit(subredditToLeave); err == nil {
                log.Printf("%s User %s leaving subreddit: %s\n",
                    time.Now().Format("2006/01/02 15:04:05"),
                    username,
                    subredditToLeave)
                metrics.Lock()
                metrics.leaveSubredditCount++
                metrics.Unlock()
            }
			// Post Operations
			randomIndex = rand.Intn(len(subredditNames))
            subredditName = subredditNames[randomIndex]
            postID := fmt.Sprintf("post%d", id)
            if err := client.CreatePost(postID, fmt.Sprintf("Post content from %s", username), subredditName); err == nil {
                log.Printf("User %s posting in subreddit %s: Post content from %s\n", username, subredditName, username)
                metrics.Lock()
                metrics.totalPosts++
                metrics.Unlock()
            }

            // Comment Operations
            var parentCommentID *string // nil for top-level comment
            if err := client.CreateComment(postID, "Test comment", parentCommentID); err == nil {
                log.Printf("User %s commenting on post %s\n", username, postID)
                metrics.Lock()
                metrics.totalComments++
                metrics.Unlock()
            }

            // Voting Operations
            start = time.Now() 
            postID = fmt.Sprintf("post%d", rand.Intn(10))
            isUpvote := rand.Intn(2) == 0
            voteType := "upvoting"
            if !isUpvote {
                voteType = "downvoting"
            }
            
            if err := client.VotePost(postID, isUpvote); err == nil {
                log.Printf("User %s %s post %s\n", username, voteType, postID)
                if isUpvote {
                    log.Printf("User %s will gain karma for upvoting post %s\n", username, postID)
                } else {
                    log.Printf("User %s will lose karma for downvoting post %s\n", username, postID)
                }
                
                karma, err := client.GetKarma()
                if err != nil {
                    log.Printf("Error retrieving karma for %s: %v\n", username, err)
                } else {
                    log.Printf("Current karma for user %s: %d\n", username, karma)
                }
                
                metrics.Lock()
                metrics.totalVotes++
                metrics.Unlock()
            } else {
                log.Printf("Failed to vote on post: %v\n", err)
            }
            
            elapsed := time.Since(start)
            log.Printf("Voting process for user %s took %s", username, elapsed)


            if err := client.VotePost(fmt.Sprintf("post%d", (id+1)%3), false); err == nil {
                log.Printf("User %s downvoting post post%d\n", username, (id+1)%3)
                metrics.Lock()
                metrics.totalVotes++
                metrics.Unlock()
            }

            // Direct Messaging
            if err := client.SendDM(fmt.Sprintf("user%d", (id+1)%10), 
                fmt.Sprintf("DM from %s", username)); err == nil {
                log.Printf("User %s sending DM to user%d\n", username, (id+1)%10)
                metrics.Lock()
                metrics.sendDMCount++
                metrics.Unlock()
            }

            // Repost Operations
            if err := client.RepostContent(fmt.Sprintf("post%d", id%2)); err == nil {
                log.Printf("User %s reposting post post%d\n", username, id%2)
                metrics.Lock()
                metrics.repostCount++
                metrics.Unlock()
            }

            // Feed Operations
            if err := client.GetFeed(); err == nil {
                log.Printf("User %s retrieving feed\n", username)
                metrics.Lock()
                metrics.retrieveFeedCount++
                metrics.Unlock()
            }

            log.Printf("Operations for user %s completed in %v\n", username, time.Since(operationStart))
        }(i, clients[i])
    }

    wg.Wait()
    
    // Print metrics report
    log.Println("--- Metrics Report ---")
    log.Printf("Total Users: %d\n", metrics.totalUsers)
    log.Printf("Total Posts: %d\n", metrics.totalPosts)
    log.Printf("Total Comments: %d\n", metrics.totalComments)
    log.Printf("Total Votes: %d\n", metrics.totalVotes)
    log.Printf("Subreddits Created: %d\n", metrics.subredditsCreated)
    log.Printf("Subreddit Joins: %d\n", metrics.subredditJoins)
    log.Printf("Leave Subreddit Count: %d\n", metrics.leaveSubredditCount)
    log.Printf("Send DM Count: %d\n", metrics.sendDMCount)
    log.Printf("Retrieve Feed Count: %d\n", metrics.retrieveFeedCount)
    log.Printf("Subscribe Count: %d\n", metrics.subscribeCount)
    log.Printf("Repost Count: %d\n", metrics.repostCount)
    log.Printf("Total execution time: %v\n", time.Since(start))
    log.Println("Exiting gracefully...")
}
