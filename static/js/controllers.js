angular.module('myApp.controllers', [])

    .controller('ProjectController', function ($scope, $routeParams, API) {
        $scope.projectLoaded = false;
        $scope.project = {};

        $scope.alertsLoaded = false;
        $scope.alerts = [];

        $scope.project_id = $routeParams.project_id;
        $scope.events = [];
        $scope.alert = {};

        $scope.start = function() {
            API.getProject($scope.project_id)
                .then(function (result) {
                    $scope.projectLoaded = true;
                    $scope.project = result.data;
                });
            API.getProjectAlerts($scope.project_id)
                .then(function (result) {
                    $scope.alertsLoaded = true;
                    $scope.alerts = result.data;
                    if ($scope.alerts.length > 0) {
                        $scope.showEvents($scope.alerts[0]);
                    }
                });
        }

        $scope.showEvents = function(alert) {
            $scope.events = [];
            $scope.alert = alert;
            API.getAlertEvents(alert.id)
                .then(function (result) {
                    $scope.events = result.data;
                });
        }


        $scope.start();
    })

    .controller('HomeController', function ($scope, API) {
        $scope.loaded = false;
        $scope.projects = [];

        $scope.refresh = function() {
            API.getProjects()
                .then(function (result) {
                    $scope.loaded = true;
                    $scope.projects = result.data;
                });
        }

        $scope.refresh();
    });
