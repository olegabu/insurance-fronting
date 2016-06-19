/**
 * @class TransactionListController
 * @classdesc
 * @ngInject
 */
function TransactionListController($scope, $log, $interval, 
    PeerService, IdentityService) {

  var ctl = this;
  
  var init = function() {
    PeerService.getTransactions().then(function(list) {
      ctl.list = list;
    });
    
    ctl.user = IdentityService.getCurrent();
  };
  
  $scope.$on('$viewContentLoaded', init);
  
  $interval(init, 1000);
  
}

angular.module('transactionListController', [])
.controller('TransactionListController', TransactionListController);