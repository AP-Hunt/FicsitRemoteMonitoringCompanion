import { GameMap } from './map';

function init()
{
    let map = new GameMap(document.getElementById("map")!);
    map.plotBuildings("http://localhost:8080/factory-geojson.json");
}
init();