angular.module('config', [])
.constant('cfg', {
  endpoint: 'http://localhost:5000/chaincode',
  chaincodeID: 'e0588c90defacf615671d8bddc960bf5fc4a540ffa0039fca30fbe2be8c75aeae40272b207c436ebe5be9657db60c52c99c3a1abd94f74f3beff6f2f6bb43177',
  users: [{
    id: 'bermuda', 
    role: 'captive',
    company: 'Bermuda',
    balance: 10000000
  },
  {
    id: 'art', 
    role: 'reinsurer',
    company: 'Art',
    balance: 0
  },
  {
    id: 'allianz', 
    role: 'fronter',
    company: 'Allianz',
    balance: 0
  },
  {
    id: 'nigeria', 
    role: 'affiliate',
    company: 'Nigeria',
    balance: 10000
  },
  {
    id: 'citi',  
    company: 'Citi',
    role: 'bank',
    balance: 100000000
  },
  {
    id: 'auditor',  
    company: 'Auditor',
    role: 'auditor',
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
