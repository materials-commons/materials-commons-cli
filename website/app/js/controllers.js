
function HomeController($scope) {

}

function ProjectsController($scope, Restangular, $http) {
    var allProjects = Restangular.all('projects');
    allProjects.getList().then(function(projects) {
        $scope.projects = projects;
    });

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