#!/bin/bash

# Test user registration
curl -X POST http://localhost:8081/api/register \
    -H "Content-Type: application/json" \
    -d '{"username": "testuser1"}'

# Test post creation
curl -X POST http://localhost:8081/api/posts \
    -H "Content-Type: application/json" \
    -d '{"title": "Test Post", "content": "This is a test post", "subreddit": "testsubreddit"}'

# Test voting
curl -X POST http://localhost:8081/api/posts/post1/vote \
    -H "Content-Type: application/json" \
    -d '{"vote": 1}'