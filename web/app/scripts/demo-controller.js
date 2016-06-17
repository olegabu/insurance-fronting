/**
 * @class DemoController
 * @classdesc
 * @ngInject
 */
function DemoController($log, $state, 
    cfg, TimeService, UserService, PeerService) {

  var ctl = this;

  ctl.now = function() {
    return TimeService.now;
  };

  ctl.tick = function() {
    TimeService.tick();
  };

  ctl.clock = function() {
    TimeService.clock();
  };
  
  ctl.user = UserService.getUser();

  ctl.users = UserService.getUsers();
  
  ctl.setUser = function() {
    UserService.setUser(ctl.user);

    if(ctl.user.role === 'captive') {
      $state.go('demo.captivePolicyList');
    }
    else if(ctl.user.role === 'affiliate') {
      $state.go('demo.affiliatePolicyList');
    }
    else if(ctl.user.role === 'bank') {
      $state.go('demo.transactionList');
    }
    else {
      $state.go('demo.policyList');
    }
  };

}

angular.module('demoController', [])
.controller('DemoController', DemoController);