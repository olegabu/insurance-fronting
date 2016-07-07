/**
 * @class PeerService
 * @classdesc
 * @ngInject
 */
function PeerService($log, $q, $http, cfg, IdentityService) {

  // jshint shadow: true
  var PeerService = this;
  
  var getPayload = function(method, functionName, functionArgs) {
    return {
      jsonrpc: '2.0',
      method: method,
      params: {
        type: 1,
        chaincodeID: {
          name: cfg.chaincodeID
        },
        ctorMsg: {
          'function': functionName,
          args: functionArgs
        },
        secureContext: IdentityService.getCurrent().id,
        attributes: ['role', 'company']
      },
      id: 1
    };
  };
  
  var logReject = function(d, o) {
    $log.error(o);
    d.reject(o);
  };

  var invoke = function(functionName, functionArgs) {
    var payload = getPayload('invoke', functionName, functionArgs);
    //$log.debug('payload', JSON.stringify(payload));

    return $http.post(cfg.endpoint, angular.copy(payload)).then(function(res) {
      //$log.debug('result', res.data.result);
    });
  };

  var query = function(functionName, functionArgs) {
    var d = $q.defer();

    var payload = getPayload('query', functionName, functionArgs);
    //$log.debug('payload', JSON.stringify(payload));

    $http.post(cfg.endpoint, angular.copy(payload)).then(function(res) {
      //$log.debug('result', res.data.result);
      
      if(res.data.error) {
        logReject(d, res.data.error);
      }
      else if(res.data.result.status === 'OK') {
        d.resolve(JSON.parse(res.data.result.message));
      }
      else {
        logReject(d, res.data.result);
      }
    });

    return d.promise;
  };
  
  PeerService.getContracts = function() {
    return query('getContracts', []);
  };
  
  PeerService.getPolicies = function() {
    return query('getPolicies', []);
  };
  
  PeerService.getClaims = function() {
    return query('getClaims', []);
  };
  
  PeerService.getTransactions = function() {
    return query('getTransactions', []);
  };
  
  PeerService.getBalance = function() {
    return query('getBalance', []);
  };
  
  PeerService.join = function(policyId) {
    return invoke('join', [policyId]);
  };
  
  PeerService.pay = function(policyId) {
    return invoke('pay', [policyId]);
  };
  
  PeerService.createClaim = function(claim) {
    return invoke('claim', [claim.policyId, '' + claim.amt]);
  };
  
  PeerService.approve = function(claim) {
    return invoke('approve', [claim.policyId, claim.id]);
  };
  
  PeerService.createPolicy = function(p) {
    return invoke('createPolicy', [p.contractId, 
                                   '' + p.coverage, '' + p.premium]);
  };

}


angular.module('peerService', []).service('PeerService', PeerService);