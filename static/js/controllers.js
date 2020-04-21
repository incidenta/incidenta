angular.module('myApp.controllers', [])

    .controller('HomeController', function ($scope, API) {
        $scope.loaded = false;
        $scope.projects = [];

        $scope.start = function() {
            API.getProjects()
                .then(function (result) {
                    $scope.loaded = true;
                    $scope.projects = result.data;
                });
        }

        $scope.start();
    });
