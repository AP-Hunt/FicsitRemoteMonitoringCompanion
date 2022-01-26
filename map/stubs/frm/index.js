var http = require("http");
var factoryJson = require("./factory.json");

http.createServer((req, res) => {
    switch(req.url) {
        case "/getFactory":
            getFactoryHandler(req, res);
            break;

        default:
            defaultRoute(req, res);
            break;
    }
}).listen(8080, () => {
    console.log("FRM stub is listening on 8080");
});

function getFactoryHandler(req, res) {
    res.writeHead(200, { "Content-Type": "application/json", "Access-Control-Allow-Origin": "*"});
    res.write(JSON.stringify(factoryJson));
    res.end();
}

function defaultRoute(req, res) {
    res.writeHead(404, { "Access-Control-Allow-Origin": "*"});
    res.end();
}