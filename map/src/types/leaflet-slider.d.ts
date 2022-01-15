import 'leaflet';

declare module 'leaflet' {

    type OnUpdateHandler = (value : string) => void;
    type Orientation = "horizontal" | "vertical";

    interface SliderOptions {
        size: string;
        position?: string;
        min: number;
        max: number;
        step: number;
        id?: string;
        value: number;
        collapsed: boolean;
        title: string;
        logo?: string;
        orientation: Orientation;
        increment: boolean;
        showValue?: boolean;
        syncSlider?: boolean;

        getValue: (value: string) => string
    }

    class Slider extends L.Control {
        constructor(callback: OnUpdateHandler, options: SliderOptions)
    }

    declare namespace control {
        function slider(f: OnUpdateHandler, options: SliderOptions): Slider;
    }
}
