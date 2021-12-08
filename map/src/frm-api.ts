import { FactoryBuilding } from "./api-types";
export async function fetchFactory() : Promise<FactoryBuilding[]> {

    let response = await fetch("http://localhost:8080/factory.json");
    return response.json();
}