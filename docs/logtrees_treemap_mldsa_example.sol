// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

/// @notice Documentation-only adaptation of LogtreesTreemap.logN using ML-DSA-65.
/// @dev The Treemap state, LogN entries, validation, and event are retained from
///      https://github.com/GravityLabLLC/logtrees/blob/main/contracts/LogtreesTreemap.sol.
///      The added code is limited to key registration, action binding, replay
///      protection, and verification through the draft precompile at 0x0101.
contract LogtreesTreemapMLDSAExample {
    address internal constant MLDSA65_PRECOMPILE = address(0x0101);
    uint256 internal constant MLDSA65_PUBLIC_KEY_LENGTH = 1952;

    uint32[30][611] public counts;

    mapping(address => bytes32) public mldsaKeyHash;
    mapping(address => uint256) public mldsaNonce;

    struct LogN {
        uint8 activityId;
        string date;
        string proofOfWorkId;
    }

    event LoggedN(address indexed walletAddress, uint8 activityId, string date, string proofOfWorkId);
    event MLDSAKeyRegistered(address indexed account, bytes32 indexed keyHash);

    /// @notice Registers one immutable ML-DSA-65 key hash for the caller.
    /// @dev Registration must happen before the ECDSA account is compromised.
    ///      A production rotation or recovery path must itself require PQ authorization.
    function registerMLDSAKey(bytes calldata publicKey) external {
        require(publicKey.length == MLDSA65_PUBLIC_KEY_LENGTH, "Invalid ML-DSA-65 public key length");
        require(mldsaKeyHash[msg.sender] == bytes32(0), "ML-DSA key already registered");

        bytes32 keyHash = keccak256(publicKey);
        mldsaKeyHash[msg.sender] = keyHash;
        emit MLDSAKeyRegistered(msg.sender, keyHash);
    }

    /// @notice Applies the original Logtrees logN update after ML-DSA authorization.
    /// @dev Any EOA or relayer may submit the outer transaction. The ML-DSA signature
    ///      authorizes `account` and the exact action rather than trusting msg.sender.
    function logN(
        address account,
        uint16 indexToUpdate,
        LogN[] calldata entries,
        uint256 deadline,
        bytes calldata publicKey,
        bytes calldata signature
    ) external {
        require(account != address(0), "Account is required");
        require(block.timestamp <= deadline, "Authorization expired");

        uint256 nonce = mldsaNonce[account];
        bytes32 digest = keccak256(
            abi.encode(
                block.chainid,
                address(this),
                account,
                this.logN.selector,
                indexToUpdate,
                keccak256(abi.encode(entries)),
                nonce,
                deadline
            )
        );

        require(_verifyMLDSA65(account, publicKey, signature, digest), "Invalid ML-DSA authorization");
        mldsaNonce[account] = nonce + 1;

        for (uint256 i = 0; i < entries.length; i++) {
            LogN calldata entry = entries[i];
            require(entry.activityId >= 3 && entry.activityId < 30, "Invalid activityId");
            counts[indexToUpdate][entry.activityId] += 1;
            emit LoggedN(account, entry.activityId, entry.date, entry.proofOfWorkId);
        }
    }

    function _verifyMLDSA65(
        address account,
        bytes calldata publicKey,
        bytes calldata signature,
        bytes32 digest
    ) internal view returns (bool) {
        if (keccak256(publicKey) != mldsaKeyHash[account]) {
            return false;
        }

        bytes memory input = abi.encodePacked(publicKey, signature, digest);
        (bool success, bytes memory out) = MLDSA65_PRECOMPILE.staticcall(input);
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
