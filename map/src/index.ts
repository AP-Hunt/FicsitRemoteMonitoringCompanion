import { GameMap } from './map';

/// <reference types="prometheus-query/dist/types" />

function init()
{
    let map = new GameMap(document.getElementById("map")!);
    map.plotBuildings("http://localhost:8090/getFactory");

    // @ts-ignore TS2304
    const prom = new Prometheus.PrometheusDriver({
        endpoint: "http://localhost:9090",
        baseURL: "/api/v1"
    });

    prom.rangeQuery(
        'machine_items_produced_per_min{item_name="Iron Ingot"} > 20',
        new Date("2022-01-06T19:40:00Z"),
        new Date("2022-01-06T19:44:00Z"),
        60,
    ).then((res : any) => {
        const series = res.result;
        series.forEach((serie:any) => {
            console.log("Serie:", serie.metric.toString());
            console.log("Values:\n" + serie.values.join('\n'));
        });
    });
}
init();