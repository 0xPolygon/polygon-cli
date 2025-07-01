# Bridging assets using polycli

## Requirements

The instructions in this doc assume you have `Foundry` installed.
It's also required to have a funded account in the network that originally has the funds.

## General setup

All the steps below will require you to interact with the RPC and also to provide an account.

Lets prepare some environment variables to use along the steps.

```bash
export rpc_url_l1 = <RPC_URL>
export acc_private_key_l1 = <0xPRIVATE_KEY>
export acc_addr_l1 = $(cast wallet address --private-key $acc_private_key_l1)
```
> replace the placeholders by the correct value 

---

## ERC20

To bridge ERC20 tokens we first need to have a ERC20 token deployed and funds to the account that
will bridge them to the other network.

In case you don't have the token contract deployed, the following steps will help you deploying one.

Firstly, lets use forge to create all the resources needed for the token contract deploy.

```bash
forge init erc20-example
cd erc20-example
```

Inside of the newly created directory `erc20-example` let's clean-up some unnecessary files and
folders

```bash
rm -Rf ./script
rm -Rf ./test
rm ./src/Counter.sol
```

Lets create our token contract file

```bash
touch ./src/MyToken.sol
```

Then add this code to `./src/MyToken.sol`

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyToken is ERC20, Ownable {
    constructor(
        string memory name,
        string memory symbol,
        address initialOwner
    ) ERC20(name, symbol) Ownable(initialOwner) {}

    function mint(address to, uint256 amount) public onlyOwner {
        _mint(to, amount);
    }
}
```

Next step is to use forge to install our contract dependencies

```bash
forge install OpenZeppelin/openzeppelin-contracts
```

Then build the contract

```bash
forge build
```

We are done for now, we have an ERC20 contract ready to be deployed to a network in the next steps.

### L1 to L2

Well, to bridge an ERC20 asset from L1 to L2, we first need to deploy ou ERC20 contract to the
network.

```bash
forge create --broadcast \
  --rpc-url $rpc_url_l1 \
  --private-key $acc_private_key_l1 \
    src/MyToken.sol:MyToken \
  --constructor-args "MyToken" "MTK" $acc_addr_l1
```

The response to the command above should be something like this:

```log
Deployer: 0x...
Deployed to: 0x...
Transaction hash: 0x...
```

Store the value from `Deployed to:` to be used further, this is the address of our token contract
deployed in the network. Lets add it to a env var, so we can use it later easily.

```bash
export token_address = <DEPLOYED_TO>
```

The first thing we need to do after the contract deployment, it to mint some funds to the account
that will be used to bridge the funds to the other network. For convenience and to make it simples,
we will use the same account that we used to deploy the SC.

```bash
cast send $token_address \
  "mint(address,uint256)" $acc_addr_l1 100000000000000000000 \
  --rpc-url $rpc_url_l1 \
  --private-key $acc_private_key_l1
```

We can check if the mint worked as expected, by checking the balance for the account we minted the
funds

```bash
cast call $token_address \
  "balanceOf(address)(uint256)" $acc_addr_l1 \
  --rpc-url $rpc_url_l1
```

Now we are ready to bridge the funds from L1 to L2.

Next step is to use polycli to bridge the asset from L1 to L2. Some information from the network is
required to perform this action, we need:

- `bridge-address`: this is the address of the bridge smart contract that is deployed to the network
that has the funds, in this case is the address of the bridge contract in the L1 network. if you are
running a local kurtosis-cdk environment, you can find this address in the JSON logged when the env
starts in the field `"polygonZkEVML2BridgeAddress"`.
- `destination-network`: this is the ID of the network that will receive the asset, in this case the ID
must correspond to the L2
- `value`: the amount of tokens we want to move from L1 to L2
- `token-address`: the address of the token contract, in this case the one we just created
- `destination-address`: the address of the account that will receive the tokens in the L2 network.

more details can be found running the following help command

```bash
polycli ulxly bridge asset --help
```

and here is the command that will bridge the assets from l1 to l2

```bash
polycli ulxly bridge asset \
    --bridge-address <BRIDGE_ADDRESS> \
    --rpc-url $rpc_url \
    --private-key $private_key \
    --destination-network 1 \
    --value 10000000000000000 \
    --token-address $token_address \
    --destination-address 
```
