import { FactoryBuilding } from './api-types';
import { fetchFactory } from './frm-api'
import { BuildingMarker } from './marker';
import { MarkerSet } from './set';

class XPopup extends HTMLElement {
    constructor(slotText: string) {
        super();
        const template = (document
            .getElementById('test-tpl')! as HTMLTemplateElement)
            .content;

        let textSlot = template.querySelector('slot[name="some-text"]')!;
        textSlot.textContent = slotText;
        const shadowRoot = 
            this
            .attachShadow({mode: 'open'})
            .appendChild(template.cloneNode(true));
    }

    loadContent(){
        console.log("I'm loading content");
    }
}

customElements.define('x-popup', XPopup);

const markerGroup : L.FeatureGroup = new L.FeatureGroup();

function gameToWorldCoords(coords : L.LatLng) : L.LatLngExpression{
    return [-coords.lat, coords.lng]
}


function plotMarkers(map : L.Map, markers : MarkerSet<BuildingMarker>) {
    markerGroup.clearLayers();

    for(let m of markers) {
        markerGroup.addLayer(
            new L.Marker(
                gameToWorldCoords(m.coordinates())
            )
        );
    }

    if(!map.hasLayer(markerGroup)){
        markerGroup.addTo(map);
    }
}

function init()
{
    const bounds : L.LatLngBoundsLiteral = [
        [-375e3, -324698.832031],
        [375e3, 425301.832031],
    ];

    const map = new L.Map("map", {
        crs: L.CRS.Simple,
        zoom: -9,
        maxZoom: -5,
        minZoom: -10,
        maxBounds: bounds
    });

    map.setMinZoom(-10);
    map.setMaxZoom(-5);

    map.fitBounds(bounds);
    map.setView(map.getCenter(), -10);

    let imgOverlayLayer = new L.ImageOverlay("map-16k.png", bounds);
    imgOverlayLayer.addTo(map);

    let popup = new XPopup("I've changed");
    new L.Marker([-1000, 1000], { "title": "Random"}).bindPopup(popup).on("popupopen", popup.loadContent).addTo(map);

    setInterval(function(){
        fetchFactory()
        .then((factory : FactoryBuilding[]) => {
            let markers = 
            factory.map(b => new BuildingMarker(b));
            let set = MarkerSet.fromArray(markers);

            plotMarkers(map, set);
        });

    }, 10000)
}
init();