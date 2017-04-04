package commands

import (
	"log"

	consul "github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/spf13/viper"
)

type config struct {
	Cluster    string
	Consul     consulConfig     `json:"consul"`
	Prometheus prometheusConfig `json:"prometheus"`
	Verbose    bool             `json:"verbose"`
}

type consulConfig struct {
	Address    string `json:"address"`
	client     *consul.Client
	Datacenter string `json:"datacenter"`
	Schema     string `json:"schema"`
	Port       string `json:"port"`
}

type prometheusConfig struct {
	Address string `json:"address"`
	client  prometheus.Client
	Schema  string `json:"schema"`
	Port    string `json:"port"`
}

var (
	commandCfg = config{
		Cluster: "local",
		Consul: consulConfig{
			Address:    defaultConsulAddress,
			Datacenter: defaultConsulDatacenter,
			Port:       "8500",
		},
		Prometheus: prometheusConfig{
			Address: "localhost",
			Port:    "9090",
			Schema:  "http",
		},
		Verbose: false,
	}
	parsedCfg config

	clusterFlag string
)

const (
	defaultCluster          = ""
	defaultConsulAddress    = "localhost"
	defaultConsulDatacenter = "promstack"
	defaultVerbose          = false
)

func globalCmdFlags() {
	PromStackCmd.PersistentFlags().StringVarP(&parsedCfg.Consul.Address, "consul.address", "", defaultConsulAddress, "Address to Consul API. [env:PROMSTACK_CONSUL_ADDRESS]")
	PromStackCmd.PersistentFlags().StringVar(&parsedCfg.Consul.Datacenter, "consul.datacenter", defaultConsulDatacenter, "Datacenter to reference in Consul API. [env:PROMSTACK_CONSUL_DATACENTER]")
	PromStackCmd.PersistentFlags().StringVar(&parsedCfg.Consul.Port, "consul.port", "8500", "Port to Consul API. [env:PROMSTACK_CONSUL_PORT]")
	PromStackCmd.PersistentFlags().StringVar(&parsedCfg.Consul.Schema, "consul.schema", "http", "Schema to access Consul API on. [env:PROMSTACK_CONSUL_SCHEMA]")

	PromStackCmd.PersistentFlags().StringVar(&parsedCfg.Prometheus.Address, "prometheus.address", "localhost", "Address to Prometheus. [env:PROMSTACK_PROMETHEUS_ADDRESS]")
	PromStackCmd.PersistentFlags().StringVar(&parsedCfg.Prometheus.Port, "prometheus.port", "9090", "Port to Prometheus API. [env:PROMSTACK_PROMETHEUS_PORT]")
	PromStackCmd.PersistentFlags().StringVar(&parsedCfg.Prometheus.Schema, "promtheus.schema", "http", "Schema to access Prometheus API on. [env:PROMSTACK_PROMETHEUS_SCHEMA]")

	PromStackCmd.PersistentFlags().StringVarP(&parsedCfg.Cluster, "cluster", "", defaultCluster, "Environment to use when loading in configuration variables. [env:PROMSTACK_ENVIRONEMNT]")
	//PromStackCmd.PersistentFlags().StringVar(&clusterFlag, "cluster", "local", "Environment to use when loading in configuration variables. [env:PROMSTACK_ENVIRONEMNT]")
	PromStackCmd.PersistentFlags().BoolVar(&parsedCfg.Verbose, "verbose", defaultVerbose, "Enable Verbose output of application. [env:PROMSTACK_VERBOSE]")
}

func initConfig() {
	// set config file defaults
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.promstack")
	viper.AddConfigPath("/etc/promstack")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	// set env defaults
	viper.SetEnvPrefix("PROMSTACK")

	// setting priority variables - Verbose
	viper.BindEnv("VERBOSE")
	if viper.GetBool("VERBOSE") {
		commandCfg.Verbose = true
	}
	if parsedCfg.Verbose != defaultVerbose {
		commandCfg.Verbose = parsedCfg.Verbose
	}
	// setting priority variables - Cluster
	viper.BindEnv("CLUSTER")
	if viper.GetString("CLUSTER") != "" {
		commandCfg.Cluster = viper.GetString("CLUSTER")
	}
	if parsedCfg.Cluster != "" {
		commandCfg.Cluster = parsedCfg.Cluster
	}

	err := viper.ReadInConfig()
	if err != nil {
		if commandCfg.Verbose {
			log.Printf("[DEBUG] No Configuration File Found (%s). Loading defaults.", err)
		}
	} else {
		// set prefix. If no prefix is provided, look for top-level configurations
		var prefix string
		if commandCfg.Cluster != "" {
			prefix = commandCfg.Cluster + "."
		}

		// set Consul from Config File
		if viper.GetString(prefix+"consul.address") != "" {
			commandCfg.Consul.Address = viper.GetString(prefix + "consul.address")
		}
		if viper.GetString(prefix+"consul.datacenter") != "" {
			commandCfg.Consul.Datacenter = viper.GetString(prefix + "consul.datacenter")
		}
		if viper.GetString(prefix+"prometheus.address") != "" {
			commandCfg.Prometheus.Address = viper.GetString(prefix + "prometheus.address")
		}
	}

	// ENV mappings (take priority over FILE)
	//-- CONSUL MAPS
	viper.BindEnv("CONSUL_ADDRESS")
	if viper.GetString("CONSUL_ADDRESS") != "" {
		commandCfg.Consul.Address = viper.GetString("CONSUL_ADDRESS")
	}
	viper.BindEnv("CONSUL_DATACENTER")
	if viper.GetString("CONSUL_DATACENTER") != "" {
		commandCfg.Consul.Datacenter = viper.GetString("CONSUL_DATACENTER")
	}
	viper.BindEnv("CONSUL_PORT")
	if viper.GetString("CONSUL_PORT") != "" {
		commandCfg.Consul.Port = viper.GetString("CONSUL_PORT")
	}
	//-- PROMETHEUS MAPS
	viper.BindEnv("PROMETHEUS_ADDRESS")
	if viper.GetString("PROMETHEUS_ADDRESS") != "" {
		commandCfg.Prometheus.Address = viper.GetString("PROMETHEUS_ADDRESS")
	}
	viper.BindEnv("PROMETHEUS_SCHEMA")
	if viper.GetString("PROMETHEUS_SCHEMA") != "" {
		commandCfg.Prometheus.Schema = viper.GetString("PROMETHEUS_SCHEMA")
	}
	viper.BindEnv("PROMETHEUS_PORT")
	if viper.GetString("PROMETHEUS_PORT") != "" {
		commandCfg.Prometheus.Port = viper.GetString("PROMETHEUS_PORT")
	}
	//-- GENERAL MAPS

	// Script ARG mappings (take priority of ENV/FILE)
	//-- CONSUL MAPS
	if parsedCfg.Consul.Address != defaultConsulAddress {
		commandCfg.Consul.Address = parsedCfg.Consul.Address
	}
}
