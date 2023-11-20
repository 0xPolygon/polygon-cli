// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import {ERC20 as OpenzeppelinERC20} from "@openzeppelin/token/ERC20/ERC20.sol";

contract Swapper is OpenzeppelinERC20 {
    constructor(string memory name, string memory symbol, uint256 amount, address recipient) OpenzeppelinERC20(name, symbol) {
        _mint(recipient, amount * 10 ** decimals());
    }
}
