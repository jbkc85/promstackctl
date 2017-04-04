package commands

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

var healthRootCmd = &cobra.Command{
	Use:   "health",
	Short: "Check Health of Endpoints for PromStack",
	Long:  `Check all endpoints in the PromStack Suite for connectivity.`,
	Run:   healthCheck,
}

func healthCheck(cmd *cobra.Command, args []string) {
	table := uitable.New()
	table.Wrap = true

	consulHealth, consulErr := commandCfg.Consul.health()
	prometheusHealth, prometheusErr := commandCfg.Prometheus.health()

	table.AddRow("COMPONENT", "URL", "STATUS", "MESSAGE")
	table.AddRow("consul", commandCfg.Consul.connectionString(), consulHealth, consulErr)
	table.AddRow("prometheus", commandCfg.Prometheus.connectionString(), prometheusHealth, prometheusErr)

	fmt.Println(table)
}
