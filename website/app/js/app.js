var app = angular.module('materials', ['ngRoute', 'restangular', 'ngTable', 'mcdirectives']);

app.config(['$routeProvider', '$locationProvider', '$httpProvider', function ($routeProvider) {
    $routeProvider.
        when('/home', {templateUrl: 'partials/home.html', controller: HomeController}).
        when('/projects', {templateUrl: 'partials/projects.html', controller: ProjectsController}).
        when('/changes', {templateUrl: 'partials/changes.html', controller: ChangesController}).
        when('/provenance', {templateUrl: 'partials/provenance.html', controller: ProvenanceController}).
        when('/about', {templateUrl: 'partials/about.html', controller: AboutController}).
        when('/contact', {templateUrl: 'partials/contact.html', controller: ContactController}).
        otherwise({redirectTo: '/home'});
}]);

app.run(function ($rootScope) {

    $rootScope.$on('$routeChangeStart', function (event, next) {
        if (matchesPartial(next, "partials/home", "HomeController")) {
            setActiveMainNav('#home-nav');
        } else if (matchesPartial(next, "partials/projects", "ProjectsController")) {
            setActiveMainNav('#projects-nav');
        } else if (matchesPartial(next, "partials/changes", "ChangesController")) {
            setActiveMainNav('#changes-nav');
        } else if (matchesPartial(next, "partials/provenance", "ProvenanceController")) {
            setActiveMainNav('#prov-nav');
        } else if (matchesPartial(next, "partials/about", "AboutController")) {
            setActiveMainNav('#about-nav');
        } else if (matchesPartial(next, "partials/contact", "ContactController")) {
            setActiveMainNav('#contact-nav');
        }
    });
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

