package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	mkr "github.com/mackerelio/mackerel-client-go"
)

type Items struct {
	Item []Item `json:"items"`
}

type Item struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Arg      string `json:"arg"`
	Icon     icon   `json:"icon"`
}

type icon struct {
	Path string `json:"path"`
}

func parseFlag() (string, string) {
	var hostname *string = flag.String(
		"h", "example", "Hostname")
	var apikey *string = flag.String(
		"i", "XXXXXXXX", "Mackerel API Key")

	flag.Parse()
	return *apikey, *hostname
}

func getOrg(a string) string {
	client := mkr.NewClient(a)
	org, _ := client.GetOrg()

	return org.Name
}

func getHosts(a string) []*mkr.Host {
	client := mkr.NewClient(a)
	hosts, _ := client.FindHosts(&mkr.FindHostsParam{
		Statuses: []string{
			"working",
			"standby",
			"maintenance",
			"poweroff"},
	})

	return hosts
}

func collectItem(hosts []*mkr.Host, h string, org string) []Item {
	var items []Item
	for _, v := range hosts {
		if strings.Contains(v.Name, h) {
			url := "https://mackerel.io/orgs/" + org +
				"/hosts/" + v.ID
			items = append(items, Item{
				Title:    v.Name,
				Subtitle: v.Status,
				Arg:      url,
				Icon:     icon{Path: v.Status + ".png"}})
		}
	}

	return items
}

func main() {
	a, h := parseFlag()
	org := getOrg(a)
	hosts := getHosts(a)
	items := collectItem(hosts, h, org)

	jsonBytes, _ := json.Marshal(Items{Item: items})
	fmt.Println(string(jsonBytes))
}
