# NunDB Golang Client

NunDB is a realtime key-value database written in Rust by
[@mateusfreira](https://github.com/mateusfreira).\
This repository provides a Golang client for interacting with a NunDB
server.

> For more details, check out the official [NunDB
> repository](https://github.com/mateusfreira/nun-db).

------------------------------------------------------------------------

## Installation

``` bash
go get github.com/kqnd/nun-db-go@v0.1.2
```

------------------------------------------------------------------------

## Basic Usage

``` go
package main

import (
    "fmt"
    nundb "github.com/kqnd/nun-db-go"
)

func main() {
    client, err := nundb.NewClient("ws://localhost:3012", "username", "password")
    if err != nil {
        panic(err)
    }

    // Create or use a database
    client.CreateDatabase("mydb", "mypass")
    client.UseDatabase("mydb", "mypass")

    // Set a value
    client.Set("foo", "bar")

    // Get a value
    val, _ := client.Get("foo")
    fmt.Println("foo:", val)

    // Watch for changes
    client.Watch("foo", func(v interface{}) {
        fmt.Println("foo changed:", v)
    })

    // Increment a numeric value
    client.Increment("counter", 1)
}
```

------------------------------------------------------------------------

## Features

-   Authentication (`auth` command)
-   Database management:
    -   `create-db`
    -   `use-db`
    -   `debug list-dbs`
-   Key-value operations:
    -   `set`
    -   `get`
    -   `remove`
    -   `increment`
-   Key listing:
    -   `keys`
    -   Filter by prefix or suffix
-   Watchers:
    -   `watch`
    -   `unwatch`
    -   `unwatch-all`

------------------------------------------------------------------------

## Example Commands

``` go
// Set and Get
client.Set("name", "Alice")
val, _ := client.Get("name")
fmt.Println(val) // Alice

// Remove key
client.Remove("name")

// Get all keys
keys, _ := client.GetAllKeys()
fmt.Println(keys)

// Get keys starting with "user:"
keys, _ := client.GetKeysStartingWith("user:")

// Get keys ending with ":id'
keys, _ := client.GetKeysEndingWith(":id")
```

------------------------------------------------------------------------

## License

MIT
