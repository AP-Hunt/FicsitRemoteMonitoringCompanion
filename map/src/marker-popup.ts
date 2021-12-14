import { LatLngExpression, Map } from "leaflet";
import { BuildingFeature } from "./feature-types";
import MarkerTemplate from "./marker-template";

class MarkerPopupViewModel {
    private _feature: BuildingFeature;

    public buildingType: KnockoutObservable<string>;
    public recipe: KnockoutObservable<string>;

    constructor(feature: BuildingFeature) {
        this._feature = feature;

        this.buildingType = ko.observable(feature.properties.building);
        this.recipe = ko.observable(feature.properties.Recipe);
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
        ko.applyBindings(this._vm, this._shadowRoot.querySelector("[data-root=true]"));
    }

    loadContent(){
        console.log("I'm loading content");
    }
}

export class MarkerPopup extends L.Popup {
    private _element: MarkerPopupElement;

    constructor(feature: BuildingFeature, options?: L.PopupOptions, source?: L.Layer){
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