var http = require("http");
var factoryJson = require("../frm/factory.json");
var prom = require("prom-client");

http.createServer(async (req, res) => {
    switch(req.url) {
        case "/metrics":
            await getMetricsHandler(req, res);
            break;

        default:
            defaultRoute(req, res);
            break;
    }
}).listen(9091, () => {
    console.log("Prometheus exporter stub is listening on 9091");
});

const machineItemsProducedPerMinute = new prom.Gauge({
    name: "machine_items_produced_per_min", 
    help: "How much of an item a building is producing",
    labelNames: ["item_name", "machine_name", "x", "y", "z"]
});
const machineItemsProducedEfficiency = new prom.Gauge({
    name: "machine_items_produced_pc", 
    help: "The efficiency with which a building is producing an item",
    labelNames: ["item_name", "machine_name", "x", "y", "z"]
});

function randIntBetween(min, max) {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min) + min); //The maximum is exclusive and the minimum is inclusive
  }
  

async function getMetricsHandler(req, res) {
    factoryJson.forEach(building => {
        building.production.forEach(output => {
            machineItemsProducedPerMinute.labels(
                output.Name,
                building.building,
                building.location.x,
                building.location.y,
                building.location.z,
            ).set(randIntBetween(0, 51));

            machineItemsProducedEfficiency.labels(
                output.Name,
                building.building,
                building.location.x,
                building.location.y,
                building.location.z,
            ).set(randIntBetween(0, 101)/100);
        })
    });


    res.writeHead(200, { "Content-Type": "text/plain", "Access-Control-Allow-Origin": "*"});
    const metricOutput = await prom.register.metrics();
    res.write(metricOutput);
    res.end();
}

function defaultRoute(req, res) {
    res.writeHead(404, { "Access-Control-Allow-Origin": "*"});
    res.end();
}