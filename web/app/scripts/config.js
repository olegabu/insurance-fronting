angular.module('config', [])
.constant('cfg', {
  refresh: 2000,
  chaincodeID: '61fc9184862f93176c08f049d3e512b6909edf39e1fcf0ab2d8bf7529ac9993125eba0b6faafcdacba99d903969a3031d9b0e2d1b0dce55bb43ec82a90ff82a1',
  users: [{
    id: 'bermuda', 
    role: 'captive',
    company: 'Bermuda',
    endpoint: 'http://54.210.64.11:5000/chaincode',
    balance: 10000000
  },
  {
    id: 'art', 
    role: 'reinsurer',
    company: 'Art',
    endpoint: 'http://52.91.72.177:5000/chaincode',
    balance: 0
  },
  {
    id: 'allianz', 
    role: 'fronter',
    company: 'Allianz',
    endpoint: 'http://54.198.39.96:5000/chaincode',
    balance: 0
  },
  {
    id: 'nigeria', 
    role: 'affiliate',
    company: 'Nigeria',
    endpoint: 'http://54.175.130.125:5000/chaincode',
    balance: 10000
  },
  {
    id: 'citi',  
    company: 'Citi',
    role: 'bank',
    endpoint: 'http://54.210.64.11:5000/chaincode',
    balance: 100000000
  },
  {
    id: 'auditor',  
    company: 'Auditor',
    role: 'auditor',
    endpoint: 'http://54.210.64.11:5000/chaincode',
    balance: 0
  }],
  contracts: [{
    id: 10, 
    maxCoverage: 10000000,
    maxPremium: 1000,
    currentTotalCoverage: 0,
    currentTotalPremium: 0,
    currentPaidClaim: 0,
    currentPaidPremium: 0,
    captive: 'bermuda',
    reinsurer: 'art'
  }],
  policies: [/*{
    id: 1000, 
    contractId: 10,
    coverage: 1000000,
    premium: 100,
    frontingChain: {
      captive: 'bermuda',
      reinsurer: 'art',
      fronter: 'allianz',
      affiliate: 'nigeria'
    }
  }*/],
  claims: [/*{
    id: 2000,
    policyId: 1000,
    amt: 100000,
    approvalChain: {
      captive: 'bermuda',
      reinsurer: 'art'}
  },
  {
    id: 2001,
    policyId: 1000,
    amt: 50000,
    approvalChain: {
      captive: 'bermuda',
      reinsurer: 'art'}
  }*/],
  transactions: [{
    from: 'citi', 
    to: 'bermuda', 
    amt: 10000000,
    purpose: 'buy coins'
  },
  {
    from: 'citi', 
    to: 'nigeria', 
    amt: 10000,
    purpose: 'buy coins'
  },
  /*{
    from: 'nigeria', 
    to: 'aliianz', 
    amt: 100,
    purpose: 'premium.1000'
  },
  {
    from: 'aliianz', 
    to: 'art', 
    amt: 100,
    purpose: 'premium.1000'
  },
  {
    from: 'art', 
    to: 'nigeria', 
    amt: 100,
    purpose: 'premium.1000'
  }*/],
});
