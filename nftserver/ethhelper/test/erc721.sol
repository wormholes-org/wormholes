pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "./TokenURI.sol";

contract NFT721 is ERC721{
    mapping(uint256 => string) private tokenURIs;
    constructor(string memory name_, string memory symbol_) ERC721(name_,symbol_){

    }
    function mint(address to, uint256 tokenId,string memory url) external virtual {
        _mint(to,tokenId);
        _setTokenURI(tokenId,url);
    }


    function tokenURI(uint256 _tokenId) public override view returns (string memory) {
        return tokenURIs[_tokenId];
    }

    function _setTokenURI(uint256 _tokenId, string memory _tokenURI) internal {
        tokenURIs[_tokenId] = _tokenURI;
    }
}