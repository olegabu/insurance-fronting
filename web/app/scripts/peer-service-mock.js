/**
 * @class PeerService
 * @classdesc
 * @ngInject
 */
function PeerService($log, $q, $http, $localStorage, 
    cfg, UserService, IdentityService) {

  // jshint shadow: true
  var PeerService = this;

  var nextPolicyId = 1001;
  var nextClaimId = 2001;
  
  var $storage = $localStorage.$default(cfg);
  
  var createTransaction = function(t) {
    $storage.transactions.push(t);
  };
  
  var getClaim = function(claimId) {
    return _.find($storage.claims, function(o) {
      return o.id === claimId;
    });
  };

  var getContract = function(contractId) {
    return _.find($storage.contracts, function(o) {
      return o.id === contractId;
    });
  };

  PeerService.getClaims = function(policy) {
    return $q(function(resolve) {
      resolve(!policy ? $storage.claims : _.find($storage.claims, function(o) {
        return o.policyId === policy.id;
      }));
    });
  };
  
  PeerService.getContracts = function() {
    return $q(function(resolve) {
      resolve($storage.contracts);
    });
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
        return user.role === 'auditor' || 
        o.from === user.id || o.to === user.id;
      }));
    });
  };
  
  PeerService.getPolicy = function(policyId) {
    return _.find($storage.policies, function(o) {
      return o.id === policyId;
    });
  };
  
  PeerService.createPolicy = function(p) {
    var contract = getContract(p.contractId);
    
    p.id = nextPolicyId++;
    p.supplyChain = angular.copy(contract.supplyChain);
    p.totalPremium = 0;
    p.totalClaim = 0;
    
    $storage.policies.push(p);
    
    contract.totalPolicyCoverage += p.coverage;
    
    contract.totalPolicyPremium += p.premium;
  };

  PeerService.join = function(policyId) {
    var user = IdentityService.getCurrent();

    var policy = PeerService.getPolicy(policyId);

    policy.supplyChain[user.role] = user.id;
  };

  PeerService.approve = function(claimId) {
    var user = IdentityService.getCurrent();

    var claim = getClaim(claimId);

    claim.approvalChain[user.role] = user.id;
    
    if(claim.approvalChain.captive && 
        claim.approvalChain.reinsurer && 
        claim.approvalChain.fronter) {
      var policy = PeerService.getPolicy(claim.policyId);
      var contract = getContract(policy.contractId);
      var purpose = 'claim.' + claim.policyId + '.' + claim.id;
      
      transfer(policy.supplyChain.captive, policy.supplyChain.reinsurer, 
          claim.amt, purpose);
      
      transfer(policy.supplyChain.reinsurer, policy.supplyChain.fronter, 
          claim.amt, purpose);
      
      transfer(policy.supplyChain.fronter, policy.supplyChain.affiliate, 
          claim.amt, purpose);
      
      policy.totalClaim += claim.amt;
      
      contract.totalClaim += claim.amt;
    }
  };

  PeerService.createClaim = function(claim) {
    claim.approvalChain = {};
    
    claim.id = nextClaimId++;
    
    $storage.claims.push(claim);
  };
  
  var transfer = function(from, to, amt, purpose) {
    var fromUser = UserService.getUser(from);
    
    var toUser = UserService.getUser(to);
    
    fromUser.balance -= amt;
    toUser.balance += amt;
    
    createTransaction({
      from: from,
      to: to,
      amt: amt,
      purpose: purpose
    });
  };
  
  PeerService.pay = function(policyId) {
    var policy = PeerService.getPolicy(policyId);
    var contract = getContract(policy.contractId);
    var purpose = 'premium.' + policy.id;
    
    transfer(policy.supplyChain.affiliate, policy.supplyChain.fronter, 
        policy.premium, purpose);
    
    transfer(policy.supplyChain.fronter, policy.supplyChain.reinsurer, 
        policy.premium, purpose);
    
    transfer(policy.supplyChain.reinsurer, policy.supplyChain.captive, 
        policy.premium, purpose);
    
    policy.totalPremium += policy.premium;
    
    contract.totalPremium += policy.premium;
  };

}


angular.module('peerService', []).service('PeerService', PeerService);