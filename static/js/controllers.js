angular.module('myApp.controllers', [])

    .controller('HomeController', function ($scope, API) {
        $scope.loaded = false;
        $scope.receivers = [];

        $scope.start = function() {
            API.getReceivers()
                .then(function (result) {
                    $scope.loaded = true;
                    $scope.receivers = result.data;
                });
        }

        $scope.start();
    });
