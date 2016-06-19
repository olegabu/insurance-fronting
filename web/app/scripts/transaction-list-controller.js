/**
 * @class TransactionListController
 * @classdesc
 * @ngInject
 */
function TransactionListController($scope, $log, $interval, 
    PeerService, IdentityService) {

  var ctl = this;
  
  var init = function() {
    ctl.list = PeerService.getTransactions();
    
    ctl.user = IdentityService.getCurrent();
  };
  
  $scope.$on('$viewContentLoaded', init);
  
  $interval(init, 1000);
  
}

angular.module('transactionListController', [])
.controller('TransactionListController', TransactionListController);