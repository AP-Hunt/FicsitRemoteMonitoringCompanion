# Ficsit Remote Monitoring Companion

Ficsit Remote Monitoring Companion is a companion application for the 
[Ficsit Remote Montioring mod](https://ficsit.app/mod/B9bEiZFtaaQZHU) 
for [Satisfactory](https://www.satisfactorygame.com/).

The Ficsit Remote Monitoring mod exports metrics about the loaded save in 
Satisfactory, via a web server serving JSON. This companion application 
provides a visulisation tool on top of those metrics.

![](./blob/main/images/example-dashboard.png)

## Installation instructions
See [installation instructions](./blob/main/InstallationInstructions.md)

## How do I use it?
Once installed, sign in (button in bottom left), using the username and password `ficsit:pioneer`.

![](./blob/main/images/menu-login.png)

Ficsit Remote Monitoring Companion comes with some preconfigured dashboards, under the Dashboards > Manage menu.

![](./blob/main/images/menu-dashboards.png)

![](./blob/main/images/page-dashboards.png)

From the manage dashboards screen, you can create your own dashboards. How to do that is out of scope here, but 
[the Grafana documentation](https://grafana.com/docs/grafana/latest/) is a good place to start.

## How it works
Ficsit Remote Monitoring Companion runs, configures, and coordinates the following components:

* [Prometheus](https://prometheus.io/) for storing metric data
* [Grafana](https://grafana.com/grafana/) for creating visualisations of the metrics stored
  in prometheus 
* [PrometheusExporter](./tree/main/PrometheusExporter) for converting the information exposed 
  by Ficsit Remote Monitoring in to the text-based [Prometheus exposition format](https://prometheus.io/docs/instrumenting/exposition_formats/)
* [Companion](./tree/main/Companion) for coordinating all of the above, and providing easy access to Grafana

It takes care of all the heavy lifting, so you can get on with building factories.