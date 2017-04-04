package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

type exporter struct {
	Name       string   `json:"name"`
	Port       int      `json:"port"`
	Repository string   `json:"repository"`
	Tags       []string `json:"tags"`
}

type exporterHub struct {
	Client    *http.Client
	URL       string     `json:"url"`
	Exporters []exporter `json:"exporters"`
}

var exporterGetCmd = &cobra.Command{
	Use:   "exporters",
	Short: "list available metadata for exporters",
	Long:  `list available metadata for exporters from the Consul KV store.`,
	Run:   getExporters,
}

var exporterHubFlag string

func init() {
	exporterGetCmd.PersistentFlags().StringVarP(&exporterHubFlag, "exporter.hub", "", "", "If presented will skip checking Consul for Inventory of Exporters and check the given URL at $URL/exporters.json")
}

func getExporters(cmd *cobra.Command, args []string) {
	table := uitable.New()
	table.Wrap = true

	table.AddRow("EXPORTER", "PORT", "TAGS")
	// if official exporters, skip on the Consul listings
	if exporterHubFlag != "" {
		hub := exporterHub{
			Client: &http.Client{Timeout: 3 * time.Second},
			URL:    exporterHubFlag,
		}
		file, err := hub.Client.Get(hub.URL + "/exporters.json")
		if err != nil {
			log.Printf("[ERROR] Unable to access %s, message: %s", hub.URL+"/exporters.json", err)
			return
		}
		defer file.Body.Close()
		buf, _ := ioutil.ReadAll(file.Body)

		json.Unmarshal(buf, &hub)

		for _, exp := range hub.Exporters {
			table.AddRow(exp.Name, exp.Port, exp.Tags)
		}
	} else {
		pairs := getKVPath("promstack/exporters")

		for _, pair := range pairs {
			kvExporter := exporter{}
			if err := json.Unmarshal(pair.Value, &kvExporter); err != nil {
				log.Printf("[ERROR] Unable to unmarshal %v into structure, message: %s.", pair.Value, err)
			} else {
				table.AddRow(strings.Replace(pair.Key, "promstack/exporters/", "", -1), kvExporter.Port, kvExporter.Tags)
			}
		}
	}
	fmt.Println(table)
}

func unmonitorExporter(cmd *cobra.Command, args []string) {

}
