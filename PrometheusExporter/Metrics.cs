using System;
using System.Collections.Generic;
using System.Text;

namespace PrometheusExporter
{
    internal class Metrics
    {
        internal static readonly Prometheus.Gauge PowerCapacity = Prometheus.Metrics.CreateGauge(
            "power_capacity_mw", "Total capacity of all generators in a circuit",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = new string[] {
                        "circuit_id"
                },
            });

        internal static readonly Prometheus.Gauge PowerConsumed = Prometheus.Metrics.CreateGauge(
            "power_consumed_mw", "Amount of power being consumed in a circuit",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = new string[] {
                        "circuit_id"
                },
            });

        internal static readonly Prometheus.Gauge PowerMaxConsumed = Prometheus.Metrics.CreateGauge(
            "power_max_consumed_mw", "Highest amount of power consumed in a circuit",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = new string[] {
                        "circuit_id"
                },
            });

        internal static readonly Prometheus.Gauge PowerProduced = Prometheus.Metrics.CreateGauge(
            "power_produced_mw", "Amount of power being produced in a circuit",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = new string[] {
                        "circuit_id"
                },
            });

        internal static readonly Prometheus.Gauge PowerBatteryDifferential = Prometheus.Metrics.CreateGauge(
             "power_battery_differential_pc", "Amount of battery capacity as a percentage of total capacity in a circuit",
             new Prometheus.GaugeConfiguration()
             {
                 LabelNames = new string[] {
                        "circuit_id"
                 },
             });

        internal static readonly Prometheus.Gauge PowerBatteryPercent = Prometheus.Metrics.CreateGauge(
             "power_battery_pc", "Amount of battery power available as a percentage of total battery capacity in a circuit",
             new Prometheus.GaugeConfiguration()
             {
                 LabelNames = new string[] {
                        "circuit_id"
                 },
             });

        internal static readonly Prometheus.Gauge PowerBatteryCapacity = Prometheus.Metrics.CreateGauge(
             "power_battery_capacity_mw", "Total battery capacity in a circuit",
             new Prometheus.GaugeConfiguration()
             {
                 LabelNames = new string[] {
                        "circuit_id"
                 },
             });

        internal static readonly Prometheus.Gauge PowerBatteryTimeEmpty = Prometheus.Metrics.CreateGauge(
             "power_battery_time_empty_seconds", "Amount of time batteries have been empty in a circuit",
             new Prometheus.GaugeConfiguration()
             {
                 LabelNames = new string[] {
                        "circuit_id"
                 },
             });

        internal static readonly Prometheus.Gauge PowerFuseBlown = Prometheus.Metrics.CreateGauge(
             "power_fuse_blown", "Whether the fuse has been blown in a circuit",
             new Prometheus.GaugeConfiguration()
             {
                 LabelNames = new string[] {
                        "circuit_id"
                 },
             });
    }
}
