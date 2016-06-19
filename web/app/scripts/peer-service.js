/**
 * @class PeerService
 * @classdesc
 * @ngInject
 */
function PeerService($log, $q, $http, cfg, IdentityService) {

  // jshint shadow: true
  var PeerService = this;

  var payload = {
      'jsonrpc': '2.0',
      'params': {
        'type': 1,
        'chaincodeID': {
          name: cfg.chaincodeID
        },
        'ctorMsg': {}
      },
      'id': 0
  };

  var invoke = function(functionName, functionArgs) {
    $log.debug('PeerService.invoke');

    payload.method = 'invoke';
    payload.params.ctorMsg['function'] = functionName;
    payload.params.ctorMsg.args = functionArgs;
    payload.secureContext = IdentityService.getCurrent().id;

    $log.debug('payload', payload);

    return $http.post(cfg.endpoint, angular.copy(payload)).then(function(res) {
      $log.debug('result', res.data.result);
    });
  };

  var query = function(functionName, functionArgs) {
    $log.debug('PeerService.query');
    
    var d = $q.defer();

    payload.method = 'query';
    payload.params.ctorMsg['function'] = functionName;
    payload.params.ctorMsg.args = functionArgs;
    payload.secureContext = IdentityService.getCurrent().id;

    $log.debug('payload', payload);

    $http.post(cfg.endpoint, angular.copy(payload)).then(function(res) {
      $log.debug('result', res.data.result);
      d.resolve(res.data.result);
    });

    return d.promise;
  };
  
  PeerService.getPolicies = function() {
    return query('getPolicies', []);
  };
  
  PeerService.getTransactions = function() {
    return query('getTransactions', []);
  };
  
  PeerService.join = function(policyId) {
    return invoke('join', [policyId]);
  };
  
  PeerService.pay = function(policyId) {
    return invoke('pay', [policyId]);
  };
  
  PeerService.claim = function(policyId, claim) {
    return invoke('claim', [policyId, claim]);
  };

}


angular.module('peerService', []).service('PeerService', PeerService);