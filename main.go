
package main

import (
    "flag"
    "log"
    "net/http"
    "os"
    "github.com/chooper/steamstatus-api/poller"
    "github.com/chooper/steamstatus-api/web"
    "strings"
    "time"
)

func usage() {
    log.Printf("Usage: %s [web|poller]")
}

func launch_web() {
    var assigned_port string
    if assigned_port = os.Getenv("PORT"); assigned_port == "" {
        assigned_port = "5000"
    }

    http.HandleFunc("/status", web.ProfileHandler)

    s := &http.Server{
        Addr:           ":" + assigned_port,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    log.Fatal(s.ListenAndServe())
}

func launch_poller() {
    var usernames []string
    if usernames = strings.Split(os.Getenv("POLL_USERNAMES"), ","); len(usernames) == 0 {
        log.Fatalf("Must set POLL_USERNAMES env var!")
    }

    c := make(chan poller.Notification)
    p := &poller.Poller{
        Usernames:      usernames,
        NotifyChan:     c,
    }
    go p.Loop()
    for {
        select {
            case n := <- c:
                log.Printf("N: %v", n)
        }
    }
}

func main() {
    flag.Parse()

    if len(os.Args) < 2 {
        usage()
        log.Fatal("See usage")
    }

    if os.Args[1] == "web" {
        launch_web()
    }
    if os.Args[1] == "poller" {
        launch_poller()
    }
}

