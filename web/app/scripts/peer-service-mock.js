/**
 * @class PeerService
 * @classdesc
 * @ngInject
 */
function PeerService($log, $q, $http, cfg, UserService) {

  // jshint shadow: true
  var PeerService = this;

  PeerService.getPolicies = function() {
    return cfg.policies;
  };

  PeerService.getTransactions = function() {
    var user = UserService.getUser();
    
    return _.filter(cfg.transactions, function(o) {
      return o.from === user.id || o.to === user.id;
    });
  };

  PeerService.canJoin = function(policy) {
    var user = UserService.getUser();

    return (user.role === 'reinsurer' && 
        policy.supplyChain.captive && !policy.supplyChain.reinsurer) || 
    (user.role === 'fronter' && 
        policy.supplyChain.reinsurer && !policy.supplyChain.fronter) || 
    (user.role === 'affiliate' &&
        policy.supplyChain.fronter && !policy.supplyChain.affiliate);
  };

  PeerService.join = function(policyId) {
    var user = UserService.getUser();

    var policy = _.find(cfg.policies, function(o) {
      return o.id === policyId;
    });

    policy.supplyChain[user.role] = user.id;
  };

  PeerService.canApprove = function(policy, claim) {
    var user = UserService.getUser();

    return (user.role === 'reinsurer' && 
        claim.approvalChain.captive && !claim.approvalChain.reinsurer && 
        policy.supplyChain.reinsurer === user.id) || 
    (user.role === 'captive' && !claim.approvalChain.captive && 
        policy.supplyChain.captive === user.id);
  };

  PeerService.approve = function(policyId, claimId) {
    var user = UserService.getUser();

    var policy = _.find(cfg.policies, function(o) {
      return o.id === policyId;
    });
    
    var claim = policy.claims[claimId];

    claim.approvalChain[user.role] = user.id;
    
    if(claim.approvalChain.captive && claim.approvalChain.reinsurer) {
      transferDown(policy, claim);
    }
  };
  
  PeerService.getMaxClaimAmount = function(policy) {
    var sum = _.sumBy(policy.claims, function(o) {
      return o.amt;
    });
    
    return policy.coverage - sum;
  };

  PeerService.canClaim = function(policy) {
    var user = UserService.getUser();

    return user.role === 'affiliate' && 
    policy.supplyChain.affiliate === user.id &&
    PeerService.getMaxClaimAmount(policy) > 0;
  };

  PeerService.claim = function(policyId, claim) {
    var user = UserService.getUser();

    var policy = _.find(cfg.policies, function(o) {
      return o.id === policyId;
    });
    
    claim.approvalChain = {};
    
    policy.claims.push(claim);
  };

  PeerService.canPay = function(policy) {
    var user = UserService.getUser();

    return user.role === 'affiliate' &&
    policy.supplyChain.affiliate === user.id;
  };
  
  var transfer = function(from, to, amt, purpose) {
    var fromUser = _.find(cfg.users, function(o) {
      return o.id === from;
    });
    
    var toUser = _.find(cfg.users, function(o) {
      return o.id === to;
    });
    
    fromUser.balance -= amt;
    toUser.balance += amt;
    
    cfg.transactions.push({
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
    var claimId = 1 + _.indexOf(policy.claims, claim);
    var purpose = 'claim.' + policy.id + '.' + claimId;
    
    transfer(policy.supplyChain.captive, policy.supplyChain.reinsurer, 
        claim.amt, purpose);
    
    transfer(policy.supplyChain.reinsurer, policy.supplyChain.fronter, 
        claim.amt, purpose);
    
    transfer(policy.supplyChain.fronter, policy.supplyChain.affiliate, 
        claim.amt, purpose);
  };

  PeerService.pay = function(policyId) {
    var user = UserService.getUser();

    var policy = _.find(cfg.policies, function(o) {
      return o.id === policyId;
    });
    
    transferUp(policy);
  };

}


angular.module('peerService', []).service('PeerService', PeerService);