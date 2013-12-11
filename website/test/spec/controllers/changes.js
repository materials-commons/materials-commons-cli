'use strict';

describe('Controller: ChangesCtrl', function () {

  // load the controller's module
  beforeEach(module('materialsApp'));

  var ChangesCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    ChangesCtrl = $controller('ChangesCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });
});
