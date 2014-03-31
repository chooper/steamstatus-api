
package main

import (
    "flag"
    "log"
    "net/http"
    "os"
    "github.com/chooper/steam-status/web"
    "time"
)

func main() {
    flag.Parse()
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

