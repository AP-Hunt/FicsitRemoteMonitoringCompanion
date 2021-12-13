import 'leaflet';

declare module 'leaflet' {
    type RealtimeSource = string | FetchOptions | function

    interface RealtimeOptions extends L.GeoJSONOptions {
        start: boolean;
        interval: number;
        getFeatureId: (featureData: L.GeoJSON) => TFeatureId;
        updateFeature: (featureData: L.GeoJSON, Layer: L.Layer) => void;
        container: L.LayerGroup;
        removeMissing: boolean;
    }

    class Realtime<TFeatureId = any> extends L.Layer {
        constructor(source: RealtimeSource, options: RealtimeOptions<TFeatureId>)

        start(): this;
        stop(): this;
        isRunning(): boolean;
        update(featureData?: L.GeoJSON): this;
        remove(featureData: L.GeoJSON): this;
        getLayer(featureId: TFeatureId): L.Layer;
        getFeature(featureId: TFeatureId): L.Layer;
    }
}