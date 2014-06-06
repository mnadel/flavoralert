package flavor

import (
    "time"
    "strings"

    "github.com/PuerkitoBio/goquery"

    "appengine"
    "appengine/datastore"
)

type Alert struct {
    Created time.Time
    Flavor string
    User string
    AlertedOn time.Time
}

type Data struct {
    Location string `json:"location"`
    Flavors []string `json:"flavors"`
    Url string `json:"-"` // ignore
}

func (alert *Alert) Create(ctx appengine.Context) error {
    _, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "alert", nil), alert)
    if err != nil {
        ctx.Errorf("Error creating alert: %s", err.Error())
        return err
    }

    return nil
}

func (alert *Alert) Delete(ctx appengine.Context) error {
    q := datastore.NewQuery("alert").
        Filter("Flavor =", alert.Flavor).
        Filter("AlertedOn =", 0).
        Filter("User =", alert.User).
        KeysOnly()

    var matches []datastore.Key
    _, err := q.GetAll(ctx, &matches)
    if err != nil {
        ctx.Errorf("Error querying %v: %s", q, err.Error())
        return err
    }

    var anError error

    for _, key := range matches {
        err := datastore.Delete(ctx, &key)
        if err != nil {
            ctx.Errorf("Error deleting %v: %s", key, err.Error())
            anError = err
        }
    }

    return anError
}

func getAllFlavors(ctx appengine.Context) []string {
    url := "http://www.yogurtlabs.com/flavors/"

    html, err := fetch(ctx, url)
    if err != nil {
        ctx.Errorf("Unable to fetch [%s]: %s", url, err.Error())
        return nil
    }

    doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
    if err != nil {
        ctx.Errorf("Unable to parse [%s]: %s", url, err.Error())
        return nil
    }

    flavors := make([]string, 0)

    doc.Find(".flvName").Each(func(i int, s *goquery.Selection) {
        flavors = append(flavors, s.Text())
    })

    return flavors
}

func getCurrentFlavors(ctx appengine.Context, data Data, ch chan<- Data) {
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