angular.module('config', [])
.constant('cfg', {
  endpoint: 'https://deadbeef-6d99-41a8-988a-c31ce4ae14dc_vp1-api.blockchain.ibm.com:443/chaincode',
  chaincodeID: 'badeadbeeff1ec5ab1d077cb38907780c79260420fd94c288d6bddb710114c44969b9d560d365e38b62d379681e2acaa5b88536ebd22d6ed56c0605736349',
  users: [{
    id: 'bermuda', 
    role: 'captive',
    balance: 10000000
  },
  {
    id: 'art', 
    role: 'reinsurer',
    balance: 0
  },
  {
    id: 'allianz', 
    role: 'fronter',
    balance: 0
  },
  {
    id: 'nigeria', 
    role: 'affiliate',
    balance: 1000
  },
  {
    id: 'citi', 
    role: 'bank',
    balance: 100000000
  }],
  policies: [{
    id: 1000, 
    coverage: 1000000,
    premium: 100,
    supplyChain: {
      captive: 'bermuda',
      /*reinsurer: 'art',
      fronter: 'allianz',
      affiliate: 'nigeria'*/
    },
    claims: [/*{
      amt: 100000,
      approvalChain: {
        captive: 'bermuda',
        reinsurer: 'art'}
    },
    {
      amt: 50000,
      approvalChain: {
        captive: 'bermuda',
        reinsurer: 'art'}
    }*/]
  }],
  transactions: [{
    from: 'citi', 
    to: 'bermuda', 
    amt: 10000000,
    purpose: 'cash'
  },
  {
    from: 'citi', 
    to: 'nigeria', 
    amt: 10000,
    purpose: 'cash'
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
