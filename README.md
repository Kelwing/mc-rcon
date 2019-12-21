# mc-rcon

A Minecraft RCON library written in Go.

Forked from [github.com/bearbin/mcgorcon](https://github.com/bearbin/mcgorcon), cleaned up auth, and added tests.

## Installation
```
go get github.com/Kelwing/mc-rcon
```

## Basic usage
```go
package main

import (
    "log"

    mcrcon "github.com/Kelwing/mc-rcon"
)

func main() {
    conn := new(mcrcon.MCConn)
    err := conn.Open("localhost:25575", "testpw")
    if err != nil {
        log.Fatalln("Open failed", err)
    }
    defer conn.Close()

    err = conn.Authenticate()
    if err != nil {
        log.Fatalln("Auth failed", err)
    }
    
    resp, err = conn.SendCommand("tps")
    if err != nil {
        log.Fatalln("Command failed", err)
    }
    log.Println(resp)
}
```