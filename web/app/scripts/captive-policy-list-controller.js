/**
 * @class CaptivePolicyListController
 * @classdesc
 * @ngInject
 */
function CaptivePolicyListController($scope, $log, $interval, $uibModal, 
    cfg, PeerService, IdentityService) {
  
  /*global PolicyListController*/
  PolicyListController.call(this, $scope, $log, $interval, $uibModal, 
      cfg, PeerService, IdentityService);

  var ctl = this;
  
  ctl.open = function() {
    var modalInstance = $uibModal.open({
      templateUrl: 'policy-modal.html',
      controller: 'PolicyModalController as ctl',
      resolve: {
        contracts: function() {
          return PeerService.getContracts();
        }
      }
    });

    modalInstance.result.then(function(policy) {
      PeerService.createPolicy(policy);
    });
  };

}

function PolicyModalController($uibModalInstance, PeerService, contracts) {

  var ctl = this;
  
  ctl.contracts = contracts;
  
  ctl.ok = function () {
    ctl.policy.contractId = ctl.contract.id;
    $uibModalInstance.close(ctl.policy);
  };

  ctl.cancel = function () {
    $uibModalInstance.dismiss('cancel');
  };
}

angular.module('captivePolicyListController', [])
.controller('CaptivePolicyListController', CaptivePolicyListController)
.controller('PolicyModalController', PolicyModalController);