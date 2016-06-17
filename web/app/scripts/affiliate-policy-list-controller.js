/**
 * @class AffiliatePolicyListController
 * @classdesc
 * @ngInject
 */
function AffiliatePolicyListController($scope, $log, $interval, $uibModal, 
    PeerService) {
  
  /*global PolicyListController*/
  PolicyListController.call(this, $scope, $log, $interval, $uibModal, 
    PeerService);

  var ctl = this;
  
  ctl.canClaim = function(policy) {
    return PeerService.canClaim(policy);
  };
  
  ctl.openClaim = function(policy) {
    var modalInstance = $uibModal.open({
      templateUrl: 'claim-modal.html',
      controller: 'ClaimModalController as ctl',
      resolve: {
        policy: function() {
          return policy;
        }
      }
    });

    modalInstance.result.then(function(claim) {
      PeerService.claim(policy.id, claim);
    });
  };
  
  ctl.canPay = function(policy) {
    return PeerService.canPay(policy);
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

function ClaimModalController($uibModalInstance, PeerService, policy) {

  var ctl = this;
  
  ctl.policy = policy;
  
  ctl.maxClaim = PeerService.getMaxClaimAmount(policy);
  
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