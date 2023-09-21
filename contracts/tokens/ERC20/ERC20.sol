// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.4;

// https://solidity-by-example.org/app/erc20/

interface IERC20 {
    function totalSupply() external view returns (uint);

    function balanceOf(address account) external view returns (uint);

    function transfer(address recipient, uint amount) external returns (bool);

    function allowance(address owner, address spender) external view returns (uint);

    function approve(address spender, uint amount) external returns (bool);

    function transferFrom(
        address sender,
        address recipient,
        uint amount
    ) external returns (bool);

    event Transfer(address indexed from, address indexed to, uint value);
    event Approval(address indexed owner, address indexed spender, uint value);
}

contract ERC20 is IERC20 {
    uint override public totalSupply;
    mapping(address => uint) public balances;
    mapping(address => mapping(address => uint)) public allowances;
    string public name;
    string public symbol;
    uint8 public decimals = 18;

    constructor(string memory _name, string memory _symbol) {
        name = _name;
        symbol = _symbol;
    }

    function balanceOf(address account) override external view returns (uint) {
        return balances[account];
    }

    function transfer(address recipient, uint amount) override external returns (bool) {
        balances[msg.sender] -= amount;
        balances[recipient] += amount;
        emit Transfer(msg.sender, recipient, amount);
        return true;
    }
    function allowance(address owner, address spender) override external view returns (uint) {
        return allowances[owner][spender];
    }

    function approve(address spender, uint amount) override external returns (bool) {
        allowances[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }

    function transferFrom(
        address sender,
        address recipient,
        uint amount
    ) override external returns (bool) {
        allowances[sender][msg.sender] -= amount;
        balances[sender] -= amount;
        balances[recipient] += amount;
        emit Transfer(sender, recipient, amount);
        return true;
    }

    function mint(uint amount) external {
        balances[msg.sender] += amount;
        totalSupply += amount;
        emit Transfer(address(0), msg.sender, amount);
    }

    function burn(uint amount) external {
        balances[msg.sender] -= amount;
        totalSupply -= amount;
        emit Transfer(msg.sender, address(0), amount);
    }
}