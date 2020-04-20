angular.module('myApp.services', [])

    .service('API', function ($http) {
        this.getReceivers = function() {
            return $http.get(
                '/v1/receivers'
            );
        };
    });
