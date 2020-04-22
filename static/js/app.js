angular.module('myApp', [
    'myApp.controllers',
    'myApp.services',
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

            .otherwise('/');

        $locationProvider.html5Mode(false);
        $locationProvider.hashPrefix('!');
    });
