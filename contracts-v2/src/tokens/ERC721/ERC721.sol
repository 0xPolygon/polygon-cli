// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.21;

import {ERC721 as OZ_ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";

contract ERC721 is OZ_ERC721 {
    constructor() OZ_ERC721("MyNFT", "MNFT") {
        _batchMint(msg.sender, 10000);
    }

    function _batchMint(address to, uint256 amount) internal {
        for (uint256 i = 1; i <= amount; i++) {
            _safeMint(to, i);
        }
    }
}
