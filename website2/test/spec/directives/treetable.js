'use strict';

describe('Directive: treetable', function () {

  // load the directive's module
  beforeEach(module('materialsApp'));

  var element,
    scope;

  beforeEach(inject(function ($rootScope) {
    scope = $rootScope.$new();
  }));

  it('should make hidden element visible', inject(function ($compile) {
    element = angular.element('<treetable></treetable>');
    element = $compile(element)(scope);
    expect(element.text()).toBe('this is the treetable directive');
  }));
});
