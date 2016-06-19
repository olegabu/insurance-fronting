/**
 * @class AffiliatePolicyListController
 * @classdesc
 * @ngInject
 */
function AffiliatePolicyListController($scope, $log, $interval, $uibModal, 
    PeerService, RoleService) {
  
  /*global PolicyListController*/
  PolicyListController.call(this, $scope, $log, $interval, $uibModal, 
    PeerService, RoleService);

  var ctl = this;
  
  var getMaxClaimAmount = function(policy) {
    var sum = _.sumBy(policy.claims, function(o) {
      return o.amt;
    });
    
    return policy.coverage - sum;
  };
  
  ctl.canClaim = function(policy) {
    return RoleService.canClaim(policy) && getMaxClaimAmount(policy) > 0;
  };
  
  ctl.openClaim = function(policy) {
    var modalInstance = $uibModal.open({
      templateUrl: 'claim-modal.html',
      controller: 'ClaimModalController as ctl',
      resolve: {
        policy: function() {
          return policy;
        },
        maxClaimAmount: function() {
          return getMaxClaimAmount(policy);
        }
      }
    });

    modalInstance.result.then(function(claim) {
      PeerService.claim(policy.id, claim);
    });
  };
  
  ctl.canPay = function(policy) {
    return RoleService.canPay(policy);
  };
  
  ctl.openPay = function(policy) {
    var modalInstance = $uibModal.open({
      templateUrl: 'pay-modal.html',
      controller: 'PayModalController as ctl',
      resolve: {
        policy: function() {
          return policy;
        }
      }
    });

    modalInstance.result.then(function() {
      PeerService.pay(policy.id);
    });
  };
}

function ClaimModalController($uibModalInstance, PeerService, 
    policy, maxClaimAmount) {

  var ctl = this;
  
  ctl.policy = policy;
  ctl.maxClaimAmount = maxClaimAmount;
  
  ctl.ok = function () {
    $uibModalInstance.close(ctl.claim);
  };

  ctl.cancel = function () {
    $uibModalInstance.dismiss('cancel');
  };
}

function PayModalController($uibModalInstance, policy) {

  var ctl = this;
  
  ctl.policy = policy;
  
  ctl.ok = function () {
    $uibModalInstance.close();
  };

  ctl.cancel = function () {
    $uibModalInstance.dismiss('cancel');
  };
}

angular.module('affiliatePolicyListController', [])
.controller('AffiliatePolicyListController', AffiliatePolicyListController)
.controller('ClaimModalController', ClaimModalController)
.controller('PayModalController', PayModalController);