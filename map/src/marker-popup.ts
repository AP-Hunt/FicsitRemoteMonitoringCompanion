import MarkerTemplate from "./marker-template";

export class MarkerPopupElement extends HTMLElement {
    constructor(slotText: string) {
        super();
        const template = MarkerTemplate.content;

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

customElements.define('x-marker-popup', MarkerPopupElement);