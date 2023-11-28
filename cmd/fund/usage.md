Bulk fund 20 wallets.

```bash
$ polycli fund --wallets 20 --verbosity 500
9:09AM INF Starting bulk funding wallets
9:09AM DBG Funder contract deployed address=0x0775aafb6dd38417581f7c583053fa3b78fd4fd1
9:09AM DBG Funder contract funded weiBalance=1000000000000000000000
9:09AM DBG Address(es) of newly generated wallet(s) addresses=["0xc88d3686b71339874b1a620aad97fd3779ac2473","0xdd7c38af9a7ebd359e3caabd785c01070ed81ab5","0x29b6d4d6a6f44e86895f93f6e10fd29d44e58efb","0xd6f861e41548262efbf4cb33bdc835be45b3ba7a","0xdc09658e1cc143be77efef47e4f62565fd83e4a5","0xe32e123028fd23bbfc984726b1e8c4f4c4de5be3","0x94535383adca5756c18622b7efd4572a8425f038","0x37b5ded15b1ffcf27f35d031da0ce0d35a68d23b","0xdbd386fb367b264490455b6942d2fe2996fdda2f","0x54376b36a718ff3d92b4d3f890eb131d892d95ba","0xf6a8de5b601c54dd95615b67ef99b8164d35591e","0x84793a9e49a842a20799ea45deee35ccaf54cc46","0x95ca4510c973b516e0e873b35949f83154ecf562","0xa41493dac048ed151c5424b23038c504ae0e08cf","0xd0ad2bb3ddd47921d8a4f8dcc165841abdcecd19","0xc771ab4d61e8a4f0b001abcfd396212828152448","0x2fd80f9296b89f645ea15e2908c654737497fab5","0xaffa3cb0e6802bfefcae2ed4cfc6dfd196a131a0","0x7cc77413c836fb78be80af25c5352b9876f09bd3","0x4be625cfdeff7b9333c52e75d4756216a04a1f51"]
9:09AM INF Wallet address(es) and private key(s) saved to file fileName=wallets.json
9:09AM INF Wallet(s) funded! 💸
9:09AM INF Total execution time: 1.027154333s
```

Bulk fund 20 wallets using a pre-deployed contract address.

```bash
$ polycli fund --wallets 20 --verbosity 500 --funder-address 0x0775aafb6dd38417581f7c583053fa3b78fd4fd1
9:09AM INF Starting bulk funding wallets
9:09AM DBG Address(es) of newly generated wallet(s) addresses=["0xc1a44b1e37ee1fca4c6fd5562c730d5b8525e4c6","0x5d8121cf716b70d3e345adb58157752304eed5c3","0xf576372fbabd14d8574c4bf54f1c666b078e76d2","0xea2231483fa6f5d6ca4ce943735e97f29dc4d2ba","0xd0dafc79ab3ec231e80b4207351f161480c4c86a","0x82c609438bfbbedb5e83b092cb7b21f67de8355d","0x17b97eb872c08f6521e175a07a41608099b304e2","0x72fb75d543b27b003113ecae7f46c8bb05864caf","0x78d7b672cab700e3cfda382f9290cc3a6f5d5daf","0xa738d3aca4cd35ae69dd09b764445538df38c142","0xb187e5e4d9b4fe2ba0f23716e795b74887f6ec95","0xdf92375ac934c768eff96154716f554760d70cb7","0xb281194c695a6ef5208b760030421f2ae29b65f7","0x77b68833c0861bb8bb1925ef4086db2c67e6d0e3","0xdacf0b0cdc6c99f03e639fa4797f26e31cf3a2c9","0x980354e6a8cb69cd8b7e9d8c23b25df4c140b27f","0x9f82b0e526fed31bbd90a544b408ee9918cd90b4","0x8d9458175f8f5406585dca086a088e7a1aae3118","0x4179f6d6a91a76b5e83afe4c996c3572d4d29b33","0xebb69f1a09dde3ec5ce2929e270cda4b60510da7"]
9:09AM INF Wallet address(es) and private key(s) saved to file fileName=wallets.json
9:09AM INF Wallet(s) funded! 💸
9:09AM INF Total execution time: 10.483125ms
```

Extract from `wallets.json`.

```json
[
  {
    "Address": "0xc1A44B1e37EE1fca4C6Fd5562c730d5b8525e4C6",
    "PrivateKey": "c1a8f737fd9f78aee361bfd856f9b2e99f853a5fe5efa2131fb030acdcee762b"
  },
  {
    "Address": "0x5D8121cf716B70d3e345adB58157752304eED5C3",
    "PrivateKey": "fbc57de542cef10fdcdf99e5578ffb5508992e9a8623ea4a39ab957d77e9b849"
  },
  ...
]
```

Check the balances of the wallets.

```bash
$ cast balance 0xc1A44B1e37EE1fca4C6Fd5562c730d5b8525e4C6
50000000000000000

$ cast balance 0x5D8121cf716B70d3e345adB58157752304eED5C3
50000000000000000
...
```
