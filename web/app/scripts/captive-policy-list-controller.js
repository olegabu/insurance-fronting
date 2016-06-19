/**
 * @class CaptivePolicyListController
 * @classdesc
 * @ngInject
 */
function CaptivePolicyListController($scope, $log, $interval, $uibModal, 
    PeerService) {
  
  /*global PolicyListController*/
  PolicyListController.call(this, $scope, $log, $interval, $uibModal, 
    PeerService);

//  var ctl = this;
  

}

angular.module('captivePolicyListController', [])
.controller('CaptivePolicyListController', CaptivePolicyListController);