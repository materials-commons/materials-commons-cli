function HomeController($scope) {
    'use strict';
}

function ProjectsController($scope, Restangular, $http) {
    'use strict';

    $scope.projectsData = [];
    var allProjects = Restangular.all('projects');
    allProjects.getList().then(function (projects) {
        $scope.projects = projects;
    });

    $scope.selected = [];

    $scope.statusButtonName = function (status) {
        if (status === "Unloaded") {
            return "Upload";
        }
        return status;
    };

    $scope.statusButtonAction = function (val) {
        console.log("uploading...");
        console.dir(val);
    };

    $scope.uploadProject = function (what) {
        Restangular.one("projects", what.name).customGET("upload").then(function () {
            allProjects.getList().then(function (projects) {
                $scope.projects = projects;
            });
        }, function () {
            console.log("Upload failed");
        });
    };

    $scope.newProject = function () {
        console.log("Creating project: " + $scope.newProjectName);
        console.log("  Located at: " + $scope.newProjectPath);
        var proj = {
            name: $scope.newProjectName,
            path: $scope.newProjectPath,
            status: "Unloaded"
        }
        allProjects.post(proj).then(function () {
            console.log("Project created");
            allProjects.getList().then(function (projects) {
                $scope.projects = projects;
            });
        }, function () {
            console.log("Project creation failed!");
        });
        $scope.newProjectName = "";
        $scope.newProjectPath = "";
    }

    $scope.projectGridOptions = {
        data: 'projects',
        multiSelect: false,
        columnDefs: [
            {field: 'name', displayName: 'Name'},
            {field: 'path', displayName: 'Path'},
            {field: 'status', displayName: 'Status', cellTemplate: 'partials/projects/status_cell.html'}
        ],
        selectedItems: $scope.selected,
        afterSelectionChange: function (rowItem) {
            $scope.project = rowItem.entity.name;
            $scope.projectStatus = rowItem.entity.status;
            Restangular.one("projects", $scope.project).customGET("tree").then(function (tree) {
                $scope.projectTree = tree;
            });
        }
    };


}

function ChangesController($scope, Restangular, $timeout) {
    $scope.events = [];
//    (function tick() {
//        console.log("tick")
//        Restangular.all('projects/changes').getList().then(function(eventsInfo) {
//            //console.dir(eventsInfo);
//            var found = false;
//            $scope.events.forEach(function(event) {
//                if (event.filepath == eventsInfo.filepath) {
//                    found = true;
//                }
//            });
//
//            if (! found) {
//                $scope.events.push(eventsInfo);
//            }
//        })
//        $timeout(tick, 3000);
//    })();
}

function ProvenanceController($scope) {

}

function AboutController($scope) {

}

function ContactController($scope) {

}

function EventController($scope) {

}