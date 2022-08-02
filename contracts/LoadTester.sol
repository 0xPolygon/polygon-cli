// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.4;

contract LoadTester {
    uint256 callCounter;

    function getCallCounter() public view returns (uint256){
        return callCounter;
    }
    function inc() public returns (uint256){
        callCounter = callCounter + 1;
        return callCounter;
    }

    // A few op codes that aren't being tested specifically
    // 0x00 STOP - 0 Gas and doesn't do anything
    // 0x50 POP - exists in the Yul instructions, but is so low level that an empty for loop would test it
    // 0x56 JUMP
    // 0x57 JUMPI
    // 0x58 PC - Is disallowed
    // 0x5B - JUMPDEST
    // 0x60 to 0x7F - PUSHi
    // 0x80 to 0x8F - DUPi
    // 0x90 to 0x9F - SWAPi
    // 0xF0 to 0xFF - These contract level functions are bit tricky to test in isolation. It's easier to test them with contact calls
    function testADD(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0001;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := add(result, 0)
            }
        }
        return result;
    }

    function testMUL(uint x) public returns(uint) {
        inc();
        uint result = 0x0DEADBEEF0002;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := mul(result, 1)
            }
        }
        return result;
    }
    function testSUB(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0003;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := sub(result, 0)
            }
        }
        return result;
    }
    function testDIV(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0004;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := div(result, 1)
            }
        }
        return result;
    }
    function testSDIV(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0005;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := sdiv(result, 1)
            }
        }
        return result;
    }
    function testMOD(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0006;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := mod(result, hex"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
            }
        }
        return result;
    }
    function testSMOD(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0007;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := smod(result, hex"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
            }
        }
        return result;
    }
    function testADDMOD(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0008;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := addmod(result, 0, hex"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
            }
        }
        return result;
    }
    function testMULMOD(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0009;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := mulmod(result, 1, hex"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
            }
        }
        return result;
    }
    function testEXP(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF000A;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := exp(result, 1)
            }
        }
        return result;
    }
    function testSIGNEXTEND(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF000B;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := signextend(32, result)
            }
        }
        return result;
    }
    function testLT(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0010;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := lt(1, result)
            }
        }
        return result;
    }
    function testGT(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0011;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := gt(result, 1)
            }
        }
        return result;
    }
    function testSLT(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0012;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := slt(1, result)
            }
        }
        return result;
    }
    function testSGT(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0013;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := sgt(result, 1)
            }
        }
        return result;
    }
    function testEQ(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0014;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := eq(result, result)
            }
        }
        return result;
    }
    function testISZERO(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0015;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := iszero(result)
            }
        }
        return result;
    }    
    function testAND(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0016;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := and(result, result)
            }
        }
        return result;
    }
    function testOR(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0017;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := or(result, 0)
            }
        }
        return result;
    }
    function testXOR(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0018;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := xor(result, 0)
            }
        }
        return result;
    }
    function testNOT(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0019;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := not(result)
            }
            if iszero(eq(result, 0xDEADBEEF0019)) {result := not(result)}
        }
        return result;
    }
    function testBYTE(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF001A;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := byte(0, result)
            }
        }
        return result;
    }
    function testSHL(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF001B;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := shl(0, result)
            }
        }
        return result;
    }
    function testSHR(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF001C;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := shr(0, result)
            }
        }
        return result;
    }
    function testSAR(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF001D;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := sar(0, result)
            }
        }
        return result;
    }
    function testSHA3(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0020;
        assembly {
            mstore(0x00, hex"FFFFFFFF00000000000000000000000000000000000000000000000000000000")
            let out := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                out := keccak256(0x00, 0x04)
            }
            if iszero(eq(out, hex"29045a592007d0c246ef02c2223570da9522d0cf0f73282c79a1bc8f0bb2c238")) { result := 0 }
        }
        return result;
    }
    function testADDRESS(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0030;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := address()
            }
        }
        return result;
    }
    function testBALANCE(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0031;
        assembly {
            let v := 0
            let addr := address()
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := balance(addr)
            }
        }
        return result;
    }
    function testORIGIN(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0032;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := origin()
            }
        }
        return result;
    }
    function testCALLER(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0033;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := caller()
            }
        }
        return result;
    }
    function testCALLVALUE(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0034;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := callvalue()
            }
        }
        return result;
    }
    function testCALLDATALOAD(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0035;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := calldataload(0)
            }
        }
        return result;
    }
    function testCALLDATASIZE(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0036;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := calldatasize()
            }
        }
        return result;
    }
    function testCALLDATACOPY(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0037;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                // not sure how this behaves if there is no call data...??
                calldatacopy(0x00, 0x00, 32)
            }
        }
        return result;
    }
    function testCODESIZE(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0038;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := codesize()
            }
        }
        return result;
    }
    function testCODECOPY(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0039;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                codecopy(0x00,0x00,32)
            }
        }
        return result;
    }
    function testGASPRICE(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF003A;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := gasprice()
            }
        }
        return result;
    }
    function testEXTCODESIZE(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF003B;
        assembly {
            let v := 0
            let addr := address()
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := extcodesize(addr)
            }
        }
        return result;
    }
    function testRETURNDATASIZE(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF003D;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := returndatasize()
            }
        }
        return result;
    }
    function testRETURNDATACOPY(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF003E;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                returndatacopy(0x00,0x00,32)
            }
        }
        return result;
    }

    function testBLOCKHASH(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0040;
        assembly {
            let v := 0
            let n := sub(number(), 1)
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := blockhash(n)
            }
        }
        return result;
    }
    function testCOINBASE(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0041;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := coinbase()
            }
        }
        return result;
    }
    function testTIMESTAMP(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0042;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := timestamp()
            }
        }
        return result;
    }
    function testNUMBER(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0043;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := number()
            }
        }
        return result;
    }
    function testDIFFICULTY(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0044;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := difficulty()
            }
        }
        return result;
    }
    function testGASLIMIT(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0045;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := gaslimit()
            }
        }
        return result;
    }
    function testCHAINID(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0046;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := chainid()
            }
        }
        return result;
    }
    function testSELFBALANCE(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0047;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := selfbalance()
            }
        }
        return result;
    }
    function testBASEFEE(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0048;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := basefee()
            }
        }
        return result;
    }
    function testMLOAD(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0051;
        assembly {
            let v := 0
            mstore(0x00, result)
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := mload(0x00)
            }
            result := v
        }
        return result;
    }
    function testMSTORE(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0052;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                mstore(0x00, result)
            }
        }
        return result;
    }
    function testMSTORE8(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0053;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                mstore(0x00, 0xDEADBEEF)
            }
        }
        return result;
    }
    function testSLOAD(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0054;
        assembly {
            sstore(0x00, result)
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                result := sload(0x00)
            }
        }
        return result;
    }
    function testSSTORE(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0055;
        assembly {
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                sstore(0x00, result)

            }
        }
        return result;
    }
    function testMSIZE(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF0059;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := msize()
            }
        }
        return result;
    }
    function testGAS(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF005A;
        assembly {
            let v := 0
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                v := gas()
            }
        }
        return result;
    }
    function testLOG0(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF00A0;
        assembly {
            mstore(0x10, result)
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                log0(0x10, 6)
            }
        }
        return result;
    }
    function testLOG1(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF00A1;
        assembly {
            mstore(0x10, result)
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                log1(0x10, 6, i)
            }
        }
        return result;
    }
    function testLOG2(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF00A2;
        assembly {
            mstore(0x10, result)
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                log2(0x10, 6, i, 0x02)
            }
        }
        return result;
    }
    function testLOG3(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF00A3;
        assembly {
            mstore(0x10, result)
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                log3(0x10, 6, i, 0x02, 0x03)
            }
        }
        return result;
    }
    function testLOG4(uint x) public returns(uint) {
        inc();
        uint result = 0xDEADBEEF00A4;
        assembly {
            mstore(0x10, result)
            for { let i := 0 } lt(i, x) { i := add(i, 1) }
            {
                log4(0x10, 6, i, 0x02, 0x03, 0x04)
            }
        }
        return result;
    }
}
