// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

contract ConformanceTester {
    string public name;
    mapping(address => uint) public balances;
    string public constant RevertErrorMessage = "Test Revert Error Message";

    constructor(string memory _name) {
        name = _name;
    }

    function deposit(uint amount) external {
        balances[msg.sender] += amount;
    }

    function testRevert() public pure{
        revert(RevertErrorMessage);
    }
}
