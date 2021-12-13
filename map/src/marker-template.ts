const template = `
<template id="marker-popup-template">
    <div>
        <p style="font-weight: bold"><slot name="some-text">Default </slot></p>
    </div>
</template>`;

let domParser = new DOMParser();
let TemplateDomNode = domParser.parseFromString(template, "text/html");

export default TemplateDomNode.getElementById('marker-popup-template')! as HTMLTemplateElement;