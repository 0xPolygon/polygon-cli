// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.21;

contract ConformanceTester {
    string public name;
    string public constant RevertErrorMessage = "Test Revert Error Message"; 

    constructor(string memory _name) {
        name = _name;
    }

    function testRevert() public pure{
        revert(RevertErrorMessage);
    }
}
