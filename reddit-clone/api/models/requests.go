package models

type RegisterRequest struct {
    Username string `json:"username"`
}

type PostRequest struct {
    Title     string `json:"title"`
    Content   string `json:"content"`
    Subreddit string `json:"subreddit"`
}

type VoteRequest struct {
    PostID string `json:"post_id"`
    Vote   int    `json:"vote"` // 1 for upvote, -1 for downvote
}

type CommentRequest struct {
    PostID  string `json:"post_id"`
    Content string `json:"content"`
}

type SubredditRequest struct {
    Name string `json:"name"`
}

type DMRequest struct {
    ToUser  string `json:"to_user"`
    Message string `json:"message"`
}