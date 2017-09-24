# King's College Chapel Ticketing API
> Krishna Kumar

[![CircleCI](https://circleci.com/gh/kings-cam/ticketing-api.svg?style=svg)](https://circleci.com/gh/kings-cam/ticketing-api)

# Clone repository
* Clone the repo: `git clone git@github.com:kings-cam/ticketing-api tickets`
* `cd` to the repo: `cd tickets`

# Install dependencies

```
              go get github.com/rs/cors && \
              go get github.com/gorilla/mux && \
              go get gopkg.in/mgo.v2 && \
              go get github.com/urfave/negroni && \
              go get github.com/stretchr/testify/assert
```

# Compile

`go build cmd/api/server.go`

# Run

`./server`

# Stats

```
curl http://localhost:4000/api/v1/stats | python -m "json.tool"
```
