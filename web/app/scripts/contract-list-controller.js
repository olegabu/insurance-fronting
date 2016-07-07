/**
 * @class ContractListController
 * @classdesc
 * @ngInject
 */
function ContractListController($scope, $log, $interval, 
    cfg, PeerService, IdentityService) {

  var ctl = this;
  
  var init = function() {
    PeerService.getContracts().then(function(list) {
      ctl.list = list;
    });
  };
  
  $scope.$on('$viewContentLoaded', init);
  
  $interval(init, cfg.refresh);
  
}

angular.module('contractListController', [])
.controller('ContractListController', ContractListController);