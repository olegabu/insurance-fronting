/**
 * @class PolicyListController
 * @classdesc
 * @ngInject
 */
function PolicyListController($scope, $log, $interval, $uibModal, 
    PeerService) {

  var ctl = this;
  
  var init = function() {
    PeerService.getPolicies().then(function(list) {
      ctl.list = list;
    });
  };
  
  $scope.$on('$viewContentLoaded', init);
  
  $interval(init, 1000);
  
  ctl.canJoin = function(policy) {
    return PeerService.canJoin(policy);
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
  
  ctl.canApprove = function(policy, claim) {
    return PeerService.canApprove(policy, claim);
  };
  
  ctl.openApprove = function(policy, claim) {
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
      PeerService.approve(policy.id, _.indexOf(policy.claims, claim));
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

angular.module('policyListController', [])
.controller('PolicyListController', PolicyListController)
.controller('JoinModalController', JoinModalController)
.controller('ApproveModalController', ApproveModalController);