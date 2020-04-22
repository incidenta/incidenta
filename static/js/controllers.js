angular.module('myApp.controllers', [])

    .controller('ProjectAddController', function ($scope, API) {
        $scope.adding = false;
        $scope.opts = {};
        $scope.add = function() {
            $scope.adding = true;
            API.addProject($scope.opts)
                .then(function (){
                    $scope.adding = false;
                    window.location = "#!/";
                });
        };
    })

    .controller('ProjectViewController', function ($scope, $routeParams, API) {
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

        $scope.delete = function(project) {
            API.deleteProject(project.id)
                .then(function () {
                    window.location = "#!/";
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
