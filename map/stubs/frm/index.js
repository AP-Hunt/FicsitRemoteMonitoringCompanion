var http = require("http");
var factoryJson = require("./factory.json");

const HEADERS = { "Content-Type": "application/json", "Access-Control-Allow-Origin": "*"}

http.createServer((req, res) => {
    switch(req.url) {
        case "/getFactory":
            getFactoryHandler(req, res);
            break;

        case "/getTrains":
            getTrainsHandler(req, res);
            break;

        default:
            defaultRoute(req, res);
            break;
    }
}).listen(8080, () => {
    console.log("FRM stub is listening on 8080");
});

function getFactoryHandler(req, res) {
    res.writeHead(200, HEADERS);
    res.write(JSON.stringify(factoryJson));
    res.end();
}

function getTrainsHandler(req, res) {
    let timestamp = new Date().valueOf() // UNIX timestamp in milliseconds

    let z = 0.0;
    let y = 0.0;
    let x = (timestamp % 30000); // Base the X value on the timestamp so that the train will appear to move over time
    let rotation = timestamp % 360; // Base rotation on timestamp to make the train appear to rotate over time (to test rotation)

    let trainDef = {
        TrainName: "choo_choo",
        location: {
            x: x,
            y: y,
            z: z,
            Rotation: rotation
        },
        ForwardSpeed: 55.1,
        TotalMass: 99.0,
        PayloadMass: 50.0,
        PowerConsumed: 5.0,
        TrainStation: "Home",
        ThrottlePercent: 100.0,
        Derailed: false,
        PendingDerail: false,
        Status: "existing"
    };

    res.writeHead(200, HEADERS);
    res.write(JSON.stringify([trainDef]));
    res.end();
}

function defaultRoute(req, res) {
    res.writeHead(404, { "Access-Control-Allow-Origin": "*"});
    res.end();
}