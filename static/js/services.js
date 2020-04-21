angular.module('myApp.services', [])

    .service('API', function ($http) {
        this.getProjects = function() {
            return $http.get(
                '/v1/projects'
            );
        };
    });
