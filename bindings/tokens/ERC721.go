// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tokens

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ERC721MetaData contains all meta data concerning the ERC721 contract.
var ERC721MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentTokenId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mintBatch\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x60806040525f60065534801562000014575f80fd5b506040518060400160405280600581526020017f4d794e46540000000000000000000000000000000000000000000000000000008152506040518060400160405280600481526020017f4d4e465400000000000000000000000000000000000000000000000000000000815250815f908162000091919062000310565b508060019081620000a3919062000310565b505050620003f4565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806200012857607f821691505b6020821081036200013e576200013d620000e3565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f60088302620001a27fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8262000165565b620001ae868362000165565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f620001f8620001f2620001ec84620001c6565b620001cf565b620001c6565b9050919050565b5f819050919050565b6200021383620001d8565b6200022b6200022282620001ff565b84845462000171565b825550505050565b5f90565b6200024162000233565b6200024e81848462000208565b505050565b5b818110156200027557620002695f8262000237565b60018101905062000254565b5050565b601f821115620002c4576200028e8162000144565b620002998462000156565b81016020851015620002a9578190505b620002c1620002b88562000156565b83018262000253565b50505b505050565b5f82821c905092915050565b5f620002e65f1984600802620002c9565b1980831691505092915050565b5f620003008383620002d5565b9150826002028217905092915050565b6200031b82620000ac565b67ffffffffffffffff811115620003375762000336620000b6565b5b62000343825462000110565b6200035082828562000279565b5f60209050601f83116001811462000386575f841562000371578287015190505b6200037d8582620002f3565b865550620003ec565b601f198416620003968662000144565b5f5b82811015620003bf5784890151825560018201915060208501945060208101905062000398565b86831015620003df5784890151620003db601f891682620002d5565b8355505b6001600288020188555050505b505050505050565b6124dd80620004025f395ff3fe608060405234801561000f575f80fd5b50600436106100f2575f3560e01c806342842e0e11610095578063a22cb46511610064578063a22cb46514610280578063b88d4fde1461029c578063c87b56dd146102b8578063e985e9c5146102e8576100f2565b806342842e0e146101e65780636352211e1461020257806370a082311461023257806395d89b4114610262576100f2565b8063081812fc116100d1578063081812fc14610162578063095ea7b31461019257806323b872dd146101ae578063248b71fc146101ca576100f2565b80629a9b7b146100f657806301ffc9a71461011457806306fdde0314610144575b5f80fd5b6100fe610318565b60405161010b9190611786565b60405180910390f35b61012e60048036038101906101299190611805565b61031e565b60405161013b919061184a565b60405180910390f35b61014c6103ff565b60405161015991906118ed565b60405180910390f35b61017c60048036038101906101779190611937565b61048e565b60405161018991906119a1565b60405180910390f35b6101ac60048036038101906101a791906119e4565b6104d0565b005b6101c860048036038101906101c39190611a22565b6105e6565b005b6101e460048036038101906101df91906119e4565b610646565b005b61020060048036038101906101fb9190611a22565b610684565b005b61021c60048036038101906102179190611937565b6106a3565b60405161022991906119a1565b60405180910390f35b61024c60048036038101906102479190611a72565b610727565b6040516102599190611786565b60405180910390f35b61026a6107db565b60405161027791906118ed565b60405180910390f35b61029a60048036038101906102959190611ac7565b61086b565b005b6102b660048036038101906102b19190611c31565b610881565b005b6102d260048036038101906102cd9190611937565b6108e3565b6040516102df91906118ed565b60405180910390f35b61030260048036038101906102fd9190611cb1565b610948565b60405161030f919061184a565b60405180910390f35b60065481565b5f7f80ac58cd000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614806103e857507f5b5e139f000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916145b806103f857506103f7826109d6565b5b9050919050565b60605f805461040d90611d1c565b80601f016020809104026020016040519081016040528092919081815260200182805461043990611d1c565b80156104845780601f1061045b57610100808354040283529160200191610484565b820191905f5260205f20905b81548152906001019060200180831161046757829003601f168201915b5050505050905090565b5f61049882610a3f565b60045f8381526020019081526020015f205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b5f6104da826106a3565b90508073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361054a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161054190611dbc565b60405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff16610569610a8a565b73ffffffffffffffffffffffffffffffffffffffff161480610598575061059781610592610a8a565b610948565b5b6105d7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105ce90611e4a565b60405180910390fd5b6105e18383610a91565b505050565b6105f76105f1610a8a565b82610b47565b610636576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161062d90611ed8565b60405180910390fd5b610641838383610bdb565b505050565b5f5b8181101561067f575f60065f815461065f90611f23565b91905081905590506106718482610ec7565b508080600101915050610648565b505050565b61069e83838360405180602001604052805f815250610881565b505050565b5f806106ae83610ee4565b90505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361071e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161071590611fb4565b60405180910390fd5b80915050919050565b5f8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610796576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161078d90612042565b60405180910390fd5b60035f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b6060600180546107ea90611d1c565b80601f016020809104026020016040519081016040528092919081815260200182805461081690611d1c565b80156108615780601f1061083857610100808354040283529160200191610861565b820191905f5260205f20905b81548152906001019060200180831161084457829003601f168201915b5050505050905090565b61087d610876610a8a565b8383610f1d565b5050565b61089261088c610a8a565b83610b47565b6108d1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108c890611ed8565b60405180910390fd5b6108dd84848484611084565b50505050565b60606108ee82610a3f565b5f6108f76110e0565b90505f8151116109155760405180602001604052805f815250610940565b8061091f846110f6565b60405160200161093092919061209a565b6040516020818303038152906040525b915050919050565b5f60055f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905092915050565b5f7f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149050919050565b610a48816111c0565b610a87576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a7e90611fb4565b60405180910390fd5b50565b5f33905090565b8160045f8381526020019081526020015f205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550808273ffffffffffffffffffffffffffffffffffffffff16610b01836106a3565b73ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45050565b5f80610b52836106a3565b90508073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161480610b945750610b938185610948565b5b80610bd257508373ffffffffffffffffffffffffffffffffffffffff16610bba8461048e565b73ffffffffffffffffffffffffffffffffffffffff16145b91505092915050565b8273ffffffffffffffffffffffffffffffffffffffff16610bfb826106a3565b73ffffffffffffffffffffffffffffffffffffffff1614610c51576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c489061212d565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610cbf576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610cb6906121bb565b60405180910390fd5b610ccc8383836001611200565b8273ffffffffffffffffffffffffffffffffffffffff16610cec826106a3565b73ffffffffffffffffffffffffffffffffffffffff1614610d42576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d399061212d565b60405180910390fd5b60045f8281526020019081526020015f205f6101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600160035f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282540392505081905550600160035f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055508160025f8381526020019081526020015f205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550808273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a4610ec28383836001611206565b505050565b610ee0828260405180602001604052805f81525061120c565b5050565b5f60025f8381526020019081526020015f205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610f8b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f8290612223565b60405180910390fd5b8060055f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c3183604051611077919061184a565b60405180910390a3505050565b61108f848484610bdb565b61109b84848484611266565b6110da576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110d1906122b1565b60405180910390fd5b50505050565b606060405180602001604052805f815250905090565b60605f6001611104846113e8565b0190505f8167ffffffffffffffff81111561112257611121611b0d565b5b6040519080825280601f01601f1916602001820160405280156111545781602001600182028036833780820191505090505b5090505f82602001820190505b6001156111b5578080600190039150507f3031323334353637383961626364656600000000000000000000000000000000600a86061a8153600a85816111aa576111a96122cf565b5b0494505f8503611161575b819350505050919050565b5f8073ffffffffffffffffffffffffffffffffffffffff166111e183610ee4565b73ffffffffffffffffffffffffffffffffffffffff1614159050919050565b50505050565b50505050565b6112168383611539565b6112225f848484611266565b611261576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611258906122b1565b60405180910390fd5b505050565b5f6112868473ffffffffffffffffffffffffffffffffffffffff1661174c565b156113db578373ffffffffffffffffffffffffffffffffffffffff1663150b7a026112af610a8a565b8786866040518563ffffffff1660e01b81526004016112d1949392919061234e565b6020604051808303815f875af192505050801561130c57506040513d601f19601f8201168201806040525081019061130991906123ac565b60015b61138b573d805f811461133a576040519150601f19603f3d011682016040523d82523d5f602084013e61133f565b606091505b505f815103611383576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161137a906122b1565b60405180910390fd5b805181602001fd5b63150b7a0260e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149150506113e0565b600190505b949350505050565b5f805f90507a184f03e93ff9f4daa797ed6e38ed64bf6a1f0100000000000000008310611444577a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000838161143a576114396122cf565b5b0492506040810190505b6d04ee2d6d415b85acef81000000008310611481576d04ee2d6d415b85acef81000000008381611477576114766122cf565b5b0492506020810190505b662386f26fc1000083106114b057662386f26fc1000083816114a6576114a56122cf565b5b0492506010810190505b6305f5e10083106114d9576305f5e10083816114cf576114ce6122cf565b5b0492506008810190505b61271083106114fe5761271083816114f4576114f36122cf565b5b0492506004810190505b606483106115215760648381611517576115166122cf565b5b0492506002810190505b600a8310611530576001810190505b80915050919050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036115a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161159e90612421565b60405180910390fd5b6115b0816111c0565b156115f0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016115e790612489565b60405180910390fd5b6115fd5f83836001611200565b611606816111c0565b15611646576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161163d90612489565b60405180910390fd5b600160035f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055508160025f8381526020019081526020015f205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550808273ffffffffffffffffffffffffffffffffffffffff165f73ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a46117485f83836001611206565b5050565b5f808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b5f819050919050565b6117808161176e565b82525050565b5f6020820190506117995f830184611777565b92915050565b5f604051905090565b5f80fd5b5f80fd5b5f7fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b6117e4816117b0565b81146117ee575f80fd5b50565b5f813590506117ff816117db565b92915050565b5f6020828403121561181a576118196117a8565b5b5f611827848285016117f1565b91505092915050565b5f8115159050919050565b61184481611830565b82525050565b5f60208201905061185d5f83018461183b565b92915050565b5f81519050919050565b5f82825260208201905092915050565b5f5b8381101561189a57808201518184015260208101905061187f565b5f8484015250505050565b5f601f19601f8301169050919050565b5f6118bf82611863565b6118c9818561186d565b93506118d981856020860161187d565b6118e2816118a5565b840191505092915050565b5f6020820190508181035f83015261190581846118b5565b905092915050565b6119168161176e565b8114611920575f80fd5b50565b5f813590506119318161190d565b92915050565b5f6020828403121561194c5761194b6117a8565b5b5f61195984828501611923565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61198b82611962565b9050919050565b61199b81611981565b82525050565b5f6020820190506119b45f830184611992565b92915050565b6119c381611981565b81146119cd575f80fd5b50565b5f813590506119de816119ba565b92915050565b5f80604083850312156119fa576119f96117a8565b5b5f611a07858286016119d0565b9250506020611a1885828601611923565b9150509250929050565b5f805f60608486031215611a3957611a386117a8565b5b5f611a46868287016119d0565b9350506020611a57868287016119d0565b9250506040611a6886828701611923565b9150509250925092565b5f60208284031215611a8757611a866117a8565b5b5f611a94848285016119d0565b91505092915050565b611aa681611830565b8114611ab0575f80fd5b50565b5f81359050611ac181611a9d565b92915050565b5f8060408385031215611add57611adc6117a8565b5b5f611aea858286016119d0565b9250506020611afb85828601611ab3565b9150509250929050565b5f80fd5b5f80fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b611b43826118a5565b810181811067ffffffffffffffff82111715611b6257611b61611b0d565b5b80604052505050565b5f611b7461179f565b9050611b808282611b3a565b919050565b5f67ffffffffffffffff821115611b9f57611b9e611b0d565b5b611ba8826118a5565b9050602081019050919050565b828183375f83830152505050565b5f611bd5611bd084611b85565b611b6b565b905082815260208101848484011115611bf157611bf0611b09565b5b611bfc848285611bb5565b509392505050565b5f82601f830112611c1857611c17611b05565b5b8135611c28848260208601611bc3565b91505092915050565b5f805f8060808587031215611c4957611c486117a8565b5b5f611c56878288016119d0565b9450506020611c67878288016119d0565b9350506040611c7887828801611923565b925050606085013567ffffffffffffffff811115611c9957611c986117ac565b5b611ca587828801611c04565b91505092959194509250565b5f8060408385031215611cc757611cc66117a8565b5b5f611cd4858286016119d0565b9250506020611ce5858286016119d0565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680611d3357607f821691505b602082108103611d4657611d45611cef565b5b50919050565b7f4552433732313a20617070726f76616c20746f2063757272656e74206f776e655f8201527f7200000000000000000000000000000000000000000000000000000000000000602082015250565b5f611da660218361186d565b9150611db182611d4c565b604082019050919050565b5f6020820190508181035f830152611dd381611d9a565b9050919050565b7f4552433732313a20617070726f76652063616c6c6572206973206e6f7420746f5f8201527f6b656e206f776e6572206f7220617070726f76656420666f7220616c6c000000602082015250565b5f611e34603d8361186d565b9150611e3f82611dda565b604082019050919050565b5f6020820190508181035f830152611e6181611e28565b9050919050565b7f4552433732313a2063616c6c6572206973206e6f7420746f6b656e206f776e655f8201527f72206f7220617070726f76656400000000000000000000000000000000000000602082015250565b5f611ec2602d8361186d565b9150611ecd82611e68565b604082019050919050565b5f6020820190508181035f830152611eef81611eb6565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f611f2d8261176e565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611f5f57611f5e611ef6565b5b600182019050919050565b7f4552433732313a20696e76616c696420746f6b656e20494400000000000000005f82015250565b5f611f9e60188361186d565b9150611fa982611f6a565b602082019050919050565b5f6020820190508181035f830152611fcb81611f92565b9050919050565b7f4552433732313a2061646472657373207a65726f206973206e6f7420612076615f8201527f6c6964206f776e65720000000000000000000000000000000000000000000000602082015250565b5f61202c60298361186d565b915061203782611fd2565b604082019050919050565b5f6020820190508181035f83015261205981612020565b9050919050565b5f81905092915050565b5f61207482611863565b61207e8185612060565b935061208e81856020860161187d565b80840191505092915050565b5f6120a5828561206a565b91506120b1828461206a565b91508190509392505050565b7f4552433732313a207472616e736665722066726f6d20696e636f7272656374205f8201527f6f776e6572000000000000000000000000000000000000000000000000000000602082015250565b5f61211760258361186d565b9150612122826120bd565b604082019050919050565b5f6020820190508181035f8301526121448161210b565b9050919050565b7f4552433732313a207472616e7366657220746f20746865207a65726f206164645f8201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b5f6121a560248361186d565b91506121b08261214b565b604082019050919050565b5f6020820190508181035f8301526121d281612199565b9050919050565b7f4552433732313a20617070726f766520746f2063616c6c6572000000000000005f82015250565b5f61220d60198361186d565b9150612218826121d9565b602082019050919050565b5f6020820190508181035f83015261223a81612201565b9050919050565b7f4552433732313a207472616e7366657220746f206e6f6e2045524337323152655f8201527f63656976657220696d706c656d656e7465720000000000000000000000000000602082015250565b5f61229b60328361186d565b91506122a682612241565b604082019050919050565b5f6020820190508181035f8301526122c88161228f565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f81519050919050565b5f82825260208201905092915050565b5f612320826122fc565b61232a8185612306565b935061233a81856020860161187d565b612343816118a5565b840191505092915050565b5f6080820190506123615f830187611992565b61236e6020830186611992565b61237b6040830185611777565b818103606083015261238d8184612316565b905095945050505050565b5f815190506123a6816117db565b92915050565b5f602082840312156123c1576123c06117a8565b5b5f6123ce84828501612398565b91505092915050565b7f4552433732313a206d696e7420746f20746865207a65726f20616464726573735f82015250565b5f61240b60208361186d565b9150612416826123d7565b602082019050919050565b5f6020820190508181035f830152612438816123ff565b9050919050565b7f4552433732313a20746f6b656e20616c7265616479206d696e746564000000005f82015250565b5f612473601c8361186d565b915061247e8261243f565b602082019050919050565b5f6020820190508181035f8301526124a081612467565b905091905056fea264697066735822122031e292bab14f9ec6e42477a074007dc3477ef85cd5b1cecabc61474f178d11c064736f6c63430008170033",
}

