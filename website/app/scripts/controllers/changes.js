'use strict';

angular.module('materialsApp')
    .controller('ChangesCtrl', function ($scope) {
        $scope.alerts = [];
        var filepathLookup = [],
            o,
            obj;
        $scope.$on('socket:connect', function (ev, data) {
            console.log("on connect");
            //console.dir(data);
            //console.log(data);
        });

        $scope.$on('socket:file', function (ev, data) {
            console.dir(data);
            if ($scope.alerts.length >= 100) {
                $scope.alerts.splice(0, 1);
            }
            if (filepathLookup[data.filepath] === undefined) {
                obj = {
                    type: 'success',
                    msg: "File changed: " + data.filepath,
                    count: 1
                };
                filepathLookup[data.filepath] = obj;
                $scope.alerts.push(obj);
            } else {
                console.log("Already saw: " + data.filepath);
                $scope.$apply(function () {
                    o = filepathLookup[data.filepath];
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
