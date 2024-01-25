// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

/// @title Funder
/// @notice A simple smart contract for funding multiple wallets at once.
/// @dev You can use this contract to fund individual wallets or a list of wallets.
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
  /// @param _address The address of the wallet to be funded.
  function fund(address _address) public {
    require(_address != address(0), "The funded address should be different than the zero address");
    require(address(this).balance >= amount, "Insufficient contract balance for funding");

    (bool success, ) = _address.call{value: amount}("");
    require(success, "Funding failed");
  }

  /// @notice Fund multiple wallets with the predefined funding amount.
  /// @param _addresses The addresses of the wallets to be funded.
  function bulkFund(address[] calldata _addresses) external {
    require(address(this).balance >= amount * _addresses.length, "Insufficient contract balance for batch funding");
    for (uint256 i = 0; i < _addresses.length; i++) {
      fund(_addresses[i]);
    }
  }

  /// @notice Allows the contract to accept direct Ether transfers.
  receive() external payable {}
}
