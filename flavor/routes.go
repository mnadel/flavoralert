package flavor

import (
    "time"
    "net/http"

    "github.com/gorilla/mux"

    "appengine"
    "appengine/user"
)

const VERSION_0 = "/version/0"

var locations map[string]string

func init() {
    locations = make(map[string]string)
    locations["Uptown"] = "http://www.yogurtlabs.com/locations/uptown/"
    locations["Apple Valley"] = "http://www.yogurtlabs.com/locations/apple-valley/"
    locations["Hopkins"] = "http://www.yogurtlabs.com/locations/hopkins/"
    locations["Wayzata"] = "http://www.yogurtlabs.com/locations/wayzata/"
    locations["Eagan"] = "http://www.yogurtlabs.com/locations/eagan/"
    locations["U of M"] = "http://www.yogurtlabs.com/locations/u-of-m-stadium-village/"
    locations["IDS"] = "http://www.yogurtlabs.com/locations/ids-tower-downtown-mpls/"
    locations["Edina"] = "http://www.yogurtlabs.com/locations/ids-tower-downtown-mpls/"
    locations["Lake Calhoun"] = "http://www.yogurtlabs.com/locations/lake-calhoun/"

    router := mux.NewRouter()
    router.HandleFunc("/", redirectHandler)
    router.HandleFunc(VERSION_0 + "/current", currentHandler)
    router.HandleFunc(VERSION_0 + "/all", allHandler)
    router.HandleFunc(VERSION_0 + "/meta", metaHandler)
    router.HandleFunc(VERSION_0 + "/alert/add/{flavor}", alertCreateHandler)
    router.HandleFunc(VERSION_0 + "/alert/remove/{flavor}", alertDeleteHandler)

    http.Handle("/", router)
}

func redirectHandler(res http.ResponseWriter, req *http.Request) {
    http.Redirect(res, req, "/web/index.html", 301)
}

func metaHandler(res http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)

    meta := make(map[string]string)
    meta["login_url"], _ = user.LoginURL(ctx, "/web/index.html")
    meta["logout_url"], _ = user.LogoutURL(ctx, "/web/index.html")
    usr := user.Current(ctx); if usr == nil {
        meta["authenticated"] = "false"
    } else {
        meta["authenticated"] = "true"
    }

    emit(ctx, res, meta, "success")
}

func alertDeleteHandler(res http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)
    usr := user.Current(ctx)

    vars := mux.Vars(req)
    flavor := vars["flavor"]

    alert := Alert {
        Flavor: flavor,
        User: usr.String(),
    }

    count := alert.Delete(ctx)
    emit(ctx, res, count, "success")
}

func alertCreateHandler(res http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)
    usr := user.Current(ctx)

    vars := mux.Vars(req)
    flavor := vars["flavor"]

    alert := Alert {
        Created: time.Now(),
        Flavor: flavor,
        User: usr.String(),
    }

    err := alert.Create(ctx)
    if err != nil {
        emit(ctx, res, err.Error(), "error")
    } else {
        emit(ctx, res, nil, "success")
    }
}

func allHandler(res http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)

    flavors := getAllFlavors(ctx)

    emit(ctx, res, flavors, "success")
}

func currentHandler(res http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)

    flavorChan := make(chan Data)

    for location, url := range locations {
        go getCurrentFlavors(ctx, Data { Location: location, Url: url, }, flavorChan)
    }

    flavors := make(map[string][]string)

    for i := 0; i < len(locations); i++ {
        data := <-flavorChan
        flavors[data.Location] = data.Flavors
    }

    emit(ctx, res, flavors, "success")
}