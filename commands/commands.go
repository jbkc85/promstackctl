package commands

import "github.com/spf13/cobra"

// PromStackCmd ...
var PromStackCmd = &cobra.Command{
	Use:           "promstackctl",
	Short:         "promstackctl is a command line script to interact with the PromStack Monitoring Suite.",
	Long:          `PromStack is a collection of software solutions to simplify ones introduction into monitoring, logging and alerting by providing a pre-configured stack of default software around Prometheus.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

// where should this go...

var monitorRootCmd = &cobra.Command{
	Use:   "monitor",
	Short: "monitor",
	Long:  `monitor`,
}

var removeRootCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove",
	Long:  `remove`,
}

func init() {
	globalCmdFlags()

	PromStackCmd.AddCommand(healthRootCmd)
	PromStackCmd.AddCommand(getRootCmd)
	PromStackCmd.AddCommand(describeRootCmd)

	// need to find a place to put these..
	PromStackCmd.AddCommand(monitorRootCmd)
	PromStackCmd.AddCommand(removeRootCmd)
	monitorRootCmd.AddCommand(serversMonitorCmd)
	removeRootCmd.AddCommand(serversUnmonitorCmd)

	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initializeConsul)
	cobra.OnInitialize(initializePrometheus)
}
