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

library SafeMath {
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

contract DOTC is Ownable {
    using SafeMath for uint;

    struct Order {
        bytes32 orderID;
        uint tradeType;  // 0 buy, 1 sell
        uint amountToken;
        uint amountETH;
        address orderOwner;
        ERC20 token;
        uint status;  // 0 undeal, 1 dealed, 2 cancel
    }

    mapping(address => bool) public tokens;
    mapping(bytes32 => Order) public orders;
    mapping(address => bytes32[]) public addr2orderID;
    bytes32[] public orderIDs;

    constructor() public {

    }

    event CreateOrder(
        uint _type,
        ERC20 indexed _token,
        uint amountToken,
        uint amountETH
    );
    event Trade(bytes32 orderID);
    event WithdrawByOwner(address indexed target, uint256 value, address token);
    event CacelOrder(bytes32 orderID);

    function getUserOrderID(address user, uint index) public view returns (bytes32) {
        bytes32[] memory array = addr2orderID[user];
        return array[index];
    }

    function () public payable {}

    function createOrder(uint _type, ERC20 _token, uint amountToken, uint amountETH) public payable {
        require(tokens[_token], "Iellgal token");

        if (_type == 0) {
            // buy
            require(msg.value > 0, "Not enough eth for buy order");
            bytes32 id = keccak256(abi.encodePacked(now, _type, _token, amountToken, msg.value, msg.sender));
            Order memory order = Order({orderID: id, tradeType: _type, amountToken: amountToken, amountETH: msg.value, orderOwner: msg.sender, token: _token, status: 0});
            orderIDs.push(id);
            orders[id] = order;
            addr2orderID[msg.sender].push(id);
        } else {
            // sell
            require(_token.balanceOf(msg.sender) >= amountToken, "Not enought token for sell order");
            if (!_token.transferFrom(msg.sender, address(this), amountToken)) {
                revert("Token transfer error");
            }
            bytes32 id2 = keccak256(abi.encodePacked(now, _type, _token, amountToken, amountETH, msg.sender));
            Order memory order2 = Order({orderID: id2, tradeType: _type, amountToken: amountToken, amountETH: amountETH, orderOwner: msg.sender, token: _token, status: 0});
            orderIDs.push(id2);
            orders[id2] = order2;
            addr2orderID[msg.sender].push(id2);
        }
        
        
        emit CreateOrder(_type, _token, amountToken, msg.value);
    }

    function trade(bytes32 orderID) public payable {
        require(orders[orderID].amountToken > 0, "Order not exist");
        require(orders[orderID].status == 0, "Order has been traded");
        Order storage order = orders[orderID];
        if (order.tradeType == 0) {
            // buy order, so trader need give her to buyer and withdraw eth
            require(address(this).balance >= order.amountETH, "Contract not enough eth balance");
            require(order.token.balanceOf(msg.sender) >= order.amountToken, "Not enough token for trader to sell");

            if (!order.token.transferFrom(msg.sender, order.orderOwner, order.amountToken)) {
                revert("token transfer error");
            }
            msg.sender.transfer(order.amountETH);
            order.status = 1;
        } else {
            // sell order, so trader need give eth to seller and withdraw her
            require(msg.value >= order.amountETH, "Not enough eth for trader to buy");
            require(order.token.balanceOf(address(this)) >= order.amountToken, "Contract not enought token");
            require(address(this).balance >= order.amountETH, "Not enought eth for contract");
            if (!order.token.transfer(msg.sender, order.amountToken)) {
                revert("token transfer error");
            }
            order.orderOwner.transfer(order.amountETH);
            order.status = 1;
        }
        emit Trade(orderID);
    }

    function cancelOrder(bytes32 orderID) public {
        Order storage order = orders[orderID];
        require(msg.sender == order.orderOwner, "Can not cancel others order");
        if (order.tradeType == 0) {
            require(address(this).balance >= order.amountETH);
            msg.sender.transfer(order.amountETH);
        } else {
            require(order.token.balanceOf(address(this)) >= order.amountToken);
            if (!order.token.transfer(msg.sender, order.amountToken)) {
                revert("token transfer error");
            }
        }
        order.status = 2;

        emit CacelOrder(orderID);
    }

    function addToken(ERC20 _token) public onlyOwner {
        tokens[_token] = true;
    }

    function removeToken(ERC20 _token) public onlyOwner {
        tokens[_token] = false;
    }

    // owner function
    function ownerWithdrawEth() public onlyOwner  {
        owner.transfer(address(this).balance);
        emit WithdrawByOwner(owner, owner.balance, address(0));
    }

    function ownerWithdrawToken(ERC20 token) public onlyOwner {
        require(tokens[token], "Token not available");
        uint256 balance = token.balanceOf(this);
        if (token.transfer(owner, balance)) {
            emit WithdrawByOwner(owner, balance, token);
        } else {
            revert("Transfer token error");
        }
    }
}