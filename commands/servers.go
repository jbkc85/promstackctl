package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gosuri/uitable"
	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

var serversGetCmd = &cobra.Command{
	Use:   "servers",
	Short: "get servers from Consul Catalog",
	Long:  `get a list of Servers registered to the Consul Catalog.`,
	Run:   getServers,
}

var serversDescribeCmd = &cobra.Command{
	Use:   "server",
	Short: "describe a server from Consul Catalog",
	Long:  `describe a specific Server registered in the Consul Catalog.`,
	Run:   describeServer,
}

var serversMonitorCmd = &cobra.Command{
	Use:   "server",
	Short: "monitor a given server",
	Long:  `monitor a given server from a predefined exporter in the Consul Catalog.`,
	Example: `
	$ promstackctl monitor server --exporter.name node-exporter --node.name example.com --node.address 1.1.1.1`,
	Run: monitorServer,
}

var serversUnmonitorCmd = &cobra.Command{
	Use:   "server",
	Short: "unmonitor a given server",
	Long:  `unmonitor a previously registered server in the Consul Catalog.`,
	Example: `
	$ promstackctl remove server --node.name`,
	Run: unmonitorServer,
}

var (
	nodeNameFlag     string
	nodeAddressFlag  string
	exporterNameFlag string
)

func init() {
	serversMonitorCmd.PersistentFlags().StringVar(&nodeAddressFlag, "node.address", "", "IPv4 Address (or DNS if website) of node to monitor")
	serversMonitorCmd.PersistentFlags().StringVar(&nodeNameFlag, "node.name", "", "name of node to monitor")
	serversMonitorCmd.PersistentFlags().StringVar(&exporterNameFlag, "exporter.name", "", "name of exporter to monitor on node")

	serversUnmonitorCmd.PersistentFlags().StringVar(&nodeNameFlag, "node.name", "", "name of node to remove from PromStack monitoring.")
}

func getServers(cmd *cobra.Command, args []string) {
	catalog := commandCfg.Consul.client.Catalog()

	consulChan := make(chan []*consul.Node)
	go func() {
		nodes, _, err := catalog.Nodes(&consul.QueryOptions{})
		if err != nil {
			log.Printf("[ERROR] Unable to connect to Consul Catalog, message: %s", err)
		}
		consulChan <- nodes
	}()

	nodes := <-consulChan

	table := uitable.New()
	table.Wrap = true
	table.AddRow("SERVER", "ADDRESS")
	for _, n := range nodes {
		table.AddRow(n.Node, n.Address)
	}
	fmt.Println(table)
}

func describeServer(cmd *cobra.Command, args []string) {
	if args[0] == "" {
		log.Printf("[ERROR] Unable to describe Server, please use the server in question as an argument to the script.")
	} else {
		node := getNode(args[0])

		table := uitable.New()
		table.Wrap = true
		table.AddRow("Name:", node.Node.Node)
		table.AddRow("Address:", node.Node.Address)
		table.AddRow("Services:")
		table.AddRow("  ExporterName", "Endpoint")
		table.AddRow("  ============", "========")
		for _, s := range node.Services {
			table.AddRow("  "+s.Service, node.Node.Address+":"+strconv.Itoa(s.Port))
		}
		fmt.Println(table)
	}
}

func monitorServer(cmd *cobra.Command, args []string) {
	cont := true
	if nodeAddressFlag == "" {
		log.Fatalf("MISSING REQUIRED --node.address argument to provide a node address to monitor.")
		cont = false
	}
	if nodeNameFlag == "" {
		log.Fatalf("MISSING REQUIRED --node.name argument to provide a node name to monitor.")
		cont = false
	}
	if exporterNameFlag == "" {
		log.Fatalf("MISSING REQUIRED --exporter.name argument to provide an exporter to monitor.")
		cont = false
	}

	if cont != true {
		log.Fatalf("[ERROR] One or more required arguments are missing.  Please refer to the log output.")
	}

	exporterDetails := getKV("promstack/exporters/" + exporterNameFlag)

	if exporterDetails == nil {
		log.Fatalf("[ERROR] Exporter %s does not appear to be in the Consul KeyVal store.  Please add it, or to see an available list of exporters run:\n\t$ promstackctl get exporters", exporterNameFlag)
	} else {
		kvExporter := exporter{}
		if err := json.Unmarshal(exporterDetails.Value, &kvExporter); err != nil {
			log.Printf("[ERROR] Unable to unmarshal %v into structure, message: %s.", exporterDetails.Value, err)
		} else {
			newNode := &consul.CatalogRegistration{
				Node:    nodeNameFlag,
				Address: nodeAddressFlag,
				Service: &consul.AgentService{
					Service: exporterNameFlag,
					Port:    kvExporter.Port,
					Tags:    kvExporter.Tags,
				},
			}

			catalog := commandCfg.Consul.client.Catalog()

			meta, err := catalog.Register(newNode, &consul.WriteOptions{})

			if err != nil {
				log.Printf("[ERROR] Unable to Register %s to Catalog, message: %s", nodeNameFlag, err)
			}

			log.Printf("Exporter %s added to Node %s (request time %v).", newNode.Node, newNode.Service.Service, meta.RequestTime)

			node := getNode(newNode.Node)
			table := uitable.New()
			table.Wrap = true
			table.AddRow("Name:", node.Node.Node)
			table.AddRow("Address:", node.Node.Address)
			table.AddRow("Services:")
			table.AddRow("  ExporterName", "Endpoint")
			table.AddRow("  ============", "========")
			for _, s := range node.Services {
				table.AddRow("  "+s.Service, node.Node.Address+":"+strconv.Itoa(s.Port))
			}
			table.AddRow("")
			fmt.Println(table)

		}
	}

}

func unmonitorServer(cmd *cobra.Command, args []string) {
	if nodeNameFlag == "" {
		log.Fatalf("MISSING REQUIRED --node.name argument to provide a node name to monitor.")
		return
	}

	dereg := &consul.CatalogDeregistration{
		Node: nodeNameFlag,
	}

	catalog := commandCfg.Consul.client.Catalog()

	if _, err := catalog.Deregister(dereg, nil); err != nil {
		log.Fatalf("[ERROR] Unable to deregister %s, message: %s", nodeNameFlag, err)
	} else {
		fmt.Printf("remove %s from PromStack monitoring.\n", nodeNameFlag)
	}
}
