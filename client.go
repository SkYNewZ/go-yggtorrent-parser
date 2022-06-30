package yggtorrent

//go:generate ifacemaker --file=client.go --struct=client --iface=Client --pkg=yggtorrent -y "Client interface describes wrapped YggTorrent client." --doc=true --output=generated.go
//go:generate stringer --type Category --linecomment
//go:generate stringer --type SubCategory --linecomment

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type client struct {
	baseURL string
}

// Category describes a main category.
type Category int

// SubCategory describes a sub category with a Category.
type SubCategory int

const (
	Movie SubCategory = iota // 2183
	TV                       // 2184
)

const (
	Video Category = iota // 2145
)

// New creates a new YggTorrent client.
func New(baseURL string) Client {
	return &client{baseURL: baseURL}
}

// Result is a search result.
type Result struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	PublishedAt time.Time `json:"published_at"`
	Size        string    `json:"size"`
	Seeders     uint      `json:"seeders"`
	Leechers    uint      `json:"leechers"`
	InfoURL     string    `json:"uri"`
	DownloadURL string    `json:"download_url"`
}

func dateStringToTime(str string) time.Time {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}

	return time.Unix(i, 0)
}

func strToUInt(str string) uint {
	r, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		panic(err)
	}

	return uint(r)
}

// ParseResults read given HTML data to search for torrents results.
func (c *client) ParseResults(data io.Reader) ([]*Result, error) {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return nil, fmt.Errorf("cannot parse HTML data: %w", err)
	}

	trim := func(str string) string {
		str = strings.ReplaceAll(str, " ", ".")
		str = strings.TrimRight(str, ".")
		return str
	}

	var results []*Result
	doc.Find("#\\#torrents").Each(func(_ int, sectionHTML *goquery.Selection) {
		sectionHTML.Find("table.table").Each(func(_ int, table *goquery.Selection) {
			table.Find("tbody").Each(func(_ int, tbody *goquery.Selection) {
				tbody.Find("tr").Each(func(_ int, row *goquery.Selection) {
					cells := row.Find("td")
					// 0: category
					// 1: name (contains a <a>)
					// 2: nfo
					// 3: comments
					// 4: age (containers a hidden div with the timestamp and a spam with the relative date)
					// 5: size
					// 6: completed
					// 7: seed
					// 8: leech

					link, _ := cells.Eq(1).Children().Attr("href")
					id := getIDFromLink(link)
					results = append(results, &Result{
						ID:          id,
						Name:        trim(cells.Eq(1).Text()),
						PublishedAt: dateStringToTime(trim(cells.Eq(4).Children().Eq(0).Text())),
						Size:        trim(cells.Eq(5).Text()),
						Seeders:     strToUInt(trim(cells.Eq(7).Text())),
						Leechers:    strToUInt(trim(cells.Eq(8).Text())),
						InfoURL:     link,
						DownloadURL: c.makeDownloadURL(id),
					})
				})
			})
		})
	})

	return results, nil
}

// SearchURL makes a YggTorrent search URL ready to be used through flaresolverr.
func (c *client) SearchURL(query string, category Category, subcategory SubCategory) string {
	return fmt.Sprintf("%s/engine/search?name=%s&category=%d&sub_category=%s&do=search", c.baseURL, query, category, subcategory)
}

func getIDFromLink(u string) string {
	idx := strings.LastIndex(u, "/")
	if idx == -1 {
		return ""
	}

	end := u[idx+1:]
	return strings.Split(end, "-")[0]
}

func (c *client) makeDownloadURL(id string) string {
	return fmt.Sprintf("%s/engine/download_torrent?id=%s", c.baseURL, id)
}
