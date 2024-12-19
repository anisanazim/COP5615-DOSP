package handlers

import (
	"encoding/json"
	"net/http"
	"reddit-clone/api/models"
	"reddit-clone/engine/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type PostHandler struct {
	system *actor.ActorSystem
	root   *actor.PID
}

func NewPostHandler(system *actor.ActorSystem, root *actor.PID) *PostHandler {
	return &PostHandler{system: system, root: root}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("Username")
	var req models.PostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	postMsg := &messages.Post{
		SubredditName: req.Subreddit,
		Username:      username,
		Content:       req.Content,
	}
	h.system.Root.Send(h.root, postMsg)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Post created successfully",
	})
}

func (h *PostHandler) Vote(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("Username")
	var req models.VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	voteMsg := &messages.Vote{
		PostID:   req.PostID,
		Username: username,
		IsUpvote: req.Vote > 0,
	}
	h.system.Root.Send(h.root, voteMsg)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Vote recorded successfully",
	})
}
