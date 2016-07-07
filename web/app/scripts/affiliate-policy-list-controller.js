/**
 * @class AffiliatePolicyListController
 * @classdesc
 * @ngInject
 */
function AffiliatePolicyListController($scope, $log, $interval, $uibModal, 
    cfg, PeerService, IdentityService) {
  
  /*global PolicyListController*/
  PolicyListController.call(this, $scope, $log, $interval, $uibModal, 
      cfg, PeerService, IdentityService);

  var ctl = this;
  
  ctl.canPay = function(policy) {
    return policy.frontingChain.affiliate === ctl.user.company;
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
.controller('PayModalController', PayModalController);