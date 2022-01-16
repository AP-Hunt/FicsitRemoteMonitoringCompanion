import { MarkerCluster, MarkerClusterGroup } from "leaflet";
import { BuildingFeature } from "./feature-types";

export class BuildingClusterGroup extends L.MarkerClusterGroup {
    constructor() {
        super({
            maxClusterRadius: 100,
            disableClusteringAtZoom: -6,
            iconCreateFunction: BuildingClusterGroup.createIconCluster,
        });
    }

    // The basis for this icon creation is taken from the Leaflet Marker Cluster source:
    // https://github.com/Leaflet/Leaflet.markercluster/blob/31360f226e1a40c03c71d68b016891beb5e63370/src/MarkerClusterGroup.js#L821
    private static createIconCluster(cluster: MarkerCluster): L.Icon | L.DivIcon {
        var childCount = cluster.getChildCount();


		var c = ' marker-cluster-';
		if (childCount < 10) {
			c += 'small';
		} else if (childCount < 100) {
			c += 'medium';
		} else {
			c += 'large';
		}

        let children = cluster.getAllChildMarkers();
        let style: string[] = [];
        if(allAreNotProducing(children)) {
            style.push('background-color: red');
            style.push('color: white');
        } else if(anyAreNotProducing(children)){
            style.push('background-color: orange');
        } else {
            style.push('background-color: green');
        }

		return new L.DivIcon({ html: '<div style="'+style.join("; ")+'"><span>' + childCount + '</span></div>', className: 'marker-cluster' + c, iconSize: new L.Point(40, 40) });
    }
}

function allAreNotProducing(children: L.Marker[]): boolean {
    return children.every(buildingIsNotProducing)
}

function anyAreNotProducing(children: L.Marker[]): boolean {
    return children.some(buildingIsNotProducing)
}

function buildingIsNotProducing(marker: L.Marker): boolean {
    let feature = marker.feature as BuildingFeature;
    if(!feature) {
        console.error("marker feature must be an instance of BuildingFeature");
        return false;
    }

    return !feature.properties.IsProducing;
}