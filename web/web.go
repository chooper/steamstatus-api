
package web

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/chooper/steamstatus-api/profiles"
    "strings"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: Error handling on bad input

    r.ParseForm()
    usernames := strings.Split(r.Form.Get("usernames"), ",")
    profiles := profiles.FetchProfiles(usernames)

    // Assemble and send response
    log.Printf("profiles: %v", profiles)
    profile_json, err := json.Marshal(profiles)
    if err != nil {
        http.Error(w, err.Error(), 500)
    }
    w.Write(profile_json)
}

