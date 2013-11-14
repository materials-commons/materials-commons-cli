
function HomeController($scope) {

}

function ProjectsController($scope, Restangular) {
    var allProjects = Restangular.all('projects');
    allProjects.getList().then(function(projects) {
        $scope.projects = projects;
    });

    Restangular.one("projects", "a").customGET("tree").then(function(tree) {
        console.dir(tree);
        $scope.projectTree = tree;
    });

//    Restangular.oneUrl("datadirs", "http://localhost:5000/datadirs/tree?apikey=4a3ec8f43cc511e3ba368851fb4688d4")
//        .get().then(function(tree) {
//            console.dir(tree);
//        })
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