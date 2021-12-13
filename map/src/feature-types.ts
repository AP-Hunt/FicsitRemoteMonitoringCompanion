
export interface BuildingProperties {
    building: string;
    Recipe: string;
}

export interface BuildingFeature extends GeoJSON.Feature {
    properties: BuildingProperties
}