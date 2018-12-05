pragma solidity ^0.5.0;

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
        require(msg.sender == owner);
        _;
    }

    function transferOwnership(address newOwner) public onlyOwner {
        require(newOwner != address(0));
        emit OwnershipTransferred(owner, newOwner);
        owner = newOwner;
    }

    function renounceOwnership() public onlyOwner {
        emit OwnershipRenounced(owner);
        owner = address(0);
    }
}

contract HNS is Ownable {

    mapping(bytes32 => bytes32) contents;

    function setContent(bytes32 key, bytes32 content) public {
        require(key.length > 0);
        require(content.length > 0);
        require(contents[key][0] == 0);
        contents[key] = content;
    }

    function getContent(bytes32 key) public view returns (bytes32) {
        require(key.length > 0);
        return contents[key];
    }

    function modify(bytes32 key, bytes32 content) public onlyOwner {
        require(key.length > 0);
        contents[key] = content;
    }
}