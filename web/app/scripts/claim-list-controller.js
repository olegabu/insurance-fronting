/**
 * @class ClaimListController
 * @classdesc
 * @ngInject
 */
function ClaimListController($scope, $log, $interval, $uibModal, 
    cfg, PeerService, IdentityService) {

  var ctl = this;
  
  var policies;
  
  var init = function() {
    PeerService.getClaims().then(function(list) {
      ctl.list = list;
    });
    PeerService.getPolicies().then(function(list) {
      policies = _.filter(list, function(o) {
        return o.frontingChain.fronter;
      });
    });
  };
  
  var getPolicy = function(policyId) {
    return _.find(policies, function(o) {
      return o.id === policyId;
    });
  };
  
  $scope.$on('$viewContentLoaded', init);
  
  $interval(init, cfg.refresh);

  ctl.canApprove = function(claim) {
    var user = IdentityService.getCurrent();
    var policy = getPolicy(claim.policyId);

    return policy && (
        (claim.captive && 
              !claim.reinsurer && 
              policy.frontingChain.reinsurer === user.company) || 
        (!claim.captive && 
            policy.frontingChain.captive === user.company) || 
        (claim.captive && 
            claim.reinsurer && 
            !claim.fronter && 
            policy.frontingChain.fronter === user.company)
    );
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
      PeerService.approve(claim);
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
          return policies;
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