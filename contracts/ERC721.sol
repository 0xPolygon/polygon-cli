// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.4;

interface IERC165 {
    function supportsInterface(bytes4 interfaceID) external view returns (bool);
}

interface IERC721 is IERC165 {

    function balanceOf(address owner) external view returns (uint balance);

    function ownerOf(uint tokenId) external view returns (address owner);

    function safeTransferFrom(
        address from,
        address to,
        uint tokenId
    ) external;

    function safeTransferFrom(
        address from,
        address to,
        uint tokenId,
        bytes calldata data
    ) external;

    function transferFrom(
        address from,
        address to,
        uint tokenId
    ) external;

    function approve(address to, uint tokenId) external;

    function getApproved(uint tokenId) external view returns (address operator);

    function setApprovalForAll(address operator, bool _approved) external;

    function isApprovedForAll(address owner, address operator)
    external
    view
    returns (bool);
}

interface IERC721Receiver {
    function onERC721Received(
        address operator,
        address from,
        uint tokenId,
        bytes calldata data
    ) external returns (bytes4);
}

contract ERC721 is IERC721 {

    event Transfer(address indexed from, address indexed to, uint indexed id);
    event Approval(address indexed owner, address indexed spender, uint indexed id);
    event ApprovalForAll(
        address indexed owner,
        address indexed operator,
        bool approved
    );

    struct AddressData {
         // Mapping owner address to token count
        uint128 balance;
        uint128 numberMinted;
    }

    // Mapping from token ID to owner address
    mapping(uint => address) internal _ownerOf;

    // Mapping owner address to token count
    mapping(address => uint) internal _balanceOf;

    // Mapping from token ID to approved address
    mapping(uint => address) internal _approvals;

    // Mapping from owner to operator approvals
    mapping(address => mapping(address => bool)) public isApprovedForAll;
    
    uint256 private currentIndex = 0;
    uint256 internal immutable maxBatchSize = 100;

    function supportsInterface(bytes4 interfaceId) external pure returns (bool) {
        return
        interfaceId == type(IERC721).interfaceId ||
        interfaceId == type(IERC165).interfaceId;
    }

    function ownerOf(uint id) external view returns (address owner) {
        owner = _ownerOf[id];
        require(owner != address(0), "token doesn't exist");
    }

    function balanceOf(address owner) external view returns (uint) {
        require(owner != address(0), "owner = zero address");
        return uint256(_balanceOf[owner]);
    }

    function setApprovalForAll(address operator, bool approved) external {
        isApprovedForAll[msg.sender][operator] = approved;
        emit ApprovalForAll(msg.sender, operator, approved);
    }

    function approve(address spender, uint id) external {
        address owner = _ownerOf[id];
        require(
            msg.sender == owner || isApprovedForAll[owner][msg.sender],
            "not authorized"
        );

        _approvals[id] = spender;

        emit Approval(owner, spender, id);
    }

    function getApproved(uint id) external view returns (address) {
        require(_ownerOf[id] != address(0), "token doesn't exist");
        return _approvals[id];
    }

    function _isApprovedOrOwner(
        address owner,
        address spender,
        uint id
    ) internal view returns (bool) {
        return (spender == owner ||
        isApprovedForAll[owner][spender] ||
        spender == _approvals[id]);
    }

    function transferFrom(
        address from,
        address to,
        uint id
    ) public {
        require(from == _ownerOf[id], "from != owner");
        require(to != address(0), "transfer to zero address");

        require(_isApprovedOrOwner(from, msg.sender, id), "not authorized");
        
        _balanceOf[from]--; 
        _balanceOf[to]++;
        _ownerOf[id] = to;

        delete _approvals[id];

        emit Transfer(from, to, id);
    }

    function safeTransferFrom(
        address from,
        address to,
        uint id
    ) external {
        transferFrom(from, to, id);

        require(
            to.code.length == 0 ||
            IERC721Receiver(to).onERC721Received(msg.sender, from, id, "") ==
            IERC721Receiver.onERC721Received.selector,
            "unsafe recipient"
        );
    }

    function safeTransferFrom(
        address from,
        address to,
        uint id,
        bytes calldata data
    ) external {
        transferFrom(from, to, id);

        require(
            to.code.length == 0 ||
            IERC721Receiver(to).onERC721Received(msg.sender, from, id, data) ==
            IERC721Receiver.onERC721Received.selector,
            "unsafe recipient"
        );
    }

    function _mintBatch(address to, uint quantity) internal {
        uint256 startTokenId = currentIndex;
        require(to != address(0), "mint to zero address");
        require(_ownerOf[startTokenId] == address(0), "already minted");
        require(quantity <= maxBatchSize, "ERC721A: quantity to mint too high");

        _balanceOf[to] += uint128(quantity);

        uint256 updatedIndex = startTokenId;
        for (uint256 i = 0; i < quantity; i++) {
            emit Transfer(address(0), to, updatedIndex);
            _ownerOf[updatedIndex] = to;
            updatedIndex++;
        }

        currentIndex = updatedIndex;
    }

    function mintBatch(address to, uint quantity) external {
        _mintBatch(to, quantity);
    }

    function _burn(uint id) internal {
        address owner = _ownerOf[id];
        require(owner != address(0), "not minted");

        _balanceOf[owner] -= 1;

        delete _ownerOf[id];
        delete _approvals[id];

        emit Transfer(owner, address(0), id);
    }
}