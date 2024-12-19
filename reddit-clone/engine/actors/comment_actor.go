package actors

import (
    "github.com/asynkron/protoactor-go/actor"
	"reddit-clone/engine/messages"
)

type CommentActor struct {
    content  string
    author   string
    replies  map[string]*actor.PID // Map reply ID to actor PID
}

func NewCommentActor(content, author string) *CommentActor {
    return &CommentActor{
        content: content,
        author:  author,
        replies: make(map[string]*actor.PID),
    }
}

func (c *CommentActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
    case *messages.Comment:
        // Handle nested replies by spawning new CommentActors as children
        replyProps := actor.PropsFromProducer(func() actor.Actor { return NewCommentActor(msg.Content, msg.Username) })
        replyPID := context.Spawn(replyProps)
        
        replyID := generateUniqueReplyID()
        c.replies[replyID] = replyPID
    }
}

func generateUniqueReplyID() string {
    // Implement a function to generate unique IDs for replies
    return "unique-reply-id"
}