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
    
    if(ctl.user.role === 'auditor' || ctl.user.role === 'bank') {
      ctl.balance = 0;
    }
    else {
      PeerService.getBalance().then(function(b) {
        ctl.balance = b;
      });
    }
  };
  
  $scope.$on('$viewContentLoaded', init);
  
  $interval(init, 1000);
  
  ctl.compareInt = function(v1, v2) {
    return parseInt(v1.value) < parseInt(v2.value) ? -1 : 1;
  }
  
}

angular.module('transactionListController', [])
.controller('TransactionListController', TransactionListController);