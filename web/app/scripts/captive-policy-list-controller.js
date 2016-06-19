/**
 * @class CaptivePolicyListController
 * @classdesc
 * @ngInject
 */
function CaptivePolicyListController($scope, $log, $interval, $uibModal, 
    PeerService, RoleService) {
  
  /*global PolicyListController*/
  PolicyListController.call(this, $scope, $log, $interval, $uibModal, 
    PeerService, RoleService);

//  var ctl = this;
  

}

angular.module('captivePolicyListController', [])
.controller('CaptivePolicyListController', CaptivePolicyListController);