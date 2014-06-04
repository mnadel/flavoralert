package flavor

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"

    "appengine"
    "appengine/urlfetch"
)

func emit(ctx appengine.Context, res http.ResponseWriter, obj interface{}, status string) {
    var data = map[string]interface{} {
        "op": status,
        "data": obj,
    }

    res.Header().Set("Content-Type", "application/json")
    
    json, err := json.Marshal(data)
    if err != nil {
        ctx.Criticalf("Error encoding %s: %s", data, err.Error())
        fmt.Fprint(res, "{\"op\":\"Man down!\"}")
        return
    }

    fmt.Fprint(res, string(json[:]))
}

func fetch(ctx appengine.Context, url string) (string, error) {
    fetcher := urlfetch.Client(ctx)
    req, _ := http.NewRequest("GET", url, nil)

    resp, curlerr := fetcher.Do(req)
    defer resp.Body.Close()

    if curlerr != nil {
        ctx.Errorf("Error fetching [%s]: %s", url, curlerr.Error())
        return "", curlerr
    }

    data, ioerr := ioutil.ReadAll(resp.Body)
    if ioerr != nil {
        ctx.Errorf("Error stringifying [%s]: %s", resp.Body, ioerr.Error())
        return "", ioerr
    }

    str := string(data[:])

    return str, nil
}