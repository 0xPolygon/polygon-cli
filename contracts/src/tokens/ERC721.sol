// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

import {ERC721 as OZ_ERC721} from "@openzeppelin/token/ERC721/ERC721.sol";

contract ERC721 is OZ_ERC721 {
    uint256 public currentTokenId = 0;

    constructor() OZ_ERC721("MyNFT", "MNFT") {
        mintBatch(msg.sender, 100);
    }

    function mintBatch(address to, uint256 amount) public {
       for (uint256 i = 0; i < amount; i++) {
            uint256 newTokenId = ++currentTokenId;
            _safeMint(to, newTokenId);
        }
    }
}
