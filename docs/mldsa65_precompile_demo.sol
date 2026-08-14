// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

/// @notice Minimal on-chain demo for the CoreGeth ML-DSA-65 verify precompile.
///
/// Precompile address: 0x0000000000000000000000000000000000000101
/// Input format: pk || sig || msg
/// Output: 32-byte word, 0x..01 on success, otherwise 0x00..
contract MLDSA65VerifyDemo {
    address internal constant PRECOMPILE = address(0x0101);

    function verify(bytes calldata pk, bytes calldata sig, bytes calldata msg_) external view returns (bool ok) {
        bytes memory input = bytes.concat(pk, sig, msg_);
        (bool success, bytes memory out) = PRECOMPILE.staticcall(input);
        if (!success || out.length < 32) {
            return false;
        }
        uint256 word;
        assembly {
            word := mload(add(out, 0x20))
        }
        return word == 1;
    }
}
