package handlers

import (
	"encoding/json"
	"net/http"
	"reddit-clone/api/models"
	"reddit-clone/engine/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type UserHandler struct {
	system *actor.ActorSystem
	root   *actor.PID
}

func NewUserHandler(system *actor.ActorSystem, root *actor.PID) *UserHandler {
	return &UserHandler{system: system, root: root}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	registerMsg := &messages.RegisterAccount{Username: req.Username}
	h.system.Root.Send(h.root, registerMsg)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "User registered successfully",
		"username": req.Username,
	})
}

func (h *UserHandler) SendDM(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("Username")
	var req models.DMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dmMsg := &messages.SendDirectMessage{
		From:    username,
		To:      req.ToUser,
		Content: req.Message,
	}
	h.system.Root.Send(h.root, dmMsg)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "DM sent successfully",
	})
}
