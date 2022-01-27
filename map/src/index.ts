import { GameMap } from './map';

function init()
{
    let location = window.location
    let params = new URLSearchParams(location.search);
    let frmPort = 8080;

    if(params.has("frmport")) {
        let p = parseInt(params.get("frmport")!, 10);
        if(p != undefined){
            frmPort = p
        }
    }

    let map = new GameMap(document.getElementById("map")!);
    map.plotBuildings(`http://localhost:${frmPort}/getFactory`);
}
init();