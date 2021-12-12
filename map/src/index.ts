
import { GameMap } from './map';

class XPopup extends HTMLElement {
    constructor(slotText: string) {
        super();
        const template = (document
            .getElementById('test-tpl')! as HTMLTemplateElement)
            .content;

        let textSlot = template.querySelector('slot[name="some-text"]')!;
        textSlot.textContent = slotText;
        const shadowRoot = 
            this
            .attachShadow({mode: 'open'})
            .appendChild(template.cloneNode(true));
    }

    loadContent(){
        console.log("I'm loading content");
    }
}

customElements.define('x-popup', XPopup);

function init()
{
    let map = new GameMap(document.getElementById("map")!);
    map.plotBuildings("http://localhost:8080/factory-geojson.json");
}
init();