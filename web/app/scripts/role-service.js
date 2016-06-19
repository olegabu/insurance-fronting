/**
 * @class RoleService
 * @classdesc
 * @ngInject
 */
function RoleService($log, IdentityService) {

  // jshint shadow: true
  var RoleService = this;

  RoleService.canApprove = function(policy, claim) {
    var user = IdentityService.getCurrent();

    return (user.role === 'reinsurer' && 
        claim.approvalChain.captive && !claim.approvalChain.reinsurer && 
        policy.supplyChain.reinsurer === user.id) || 
    (user.role === 'captive' && !claim.approvalChain.captive && 
        policy.supplyChain.captive === user.id);
  };

  RoleService.canClaim = function(policy) {
    var user = IdentityService.getCurrent();

    return user.role === 'affiliate' && 
    policy.supplyChain.affiliate === user.id;
  };

  RoleService.canPay = function(policy) {
    var user = IdentityService.getCurrent();

    return user.role === 'affiliate' &&
    policy.supplyChain.affiliate === user.id;
  };

  RoleService.canJoin = function(policy) {
    var user = IdentityService.getCurrent();

    return (user.role === 'reinsurer' && 
        policy.supplyChain.captive && !policy.supplyChain.reinsurer) || 
    (user.role === 'fronter' && 
        policy.supplyChain.reinsurer && !policy.supplyChain.fronter) || 
    (user.role === 'affiliate' &&
        policy.supplyChain.fronter && !policy.supplyChain.affiliate);
  };

}

angular.module('roleService', []).service('RoleService', RoleService);