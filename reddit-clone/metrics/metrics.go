package metrics

import (
    "log"
    "time"
    "runtime"
)

type Metrics struct {
    UserRegistrations   int
    Posts               int
    Comments            int
    Votes               int
    ActiveUsers         int
    TotalKarma          int64
	SubredditJoins      int
	SubredditsCreated   int
	LeaveSubredditCount int 
    SendDMCount         int 
    RetrieveFeedCount   int 
    SubscribeSubredditCount int 
    RepostCount         int 
    UserOfflineCount    int 
}

var metrics Metrics
var startTime time.Time

func InitMetrics() {
    metrics = Metrics{}
    startTime = time.Now()
}

func LogEvent(eventType string) {
    switch eventType {
    case "registration":
        metrics.UserRegistrations++
        log.Printf("%s event occurred. Total: %d", eventType, metrics.UserRegistrations)
    case "post":
        metrics.Posts++
        log.Printf("%s event occurred. Total: %d", eventType, metrics.Posts)
    case "comment":
        metrics.Comments++
        log.Printf("%s event occurred. Total: %d", eventType, metrics.Comments)
    case "vote":
        metrics.Votes++
        log.Printf("%s event occurred. Total: %d", eventType, metrics.Votes)
	case "join_subreddit":
        metrics.SubredditJoins++
        log.Printf("%s event occurred. Total: %d", eventType, metrics.SubredditJoins)
	case "create_subreddit":
        metrics.SubredditsCreated++
        log.Printf("%s event occurred. Total: %d", eventType, metrics.SubredditsCreated)
    case "leave_subreddit":
        metrics.LeaveSubredditCount++  
		log.Printf("%s event occurred. Total: %d", eventType, metrics.LeaveSubredditCount)
	case "send_dm":
        metrics.SendDMCount++
        log.Printf("%s event occurred. Total: %d", eventType, metrics.SendDMCount)
    case "retrieve_feed":
        metrics.RetrieveFeedCount++
        log.Printf("%s event occurred. Total: %d", eventType, metrics.RetrieveFeedCount)
    case "subscribe_subreddit":
        metrics.SubscribeSubredditCount++
        log.Printf("%s event occurred. Total: %d", eventType, metrics.SubscribeSubredditCount)
    case "repost":
        metrics.RepostCount++
        log.Printf("%s event occurred. Total: %d", eventType, metrics.RepostCount)
    case "user_offline":
        metrics.UserOfflineCount++
        log.Printf("%s event occurred. Total: %d", eventType, metrics.UserOfflineCount)		
    }
}

func MeasureExecutionTime(operation func()) time.Duration {
    start := time.Now()
    operation()
    return time.Since(start)
}

func GenerateReport() {
    log.Printf("--- Metrics Report ---")
    log.Printf("Total Users: %d", metrics.UserRegistrations)
    log.Printf("Total Posts: %d", metrics.Posts)
    log.Printf("Total Comments: %d", metrics.Comments)
    log.Printf("Total Votes: %d", metrics.Votes)
    log.Printf("Subreddits Created: %d", metrics.SubredditsCreated)
    log.Printf("Subreddit Joins: %d", metrics.SubredditJoins)
    log.Printf("Leave Subreddit Count: %d", metrics.LeaveSubredditCount)
    log.Printf("Send DM Count: %d", metrics.SendDMCount)
    log.Printf("Retrieve Feed Count: %d", metrics.RetrieveFeedCount)
    log.Printf("Subscribe Subreddit Count: %d", metrics.SubscribeSubredditCount)
    log.Printf("Repost Count: %d", metrics.RepostCount)

    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    log.Printf("Memory Usage: %v MB", bToMb(m.Alloc))
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}