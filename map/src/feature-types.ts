
export interface RecipeOutput {
    Name: string;
    CurrentProd: number;
    MaxProd: number;
    ProdPercent: number;
}

export interface BuildingProperties {
    building: string;
    Recipe: string;
    production: RecipeOutput[];
    IsProducing: boolean;
}

export interface TrainProperties {
    TrainName: string;
    location: {
        x: number,
        y: number, 
        z: number
        Rotation: number
    }
}

export interface BuildingFeature extends GeoJSON.Feature<GeoJSON.Point> {
    properties: BuildingProperties
}

export interface TrainFeature extends GeoJSON.Feature<GeoJSON.Point> {
    properties: TrainProperties
}