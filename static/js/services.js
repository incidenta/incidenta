angular.module('myApp.services', [])

    .service('API', function ($http) {
        this.getAlertLogs = function(alert_id) {
            return $http.get(
                '/v1/alert/' + alert_id + '/logs'
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
