const template = `
<template id="marker-popup-template">
    <div data-root="true">
        <p>
            I am a <span data-bind="text: buildingType"></span> making <span data-bind="text: recipe"></span>
        </p>
    </div>
</template>`;

let domParser = new DOMParser();
let TemplateDomNode = domParser.parseFromString(template, "text/html");

export default TemplateDomNode.getElementById('marker-popup-template')! as HTMLTemplateElement;