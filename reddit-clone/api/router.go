package api

import (
    "net/http"
    "github.com/gorilla/mux"
    corsHandlers "github.com/gorilla/handlers"
    "reddit-clone/api/handlers"
)

func SetupRouter(userHandler *handlers.UserHandler, postHandler *handlers.PostHandler, 
    subredditHandler *handlers.SubredditHandler, commentHandler *handlers.CommentHandler) http.Handler {
    r := mux.NewRouter()

    // User routes
    r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
    r.HandleFunc("/api/dm", userHandler.SendDM).Methods("POST")

    // Post routes
    r.HandleFunc("/api/posts", postHandler.CreatePost).Methods("POST")
    r.HandleFunc("/api/posts/{id}/vote", postHandler.Vote).Methods("POST")

    // Comment routes
    r.HandleFunc("/api/posts/{postId}/comments", commentHandler.CreateComment).Methods("POST")
    r.HandleFunc("/api/posts/{postId}/comments", commentHandler.GetComments).Methods("GET")
    r.HandleFunc("/api/posts/{postId}/comments/{commentId}", commentHandler.DeleteComment).Methods("DELETE")

    // Subreddit routes
    r.HandleFunc("/api/subreddits", subredditHandler.CreateSubreddit).Methods("POST")
    r.HandleFunc("/api/subreddits/{name}/join", subredditHandler.JoinSubreddit).Methods("POST")
    r.HandleFunc("/api/subreddits/{name}/leave", subredditHandler.LeaveSubreddit).Methods("POST")

    return corsHandlers.CORS(
        corsHandlers.AllowedOrigins([]string{"*"}),
        corsHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
        corsHandlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Username"}),
    )(r)
}