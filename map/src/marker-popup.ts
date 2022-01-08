import { BuildingFeature } from "./feature-types";
import MarkerTemplate from "./marker-template";; 
import { LatLngExpression, Map } from "leaflet";
import { ChartDataSets } from "chart.js";

class MarkerPopupViewModel {
    private _feature: BuildingFeature;
    private _chart!: Chart;
    private _prom: any; // The prometheus-query library's typescript definitions appear not to be usable (or I'm doing something wrong)
    private _isInited: boolean;
    private _chartUpdaterInterval!: number|null;


    public buildingType: KnockoutObservable<string>;
    public recipe: KnockoutObservable<string>;
    public recipeOutputs: KnockoutReadonlyComputed<string[]>;
    public isOpen: boolean;

    constructor(feature: BuildingFeature) {
        this._feature = feature;
        this._isInited = false;

        this.buildingType = ko.observable(feature.properties.building);
        this.recipe = ko.observable(feature.properties.Recipe);
        this.recipeOutputs = ko.computed(this._formatRecipeOutputs.bind(this));
        this.isOpen = false;

        // @ts-ignore TS2304
        this._prom = new Prometheus.PrometheusDriver({
            endpoint: "http://localhost:9090",
            baseURL: "/api/v1"
        });
    }

    init(root: ShadowRoot){
        if(this._isInited) {
            return;
        }

        ko.applyBindings(this, root.querySelector("[data-root=true]"));

        const canvas = (root.getElementById("chart") as HTMLCanvasElement).getContext("2d");
        this._chart = new Chart(canvas as CanvasRenderingContext2D, {
            type: "line",
            options: {
                plugins: {
                    title: {
                        text: this.recipeOutputs()[0]
                    }
                }
            }
        });

        this._isInited = true;
    }

    onShow(shadowRoot: ShadowRoot) {
        this.init(shadowRoot);
        this.isOpen = true;

        this._updateChart();
        this._chartUpdaterInterval = setInterval(() => {
            this._updateChart();
        }, 10*1000);
    }


    onHide(shadowRoot: ShadowRoot) {
        this.isOpen = false;

        if(this._chartUpdaterInterval != null) {
            clearInterval(this._chartUpdaterInterval);
            this._chartUpdaterInterval = null;
        }
    }

    private _updateChart(){
        let labels = [
            `machine_name="${this._feature.properties.building}"`, 
            `item_name="${this._feature.properties.production[0].Name}"`,
            `x="${this._feature.geometry.coordinates[0]}"`,
            `y="${this._feature.geometry.coordinates[1]}"`,
            `z="${this._feature.geometry.coordinates[2]}"`,
        ]

        this._prom.rangeQuery(
            `machine_items_produced_per_min{${labels.join(",")}}`,
            new Date().valueOf() - (5*60*1000),
            new Date(),
            60,
        ).then((res : any) => {
            const series = res.result;
            
            let timestamps = series[0].values.map((v:{time: Date, value: any}) => `${v.time.getHours()}:${v.time.getMinutes()}`);
            let datasets : ChartDataSets[] = [];
            
            series.forEach((s : any) => {
                let itemBeingProduced = s.metric.labels["item_name"];
                let values = s.values.map((v:{time: Date, value: any}) => v.value);
                
                datasets.push({
                    label: itemBeingProduced,
                    data: values,
                    fill: false,
                    borderColor: "#000000"
                })
            })

            this._chart.data = {
                labels: timestamps,
                datasets: datasets,
            }
            this._chart.update();
        });
    }

    private _formatRecipeOutputs(): string[] {
        return this._feature.properties.production.map(p => {
            return `${p.Name} (${p.CurrentProd}/min, ${p.ProdPercent}% efficiency)`;
        })
    }
}

export class MarkerPopupElement extends HTMLElement {
    private _vm: MarkerPopupViewModel;
    private _shadowRoot: ShadowRoot

    constructor(feature: BuildingFeature) {
        super();

        let template = MarkerTemplate.content;

        this._vm = new MarkerPopupViewModel(feature);

        this._shadowRoot = this.attachShadow({mode: 'open'});      
        this._shadowRoot.appendChild(template.cloneNode(true));
        
    }

    onShow(){
        this._vm.init(this._shadowRoot);
        this._vm.onShow(this._shadowRoot);
    }

    onHide() {
        this._vm.onHide(this._shadowRoot);
    }

    loadContent(){
        console.log("I'm loading content");
    }
}

export class MarkerPopup extends L.Popup {
    private _element: MarkerPopupElement;

    constructor(feature: BuildingFeature, options?: L.PopupOptions, source?: L.Layer){ 
        options = options || {};
        (options as any)["minWidth"] = "fit-content";
        super(options, source);

        this._element = new MarkerPopupElement(feature);

        this.setContent(this._element);

        this.on('popupclose', (p) => {
            console.log("popup closed");
        });

        this.on('popupopen', (p) => {
            console.log("popup opened");
        });

        this.on('remove', (p) => {
            console.log("popup removed");
        })
    }

    public override onAdd(map: Map): this {
        super.onAdd(map);
        this._element.onShow();
        return this;
    }

    public override onRemove(map: Map): this {
        super.onRemove(map);
        this._element.onHide();
        return this;
    }
}

customElements.define('x-marker-popup', MarkerPopupElement);

function timestamps(interval: number, count: number): string[] {
    let now = new Date();
    let current = now;
    let i = count;
    let out: string[] = [];
    while(i > 0) {
        out.push(`${current.getHours()}:${current.getMinutes()}`);
        current = new Date(current.valueOf() + interval);
        i--;
    }

    return out;
}
