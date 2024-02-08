function mOver(id) {
    var target = document.getElementById(id)
    if (target.innerText.length >= 23){
    target.classList.add("s");
}

}
function mOut(id) {
    document.getElementById(id).classList.remove("s");
}

function key() {
    mapboxgl.accessToken = '';
    const mapboxClient = mapboxSdk({ accessToken: mapboxgl.accessToken });
    return mapboxClient
}

function map() {
    const map = new mapboxgl.Map({
        container: 'map',
        style: '',
        center: [1, 1],
        zoom: 0
    });
    map.addControl(new mapboxgl.NavigationControl());
    return map
}

function marker(location) {
    mapboxClient.geocoding
        .forwardGeocode({
            query: location,
            autocomplete: false,
            limit: 1
        })
        .send()
        .then((response) => {
            if (
                !response ||
                !response.body ||
                !response.body.features ||
                !response.body.features.length
            ) {
                console.error('Invalid response:');
                console.error(response);
                return;
            }
            const feature = response.body.features[0];
            new mapboxgl.Marker().setLngLat(feature.center).setPopup(new mapboxgl.Popup().setHTML(location)).addTo(map);
        });
}

