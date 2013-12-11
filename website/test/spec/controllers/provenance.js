'use strict';

describe('Controller: ProvenanceCtrl', function () {

  // load the controller's module
  beforeEach(module('materialsApp'));

  var ProvenanceCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    ProvenanceCtrl = $controller('ProvenanceCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });
});
