
package main

import (
    "encoding/json"
    "flag"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "regexp"
    "strings"
    "time"
)

type ProfileData struct {
    Url             string  `json:"url"`
    SteamId         string  `json:"steamid"`
    PersonaName     string  `json:"personaname"`
    Summary         string  `json:"summary"`
    InGame          string  `json:"ingame"`
}

func GetProfile(username string) ProfileData {
    // Download the profile from steam
    profile_url := "http://steamcommunity.com/id/" + username + "/"
    response, err := http.Get(profile_url)
    defer response.Body.Close()
    if err != nil {
        panic(err)
    }
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        panic(err)
    }

    // Parse profile data
    var profile ProfileData
    json_regex := regexp.MustCompile(`g_rgProfileData = (.*);`)
    json_matches := json_regex.FindStringSubmatch(string(body))

    if len(json_matches) > 0 {
        if err = json.Unmarshal([]byte(json_matches[1]), &profile); err != nil {
            panic(err)
        }
    }

    // Find out if user is in a game
    ingame_regex := regexp.MustCompile(`<div class="profile_in_game_header">(.*)</div>`)
    ingame_matches := ingame_regex.FindStringSubmatch(string(body))

    var ingame bool = false
    if len(ingame_matches) > 0 && ingame_matches[1] == "Currently In-Game" {
        ingame = true

        // Find out which game
        gamename_regex := regexp.MustCompile(`<div class="profile_in_game_name">(.*)</div>`)
        gamename_matches := gamename_regex.FindStringSubmatch(string(body))

        // Add the game name to ProfileData
        if ingame && len(gamename_matches) > 0 {
            profile.InGame = gamename_matches[1]
        }
    }

    return profile
}

func FetchProfile(username string, c chan ProfileData) {
    go func() { c <- GetProfile(username) }()
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: Error handling on bad input

    r.ParseForm()
    usernames := strings.Split(r.Form.Get("usernames"), ",")
    user_count := len(usernames)
    var profiles = make([]ProfileData, user_count)

    // Request multiple users at once
    var profile_c = make(chan ProfileData)
    for _, username := range usernames {
        FetchProfile(username, profile_c)
    }

    // Wait for responses
    timeout := time.After(1000 * time.Millisecond)
    for idx := 0; idx < user_count; idx++ {
        select {
        case profile := <- profile_c:
            profiles[idx] = profile
        case <- timeout:
            log.Print("Timed out!")
            break
        }
    }

    // Assemble and send response
    log.Printf("profiles: %v", profiles)
    profile_json, err := json.Marshal(profiles)
    if err != nil {
        panic(err)
    }
    w.Write(profile_json)
}

func main() {
    flag.Parse()
    var assigned_port string
    if assigned_port = os.Getenv("PORT"); assigned_port == "" {
        assigned_port = "5000"
    }


    http.HandleFunc("/status", ProfileHandler)

    s := &http.Server{
        Addr:           ":" + assigned_port,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    log.Fatal(s.ListenAndServe())
}

