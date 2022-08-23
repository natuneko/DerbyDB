package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	var DerbyName string
	fmt.Print("DerbyName: ")
	fmt.Scan(&DerbyName)

	DerbyName, _, err := transform.String(japanese.EUCJP.NewEncoder(), DerbyName)
	if err != nil {
		panic(err)
	}

	DerbyName = url.QueryEscape(DerbyName)

	url := "https://db.netkeiba.com/?pid=horse_list&word=" + DerbyName
	req, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	query, err := goquery.NewDocumentFromReader(req.Body)
	if err != nil {
		panic(err)
	}

	query.Find("table.db_h_race_results > tbody > tr").Each(func(i int, s *goquery.Selection) {
		data := s.Find("td")
		date := data.Eq(0).Text()
		RaceName, _, _ := transform.Bytes(japanese.EUCJP.NewDecoder(), []byte(data.Eq(4).Text()))
		Arrivals, _, _ := transform.Bytes(japanese.EUCJP.NewDecoder(), []byte(data.Eq(11).Text()))

		if string(Arrivals) == "中" {
			Arrivals = []byte("中止")
		} else {
			Arrivals = []byte(string(Arrivals) + "着")
		}

		fmt.Println(fmt.Sprintf("%s %s %s", date, string(RaceName), string(Arrivals)))
	})

}
