// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.19;

import {ERC20 as OpenzeppelinERC20} from "@openzeppelin/token/ERC20/ERC20.sol";

contract Swapper is OpenzeppelinERC20 {
    constructor(string memory name, string memory symbol, address recipient) OpenzeppelinERC20(name, symbol) {
        _mint(recipient, 1000000000 * 10 ** decimals());
    }
}
