angular.module('myApp.filters', [])

    .filter('prettyJSON', function () {
        function prettyPrintJson(json) {
            return JSON ? JSON.stringify(json, null, '  ') : 'Your browser doesnt support JSON so cant pretty print';
        }
        return prettyPrintJson;
    });
