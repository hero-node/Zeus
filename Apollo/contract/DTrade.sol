pragma solidity ^0.4.24;

contract Ownable {
    address public owner;

    event OwnershipRenounced(address indexed previousOwner);
    event OwnershipTransferred(
        address indexed previousOwner,
        address indexed newOwner
    );

    constructor() public {
        owner = msg.sender;
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner");
        _;
    }

    function transferOwnership(address newOwner) public onlyOwner {
        require(newOwner != address(0), "Owner mustn't be zero address");
        emit OwnershipTransferred(owner, newOwner);
        owner = newOwner;
    }

    function renounceOwnership() public onlyOwner {
        emit OwnershipRenounced(owner);
        owner = address(0);
    }
}

contract ERC20 {
    function balanceOf(address who) public view returns (uint256);
    function transfer(address to, uint256 value) public returns (bool);
    function transferFrom(address from, address to, uint256 value) public returns (bool);
}

contract SafeMath {
    function safeAdd(uint256 a, uint256 b) internal pure returns(uint256 c) {
        c = a + b;
        assert(c >= a);
        return c;
    }
    function safeSub(uint256 a, uint256 b) internal pure returns(uint256) {
        assert(b <= a);
        return a - b;
    }
    function safeMult(uint256 a, uint256 b) internal pure returns(uint256 c) {
        if (a == 0) {
            return 0;
        }
        c = a * b;
        assert(c / a == b);
        return c;
    }
    function safeDiv(uint256 a, uint256 b) internal pure returns (uint256) {
        return a / b;
    }

    function safeAdd2(uint256 a, int256 b) internal pure returns (uint256) {
        uint256 c = 0;
        if (b >= 0) {
            c = safeAdd(a, uint256(b));
        } else {
            c = safeSub(a, uint256(-b));
        }
        return c;
    }
}

contract DTrade is SafeMath, Ownable {
    ERC20 public token;
    uint public rate;    // 1 Eth = rate Token
    int public coefficient;
    uint public limitPerRound;
    uint public blocksPerRound;
    uint public volumn;
    uint public nextRound;
    
    uint lastRound;
    uint constant weii = 1000000000000000000;

    event WithdrawByOwner(address indexed target, uint256 value, bool isToken);
    event Join(address indexed target, uint256 value, uint256 tokenValue);
    event Leave(address indexed target, uint256 tokenValue, uint256 value);

    constructor(ERC20 _token, uint256 _rate, int _coefficient) public {
        token = _token;
        rate = _rate;
        coefficient = _coefficient;
        lastRound = block.number;
        limitPerRound = safeMult(1, weii);
        blocksPerRound = 5760;
        nextRound = safeAdd(block.number, blocksPerRound);
        volumn = 0;
    }

    function() external payable {}

    function join() public payable {
        update();
        uint tokenNumber = safeMult(msg.value, rate);
        require(token.balanceOf(address(this)) >= tokenNumber, "Token balance not enough");
        token.transfer(msg.sender, tokenNumber);
        emit Join(msg.sender, msg.value, tokenNumber);
    }

    function leave(uint value) public {
        update();
        uint ethNumber = safeDiv(value, rate);
        uint newAmount = safeAdd(volumn, ethNumber);
        require(newAmount <= limitPerRound, "No more than limit");
        require(ethNumber <= address(this).balance, "Eth balance not enough");
        require(token.balanceOf(msg.sender) >= value, "Sender not enough token");
        
        if (!token.transferFrom(msg.sender, address(this), value)) {
            revert("Transfer token error");
        }
        msg.sender.transfer(ethNumber);

        volumn = newAmount;
        emit Leave(msg.sender, value, ethNumber);
    }

    // owner function
    function ownerWithdrawEth() public onlyOwner  {
        owner.transfer(address(this).balance);
        emit WithdrawByOwner(owner, owner.balance, false);
    }

    function ownerWithdrawToken() public onlyOwner {
        uint256 balance = token.balanceOf(this);
        if (token.transfer(owner, balance)) {
            emit WithdrawByOwner(owner, balance, true);
        } else {
            revert("Transfer token error");
        }
    }

    function update() public {
        if (block.number >= nextRound) {
            uint sub = safeSub(block.number, lastRound);
            lastRound = safeSub(block.number, sub % blocksPerRound);
            nextRound = safeAdd(lastRound, blocksPerRound);
            volumn = 0;
            uint rounds = safeDiv(sub, blocksPerRound);
            rate = safeAdd2(rate, int(rounds) * coefficient);
        }
    }

    function forceUpdate() private {
        lastRound = block.number;
        nextRound = safeAdd(lastRound, blocksPerRound);
    }

    // update
    function updateCoefficient(int _coefficient) public onlyOwner {
        // update round and initial rate
        update();
        coefficient = _coefficient;
        forceUpdate();
    }

    function updateLimitPerRound(uint _limitPerRound) public onlyOwner {
        limitPerRound = _limitPerRound;
    }

    function updateBlocksPerRound(uint _blocksPerRound) public onlyOwner {
        update();
        blocksPerRound = _blocksPerRound;
        forceUpdate();
    }
}