package flavor

import (
	"fmt"
	"strings"
	"net/http"
    "io/ioutil"
	"encoding/json"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"

    "appengine"
    "appengine/urlfetch"
)

const VERSION_0 = "/version/0"

func init() {
	router := mux.NewRouter()

	router.HandleFunc(VERSION_0 + "/list", listHandler)

	http.Handle("/", router)
}

type Data struct {
	Location string `json:"location"`
	Flavors []string `json:"flavors"`
	Url string `json:"-"` // ignore
}

func listHandler(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

    locations := make(map[string]string)
    locations["Uptown"] = "http://www.yogurtlabs.com/locations/uptown/"
    locations["Apple Valley"] = "http://www.yogurtlabs.com/locations/apple-valley/"
    locations["Hopkins"] = "http://www.yogurtlabs.com/locations/hopkins/"
    locations["Wayzata"] = "http://www.yogurtlabs.com/locations/wayzata/"
    locations["Eagan"] = "http://www.yogurtlabs.com/locations/eagan/"
    locations["U of M"] = "http://www.yogurtlabs.com/locations/u-of-m-stadium-village/"
    locations["IDS"] = "http://www.yogurtlabs.com/locations/ids-tower-downtown-mpls/"
    locations["Edina"] = "http://www.yogurtlabs.com/locations/ids-tower-downtown-mpls/"
    locations["Lake Calhoun"] = "http://www.yogurtlabs.com/locations/lake-calhoun/"

    var counter int
    flavorChan := make(chan Data)

    for location, url := range locations {
    	counter++
    	go getFlavors(ctx, Data { Location: location, Url: url, }, flavorChan)
    }

    flavors := make(map[string][]string)

    for i := 0; i < counter; i++ {
    	data := <-flavorChan
    	flavors[data.Location] = data.Flavors
    }

    emit(ctx, res, flavors, "success")
}

func getFlavors(ctx appengine.Context, data Data, ch chan<- Data) {
	html, err := fetch(ctx, data.Url)
	if err != nil {
		ctx.Errorf("Unable to fetch [%s]: %s", data.Url, err.Error())
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		ctx.Errorf("Unable to parse [%s]: %s", data.Url, err.Error())
		return
	}

	flavors := make([]string, 0)

	doc.Find(".flavorsOnTap").Each(func(i int, s *goquery.Selection) {
        s.Find(".bdcpy").Each(func(i2 int, s2 *goquery.Selection) {
        	flavor := s2.Text()
	        flavors = append(flavors, flavor)
	    })
	})

	data.Flavors = flavors

	ch <- data
}

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
        ctx.Errorf("Error stringing [%s]: %s", resp.Body, ioerr.Error())
        return "", ioerr
    }

    str := string(data[:])

    return str, nil
}