/**
 * @class IdentityService
 * @classdesc
 * @ngInject
 */
function IdentityService($log, UserService) {

  // jshint shadow: true
  var IdentityService = this;

  var user = UserService.getUsers()[0];
  
  IdentityService.setCurrent = function(u) {
    user = UserService.getUser(u.id);
  };
  
  IdentityService.getCurrent = function() {
    return UserService.getUser(user.id);
  };

}

angular.module('identityService', []).service('IdentityService', IdentityService);