function HomeController($scope) {
    'use strict';
}

function ProjectsController($scope, Restangular, $http) {
    'use strict';

    $scope.$on('socket:connect', function (ev, data) {
        console.log("on connect");
        //console.dir(data);
        //console.log(data);
    });

    $scope.$on('socket:file', function (ev, data) {
        console.log("socket:file event");
        console.dir(data);
    });

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
        };
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
    };

    $scope.showProject = function (project) {
        $scope.projectName = project.name;
        $scope.projectStatus = project.status;
        Restangular.one("projects", $scope.projectName).customGET("tree").then(function (tree) {
            var flattened = $scope.flattenTree(tree);
            $scope.projectTree = flattened;
        });
    };

    $scope.action1 = function (item) {
        console.log("action1");
        console.dir(item);
    };

    $scope.action2 = function (item) {
        console.log("action2");
        console.dir(item);
    };

    $scope.flattenTree = function (tree) {
        var flatTree = [],
            treeModel = new TreeModel(),
            root = treeModel.parse(tree[0]);
        root.walk({strategy: 'pre'}, function (node) {
            flatTree.push(node.model);
        });
        return flatTree;
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