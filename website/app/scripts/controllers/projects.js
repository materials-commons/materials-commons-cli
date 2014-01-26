angular.module('materialsApp')
    .controller('ProjectsCtrl', function ($scope, materials) {
        'use strict';
        $scope.projectsData = [];

        $scope.getAllProjects = function () {
            materials('/projects')
                .success(function (projects) {
                    projects.forEach(function (project) {
                        project.originalName = project.Name;
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
                Name: $scope.newProjectName,
                Path: $scope.newProjectPath,
                Status: "Unloaded"
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
                Name: project.name,
                Path: project.path
            };
            materials('/projects/%', project.originalName)
                .success(function (value) {
                    console.dir(value);
                })
                .put(proj);
        };

        $scope.trail = [];

        $scope.showProject2 = function (project) {
            $scope.projectName = project.Name;
            $scope.projectStatus = project.Status;
            materials('/projects/%/tree', $scope.projectName)
                .success(function (tree) {
                    $scope.trail.push(tree[0]);
                    $scope.tree = tree;
                    $scope.dir = $scope.tree[0].children;
                    $scope.displayProject2 = true;
                }).getJson();
        }
        $scope.showProject = function (project) {
            $scope.projectName = project.Name;
            $scope.projectStatus = project.Status;
            materials('/projects/%/tree', $scope.projectName)
                .success(function (tree) {
                    var flattened = $scope.flattenTree(tree);
                    $scope.projectTree = flattened;
                    $scope.displayProject = true;
                })
                .error(function () {
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

        $scope.openFolder = function (item) {
            var e = _.find($scope.trail, function (item1) {
                return item1.id === item.id;
            });

            if (typeof e === 'undefined') {
                // first level is 0 so we need to add 1 to our test
                if (item.level+1 <= $scope.trail.length) {
                    // Remove everything at this level and above
                    $scope.trail = $scope.trail.splice(0, item.level);
                }
                $scope.trail.push(item);
            }

            $scope.dir = item.children;
        }

        $scope.backToFolder = function (item) {
            $scope.dir = item.children;
        }
    });
