/**
 * @class PolicyListController
 * @classdesc
 * @ngInject
 */
function PolicyListController($scope, $log, $interval, $uibModal, 
    cfg, PeerService, IdentityService) {

  var ctl = this;
  
  var init = function() {
    PeerService.getPolicies().then(function(list) {
      ctl.list = list;
    });
    
    ctl.user = IdentityService.getCurrent();
  };
  
  $scope.$on('$viewContentLoaded', init);
  
  $interval(init, cfg.refresh);
  
  ctl.canCreate = function() {    
    return ctl.user && ctl.user.role === 'captive';
  };
  
  ctl.canJoin = function(policy) {
    return (ctl.user.role === 'reinsurer' && 
        policy.frontingChain.captive && !policy.frontingChain.reinsurer) || 
    (ctl.user.role === 'fronter' && 
        policy.frontingChain.reinsurer && !policy.frontingChain.fronter) || 
    (ctl.user.role === 'affiliate' &&
        policy.frontingChain.fronter && !policy.frontingChain.affiliate);
  };
  
  ctl.openJoin = function(policy) {
    var modalInstance = $uibModal.open({
      templateUrl: 'join-modal.html',
      controller: 'JoinModalController as ctl',
      resolve: {
        policy: function() {
          return policy;
        }
      }
    });

    modalInstance.result.then(function(policy) {
      PeerService.join(policy.id);
    });
  };
  
}

function JoinModalController($uibModalInstance, policy) {

  var ctl = this;
  
  ctl.policy = policy;
  
  ctl.ok = function () {
    $uibModalInstance.close(ctl.policy);
  };

  ctl.cancel = function () {
    $uibModalInstance.dismiss('cancel');
  };
}

angular.module('policyListController', [])
.controller('PolicyListController', PolicyListController)
.controller('JoinModalController', JoinModalController);