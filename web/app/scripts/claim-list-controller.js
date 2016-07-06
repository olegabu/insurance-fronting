/**
 * @class ClaimListController
 * @classdesc
 * @ngInject
 */
function ClaimListController($scope, $log, $interval, $uibModal, 
    PeerService, IdentityService) {

  var ctl = this;
  
  var init = function() {
    PeerService.getClaims().then(function(list) {
      ctl.list = list;
    });
  };
  
  $scope.$on('$viewContentLoaded', init);
  
  $interval(init, 1000);

  ctl.canApprove = function(claim) {
    var user = IdentityService.getCurrent();
    var policy = PeerService.getPolicy(claim.policyId);

    return (claim.approvalChain.captive && 
              !claim.approvalChain.reinsurer && 
              policy.frontingChain.reinsurer === user.id) || 
        (!claim.approvalChain.captive && 
            policy.frontingChain.captive === user.id) || 
        (claim.approvalChain.captive && 
            claim.approvalChain.reinsurer && 
            !claim.approvalChain.fronter && 
            policy.frontingChain.fronter === user.id);
  };
  
  ctl.openApprove = function(claim) {
    var modalInstance = $uibModal.open({
      templateUrl: 'approve-modal.html',
      controller: 'ApproveModalController as ctl',
      resolve: {
        claim: function() {
          return claim;
        }
      }
    });

    modalInstance.result.then(function() {
      PeerService.approve(claim.id);
    });
  };

  ctl.canClaim = function() {
    var user = IdentityService.getCurrent();

    return user.role === 'affiliate';
  };
  
  ctl.openClaim = function() {
    var modalInstance = $uibModal.open({
      templateUrl: 'claim-modal.html',
      controller: 'ClaimModalController as ctl',
      resolve: {
        policies: function() {
          return PeerService.getPolicies();
        }
      }
    });

    modalInstance.result.then(function(claim) {
      PeerService.createClaim(claim);
    });
  };

}

function ApproveModalController($uibModalInstance, claim) {

  var ctl = this;
  
  ctl.claim = claim;
  
  ctl.ok = function () {
    $uibModalInstance.close();
  };

  ctl.cancel = function () {
    $uibModalInstance.dismiss('cancel');
  };
}

function ClaimModalController($uibModalInstance, PeerService, policies) {

  var ctl = this;
  
  ctl.policies = policies;
  
  ctl.ok = function () {
    ctl.claim.policyId = ctl.policy.id;
    $uibModalInstance.close(ctl.claim);
  };

  ctl.cancel = function () {
    $uibModalInstance.dismiss('cancel');
  };
}

angular.module('claimListController', [])
.controller('ClaimListController', ClaimListController)
.controller('ApproveModalController', ApproveModalController)
.controller('ClaimModalController', ClaimModalController);