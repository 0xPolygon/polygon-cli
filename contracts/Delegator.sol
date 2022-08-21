// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.4;

contract Delegator {
    function call(address contractAddress, bytes calldata packedCall) public returns(bool){
        bool success;
        bytes memory data;
        (success, data) = contractAddress.call(packedCall);
        return success;
    }
    function delegateCall(address contractAddress, bytes calldata packedCall) public returns(bool){
        bool success;
        bytes memory data;
        (success, data) = contractAddress.delegatecall(packedCall);
        return success;
    }

    function loopCall(address contractAddress, bytes calldata packedCall) public returns(bool){
        bool success;
        bytes memory data;
        while(gasleft() > 1000) {
            (success, data) = contractAddress.call(packedCall);
        }
        return success;
    }
    function loopDelegateCall(address contractAddress, bytes calldata packedCall) public returns(bool){
        bool success;
        bytes memory data;
        while(gasleft() > 1000) {
            (success, data) = contractAddress.delegatecall(packedCall);
        }
        return success;
    }
}
