'use strict';

describe('Service: Materials', function () {

  // load the service's module
  beforeEach(module('materialsApp'));

  // instantiate service
  var Materials;
  beforeEach(inject(function (_Materials_) {
    Materials = _Materials_;
  }));

  it('should do something', function () {
    expect(!!Materials).toBe(true);
  });

});
