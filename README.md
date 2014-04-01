# steamstatus-api

`steamstatus-api` is a REST-like API wrapper for Steam profiles. It only supports
one method and gets its data by scraping the webpage.

`steamstatus-api` supports multiple usernames being queried at once and will request
them in parallel. In testing, each response from steam takes 400-600ms and
`steamstatus-api` is able to return all responses within that time.

## Usage

```
$ curl 'http://steamstatus-api/status?usernames=foxhop,chuckbang,japherwocky'
[{"url":"http://steamcommunity.com/id/foxhop/","steamid":"76561197960708678","personaname":"Foxhop","summary":"No information given.","ingame":""},{"url":"http://steamcommunity.com/id/japherwocky/","steamid":"76561198049551053","personaname":"japherwocky","summary":"No information given.","ingame":""},{"url":"http://steamcommunity.com/id/chuckbang/","steamid":"76561197961485911","personaname":"chuck!","summary":"No information given.","ingame":"Counter-Strike: Global Offensive"}]
```

## Configuration and running go-bot

1. `cp .env.sample .env`
1. `$EDITOR .env` and set your environment variables
1. `go run main.go`

### Environment variables

Variable | Description | Example
-------- | ----------- | -------
`PORT` | The port to listen for API requests on | `5000`

