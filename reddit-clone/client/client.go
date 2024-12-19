package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type Client struct {
    baseURL    string
    httpClient *http.Client
    username   string
}

func NewClient(baseURL string) *Client {
    return &Client{
        baseURL:    baseURL,
        httpClient: &http.Client{},
    }
}

func (c *Client) Register(username string) error {
    data := map[string]string{"username": username}
    jsonData, _ := json.Marshal(data)
    
    resp, err := c.httpClient.Post(c.baseURL+"/api/register", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("registration failed: %s", resp.Status)
    }
    
    c.username = username
    return nil
}

func (c *Client) CreatePost(title, content, subreddit string) error {
    data := map[string]string{
        "title": title,
        "content": content,
        "subreddit": subreddit,
    }
    jsonData, _ := json.Marshal(data)
    
    req, _ := http.NewRequest("POST", c.baseURL+"/api/posts", bytes.NewBuffer(jsonData))
    req.Header.Set("Username", c.username)
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    return nil
}

func (c *Client) CreateSubreddit(name string) error {
    data := map[string]string{"name": name}
    jsonData, _ := json.Marshal(data)
    
    req, _ := http.NewRequest("POST", c.baseURL+"/api/subreddits", bytes.NewBuffer(jsonData))
    req.Header.Set("Username", c.username)
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    return nil
}

func (c *Client) SendDM(toUser, message string) error {
    data := map[string]string{
        "to_user": toUser,
        "message": message,
    }
    jsonData, _ := json.Marshal(data)
    
    req, _ := http.NewRequest("POST", c.baseURL+"/api/dm", bytes.NewBuffer(jsonData))
    req.Header.Set("Username", c.username)
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    return nil
}

func (c *Client) VotePost(postID string, isUpvote bool) error {
    data := map[string]interface{}{
        "post_id": postID,
        "is_upvote": isUpvote,
    }
    jsonData, _ := json.Marshal(data)
    
    req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/posts/%s/vote", c.baseURL, postID), bytes.NewBuffer(jsonData))
    req.Header.Set("Username", c.username)
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}

func (c *Client) GetKarma() (int32, error) {
    req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/users/%s/karma", c.baseURL, c.username), nil)
    req.Header.Set("Username", c.username)
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()
    
    var karma int32
    if err := json.NewDecoder(resp.Body).Decode(&karma); err != nil {
        return 0, err
    }
    return karma, nil
}


func (c *Client) JoinSubreddit(subredditName string) error {
    req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/subreddits/%s/join", c.baseURL, subredditName), nil)
    req.Header.Set("Username", c.username)
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}

func (c *Client) LeaveSubreddit(subredditName string) error {
    req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/subreddits/%s/leave", c.baseURL, subredditName), nil)
    req.Header.Set("Username", c.username)
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}

func (c *Client) GetFeed() error {
    req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/feed", c.baseURL), nil)
    req.Header.Set("Username", c.username)
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}

func (c *Client) RepostContent(postID string) error {
    data := map[string]string{"post_id": postID}
    jsonData, _ := json.Marshal(data)
    
    req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/posts/%s/repost", c.baseURL, postID), bytes.NewBuffer(jsonData))
    req.Header.Set("Username", c.username)
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}

func (c *Client) CreateComment(postID string, content string, parentCommentID *string) error {
    data := map[string]interface{}{
        "content": content,
        "post_id": postID,
    }
    if parentCommentID != nil {
        data["parent_comment_id"] = *parentCommentID
    }
    
    jsonData, _ := json.Marshal(data)
    
    req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/posts/%s/comments", c.baseURL, postID), bytes.NewBuffer(jsonData))
    req.Header.Set("Username", c.username)
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}


func (c *Client) GetComments(postID string) error {
    req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/posts/%s/comments", c.baseURL, postID), nil)
    req.Header.Set("Username", c.username)
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}