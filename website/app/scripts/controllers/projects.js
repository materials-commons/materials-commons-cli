angular.module('materialsApp')
    .controller('ProjectsCtrl', function ($scope, materials) {
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

        $scope.getAllProjects = function () {
            materials('/projects')
                .success(function (projects) {
                    projects.forEach(function (project) {
                        project.originalName = project.name;
                    });
                    $scope.projects = projects;
                })
                .jsonp();
        };

        $scope.getAllProjects();

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
            materials('/projects/%/upload', what.name)
                .success(function () {
                    $scope.getAllProjects();
                })
                .get();
        };

        $scope.newProject = function () {
            console.log("Creating project: " + $scope.newProjectName);
            console.log("  Located at: " + $scope.newProjectPath);
            var proj = {
                name: $scope.newProjectName,
                path: $scope.newProjectPath,
                status: "Unloaded"
            };
            materials('/projects')
                .success(function () {
                    $scope.getAllProjects();
                })
                .post(proj);
            $scope.newProjectName = "";
            $scope.newProjectPath = "";
        };

        $scope.projectUpdate = function (project) {
            console.log("projectUpdate");
            console.dir(project);
            project.$edit = false;
            var proj = {
                name: project.name,
                path: project.path
            };
            materials('/projects/%', project.originalName)
                .success(function (value) {
                    console.dir(value);
                })
                .put(proj);
        };

        $scope.showProject = function (project) {
            $scope.projectName = project.name;
            $scope.projectStatus = project.status;
            materials('/projects/%/tree', $scope.projectName)
                .success(function (tree) {
                    var flattened = $scope.flattenTree(tree);
                    $scope.projectTree = flattened;
                    $scope.displayProject = true;
                })
                .jsonp();
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
    });
