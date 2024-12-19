package handlers

import (
	"encoding/json"
	"net/http"
	"reddit-clone/api/models"
	"reddit-clone/engine/messages"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/gorilla/mux"
)

type CommentHandler struct {
	system *actor.ActorSystem
	root   *actor.PID
}

func NewCommentHandler(system *actor.ActorSystem, root *actor.PID) *CommentHandler {
	return &CommentHandler{system: system, root: root}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]
	username := r.Header.Get("Username")

	var req models.CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	commentMsg := &messages.Comment{
		PostID:   postID,
		Content:  req.Content,
		Username: username,
	}
	h.system.Root.Send(h.root, commentMsg)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Comment created successfully",
	})
}

func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]

	getCommentsMsg := &messages.PrintComments{
		PostId: postID,
	}
	h.system.Root.Send(h.root, getCommentsMsg)

	// Note: In a real implementation, you'd want to wait for the response
	// and return the actual comments. This is a simplified version.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Comments retrieved successfully",
	})
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["commentId"]
	postID := vars["postId"]
	username := r.Header.Get("Username")

	deleteMsg := &messages.DeleteComment{
		PostID:    postID,
		CommentID: commentID,
		Username:  username,
	}
	h.system.Root.Send(h.root, deleteMsg)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Comment deleted successfully",
	})
}
