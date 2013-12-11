angular.module('materialsApp')
    .controller('ProjectsCtrl', function ($scope, materials) {
        'use strict';



        $scope.projectsData = [];

        $scope.getAllProjects = function () {
            materials('/projects')
                .success(function (projects) {
                    projects.forEach(function (project) {
                        project.originalName = project.name;
                    });
                    $scope.projects = projects;
                })
                .getJson();
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
                .error(function() {
                    $scope.projectTree = [];
                    $scope.displayProject = false;
                })
                .getJson();
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
