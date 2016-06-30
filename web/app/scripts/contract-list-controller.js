/**
 * @class ContractListController
 * @classdesc
 * @ngInject
 */
function ContractListController($scope, $log, $interval, 
    PeerService, IdentityService) {

  var ctl = this;
  
  var init = function() {
    PeerService.getContracts().then(function(list) {
      ctl.list = list;
    });
  };
  
  $scope.$on('$viewContentLoaded', init);
  
  $interval(init, 1000);
  
}

angular.module('contractListController', [])
.controller('ContractListController', ContractListController);