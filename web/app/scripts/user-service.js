/**
 * @class UserService
 * @classdesc
 * @ngInject
 */
function UserService($log, $localStorage, cfg) {

  // jshint shadow: true
  var UserService = this;
  
  var $storage = $localStorage.$default({
    users: cfg.users
  });
  
  UserService.getUsers = function() {
    return $storage.users;
  };
  
  UserService.getUser = function(userId) {
    return _.find($storage.users, function(o) {
      return o.id === userId;
    });
  };

}

angular.module('userService', []).service('UserService', UserService);