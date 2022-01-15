export function requestAsGeJSON(url: string) {
    return function(success: (featuers: any) => void, error: (error: object, message: string) => void) {
        return fetch(url)
            .then(response => response.json())
            .then(data => {
                let geo: GeoJSON.FeatureCollection = {
                    type: "FeatureCollection",
                    features: [] as Array<GeoJSON.Feature>
                };

                data.forEach((building: any) => {
                    let feature = {
                        type: "Feature",
                        geometry: {
                            type: "Point",
                            coordinates: [
                                building.location.x,
                                building.location.y,
                                building.location.z
                            ]
                        }
                    } as GeoJSON.Feature;

                    delete building.location;
                    feature.properties = building;

                    geo.features.push(feature)
                })

                return geo;
            })
            .then(success)
            .catch((reason) => {
                error({}, reason);
            });
        }
}
