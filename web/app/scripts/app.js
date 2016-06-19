angular.module('app', ['ui.router',
                       'ui.bootstrap',
                       'ngStorage',
                       'timeService',
                       'userService',
                       'identityService',
                       'roleService',
                       'peerService',
                       'demoController',
                       'policyListController',
                       'captivePolicyListController',
                       'affiliatePolicyListController',
                       'transactionListController',
                       'config'])
                       
.config(function($stateProvider, $urlRouterProvider, $localStorageProvider) {
  
  $localStorageProvider.setKeyPrefix('insurance-fronting-');
  
  $urlRouterProvider.otherwise('/');
  
  $stateProvider
  .state('demo', {
    url: '/',
    templateUrl: 'partials/demo.html',
    controller: 'DemoController as ctl'
  })
  .state('demo.policyList', {
    url: 'policies',
    templateUrl: 'partials/policy-list.html',
    controller: 'PolicyListController as ctl'
  })
  .state('demo.captivePolicyList', {
    url: 'captive-policies',
    templateUrl: 'partials/policy-list.html',
    controller: 'CaptivePolicyListController as ctl'
  })
  .state('demo.affiliatePolicyList', {
    url: 'affiliate-policies',
    templateUrl: 'partials/policy-list.html',
    controller: 'AffiliatePolicyListController as ctl'
  })
  .state('demo.transactionList', {
    url: 'transactions',
    templateUrl: 'partials/transaction-list.html',
    controller: 'TransactionListController as ctl'
  });

});
