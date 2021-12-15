const template = `
<template id="marker-popup-template">
    <style>
        dd {
            white-space: nowrap;
        }

        *[data-root=true] {
            width: fit-content;
        }
    </style>
    <div data-root="true">
        <h2 data-bind="text: buildingType"></h2>

        <dl>
            <dt>Current recipe</dt>
            <dd data-bind="text: recipe"></dd>

            <dt>Producing</dt>
            <!-- ko foreach: recipeOutputs -->
            <dd data-bind="text: $data"></dd>
            <!-- /ko -->
        </dl>

        <svg width="250" height="150">
            <path d="M50 0 L50 100 L250 100" stroke="red" fill="none" />

            <path d="
                M50 50 
                L100 25 
                L150 15 
                L200 100 
                L250 50" stroke="#000000" fill="none" />

            <text x="25" y="25" style="font-size: 8pt;">75%</text> 
            <text x="25" y="50" style="font-size: 8pt;">50%</text> 
            <text x="25" y="75" style="font-size: 8pt;">25%</text> 

            <path d="M50 25 L250 25 Z" stroke="#555555" fill="none" opacity="0.3" />
            <path d="M50 50 L250 50 Z" stroke="#555555" fill="none" opacity="0.3" />
            <path d="M50 75 L250 75 Z" stroke="#555555" fill="none" opacity="0.3" />

            <text x="50" y="125" width="200" style="text-align: center;">Production over last 4m</text>
        </svg>
    </div>
</template>`;

let domParser = new DOMParser();
let TemplateDomNode = domParser.parseFromString(template, "text/html");

export default TemplateDomNode.getElementById('marker-popup-template')! as HTMLTemplateElement;