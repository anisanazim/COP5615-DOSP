package actors

import (
    "github.com/asynkron/protoactor-go/actor"
    "reddit-clone/engine/messages"
	"fmt"
	"strings"
)

type PostActor struct {
    content  string
    author   string
    comments map[string][]*messages.Comment
    votes    int
}

func NewPostActor(content, author string) *PostActor {
    return &PostActor{
        content: content,
        author:  author,
		comments: make(map[string][]*messages.Comment),
		votes:    0,
    }
}

func (p *PostActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
    case *messages.Comment:
        commentID := generateUniqueCommentID()
        parentID := msg.ParentCommentID
        if parentID == nil {
            parentID = new(string) // Use an empty string for top-level comments
        }
        msg.PostID = commentID
		p.comments[*parentID] = append(p.comments[*parentID], msg)
        fmt.Printf("User %s commented on %s: %s\n", msg.Username, *parentID, msg.Content)
		
	case *messages.Vote:
    var karmaChange int32 = 1
    if msg.IsUpvote {
        p.votes++
        // Update post author's karma
        context.Send(context.Parent(), &messages.UpdateKarma{
            Username: p.author,
            Amount:   karmaChange,
        })
    } else {
        p.votes--
        context.Send(context.Parent(), &messages.UpdateKarma{
            Username: p.author,
            Amount:   -1,
        })
    }

	case *messages.PrintComments:
        fmt.Println("Printing comments for post:")
        printComments(p.comments, "", 0) // Start with top-level comments
    }
}

func generateUniqueCommentID() string {
    // Implement a function to generate unique IDs for comments
    return "unique-comment-id"
}

func printComments(comments map[string][]*messages.Comment, parentID string, level int) {
    for _, comment := range comments[parentID] {
        fmt.Printf("%sUser %s commented: %s\n", strings.Repeat("  ", level), comment.Username, comment.Content)
        printComments(comments, comment.PostID, level+1)
    }
}