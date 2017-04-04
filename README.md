# PromStackCTL

PromStackCTL interacts with endpoints from the [PromStack Implementation](https://github.com/jbkc85/promstack) to provide a simple Command Line Utility Management Script.

> This script at this time does not DEPLOY an exporter onto a machine. It simply registers said machine in Consul for Prometheus to monitor via Service Discovery.

**WARNING**: promstackctl is under development and needs review.  Certain variable names will change overtime (ie: changing node to machine for SRE terminology) and so on.

## Requirements

This script expects a Consul and Prometheus endpoint to be available.

## Examples of Usage

### Adding an Exporter

```sh
curl -XPUT consul.endpoint:8500/v1/kv/promstack/exporters/node-exporter -d '{"port":9100,"tags":["exporter"]}'
```

### Adding a Machine to Monitor

```sh
$ promstackctl monitor server --node.name example.com --node.address 1.1.1.1 --exporter.name cadvisor
```

## TODO

[ ] - implement health for entire PromStack (missing Grafana and AlertManager) using parallel requests
[ ] - add external JSON file for exporter list based on Prometheus's Wiki page
[ ] - migrate 'node' variables to 'machine' where applicable (node is Consul terminology)
