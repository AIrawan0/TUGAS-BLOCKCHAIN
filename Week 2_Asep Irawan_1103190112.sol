pragma solidity ^0.5.0;

contract helloworld{
    uint256 totalCoin;
    function addcoin (uint256 nCoin) public {
        totalCoin += nCoin;
    }
    function viewTotalCoin() public view returns (uint){
        return totalCoin;
    }
}