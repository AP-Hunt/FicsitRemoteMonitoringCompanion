import { gameToWorldCoords } from "./coordinates";
import { BuildingFeature } from "./feature-types";
import { AssemblerIcon, BlenderIcon, ConstructorIcon, FoundryIcon, ManufacturerIcon, PackagerIcon, RefineryIcon, SmelterIcon } from "./icons";
import { MarkerPopup } from "./marker-popup";

function requestAsGeJSON(url: string) {
    return function(success: (featuers: any) => void, error: (error: object, message: string) => void) {
        return fetch(url)
            .then(response => response.json())
            .then(data => {
                let geo: GeoJSON.FeatureCollection = {
                    type: "FeatureCollection",
                    features: [] as Array<GeoJSON.Feature>
                };

                data.forEach((building: any) => {
                    let feature = {
                        type: "Feature",
                        geometry: {
                            type: "Point",
                            coordinates: [
                                building.location.x,
                                building.location.y,
                                building.location.z
                            ]
                        }
                    } as GeoJSON.Feature;

                    delete building.location;
                    feature.properties = building;

                    geo.features.push(feature)
                })

                return geo;
            })
            .then(success)
            .catch((reason) => {
                error({}, reason);
            });
        }
}

export class GameMap {
    private _domTarget : HTMLElement
    private _map! : L.Map

    private _realtime!: any;

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
    }

    plotBuildings(url : string) {
        const self = this;
        this._realtime = new L.Realtime<L.LatLng>(
            requestAsGeJSON(url),
            {
                interval: 3000,
                getFeatureId(feature : GeoJSON.Feature) {
                    return (feature.geometry as GeoJSON.Point).coordinates;
                },

                updateFeature(feature: GeoJSON.Feature, marker: L.Marker) {
                    
                    // If the given (old) layer is null, return null
                    // so that leaflet-realtime will make an appropriate layer
                    // for us, which we can customie
                    // https://github.com/perliedman/leaflet-realtime/blob/88d364da9dde8aa0c8c01c5b46bc0673832c8965/src/L.Realtime.js#L202
                    if(!marker){return}

                    let addToMap = marker === undefined;
                    let m = marker || new L.Marker([0, 0]);

                    let geom = feature.geometry as GeoJSON.Point
                    
                    m.setLatLng(
                        gameToWorldCoords(new L.LatLng(
                            geom.coordinates[1], 
                            geom.coordinates[0], 
                            geom.coordinates[2]
                        ))
                    );

                    if(addToMap) {
                        m.addTo(self._map);
                    }

                    if(m.getPopup() instanceof MarkerPopup){
                        (m.getPopup() as MarkerPopup).updateFeature(feature);
                    }

                    return m;
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
                }
            }
        );

        this._realtime.addTo(this._map);
        this._realtime.start();
    }
}