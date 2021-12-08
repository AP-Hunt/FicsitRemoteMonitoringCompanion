import { FactoryBuilding } from "./api-types";

export default interface Marker<T = any> {
    X: number
    Y: number
    Z: number

    metadata: T
}

abstract class MarkerBase<T> implements Marker<T> {
    X: number;
    Y: number;
    Z: number;

    metadata!: T;

    constructor(x: number, y: number, z: number) {
        this.X = x;
        this.Y = y;
        this.Z = z;
    }

    coordinates() : L.LatLng {
        return new L.LatLng(this.Y, this.X);
    }
}

export class BuildingMarker extends MarkerBase<string> {
    constructor(building : FactoryBuilding) {
        super(building.location.x, building.location.y, building.location.z)
        this.metadata = building.building;
    }
}