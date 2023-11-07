// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.21;

import {ERC20 as OZ_ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol" ;

contract ERC20 is OZ_ERC20 {
    constructor() OZ_ERC20("MyToken", "MTK") {
        _mint(msg.sender, 1000000 * (10 ** uint256(decimals())));
    }
}
