import { GameMap } from './map';

function init()
{
    let location = window.location
    let params = new URLSearchParams(location.search);
    let frmHost = "localhost"
    let frmPort = 8080;

    if(params.has("frmport")) {
        let p = parseInt(params.get("frmport")!, 10);
        if(p != undefined){
            frmPort = p
        }
    }
    if(params.has("frmhost")) {
        let h = params.get("frmhost");
        if(h != undefined){
            frmHost = h
        }
    }

    let map = new GameMap(document.getElementById("map")!);
    map.plotBuildings(`http://${frmHost}:${frmPort}/getFactory`);
    map.plotTrains(`http://${frmHost}:${frmPort}/getTrains`);
}
init();
