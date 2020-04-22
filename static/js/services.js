angular.module('myApp.services', [])

    .service('API', function ($http) {
        this.getAlertEvents = function(alert_id) {
            return $http.get(
                '/v1/alert/' + alert_id + '/events'
            );
        }
        this.getProjectAlerts = function(project_id) {
            return $http.get(
                '/v1/project/' + project_id + '/alerts'
            );
        }
        this.getProject = function(project_id) {
            return $http.get(
                '/v1/project/' + project_id
            );
        }
        this.getProjects = function() {
            return $http.get(
                '/v1/projects'
            );
        };
    });