// ERC721ABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC721MetaData.ABI instead.
var ERC721ABI = ERC721MetaData.ABI

// ERC721Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC721MetaData.Bin instead.
var ERC721Bin = ERC721MetaData.Bin

// DeployERC721 deploys a new Ethereum contract, binding an instance of ERC721 to it.
func DeployERC721(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC721, error) {
	parsed, err := ERC721MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC721Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC721{ERC721Caller: ERC721Caller{contract: contract}, ERC721Transactor: ERC721Transactor{contract: contract}, ERC721Filterer: ERC721Filterer{contract: contract}}, nil
}

// ERC721 is an auto generated Go binding around an Ethereum contract.
type ERC721 struct {
	ERC721Caller     // Read-only binding to the contract
	ERC721Transactor // Write-only binding to the contract
	ERC721Filterer   // Log filterer for contract events
}

// ERC721Caller is an auto generated read-only Go binding around an Ethereum contract.
type ERC721Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC721Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC721Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC721Session struct {
	Contract     *ERC721           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC721CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC721CallerSession struct {
	Contract *ERC721Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ERC721TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC721TransactorSession struct {
	Contract     *ERC721Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC721Raw is an auto generated low-level Go binding around an Ethereum contract.
type ERC721Raw struct {
	Contract *ERC721 // Generic contract binding to access the raw methods on
}

// ERC721CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC721CallerRaw struct {
	Contract *ERC721Caller // Generic read-only contract binding to access the raw methods on
}

// ERC721TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC721TransactorRaw struct {
	Contract *ERC721Transactor // Generic write-only contract binding to access the raw methods on
}

// NewERC721 creates a new instance of ERC721, bound to a specific deployed contract.
func NewERC721(address common.Address, backend bind.ContractBackend) (*ERC721, error) {
	contract, err := bindERC721(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC721{ERC721Caller: ERC721Caller{contract: contract}, ERC721Transactor: ERC721Transactor{contract: contract}, ERC721Filterer: ERC721Filterer{contract: contract}}, nil
}

// NewERC721Caller creates a new read-only instance of ERC721, bound to a specific deployed contract.
func NewERC721Caller(address common.Address, caller bind.ContractCaller) (*ERC721Caller, error) {
	contract, err := bindERC721(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721Caller{contract: contract}, nil
}

// NewERC721Transactor creates a new write-only instance of ERC721, bound to a specific deployed contract.
func NewERC721Transactor(address common.Address, transactor bind.ContractTransactor) (*ERC721Transactor, error) {
	contract, err := bindERC721(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721Transactor{contract: contract}, nil
}

// NewERC721Filterer creates a new log filterer instance of ERC721, bound to a specific deployed contract.
func NewERC721Filterer(address common.Address, filterer bind.ContractFilterer) (*ERC721Filterer, error) {
	contract, err := bindERC721(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC721Filterer{contract: contract}, nil
}

// bindERC721 binds a generic wrapper to an already deployed contract.
func bindERC721(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC721MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721 *ERC721Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC721.Contract.ERC721Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721 *ERC721Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721.Contract.ERC721Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721 *ERC721Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721.Contract.ERC721Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721 *ERC721CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC721.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721 *ERC721TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721 *ERC721TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ERC721 *ERC721Caller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC721.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ERC721 *ERC721Session) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC721.Contract.BalanceOf(&_ERC721.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ERC721 *ERC721CallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC721.Contract.BalanceOf(&_ERC721.CallOpts, owner)
}

// CurrentTokenId is a free data retrieval call binding the contract method 0x009a9b7b.
//
// Solidity: function currentTokenId() view returns(uint256)
func (_ERC721 *ERC721Caller) CurrentTokenId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC721.contract.Call(opts, &out, "currentTokenId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentTokenId is a free data retrieval call binding the contract method 0x009a9b7b.
//
// Solidity: function currentTokenId() view returns(uint256)
func (_ERC721 *ERC721Session) CurrentTokenId() (*big.Int, error) {
	return _ERC721.Contract.CurrentTokenId(&_ERC721.CallOpts)
}

// CurrentTokenId is a free data retrieval call binding the contract method 0x009a9b7b.
//
// Solidity: function currentTokenId() view returns(uint256)
func (_ERC721 *ERC721CallerSession) CurrentTokenId() (*big.Int, error) {
	return _ERC721.Contract.CurrentTokenId(&_ERC721.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ERC721 *ERC721Caller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ERC721.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ERC721 *ERC721Session) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ERC721.Contract.GetApproved(&_ERC721.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ERC721 *ERC721CallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ERC721.Contract.GetApproved(&_ERC721.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ERC721 *ERC721Caller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _ERC721.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ERC721 *ERC721Session) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC721.Contract.IsApprovedForAll(&_ERC721.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ERC721 *ERC721CallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC721.Contract.IsApprovedForAll(&_ERC721.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC721 *ERC721Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC721.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC721 *ERC721Session) Name() (string, error) {
	return _ERC721.Contract.Name(&_ERC721.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC721 *ERC721CallerSession) Name() (string, error) {
	return _ERC721.Contract.Name(&_ERC721.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ERC721 *ERC721Caller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ERC721.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ERC721 *ERC721Session) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ERC721.Contract.OwnerOf(&_ERC721.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ERC721 *ERC721CallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ERC721.Contract.OwnerOf(&_ERC721.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ERC721 *ERC721Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ERC721.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ERC721 *ERC721Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC721.Contract.SupportsInterface(&_ERC721.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ERC721 *ERC721CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC721.Contract.SupportsInterface(&_ERC721.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC721 *ERC721Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC721.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC721 *ERC721Session) Symbol() (string, error) {
	return _ERC721.Contract.Symbol(&_ERC721.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC721 *ERC721CallerSession) Symbol() (string, error) {
	return _ERC721.Contract.Symbol(&_ERC721.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ERC721 *ERC721Caller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _ERC721.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ERC721 *ERC721Session) TokenURI(tokenId *big.Int) (string, error) {
	return _ERC721.Contract.TokenURI(&_ERC721.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ERC721 *ERC721CallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _ERC721.Contract.TokenURI(&_ERC721.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ERC721 *ERC721Transactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ERC721 *ERC721Session) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.Contract.Approve(&_ERC721.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ERC721 *ERC721TransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.Contract.Approve(&_ERC721.TransactOpts, to, tokenId)
}

// MintBatch is a paid mutator transaction binding the contract method 0x248b71fc.
//
// Solidity: function mintBatch(address to, uint256 amount) returns()
func (_ERC721 *ERC721Transactor) MintBatch(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC721.contract.Transact(opts, "mintBatch", to, amount)
}

// MintBatch is a paid mutator transaction binding the contract method 0x248b71fc.
//
// Solidity: function mintBatch(address to, uint256 amount) returns()
func (_ERC721 *ERC721Session) MintBatch(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC721.Contract.MintBatch(&_ERC721.TransactOpts, to, amount)
}

// MintBatch is a paid mutator transaction binding the contract method 0x248b71fc.
//
// Solidity: function mintBatch(address to, uint256 amount) returns()
func (_ERC721 *ERC721TransactorSession) MintBatch(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC721.Contract.MintBatch(&_ERC721.TransactOpts, to, amount)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721 *ERC721Transactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721 *ERC721Session) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.Contract.SafeTransferFrom(&_ERC721.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721 *ERC721TransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.Contract.SafeTransferFrom(&_ERC721.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ERC721 *ERC721Transactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC721.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ERC721 *ERC721Session) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC721.Contract.SafeTransferFrom0(&_ERC721.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ERC721 *ERC721TransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC721.Contract.SafeTransferFrom0(&_ERC721.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ERC721 *ERC721Transactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ERC721 *ERC721Session) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721.Contract.SetApprovalForAll(&_ERC721.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ERC721 *ERC721TransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721.Contract.SetApprovalForAll(&_ERC721.TransactOpts, operator, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721 *ERC721Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721 *ERC721Session) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.Contract.TransferFrom(&_ERC721.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721 *ERC721TransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721.Contract.TransferFrom(&_ERC721.TransactOpts, from, to, tokenId)
}

// ERC721ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC721 contract.
type ERC721ApprovalIterator struct {
	Event *ERC721Approval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC721ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721Approval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ERC721Approval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ERC721ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721Approval represents a Approval event raised by the ERC721 contract.
type ERC721Approval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ERC721 *ERC721Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*ERC721ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ERC721ApprovalIterator{contract: _ERC721.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ERC721 *ERC721Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC721Approval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721Approval)
				if err := _ERC721.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ERC721 *ERC721Filterer) ParseApproval(log types.Log) (*ERC721Approval, error) {
	event := new(ERC721Approval)
	if err := _ERC721.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC721ApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the ERC721 contract.
type ERC721ApprovalForAllIterator struct {
	Event *ERC721ApprovalForAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC721ApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721ApprovalForAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ERC721ApprovalForAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ERC721ApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721ApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721ApprovalForAll represents a ApprovalForAll event raised by the ERC721 contract.
type ERC721ApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ERC721 *ERC721Filterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*ERC721ApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC721.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ERC721ApprovalForAllIterator{contract: _ERC721.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ERC721 *ERC721Filterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *ERC721ApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC721.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721ApprovalForAll)
				if err := _ERC721.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ERC721 *ERC721Filterer) ParseApprovalForAll(log types.Log) (*ERC721ApprovalForAll, error) {
	event := new(ERC721ApprovalForAll)
	if err := _ERC721.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC721TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC721 contract.
type ERC721TransferIterator struct {
	Event *ERC721Transfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC721TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721Transfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ERC721Transfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ERC721TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721Transfer represents a Transfer event raised by the ERC721 contract.
type ERC721Transfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ERC721 *ERC721Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*ERC721TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ERC721TransferIterator{contract: _ERC721.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ERC721 *ERC721Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC721Transfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721Transfer)
				if err := _ERC721.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ERC721 *ERC721Filterer) ParseTransfer(log types.Log) (*ERC721Transfer, error) {
	event := new(ERC721Transfer)
	if err := _ERC721.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
