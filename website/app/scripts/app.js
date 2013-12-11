'use strict';

angular.module('materialsApp', [
        'ngCookies',
        'ngResource',
        'ngSanitize',
        'ngRoute',
        'ui.bootstrap',
        'ngTable',
        'btford.socket-io'
    ])
    .config(function ($routeProvider) {
        $routeProvider
            .when('/projects', {
                templateUrl: 'views/projects.html',
                controller: 'ProjectsCtrl'
            })
            .when('/home', {
                templateUrl: 'views/home.html',
                controller: 'HomeCtrl'
            })
            .when('/changes', {
                templateUrl: 'views/changes.html',
                controller: 'ChangesCtrl'
            })
            .when('/provenance', {
                templateUrl: 'views/provenance.html',
                controller: 'ProvenanceCtrl'
            })
            .when('/about', {
                templateUrl: 'views/about.html',
                controller: 'AboutCtrl'
            })
            .when('/contact', {
                templateUrl: 'views/contact.html',
                controller: 'ContactCtrl'
            })
            .otherwise({
                redirectTo: '/home'
            });
    })
    .run(function ($rootScope, socket) {
        $rootScope.$on('$routeChangeStart', function (event, next) {
            if (matchesPartial(next, "views/home", "HomeCtrl")) {
                setActiveMainNav('#home-nav');
            } else if (matchesPartial(next, "views/projects", "ProjectsCtrl")) {
                setActiveMainNav('#projects-nav');
            } else if (matchesPartial(next, "views/changes", "ChangesCtrl")) {
                setActiveMainNav('#changes-nav');
            } else if (matchesPartial(next, "views/provenance", "ProvenanceCtrl")) {
                setActiveMainNav('#prov-nav');
            } else if (matchesPartial(next, "views/about", "AboutCtrl")) {
                setActiveMainNav('#about-nav');
            } else if (matchesPartial(next, "views/contact", "ContactCtrl")) {
                setActiveMainNav('#contact-nav');
            }
        });

        socket.forward('connect');
        socket.forward('file');
    });

function setActiveMainNav(nav) {
    $('#main-nav li').removeClass("active");
    $(nav).addClass("active");
}

function matchesPartial(next, what, controller) {
    if (!next.templateUrl) {
        return false;
    }
    var value = next.templateUrl.indexOf(what) !== -1;
    /*
     Hack to look at controller name to figure out tab. We do this so that partials can be
     shared by controllers, but we need to show which tab is active. So, we look at the
     name of the controller (only if controller != 'ignore').
     */
    if (controller === "ignore") {
        return value;
    }

    if (value) {
        return true;
    }

    return next.controller.toString().indexOf(controller) !== -1;
}