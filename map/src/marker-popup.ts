import { BuildingFeature } from "./feature-types";
import MarkerTemplate from "./marker-template";; 
import { Map } from "leaflet";

class MarkerPopupViewModel {
    private _feature: BuildingFeature;
    private _chart!: Chart;

    public buildingType: KnockoutObservable<string>;
    public recipe: KnockoutObservable<string>;
    public recipeOutputs: KnockoutReadonlyComputed<string[]>;

    constructor(feature: BuildingFeature) {
        this._feature = feature;

        this.buildingType = ko.observable(feature.properties.building);
        this.recipe = ko.observable(feature.properties.Recipe);
        this.recipeOutputs = ko.computed(this._formatRecipeOutputs.bind(this));
    }

    init(root: ShadowRoot){
        ko.applyBindings(this, root.querySelector("[data-root=true]"));

        const canvas = (root.getElementById("chart") as HTMLCanvasElement).getContext("2d");
        this._chart = new Chart(canvas as CanvasRenderingContext2D, {
            type: "line",
            data: {
                labels: timestamps(-60 * 1000, 6),
                datasets: [
                    {
                        label: this.recipe(),
                        data: [100, 120, 80, 75, 75, 99],
                        fill: false,
                        borderColor: "#000000"
                    }
                ]
            },
            options: {
                plugins: {
                    title: {
                        text: this.recipeOutputs()[0]
                    }
                }
            }
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

    init(){
        this._vm.init(this._shadowRoot);
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
    }

    public override onAdd(map: Map): this {
        super.onAdd(map);
        this._element.init();
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
