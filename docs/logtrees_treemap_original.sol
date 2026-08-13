// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

/// @notice Baseline contract paired with LogtreesTreemapMLDSAExample for review.
/// @dev The contract body is copied without behavioral changes from
///      https://github.com/GravityLabLLC/logtrees/blob/c4aa3460606e1684fae62c3519be40207bdd9f39/contracts/LogtreesTreemap.sol.
contract LogtreesTreemap {
    address public owner;
    uint32[30][611] public counts;
    uint256 public patronPriceWei;
    uint256 public patronTreasuryBalance;
    uint256 public patronCount;

    mapping(address => bool) public isPatron;
    mapping(address => uint256) public patronContributionTotal;
    mapping(address => uint256) public patronPurchaseCount;
    mapping(address => uint256) public patronJoinedAt;

    struct LogN {
        uint8 activityId;
        string date;
        string proofOfWorkId;
    }

    event UpdatedCount(address indexed walletAddress, uint16 indexToUpdate, uint8 positionToUpdate, uint32 additional, string species, string proofOfWorkId, string date);
    event Watered(address indexed walletAddress, uint8 kind, string species, string proofOfWorkId, string date);
    event Fed(address indexed walletAddress, uint8 kind, string species, string proofOfWorkId, string date);
    event TrashTagged(address indexed walletAddress, uint8 duration, uint8 amount, string proofOfWorkId, string date);
    event Tip(address indexed walletAddress, uint256 amount, string date);
    event LoggedN(address indexed walletAddress, uint8 activityId, string date, string proofOfWorkId);
    event PatronJoined(address indexed walletAddress, uint256 amount, string date, string profileName);
    event PatronPriceUpdated(uint256 newPriceWei);
    event PatronFundsWithdrawn(address indexed recipient, uint256 amount);

    modifier onlyOwner() {
        require(msg.sender == owner, "Only the owner can call this function");
        _;
    }

    constructor() {
        owner = msg.sender;
        patronPriceWei = 42000000000000000;
    }

    function getData() public view returns (uint32[30][611] memory) {
        return counts;
    }

    function logN(uint16 indexToUpdate, LogN[] calldata entries) public {
        for (uint256 i = 0; i < entries.length; i++) {
            LogN calldata entry = entries[i];
            require(entry.activityId >= 3 && entry.activityId < 30, "Invalid activityId");
            counts[indexToUpdate][entry.activityId] += 1;
            emit LoggedN(msg.sender, entry.activityId, entry.date, entry.proofOfWorkId);
        }
    }

    function updateCount(uint16 indexToUpdate, uint8 positionToUpdate, uint32 additional, string memory species, string memory proofOfWorkId, string memory date) public {
        require(additional > 0, "The number of additional must be greater than zero");
        counts[indexToUpdate][positionToUpdate] += additional;
        emit UpdatedCount(msg.sender, indexToUpdate, positionToUpdate, additional, species, proofOfWorkId, date);
    }

    function water(uint8 kind, string memory species, string memory proofOfWorkId, string memory date) public {
        emit Watered(msg.sender, kind, species, proofOfWorkId, date);
    }

    function feed(uint8 kind, string memory species, string memory proofOfWorkId, string memory date) public {
        emit Fed(msg.sender, kind, species, proofOfWorkId, date);
    }

    function trashTag(uint16 indexToUpdate, uint8 duration, uint8 amount, string memory proofOfWorkId, string memory date) public {
        counts[indexToUpdate][2] += 1;
        emit TrashTagged(msg.sender, duration, amount, proofOfWorkId, date);
    }

    function tip(uint256 amount, string memory date) public {
        emit Tip(msg.sender, amount, date);
    }

    function becomePatron(string memory date, string memory profileName) public payable {
        require(msg.value >= patronPriceWei, "Patron price not met");

        if (!isPatron[msg.sender]) {
            isPatron[msg.sender] = true;
            patronJoinedAt[msg.sender] = block.timestamp;
            patronCount += 1;
        }

        patronContributionTotal[msg.sender] += msg.value;
        patronPurchaseCount[msg.sender] += 1;
        patronTreasuryBalance += msg.value;

        emit PatronJoined(msg.sender, msg.value, date, profileName);
    }

    function setPatronPriceWei(uint256 newPriceWei) public onlyOwner {
        require(newPriceWei > 0, "Patron price must be greater than zero");
        patronPriceWei = newPriceWei;
        emit PatronPriceUpdated(newPriceWei);
    }

    function withdrawPatronFunds(address payable recipient, uint256 amount) public onlyOwner {
        require(recipient != address(0), "Recipient is required");
        require(amount > 0, "Amount must be greater than zero");
        require(amount <= patronTreasuryBalance, "Amount exceeds patron treasury");

        patronTreasuryBalance -= amount;
        (bool success, ) = recipient.call{value: amount}("");
        require(success, "Patron withdrawal failed");

        emit PatronFundsWithdrawn(recipient, amount);
    }

    function withdrawAllPatronFunds(address payable recipient) public onlyOwner {
        withdrawPatronFunds(recipient, patronTreasuryBalance);
    }

    function setCount(uint32 indexToUpdate, uint8 positionToUpdate, uint32 newCount) public onlyOwner {
        counts[indexToUpdate][positionToUpdate] = newCount;
    }

    function updateCounts(uint32[30][611] memory newCounts) public onlyOwner {
        counts = newCounts;
    }

    function addCounts(uint32[30][611] memory additionalCounts) public onlyOwner {
        for (uint16 i = 0; i < 611; i++) {
            counts[i][0] += additionalCounts[i][0];
            counts[i][1] += additionalCounts[i][1];
            counts[i][2] += additionalCounts[i][2];
        }
    }
}
