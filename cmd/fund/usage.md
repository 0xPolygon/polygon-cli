Bulk fund 30 wallets using a pre-deployed contract address.

```bash
$ polycli fund --wallets 20 --verbosity 500
7:08PM INF Starting bulk funding wallets
7:08PM DBG Funder contract deployed address=0x0b589e1cb2457f0ba5a5eef2800d47a4d6fa9fab
7:08PM DBG Funder contract funded balance=1000000000000000000
7:08PM DBG Address(es) of newly generated wallet(s) addresses=["0xc549575af9cebd9940fd0319d2f9a68c157498f5","0x4db281de067d5473fd529530b9da124e558c1ab2","0xe6ebaca5ecb2dba2a03778f934fae6cbe42a6912","0x09cd5e5a6657b4b4f0dc48b63204254d7ded9db6","0xb38cc256482cb84351df5d9da210d726c5b8dd97","0xe66b857f0fceadd5cd015491cce1aa34432c23d3","0x52492c0f89ce6dbac9a92ac2b78d70c82585508d","0xee108e2c179069730977ee4452c4e9ab29d679bf","0xe7846734c489ec26f31a21d1778e8afc8147d329","0xdef6e1ba71a87b26a383a7f7f48aad92d9a3efcc","0x97d4f694dc1c7c99f636e49dd4aa78a486adf174","0x7059775929e96b338765175cdaf5d39ae61c8bd9","0x0cb26d04628b78bd3e7d94056edf1fa775227585","0x195bade7f0f8237e7d8ef314106cecea04f05ed6","0xcd0af0d53bdbdf2a182d3ea8f71dec4bd19378d5","0x6073a6bc815f3a0153a55a4fed3d7fa3f61d48ea","0x97cbc714a5c60c070ca36f4d8b4368c0e9266d20","0xe94ab41b7e0afcf705b68f51e92a31ef70cda280","0x335240977531ba7d1addb252d5b57cb4294e6961","0x4ed675c89cd45e341fd9e74cd7054e22cc02306a"]
7:08PM INF Wallet address(es) and private key(s) saved to file fileName=wallets.json
7:08PM INF Wallet(s) funded! 💸
7:08PM INF Total execution time: 1.042139709s

$ polycli fund --wallets 20 --verbosity 500 --funder-address 0x0b589e1cb2457f0ba5a5eef2800d47a4d6fa9fab
7:09PM INF Starting bulk funding wallets
7:09PM DBG Funder contract funded balance=1000000000000000000
7:09PM DBG Address(es) of newly generated wallet(s) addresses=["0x61568e9430118ee1439a030c28be95403d0b8aec","0x561f610ab2ae4c593def5b8df57829c846a5b493","0xadbb4860895ae926d88742b27233155280eec1c6","0x02c4ee076e5f17626c0c4434273390342134b457","0x3cc6305a743c141c90371861c97c2773a5a2707a","0xfd09346c31782f5ae152ce87ab3534129ecbc25c","0xe9c69be670e243cedcac079386c5733214176b22","0xbe3842385566d6c535af15c9a9b90a627c9e8d27","0x3fd74086309a7ade88f83a53bcd1ff970fc7525a","0x38d84b89d72a1ffa21b12ecb1540fb742148aa1c","0x03508d2cc56dfd7ca82685c4015b1babfb4abcc6","0x127c21979944504a3186cbf52985986b778769f4","0x96661b653ad1dbb1db80e4ce3a842432edf5c5f1","0xc1b85e0986e9200c2d72d42eba347c6f5899607f","0x2e4aad24a103dfb90d0e4b2f55ee343a19c8814a","0xa7fb7da023160a1fd7cc615ea42ea60121a113fc","0xe8dcc4e6454d1fa51beccfb8a4139b7b6d35c08c","0x0cdc66afb63966d4dab6bc0e4782b645839705ae","0xb4b9484609bf36e94860f63e0504c846bb14f3ff","0xeba4f939954a3710e8f3693cbc10bef40a9a6779"]
7:09PM INF Wallet address(es) and private key(s) saved to file fileName=wallets.json
7:09PM INF Wallet(s) funded! 💸
7:09PM INF Total execution time: 1.03648425s
```

Extract from `wallets.json`.

```json
[
  {
    "Address": "0x8098e0092875a89d8db66708eC1dD248D2DD4Dac",
    "PrivateKey": "7f025a5ab0a8699ca79495d8158ddbb9a6b471085a92a20ff39a274235499f22"
  },
  {
    "Address": "0x27F93e701cf7e278687FB1fA1cc9A30932E17587",
    "PrivateKey": "95b2d4484fa219f4ca74df41e04ae8557becf707824d04e4cbaab10c410ac983"
  },
  {
    "Address": "0x9174D0B938C10f89787ebBD37905593cB76a26E4",
    "PrivateKey": "8b6b5032e56dcf95e9c3ef317da5bc41b89538af652d3502973c3ccf3a36fd83"
  },
  {
    "Address": "0x50280161AA7656F57c18Be4ef558786E2c5510C1",
    "PrivateKey": "cdba7d43e36981672d5c73af82298526597a44ef20e4fef3d439283c40c2b8a1"
  },
  ...
]
```

Check the balances of the wallets.

```bash
$ cast balance 0x8098e0092875a89d8db66708eC1dD248D2DD4Dac
50000000000000000

$ cast balance 0x27F93e701cf7e278687FB1fA1cc9A30932E17587
50000000000000000
...
```
