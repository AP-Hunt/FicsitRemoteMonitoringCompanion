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
            <dd>
                <span data-bind="text: $data.Name"></span> (<span data-bind="text: $data.CurrentProd"></span>/min, <span data-bind="text: $data.ProdPercent"></span>% efficiency)
            </dd>
            <!-- /ko -->
        </dl>

        <canvas id="chart" width="200" height="150">
            <p>Loading...</p>
        </canvas>
    </div>
</template>`;

let domParser = new DOMParser();
let TemplateDomNode = domParser.parseFromString(template, "text/html");

export default TemplateDomNode.getElementById('marker-popup-template')! as HTMLTemplateElement;