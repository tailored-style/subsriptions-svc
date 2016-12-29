# Subscriptions Service

## Getting Started

This service should work with minimal setup

1. Setup your .env file
  * `cp .env.dev .env`
  * If necessary, add your `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` values
1. Install gin
  * `go get github.com/codegangsta/gin`
2. Run with gin (ensure your pwd is this cloned repo)
  * `gin`
3. Test the server
  * `curl http://127.0.0.1:3000`
