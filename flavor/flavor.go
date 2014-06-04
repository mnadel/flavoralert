package flavor

import (
    "strings"

    "github.com/PuerkitoBio/goquery"

    "appengine"
)

type Data struct {
    Location string `json:"location"`
    Flavors []string `json:"flavors"`
    Url string `json:"-"` // ignore
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