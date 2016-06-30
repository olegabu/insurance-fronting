/**
 * @class PolicyListController
 * @classdesc
 * @ngInject
 */
function PolicyListController($scope, $log, $interval, $uibModal, 
    PeerService, IdentityService) {

  var ctl = this;
  
  var init = function() {
    PeerService.getPolicies().then(function(list) {
      ctl.list = list;
    });
    
    ctl.user = IdentityService.getCurrent();
  };
  
  $scope.$on('$viewContentLoaded', init);
  
  $interval(init, 1000);
  
  ctl.canCreate = function() {    
    return ctl.user && ctl.user.role === 'captive';
  }
  
  ctl.canJoin = function(policy) {
    return (ctl.user.role === 'reinsurer' && 
        policy.supplyChain.captive && !policy.supplyChain.reinsurer) || 
    (ctl.user.role === 'fronter' && 
        policy.supplyChain.reinsurer && !policy.supplyChain.fronter) || 
    (ctl.user.role === 'affiliate' &&
        policy.supplyChain.fronter && !policy.supplyChain.affiliate);
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