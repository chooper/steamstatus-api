
package main

import (
    "encoding/json"
	"flag"
    "io/ioutil"
	"log"
    "net/http"
    "os"
    "regexp"
    "time"
)

type ProfileData struct {
    Url string `json:"url"`
    SteamId string `json:"steamid"`
    PersonaName string `json:"personaname"`
    Summary string `json:"summary"`
    InGame string
}

func GetProfile(username string) ProfileData {
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

    // Get my profile info
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
    }

    // Find out which game
    gamename_regex := regexp.MustCompile(`<div class="profile_in_game_name">(.*)</div>`)
    gamename_matches := gamename_regex.FindStringSubmatch(string(body))

    // Add the game name to my ProfileData
    if ingame && len(gamename_matches) > 0 {
        profile.InGame = gamename_matches[1]
    }
    return profile
}

func FetchProfile(username string) chan ProfileData {
    // Make 3 parallel attempts to fetch a user's profile
    c := make(chan ProfileData)
    for i := 0; i < 3; i++ {
        go func() { c <- GetProfile(username) }()
    }
    return c
}

func main() {
	flag.Parse()

    if len(os.Args) < 2 {
        panic("Not enough args")
    }

    start := time.Now()
    profile_c := FetchProfile(os.Args[1])

    select {
    case profile := <- profile_c:
        log.Print("Prof: %v", profile)
    case <- time.After(600 * time.Millisecond):
        log.Print("Timed out!")
    }

    elapsed := time.Since(start)
    log.Printf("Execution took %s", elapsed)
}
