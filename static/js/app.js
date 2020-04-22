angular.module('myApp', [
    'myApp.controllers',
    'myApp.services',
    'myApp.filters',
    'ngMaterial',
    'ngRoute',
    'ngMaterial',
])

    .config(function($httpProvider, $routeProvider, $locationProvider) {
        $httpProvider.useApplyAsync(true);

        $routeProvider
            .when('/', {
                templateUrl: 'partials/views/home.html',
                controller: 'HomeController',
            })

            .when('/project/:project_id', {
                templateUrl: 'partials/views/project.html',
                controller: 'ProjectController',
            })

            .otherwise('/');

        $locationProvider.html5Mode(false);
        $locationProvider.hashPrefix('!');
    });
