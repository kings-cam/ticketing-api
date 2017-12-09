# King's College Chapel Ticketing API
> Krishna Kumar

[![CircleCI](https://circleci.com/gh/kings-cam/ticketing-api.svg?style=svg)](https://circleci.com/gh/kings-cam/ticketing-api)

# Clone repository
* Clone the repo: `git clone git@github.com:kings-cam/ticketing-api tickets`
* `cd` to the repo: `cd tickets`

# Install dependencies

```
go get -u github.com/golang/dep/cmd/dep
dep ensure
dep ensure -update
```

# Compile

`go build cmd/api.go`

# Run

## Launch mongo
`mkdir ./db && mongod --dbpath ./db`

## Run api server
`MongoUser=<username> MongoPW=<passwd> MongoPort=27017 IP=127.0.0.1 Port=4000 ./api `

# Stats

```
curl http://localhost:4000/api/v1/stats | python -m "json.tool"
```


# Config
```
curl -X POST -H "Content-Type: application/json" -d @config.json http://localhost:4000/api/v1/config/dates
curl -X POST -H "Content-Type: application/json" -d @price.json http://localhost:4000/api/v1/config/prices
```
