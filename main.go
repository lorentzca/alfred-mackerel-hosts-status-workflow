package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	mkr "github.com/mackerelio/mackerel-client-go"
)

type Items struct {
	Items []Item `json:"items"`
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

func main() {
	var hostname *string = flag.String("h", "example", "Hostname")
	var apikey *string = flag.String("i", "XXXXXXXX", "Mackerel API Key")
	flag.Parse()

	client := mkr.NewClient(*apikey)
	hosts, _ := client.FindHosts(&mkr.FindHostsParam{
		Statuses: []string{"working", "standby", "maintenance", "poweroff"},
	})

	org, _ := client.GetOrg()

	var items []Item
	for _, v := range hosts {
		if strings.Contains(v.Name, *hostname) {
			url := "https://mackerel.io/orgs/" + org.Name + "/hosts/" + v.ID
			items = append(items, Item{Title: v.Name, Subtitle: v.Status, Arg: url, Icon: icon{Path: v.Status + ".png"}})
		}
	}

	jsonBytes, _ := json.Marshal(Items{Items: items})
	fmt.Println(string(jsonBytes))
}
