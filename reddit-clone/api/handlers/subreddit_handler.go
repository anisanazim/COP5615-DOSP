package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/asynkron/protoactor-go/actor"
    "reddit-clone/api/models"
    "reddit-clone/engine/messages"
)

type SubredditHandler struct {
    system *actor.ActorSystem
    root   *actor.PID
}

func NewSubredditHandler(system *actor.ActorSystem, root *actor.PID) *SubredditHandler {
    return &SubredditHandler{system: system, root: root}
}

func (h *SubredditHandler) CreateSubreddit(w http.ResponseWriter, r *http.Request) {
    var req models.SubredditRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    createMsg := &messages.CreateSubreddit{
        Name: req.Name,
    }
    h.system.Root.Send(h.root, createMsg)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Subreddit created successfully",
        "name": req.Name,
    })
}

func (h *SubredditHandler) JoinSubreddit(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    subredditName := vars["name"]
    username := r.Header.Get("Username") // Assuming username is passed in header

    joinMsg := &messages.JoinSubreddit{
        Username: username,
        SubredditName: subredditName,
    }
    h.system.Root.Send(h.root, joinMsg)

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Successfully joined subreddit",
    })
}

func (h *SubredditHandler) LeaveSubreddit(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    subredditName := vars["name"]
    username := r.Header.Get("Username")

    leaveMsg := &messages.LeaveSubreddit{
        Username: username,
        SubredditName: subredditName,
    }
    h.system.Root.Send(h.root, leaveMsg)

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Successfully left subreddit",
    })
}