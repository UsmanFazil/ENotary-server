pragma solidity ^0.5.0;

contract ENotary {
    
    // variables
    address constant empt= 0x0000000000000000000000000000000000000000;
    uint public constant fee=100000000000000000;
    address payable public owner;
    uint GAS_LIMIT = 4000000;
    
    //mapping of contract hashdata on users wallet address
    mapping (bytes32 => address) contractHash;
    
    //event (blockchain logger)  
    event contractHashSaved (bytes32,address);
    event Cashed (uint);

    // modifier for owner functions
    modifier onlyOwner() {
        require(msg.sender == owner);
        _;
    }
    
    //constructor sets the address of owner
    constructor (address payable addr) payable public {
        owner = addr;
    }
    
    // function to save the contract hash (reqires atleast 0.1 ethers to call)
    function saveHash (bytes32 hash)payable public returns(bool){
        require (msg.value == fee);
        
        // checks the contract is already saved by the creator.
        require (contractHash[hash] != msg.sender);
        
        contractHash[hash] = msg.sender;
        emit contractHashSaved(hash,msg.sender);
        
        return true;
        
    }
    
    //function to get the address of contract's owner
    function getContractOwner(bytes32 hash)view public returns(address){
       
        if (contractHash[hash] == empt){
            return empt;
        }
        
        else{
        return contractHash[hash];
        }
        
    }
    
    //function to check weather the contract is saved on blockchain or not
    function hashExists (bytes32 hash) view public returns(bool){
       
        if (contractHash[hash] == empt){
            return false;
        }
        else {
            return true;
        }
    }
    
    // function to withdraw amount (only for owner)
    function etherWithdraw(uint amount)payable onlyOwner public {
        require(address(this).balance >= amount);
        owner.transfer(amount);
    }
    
    // function to check total balance of the contract(only owner can call this function)
    function balancer()public view onlyOwner returns(uint){
        return address(this).balance;
    }
    
    //fall back function
    function ()external payable {

    }

}
      // require (address(this).balance >= amount);
