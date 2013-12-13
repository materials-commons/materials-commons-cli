'use strict';

angular.module('materialsApp')
    .controller('ChangesCtrl', function ($scope) {
        $scope.alerts = [];
        var filepathLookup = [],
            o,
            obj;
//        $scope.$on('socket:connect', function (ev, data) {
//
//        });

        $scope.$on('socket:file', function (ev, data) {
//            console.dir(data);
            if ($scope.alerts.length >= 100) {
                $scope.alerts.splice(0, 1);
            }
            if (filepathLookup[data.filepath] === undefined) {
                obj = {
                    type: 'success',
                    msg: "File changed: " + data.filepath,
                    event: data.event,
                    count: 1
                };
                filepathLookup[data.filepath] = obj;
                $scope.alerts.push(obj);
            } else {
                $scope.$apply(function () {
                    o = filepathLookup[data.filepath];
                    o.event = data.event;
                    o.count = o.count + 1;
                });
            }
        });

        $scope.closeAlert = function (index) {
            o = $scope.alerts[index];
            delete filepathLookup[o.filepath];
            $scope.alerts.splice(index, 1);
        };
    });
