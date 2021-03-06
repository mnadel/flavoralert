package flavor

import (
    "time"
    "strings"

    "github.com/PuerkitoBio/goquery"

    "appengine"
    "appengine/datastore"
)

const (
    AlertEntity = "Alert"
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
    _, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, AlertEntity, nil), alert)
    if err != nil {
        ctx.Errorf("Error creating alert: %s", err.Error())
        return err
    }

    return nil
}

func (alert *Alert) List(ctx appengine.Context) map[string]string {
    q := datastore.NewQuery(AlertEntity).
        Filter("User =", alert.User)

    var alerts []Alert
    _, err := q.GetAll(ctx, &alerts)
    if err != nil {
        ctx.Errorf("Error querying %v: %s", q, err.Error())
        return nil
    }

    ctx.Debugf("Found %d alerts for %v", len(alerts), q)

    data := make(map[string]string)

    for _, alert := range alerts {
        data[alert.Flavor] = "true"
    }

    return data
}

func (alert *Alert) Delete(ctx appengine.Context) int {
    q := datastore.NewQuery(AlertEntity).
        Filter("Flavor =", alert.Flavor).
        Filter("User =", alert.User).
        Filter("AlertedOn =", alert.AlertedOn).
        KeysOnly()

    keys, err := q.GetAll(ctx, nil)
    if err != nil {
        ctx.Errorf("Error querying %v: %s", q, err.Error())
        return 0
    }

    ctx.Debugf("Found %d unsent alerts for %v", len(keys), q)

    count := 0

    for _, key := range keys {
        err := datastore.Delete(ctx, key)
        if err != nil {
            ctx.Errorf("Error deleting %v: %s", key, err.Error())
        } else {
            count += 1
        }
    }

    return count
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