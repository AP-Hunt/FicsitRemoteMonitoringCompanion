import { BuildingFeature } from "./feature-types";
import { Layer, MarkerClusterGroup } from "leaflet";
import { BuildingClusterGroup } from "./building-cluster-group";

type ElevationGroup = {min: number, max: number, layer: MarkerClusterGroup}

export class ElevationLayerGroups extends L.FeatureGroup {
    private _groups: ElevationGroup[] = [];

    constructor(min: number, max: number, step: number) {
        super();
        let z = min;
        while(z < max) {
            this._groups.push({
                min: z,
                max: z + step,
                layer: new BuildingClusterGroup()
            });

            z += step;
        }
    }

    public override addLayer(layer: L.Layer): this {
        if(!(layer instanceof L.Marker)){
            console.error("can only add instances of L.Marker");
            return this;
        }

        let marker = layer as L.Marker;
        let feature = marker.feature as BuildingFeature;
        if(!feature) {
            console.error("marker feature must be an instance of BuildingFeature");
            return this;
        }

        let targetGroup: ElevationGroup = this._groups[0];
        let z = feature.geometry.coordinates[2];
        for (const group of this._groups) {
            if(z > group.min && z <= group.max){
                targetGroup = group;
                break;
            }
        }

        targetGroup.layer.addLayer(layer);

        return this;
    }

    public override removeLayer(layer: number | Layer): this {
        if(!(layer instanceof L.Marker)){
            console.error("can only add instances of L.Marker");
            return this;
        }

        let marker = layer as L.Marker;
        let feature = marker.feature as BuildingFeature;
        if(!feature) {
            console.error("marker feature must be an instance of BuildingFeature");
            return this;
        }

        let targetGroup: ElevationGroup = this._groups[0];
        let z = feature.geometry.coordinates[2];
        for (const group of this._groups) {
            if(z > group.min && z <= group.max){
                targetGroup = group;
                break;
            }
        }

        targetGroup.layer.removeLayer(layer);

        return this;
    }

    public showElevation(elevation: number) {
        let targetGroup: ElevationGroup | null = null;
        let z = elevation
        for (const group of this._groups) {
            if(z > group.min && z <= group.max){
                targetGroup = group;
                break;
            }
        }

        if(targetGroup == null) {
            console.error("no layer covers the elevation " + elevation);
            return;
        }

        this._groups.forEach((group : ElevationGroup) => {
            if(group == targetGroup) {
                group.layer.addTo(this._map);
            } else {
                group.layer.removeFrom(this._map);
            }
        });
    }

    public updateLayer(layer: L.Marker) {
        this._groups.forEach((group: ElevationGroup) => {
            if(group.layer.hasLayer(layer)) {
                group.layer.refreshClusters(layer);
            }
        })
    }

    public refresh() {
        this._groups.forEach((group: ElevationGroup) => {
            if(this._map.hasLayer(group.layer)){
                group.layer.refreshClusters();
            }
        });
    }
}