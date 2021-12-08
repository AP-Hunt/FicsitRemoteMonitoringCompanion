export interface Location {
    x: number
    y: number
    z: number
    rotation: number
}

export interface Production {
    Name: string
    CurrentProd: number
    MaxProd: number
    ProdPercent: number
}

export interface Consumption {
    Name: string
    CurrentConsumed: number
    MaxConsumed: number
    ConsPercent: number
}

export interface FactoryBuilding {
    building: string
    location: Location
    Recipe: string
    production: Production[]
    ingredients: Consumption[]
    ManuSpeed: number
    IsConfigured: boolean
    IsProducing: boolean
    IsPaused: boolean
    CircuitID: number
}
