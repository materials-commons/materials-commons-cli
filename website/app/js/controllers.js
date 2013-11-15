
function HomeController($scope) {

}

function ProjectsController($scope, Restangular, $http) {
    $scope.projectsData = [];
    var allProjects = Restangular.all('projects');
    allProjects.getList().then(function(projects) {
        angular.forEach(projects, function(project) {
            $scope.projectsData.push({Name: project.name, Path: project.path, Status: "Unloaded"})
        });
//        console.dir($scope.projects);
//        if (!$scope.$$phase) {
//            $scope.$apply()
//        }
    });

    $scope.selected = [];

    $scope.projectGridOptions = {
        data: 'projectsData',
        multiSelect: false,
        selectedItems: $scope.selected,
        afterSelectionChange: function() {
            //console.dir($scope.selected);
        }
    };

    Restangular.one("projects", "a").customGET("tree").then(function(tree) {
        $scope.projectTree = tree;
    });
}

function ChangesController($scope) {

}

function ProvenanceController($scope) {

}

function AboutController($scope) {

}

function ContactController($scope) {

}

function EventController($scope) {

}