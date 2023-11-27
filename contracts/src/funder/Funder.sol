// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

/// @title Funder
/// @notice A simple smart contract for funding multiple accounts at once.
/// @dev You can use this contract to fund individual accounts or a list of accounts.
/// @dev This contract accepts direct Ether transfers.
contract Funder {
  /// @notice The amount to be sent to each account.
  uint256 public amount;

  /// @dev Initialize the contract with the specified funding amount.
  /// @param _amount The amount to be sent to each account.
  constructor(uint256 _amount) {
    require(_amount > 0, "The funding amount should be greater than zero");
    amount = _amount;
  }

  /// @notice Fund a specific account with the predefined funding amount.
  /// @param _account The address of the account to be funded.
  function fund(address _account) public {
    require(_account != address(0), "The funded address should be different thant the zero address");
    require(address(this).balance >= amount, "Insufficient contract balance for funding");

    (bool success, ) = _account.call{value: amount}("");
    require(success, "Funding failed");
  }

  /// @notice Fund multiple accounts with the predefined funding amount.
  /// @param _accounts The addresses of the accounts to be funded.
  function bulkFund(address[] calldata _accounts) external {
    require(address(this).balance >= amount * _accounts.length, "Insufficient contract balance for batch funding");
    for (uint256 i = 0; i < _accounts.length; i++) {
      fund(_accounts[i]);
    }
  }

  /// @notice Allows the contract to accept direct Ether transfers.
  receive() external payable {}
}
