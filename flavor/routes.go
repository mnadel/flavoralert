package flavor

import (
    "net/http"

    "github.com/gorilla/mux"

    "appengine"
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
    router.HandleFunc(VERSION_0 + "/list", listHandler)

    http.Handle("/", router)
}

func listHandler(res http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)

    flavorChan := make(chan Data)

    for location, url := range locations {
    	go getFlavors(ctx, Data { Location: location, Url: url, }, flavorChan)
    }

    flavors := make(map[string][]string)

    for i := 0; i < len(locations); i++ {
    	data := <-flavorChan
	    flavors[data.Location] = data.Flavors
    }

    emit(ctx, res, flavors, "success")
}
