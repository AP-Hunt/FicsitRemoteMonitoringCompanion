import { Realtime, Slider } from "leaflet";
import { gameToWorldCoords } from "./coordinates";
import { ElevationLayerGroups } from "./elevation-layer-group";
import { BuildingFeature } from "./feature-types";
import { AssemblerIcon, BlenderIcon, ConstructorIcon, FoundryIcon, ManufacturerIcon, PackagerIcon, RefineryIcon, SmelterIcon } from "./icons";
import { MarkerPopup } from "./marker-popup";
import { requestAsGeJSON } from "./realtime-helpers";

const Z_MAX = metresToGameUnits(2000);
const Z_MIN = metresToGameUnits(-250);
const Z_STEP = metresToGameUnits(15);
const Z_DEFAULT = metresToGameUnits(-25);

function metresToGameUnits(x: number) : number {
    return x * 100;
}

function gameUnitsToMetres(x: number): number {
    return x * 0.01;
}

export class GameMap {
    private _domTarget : HTMLElement
    private _map! : L.Map
    private _realtime!: Realtime;
    private _elevationGroups!: ElevationLayerGroups;
    
    private _slider! : Slider;
    private _elevation: number = Z_DEFAULT;

    private readonly _bounds : L.LatLngBoundsLiteral = [
        [-375e3, -324698.832031],
        [375e3, 425301.832031],
    ];
    private readonly _minZoom = -10;
    private readonly _maxZoom = -5;
    private readonly _defaultZoom = this._minZoom;

    constructor(target : HTMLElement){
        this._domTarget = target;
        this._initialize();
    }

    private _initialize(){
        this._map = new L.Map(this._domTarget, {
            crs: L.CRS.Simple,
        });

        this._map.setMinZoom(this._minZoom);
        this._map.setMaxZoom(this._maxZoom);
        this._map.fitBounds(this._bounds);
        this._map.setView(this._map.getCenter(), this._defaultZoom);

        let imgOverlayLayer = new L.ImageOverlay("map-16k.png", this._bounds);
        imgOverlayLayer.addTo(this._map);

        
        this._elevationGroups = new ElevationLayerGroups(Z_MIN, Z_MAX, Z_STEP);
        this._elevationGroups.addTo(this._map);

        this._slider = L.control.slider(
            (value : string) => {
                this._updateElevation(parseInt(value));
            },
            {
                min: gameUnitsToMetres(Z_MIN),
                max: gameUnitsToMetres(Z_MAX),
                value: gameUnitsToMetres(Z_DEFAULT),
                size: "500px",
                orientation: "vertical",
                collapsed: false,
                title: "Elevation",
                increment: true,
                step: gameUnitsToMetres(Z_STEP),
                getValue(value : string): string {
                    return `From ${value}m <br /> To: ${parseInt(value) + gameUnitsToMetres(Z_STEP)}m`;
                }
            }
        );
        this._slider.addTo(this._map);
    }

    plotBuildings(url : string) {
        const self = this;
        this._realtime = new L.Realtime<L.LatLng>(
            requestAsGeJSON(url),
            {
                interval: 10 * 1000,
                container: self._elevationGroups,
                getFeatureId(feature : GeoJSON.Feature) {
                    return (feature.geometry as GeoJSON.Point).coordinates.join("/");
                },

                updateFeature(feature: GeoJSON.Feature, marker: L.Marker) {
                    
                    // If the given (old) layer is null, return null
                    // so that leaflet-realtime will make an appropriate layer
                    // for us, which we can customie
                    // https://github.com/perliedman/leaflet-realtime/blob/88d364da9dde8aa0c8c01c5b46bc0673832c8965/src/L.Realtime.js#L202
                    if(!marker){return}

                    if(marker.getPopup() instanceof MarkerPopup){
                        (marker.getPopup() as MarkerPopup).updateFeature(feature);
                    }

                    return marker;
                },

                onEachFeature(feature: BuildingFeature, marker: L.Marker) {
                    let popup = new MarkerPopup(feature);
                    marker.bindPopup(popup);

                    var icon = new L.Icon.Default();
                    switch(feature.properties.building) {
                        case "Assembler":
                            icon = new AssemblerIcon();
                            break;

                        case "Blender":
                            icon = new BlenderIcon();
                            break;

                        case "Constructor":
                            icon = new ConstructorIcon();
                            break;

                        case "Foundry":
                            icon = new FoundryIcon();
                            break;

                        case "Manufacturer":
                            icon = new ManufacturerIcon();
                            break;

                        case "Packager":
                            icon = new PackagerIcon();
                            break;

                        case "Refinery":
                            icon = new RefineryIcon();
                            break;

                        case "Smelter":
                            icon = new SmelterIcon();
                            break;

                    }
                    marker.setIcon(icon);

                    let geom = feature.geometry as GeoJSON.Point
                    
                    marker.setLatLng(
                        gameToWorldCoords(new L.LatLng(
                            geom.coordinates[1], 
                            geom.coordinates[0], 
                            geom.coordinates[2]
                        ))
                    );
                }
            }
        );

        this._realtime.addTo(this._map);
        this._realtime.on("update", (evt: L.LeafletEvent) => {
            this._elevationGroups.refresh();
        })

        this._realtime.start();
    }

    private _updateElevation(elevation : number) {
        this._elevation = metresToGameUnits(elevation);
        this._elevationGroups.showElevation(this._elevation);
    }

}