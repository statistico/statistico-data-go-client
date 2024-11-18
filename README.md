# Statistico Football Data Go Client

[![CircleCI](https://circleci.com/gh/statistico/statistico-football-data-go-grpc-client/tree/main.svg?style=shield)](https://circleci.com/gh/statistico/statistico-football-data-go-grpc-client/tree/main)

The library is a Golang wrapper around the [Statistico Football Data gRPC API](https://github.com/statistico/statistico-football-data).

## Installation
```.env
$ go get -u github.com/statistico/statistico-football-data-go-grpc-client
```
## Usage
```go
package main

import (
    "context"
    "fmt"
    "github.com/statistico/statistico-football-data-go-grpc-client/statisticofootballdata"
    "github.com/statistico/statistico-proto/go"
    "google.golang.org/grpc"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

    teamClient := statistico.NewTeamServiceClient(conn)

    client := statisticofootballdata.NewTeamClient(teamClient)
    
    team, err := client.ByID(context.Background(), 10) 

    if err != nil {
        fmt.Printf("%s\n", err.Error())
        return
    }

    // Handle team variable
}
```