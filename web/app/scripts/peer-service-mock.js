/**
 * @class PeerService
 * @classdesc
 * @ngInject
 */
function PeerService($log, $q, $http, $localStorage, 
    cfg, UserService, IdentityService) {

  // jshint shadow: true
  var PeerService = this;
  
  var $storage = $localStorage.$default(cfg);
  
  var getPolicy = function(policyId) {
    return _.find($storage.policies, function(o) {
      return o.id === policyId;
    });
  };
  
  var addTransaction = function(t) {
    $storage.transactions.push(t);
  };
  
  var addClaim = function(policy, claim) {
    policy.claims.push(claim);
  };

  var getClaim = function(policy, claimId) {
    return policy.claims[claimId];
  };

  var getClaims = function(policy) {
    return policy.claims;
  };

  var getClaimId = function(policy, claim) {
    return 1 + _.indexOf(policy.claims, claim);
  };
  
  PeerService.getPolicies = function() {
    return $q(function(resolve) {
      resolve($storage.policies);
    });
  };

  PeerService.getTransactions = function() {
    var user = IdentityService.getCurrent();

    return $q(function(resolve) {
      resolve(_.filter($storage.transactions, function(o) {
        return o.from === user.id || o.to === user.id;
      }));
    });
  };

  PeerService.canJoin = function(policy) {
    var user = IdentityService.getCurrent();

    return (user.role === 'reinsurer' && 
        policy.supplyChain.captive && !policy.supplyChain.reinsurer) || 
    (user.role === 'fronter' && 
        policy.supplyChain.reinsurer && !policy.supplyChain.fronter) || 
    (user.role === 'affiliate' &&
        policy.supplyChain.fronter && !policy.supplyChain.affiliate);
  };

  PeerService.join = function(policyId) {
    var user = IdentityService.getCurrent();

    var policy = getPolicy(policyId);

    policy.supplyChain[user.role] = user.id;
  };

  PeerService.canApprove = function(policy, claim) {
    var user = IdentityService.getCurrent();

    return (user.role === 'reinsurer' && 
        claim.approvalChain.captive && !claim.approvalChain.reinsurer && 
        policy.supplyChain.reinsurer === user.id) || 
    (user.role === 'captive' && !claim.approvalChain.captive && 
        policy.supplyChain.captive === user.id);
  };

  PeerService.approve = function(policyId, claimId) {
    var user = IdentityService.getCurrent();

    var policy = getPolicy(policyId);

    var claim = getClaim(policy, claimId);

    claim.approvalChain[user.role] = user.id;
    
    if(claim.approvalChain.captive && claim.approvalChain.reinsurer) {
      transferDown(policy, claim);
    }
  };
  
  PeerService.getMaxClaimAmount = function(policy) {
    var sum = _.sumBy(getClaims(policy), function(o) {
      return o.amt;
    });
    
    return policy.coverage - sum;
  };

  PeerService.canClaim = function(policy) {
    var user = IdentityService.getCurrent();

    return user.role === 'affiliate' && 
    policy.supplyChain.affiliate === user.id &&
    PeerService.getMaxClaimAmount(policy) > 0;
  };

  PeerService.claim = function(policyId, claim) {
    var policy = getPolicy(policyId);
    
    claim.approvalChain = {};
    
    addClaim(policy, claim);
  };

  PeerService.canPay = function(policy) {
    var user = IdentityService.getCurrent();

    return user.role === 'affiliate' &&
    policy.supplyChain.affiliate === user.id;
  };
  
  var transfer = function(from, to, amt, purpose) {
    var fromUser = UserService.getUser(from);
    
    var toUser = UserService.getUser(to);
    
    fromUser.balance -= amt;
    toUser.balance += amt;
    
    addTransaction({
      from: from,
      to: to,
      amt: amt,
      purpose: purpose
    });
  };
  
  var transferUp = function(policy) {
    var purpose = 'premium.' + policy.id;
    
    transfer(policy.supplyChain.affiliate, policy.supplyChain.fronter, 
        policy.premium, purpose);
    
    transfer(policy.supplyChain.fronter, policy.supplyChain.reinsurer, 
        policy.premium, purpose);
    
    transfer(policy.supplyChain.reinsurer, policy.supplyChain.captive, 
        policy.premium, purpose);
  };
  
  var transferDown = function(policy, claim) {
    var claimId = getClaimId(policy, claim);
    var purpose = 'claim.' + policy.id + '.' + claimId;
    
    transfer(policy.supplyChain.captive, policy.supplyChain.reinsurer, 
        claim.amt, purpose);
    
    transfer(policy.supplyChain.reinsurer, policy.supplyChain.fronter, 
        claim.amt, purpose);
    
    transfer(policy.supplyChain.fronter, policy.supplyChain.affiliate, 
        claim.amt, purpose);
  };

  PeerService.pay = function(policyId) {
    var policy = getPolicy(policyId);
    
    transferUp(policy);
  };

}


angular.module('peerService', []).service('PeerService', PeerService);