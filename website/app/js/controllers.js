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

    $scope.projectGridOptions = {
        data: 'projects',
        multiSelect: false,
        columnDefs: [
            {field: 'name', displayName: 'Name'},
            {field: 'path', displayName: 'Path'},
            {field: 'status', displayName: 'Status',
                cellTemplate: '<div>{{ row.entity[col.field] }}  <button ng-click="statusButtonAction(row.entity)">{{ statusButtonName(row.entity[col.field]) }}</button></div>'}
        ],
        selectedItems: $scope.selected,
        afterSelectionChange: function () {
            //console.dir($scope.selected);
        }
    };

    Restangular.one("projects", "a").customGET("tree").then(function (tree) {
        $scope.projectTree = tree;
    });
}

function ChangesController($scope, Restangular, $timeout) {
    $scope.events = [];
    (function tick() {
        console.log("tick")
        Restangular.all('projects/changes').getList().then(function(eventsInfo) {
            //console.dir(eventsInfo);
            var found = false;
            $scope.events.forEach(function(event) {
                if (event.filepath == eventsInfo.filepath) {
                    found = true;
                }
            });

            if (! found) {
                $scope.events.push(eventsInfo);
            }
        })
        $timeout(tick, 3000);
    })();


}

function ProvenanceController($scope) {

}

function AboutController($scope) {

}

function ContactController($scope) {

}

function EventController($scope) {

}