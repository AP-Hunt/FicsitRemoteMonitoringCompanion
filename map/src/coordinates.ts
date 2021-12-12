export function gameToWorldCoords(coords : L.LatLng) : L.LatLng{
    return new L.LatLng(-coords.lat, coords.lng, coords.alt)
}