// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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

// ConsensusValidatorEntrypointMetaData contains all meta data concerning the ConsensusValidatorEntrypoint contract.
var ConsensusValidatorEntrypointMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositCollateral\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isPermittedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerValidator\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPermittedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isPermitted\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unjail\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateExtraVotingPower\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraVotingPower\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"withdrawCollateral\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"maturesAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MsgDepositCollateral\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amountGwei\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MsgRegisterValidator\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"pubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"initialCollateralAmountGwei\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MsgUnjail\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MsgUpdateExtraVotingPower\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"extraVotingPowerGwei\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MsgWithdrawCollateral\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amountGwei\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"maturesAt\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PermittedCallerSet\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"isPermitted\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidParameter\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
	Bin: "0x30608052610120604052602c60c08181526100319161224860e03980516020918201205f19015f9081522060ff191690565b60a05234801561003f575f5ffd5b5061004861004d565b6100ff565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff161561009d5760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100fc5780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b60805160a0516120f06101585f395f8181610237015281816103dd01528181610474015281816106490152818161081401528181610b630152610d8e01525f8181610f7101528181610f9a01526111bb01526120f05ff3fe6080604052600436106100f6575f3560e01c806397d475d711610089578063d801f36d11610058578063d801f36d14610353578063e30c397814610366578063f09d31491461037a578063f2fde38b146103995761012d565b806397d475d7146102ad578063ad3cb1cc146102c0578063bfbac82b14610315578063c4d66de8146103345761012d565b8063715018a6116100c5578063715018a6146101d957806379ba5097146101ed5780637fcc389f146102015780638da5cb5b146102745761012d565b80631727b6f31461015f578063449ecfe6146101805780634f1ef2861461019f57806352d1902d146101b25761012d565b3661012d576040517fa038794000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517fa038794000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b34801561016a575f5ffd5b5061017e610179366004611c79565b6103b8565b005b34801561018b575f5ffd5b5061017e61019a366004611cb2565b61046d565b61017e6101ad366004611cf8565b610521565b3480156101bd575f5ffd5b506101c6610540565b6040519081526020015b60405180910390f35b3480156101e4575f5ffd5b5061017e61056e565b3480156101f8575f5ffd5b5061017e610581565b34801561020c575f5ffd5b5061026461021b366004611cb2565b73ffffffffffffffffffffffffffffffffffffffff165f9081527f0000000000000000000000000000000000000000000000000000000000000000602052604090205460ff1690565b60405190151581526020016101d0565b34801561027f575f5ffd5b50610288610601565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101d0565b61017e6102bb366004611cb2565b610642565b3480156102cb575f5ffd5b506103086040518060400160405280600581526020017f352e302e3000000000000000000000000000000000000000000000000000000081525081565b6040516101d09190611df9565b348015610320575f5ffd5b5061017e61032f366004611e4c565b61080d565b34801561033f575f5ffd5b5061017e61034e366004611cb2565b6109ce565b61017e610361366004611ea2565b610b5c565b348015610371575f5ffd5b50610288610d5f565b348015610385575f5ffd5b5061017e610394366004611f20565b610d87565b3480156103a4575f5ffd5b5061017e6103b3366004611cb2565b610e4a565b6103c0610f01565b73ffffffffffffffffffffffffffffffffffffffff82165f8181527f0000000000000000000000000000000000000000000000000000000000000000602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168515159081179091558251938452908301527f58b0246a79531a991271a8abe150f2c09805dec04338c87eca66ed423855d4c591015b60405180910390a15050565b335f9081527f0000000000000000000000000000000000000000000000000000000000000000602052604090205460ff166104d4576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405173ffffffffffffffffffffffffffffffffffffffff821681527fc2c03e4fbe86816915120cf64410100921065083a63a94c0a510d190bb79a893906020015b60405180910390a150565b610529610f59565b6105328261105d565b61053c8282611065565b5050565b5f6105496111a3565b507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc90565b610576610f01565b61057f5f611212565b565b338061058b610d5f565b73ffffffffffffffffffffffffffffffffffffffff16146105f5576040517f118cdaa700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024015b60405180910390fd5b6105fe81611212565b50565b5f807f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b5473ffffffffffffffffffffffffffffffffffffffff1692915050565b335f9081527f0000000000000000000000000000000000000000000000000000000000000000602052604090205460ff166106a9576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f3411610712576040517faa33ade000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6d73672e76616c7565000000000000000000000000000000000000000000000060448201526064016105ec565b610720633b9aca0034611f75565b15610787576040517faa33ade000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6d73672e76616c7565000000000000000000000000000000000000000000000060448201526064016105ec565b6040515f903480156108fc029183818181858288f193505050501580156107b0573d5f5f3e3d5ffd5b507f0ece843e2f22b53f4c2c333294e8ec7a628801a1a6238d0c7250651ec700f056816107e1633b9aca0034611fb5565b6040805173ffffffffffffffffffffffffffffffffffffffff9093168352602083019190915201610516565b335f9081527f0000000000000000000000000000000000000000000000000000000000000000602052604090205460ff16610874576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f83116108dd576040517faa33ade000000000000000000000000000000000000000000000000000000000815260206004820152600660248201527f616d6f756e74000000000000000000000000000000000000000000000000000060448201526064016105ec565b6108eb633b9aca0084611f75565b15610952576040517faa33ade000000000000000000000000000000000000000000000000000000000815260206004820152600660248201527f616d6f756e74000000000000000000000000000000000000000000000000000060448201526064016105ec565b7fa2f5c5687a98936a776591ca983c49720f98220a01257a1414160b7035b02e8084610982633b9aca0086611fb5565b6040805173ffffffffffffffffffffffffffffffffffffffff938416815260208101929092529185168183015265ffffffffffff8416606082015290519081900360800190a150505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f81158015610a185750825b90505f8267ffffffffffffffff166001148015610a345750303b155b905081158015610a42575080155b15610a79576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001660011785558315610ada5784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680100000000000000001785555b610ae2611262565b610aeb8661126a565b610af3611262565b8315610b545784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050565b335f9081527f0000000000000000000000000000000000000000000000000000000000000000602052604090205460ff16610bc3576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81818080601f0160208091040260200160405190810160405280939291908181526020018383808284375f92019190915250869250610c0691508390508261127b565b5f3411610c6f576040517faa33ade000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6d73672e76616c7565000000000000000000000000000000000000000000000060448201526064016105ec565b610c7d633b9aca0034611f75565b15610ce4576040517faa33ade000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6d73672e76616c7565000000000000000000000000000000000000000000000060448201526064016105ec565b6040515f903480156108fc029183818181858288f19350505050158015610d0d573d5f5f3e3d5ffd5b507f4981a1131eb40102a34e93fcdbbba169712bf298728708d07b83a8543866bc9f858585610d40633b9aca0034611fb5565b604051610d509493929190611fc8565b60405180910390a15050505050565b5f807f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00610625565b335f9081527f0000000000000000000000000000000000000000000000000000000000000000602052604090205460ff16610dee576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f38a463da3fce9952cb48e4bbe1b35ba1c3cfed5aec3043a4a0461a3168191d9c82610e1e633b9aca0084611fb5565b6040805173ffffffffffffffffffffffffffffffffffffffff9093168352602083019190915201610461565b610e52610f01565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081178255610ebb610601565b73ffffffffffffffffffffffffffffffffffffffff167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a35050565b33610f0a610601565b73ffffffffffffffffffffffffffffffffffffffff161461057f576040517f118cdaa70000000000000000000000000000000000000000000000000000000081523360048201526024016105ec565b3073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016148061102657507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1661100d7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1614155b1561057f576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105fe610f01565b8173ffffffffffffffffffffffffffffffffffffffff166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa9250505080156110ea575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682019092526110e79181019061203a565b60015b611138576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831660048201526024016105ec565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8114611194576040517faa1d49a4000000000000000000000000000000000000000000000000000000008152600481018290526024016105ec565b61119e838361134a565b505050565b3073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461057f576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0080547fffffffffffffffffffffffff000000000000000000000000000000000000000016815561053c826113ac565b61057f611441565b611272611441565b6105fe816114a8565b5f611285836114ff565b90508173ffffffffffffffffffffffffffffffffffffffff166112a78261175f565b73ffffffffffffffffffffffffffffffffffffffff161461119e576040517faa33ade000000000000000000000000000000000000000000000000000000000815260206004820152603060248201527f756e636d705075626b65793a207468652064657269766564206164647265737360448201527f206973206e6f742065787065637465640000000000000000000000000000000060648201526084016105ec565b61135382611799565b60405173ffffffffffffffffffffffffffffffffffffffff8316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a28051156113a45761119e8282611867565b61053c6118e8565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080547fffffffffffffffffffffffff0000000000000000000000000000000000000000811673ffffffffffffffffffffffffffffffffffffffff848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff1661057f576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6114b0611441565b73ffffffffffffffffffffffffffffffffffffffff81166105f5576040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081525f60048201526024016105ec565b6060815160211461156c576040517faa33ade000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f636d705075624b65792e6c656e6774680000000000000000000000000000000060448201526064016105ec565b815f8151811061157e5761157e612051565b6020910101517fff00000000000000000000000000000000000000000000000000000000000000167f0200000000000000000000000000000000000000000000000000000000000000148061162b5750815f815181106115e0576115e0612051565b6020910101517fff00000000000000000000000000000000000000000000000000000000000000167f0300000000000000000000000000000000000000000000000000000000000000145b611691576040517faa33ade000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f636d705075624b65795b305d000000000000000000000000000000000000000060448201526064016105ec565b5f825f815181106116a4576116a4612051565b0160200151602184015160f89190911c91505f6116e583838360077ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f611920565b60408051604180825260808201909252919250602082018180368337019050509350600460f81b845f8151811061171e5761171e612051565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690815f1a9053506021840191909152604183015250919050565b604080518181526060810182525f91829190602082018180368337019050509050604060218401602083015e805160209091012092915050565b8073ffffffffffffffffffffffffffffffffffffffff163b5f03611801576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016105ec565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60605f5f8473ffffffffffffffffffffffffffffffffffffffff1684604051611890919061207e565b5f60405180830381855af49150503d805f81146118c8576040519150601f19603f3d011682016040523d82523d5f602084013e6118cd565b606091505b50915091506118dd858383611a75565b925050505b92915050565b341561057f576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8560ff166002148061193657508560ff166003145b6119c2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603160248201527f456c6c697074696343757276653a696e6e76616c696420636f6d70726573736560448201527f6420454320706f696e742070726566697800000000000000000000000000000060648201526084016105ec565b5f82806119d1576119d1611f48565b83806119df576119df611f48565b8585806119ee576119ee611f48565b888a09088480611a0057611a00611f48565b8580611a0e57611a0e611f48565b898a098909089050611a37816004611a27866001612094565b611a319190611fb5565b85611b07565b90505f6002611a4960ff8a1684612094565b611a539190611f75565b15611a6757611a6282856120a7565b611a69565b815b98975050505050505050565b606082611a8a57611a8582611c0f565b611b00565b8151158015611aae575073ffffffffffffffffffffffffffffffffffffffff84163b155b15611afd576040517f9996b31500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff851660048201526024016105ec565b50805b9392505050565b5f815f03611b71576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f456c6c697074696343757276653a206d6f64756c7573206973207a65726f000060448201526064016105ec565b835f03611b7f57505f611b00565b825f03611b8e57506001611b00565b60017f80000000000000000000000000000000000000000000000000000000000000005b8015611c0657838186161515870a85848509099150836002820486161515870a85848509099150836004820486161515870a85848509099150836008820486161515870a8584850909915060109004611bb2565b50949350505050565b805115611c1f5780518082602001fd5b6040517fd6bda27500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b803573ffffffffffffffffffffffffffffffffffffffff81168114611c74575f5ffd5b919050565b5f5f60408385031215611c8a575f5ffd5b611c9383611c51565b915060208301358015158114611ca7575f5ffd5b809150509250929050565b5f60208284031215611cc2575f5ffd5b611b0082611c51565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f5f60408385031215611d09575f5ffd5b611d1283611c51565b9150602083013567ffffffffffffffff811115611d2d575f5ffd5b8301601f81018513611d3d575f5ffd5b803567ffffffffffffffff811115611d5757611d57611ccb565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff82111715611dc357611dc3611ccb565b604052818152828201602001871015611dda575f5ffd5b816020840160208301375f602083830101528093505050509250929050565b602081525f82518060208401528060208501604085015e5f6040828501015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011684010191505092915050565b5f5f5f5f60808587031215611e5f575f5ffd5b611e6885611c51565b935060208501359250611e7d60408601611c51565b9150606085013565ffffffffffff81168114611e97575f5ffd5b939692955090935050565b5f5f5f60408486031215611eb4575f5ffd5b611ebd84611c51565b9250602084013567ffffffffffffffff811115611ed8575f5ffd5b8401601f81018613611ee8575f5ffd5b803567ffffffffffffffff811115611efe575f5ffd5b866020828401011115611f0f575f5ffd5b939660209190910195509293505050565b5f5f60408385031215611f31575f5ffd5b611f3a83611c51565b946020939093013593505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f82611f8357611f83611f48565b500690565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f82611fc357611fc3611f48565b500490565b73ffffffffffffffffffffffffffffffffffffffff8516815260606020820152826060820152828460808301375f608084830101525f60807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116830101905082604083015295945050505050565b5f6020828403121561204a575f5ffd5b5051919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f82518060208501845e5f920191825250919050565b808201808211156118e2576118e2611f88565b818103818111156118e2576118e2611f8856fea26469706673582212200bf1adaee3bc5941879571b048b16f2ba22f7973d74c450d575b6e8ae3572a0964736f6c634300081c00336d69746f7369732e73746f726167652e436f6e73656e73757356616c696461746f72456e747279706f696e74",
}

// ConsensusValidatorEntrypointABI is the input ABI used to generate the binding from.
// Deprecated: Use ConsensusValidatorEntrypointMetaData.ABI instead.
var ConsensusValidatorEntrypointABI = ConsensusValidatorEntrypointMetaData.ABI

// ConsensusValidatorEntrypointBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConsensusValidatorEntrypointMetaData.Bin instead.
var ConsensusValidatorEntrypointBin = ConsensusValidatorEntrypointMetaData.Bin

// DeployConsensusValidatorEntrypoint deploys a new Ethereum contract, binding an instance of ConsensusValidatorEntrypoint to it.
func DeployConsensusValidatorEntrypoint(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ConsensusValidatorEntrypoint, error) {
	parsed, err := ConsensusValidatorEntrypointMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConsensusValidatorEntrypointBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConsensusValidatorEntrypoint{ConsensusValidatorEntrypointCaller: ConsensusValidatorEntrypointCaller{contract: contract}, ConsensusValidatorEntrypointTransactor: ConsensusValidatorEntrypointTransactor{contract: contract}, ConsensusValidatorEntrypointFilterer: ConsensusValidatorEntrypointFilterer{contract: contract}}, nil
}

// ConsensusValidatorEntrypoint is an auto generated Go binding around an Ethereum contract.
type ConsensusValidatorEntrypoint struct {
	ConsensusValidatorEntrypointCaller     // Read-only binding to the contract
	ConsensusValidatorEntrypointTransactor // Write-only binding to the contract
	ConsensusValidatorEntrypointFilterer   // Log filterer for contract events
}

// ConsensusValidatorEntrypointCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConsensusValidatorEntrypointCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusValidatorEntrypointTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConsensusValidatorEntrypointTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusValidatorEntrypointFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConsensusValidatorEntrypointFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusValidatorEntrypointSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConsensusValidatorEntrypointSession struct {
	Contract     *ConsensusValidatorEntrypoint // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ConsensusValidatorEntrypointCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConsensusValidatorEntrypointCallerSession struct {
	Contract *ConsensusValidatorEntrypointCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// ConsensusValidatorEntrypointTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConsensusValidatorEntrypointTransactorSession struct {
	Contract     *ConsensusValidatorEntrypointTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// ConsensusValidatorEntrypointRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConsensusValidatorEntrypointRaw struct {
	Contract *ConsensusValidatorEntrypoint // Generic contract binding to access the raw methods on
}

// ConsensusValidatorEntrypointCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConsensusValidatorEntrypointCallerRaw struct {
	Contract *ConsensusValidatorEntrypointCaller // Generic read-only contract binding to access the raw methods on
}

// ConsensusValidatorEntrypointTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConsensusValidatorEntrypointTransactorRaw struct {
	Contract *ConsensusValidatorEntrypointTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConsensusValidatorEntrypoint creates a new instance of ConsensusValidatorEntrypoint, bound to a specific deployed contract.
func NewConsensusValidatorEntrypoint(address common.Address, backend bind.ContractBackend) (*ConsensusValidatorEntrypoint, error) {
	contract, err := bindConsensusValidatorEntrypoint(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypoint{ConsensusValidatorEntrypointCaller: ConsensusValidatorEntrypointCaller{contract: contract}, ConsensusValidatorEntrypointTransactor: ConsensusValidatorEntrypointTransactor{contract: contract}, ConsensusValidatorEntrypointFilterer: ConsensusValidatorEntrypointFilterer{contract: contract}}, nil
}

// NewConsensusValidatorEntrypointCaller creates a new read-only instance of ConsensusValidatorEntrypoint, bound to a specific deployed contract.
func NewConsensusValidatorEntrypointCaller(address common.Address, caller bind.ContractCaller) (*ConsensusValidatorEntrypointCaller, error) {
	contract, err := bindConsensusValidatorEntrypoint(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointCaller{contract: contract}, nil
}

// NewConsensusValidatorEntrypointTransactor creates a new write-only instance of ConsensusValidatorEntrypoint, bound to a specific deployed contract.
func NewConsensusValidatorEntrypointTransactor(address common.Address, transactor bind.ContractTransactor) (*ConsensusValidatorEntrypointTransactor, error) {
	contract, err := bindConsensusValidatorEntrypoint(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointTransactor{contract: contract}, nil
}

// NewConsensusValidatorEntrypointFilterer creates a new log filterer instance of ConsensusValidatorEntrypoint, bound to a specific deployed contract.
func NewConsensusValidatorEntrypointFilterer(address common.Address, filterer bind.ContractFilterer) (*ConsensusValidatorEntrypointFilterer, error) {
	contract, err := bindConsensusValidatorEntrypoint(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointFilterer{contract: contract}, nil
}

// bindConsensusValidatorEntrypoint binds a generic wrapper to an already deployed contract.
func bindConsensusValidatorEntrypoint(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConsensusValidatorEntrypointMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConsensusValidatorEntrypoint.Contract.ConsensusValidatorEntrypointCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.ConsensusValidatorEntrypointTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.ConsensusValidatorEntrypointTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConsensusValidatorEntrypoint.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConsensusValidatorEntrypoint.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ConsensusValidatorEntrypoint.Contract.UPGRADEINTERFACEVERSION(&_ConsensusValidatorEntrypoint.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ConsensusValidatorEntrypoint.Contract.UPGRADEINTERFACEVERSION(&_ConsensusValidatorEntrypoint.CallOpts)
}

// IsPermittedCaller is a free data retrieval call binding the contract method 0x7fcc389f.
//
// Solidity: function isPermittedCaller(address caller) view returns(bool)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCaller) IsPermittedCaller(opts *bind.CallOpts, caller common.Address) (bool, error) {
	var out []interface{}
	err := _ConsensusValidatorEntrypoint.contract.Call(opts, &out, "isPermittedCaller", caller)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPermittedCaller is a free data retrieval call binding the contract method 0x7fcc389f.
//
// Solidity: function isPermittedCaller(address caller) view returns(bool)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) IsPermittedCaller(caller common.Address) (bool, error) {
	return _ConsensusValidatorEntrypoint.Contract.IsPermittedCaller(&_ConsensusValidatorEntrypoint.CallOpts, caller)
}

// IsPermittedCaller is a free data retrieval call binding the contract method 0x7fcc389f.
//
// Solidity: function isPermittedCaller(address caller) view returns(bool)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCallerSession) IsPermittedCaller(caller common.Address) (bool, error) {
	return _ConsensusValidatorEntrypoint.Contract.IsPermittedCaller(&_ConsensusValidatorEntrypoint.CallOpts, caller)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConsensusValidatorEntrypoint.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) Owner() (common.Address, error) {
	return _ConsensusValidatorEntrypoint.Contract.Owner(&_ConsensusValidatorEntrypoint.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCallerSession) Owner() (common.Address, error) {
	return _ConsensusValidatorEntrypoint.Contract.Owner(&_ConsensusValidatorEntrypoint.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConsensusValidatorEntrypoint.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) PendingOwner() (common.Address, error) {
	return _ConsensusValidatorEntrypoint.Contract.PendingOwner(&_ConsensusValidatorEntrypoint.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCallerSession) PendingOwner() (common.Address, error) {
	return _ConsensusValidatorEntrypoint.Contract.PendingOwner(&_ConsensusValidatorEntrypoint.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ConsensusValidatorEntrypoint.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) ProxiableUUID() ([32]byte, error) {
	return _ConsensusValidatorEntrypoint.Contract.ProxiableUUID(&_ConsensusValidatorEntrypoint.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCallerSession) ProxiableUUID() ([32]byte, error) {
	return _ConsensusValidatorEntrypoint.Contract.ProxiableUUID(&_ConsensusValidatorEntrypoint.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.AcceptOwnership(&_ConsensusValidatorEntrypoint.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.AcceptOwnership(&_ConsensusValidatorEntrypoint.TransactOpts)
}

// DepositCollateral is a paid mutator transaction binding the contract method 0x97d475d7.
//
// Solidity: function depositCollateral(address valAddr) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) DepositCollateral(opts *bind.TransactOpts, valAddr common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "depositCollateral", valAddr)
}

// DepositCollateral is a paid mutator transaction binding the contract method 0x97d475d7.
//
// Solidity: function depositCollateral(address valAddr) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) DepositCollateral(valAddr common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.DepositCollateral(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr)
}

// DepositCollateral is a paid mutator transaction binding the contract method 0x97d475d7.
//
// Solidity: function depositCollateral(address valAddr) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) DepositCollateral(valAddr common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.DepositCollateral(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner_) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) Initialize(opts *bind.TransactOpts, owner_ common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "initialize", owner_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner_) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) Initialize(owner_ common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Initialize(&_ConsensusValidatorEntrypoint.TransactOpts, owner_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner_) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) Initialize(owner_ common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Initialize(&_ConsensusValidatorEntrypoint.TransactOpts, owner_)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xd801f36d.
//
// Solidity: function registerValidator(address valAddr, bytes pubKey) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) RegisterValidator(opts *bind.TransactOpts, valAddr common.Address, pubKey []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "registerValidator", valAddr, pubKey)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xd801f36d.
//
// Solidity: function registerValidator(address valAddr, bytes pubKey) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) RegisterValidator(valAddr common.Address, pubKey []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.RegisterValidator(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, pubKey)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xd801f36d.
//
// Solidity: function registerValidator(address valAddr, bytes pubKey) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) RegisterValidator(valAddr common.Address, pubKey []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.RegisterValidator(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, pubKey)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) RenounceOwnership() (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.RenounceOwnership(&_ConsensusValidatorEntrypoint.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.RenounceOwnership(&_ConsensusValidatorEntrypoint.TransactOpts)
}

// SetPermittedCaller is a paid mutator transaction binding the contract method 0x1727b6f3.
//
// Solidity: function setPermittedCaller(address caller, bool isPermitted) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) SetPermittedCaller(opts *bind.TransactOpts, caller common.Address, isPermitted bool) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "setPermittedCaller", caller, isPermitted)
}

// SetPermittedCaller is a paid mutator transaction binding the contract method 0x1727b6f3.
//
// Solidity: function setPermittedCaller(address caller, bool isPermitted) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) SetPermittedCaller(caller common.Address, isPermitted bool) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.SetPermittedCaller(&_ConsensusValidatorEntrypoint.TransactOpts, caller, isPermitted)
}

// SetPermittedCaller is a paid mutator transaction binding the contract method 0x1727b6f3.
//
// Solidity: function setPermittedCaller(address caller, bool isPermitted) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) SetPermittedCaller(caller common.Address, isPermitted bool) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.SetPermittedCaller(&_ConsensusValidatorEntrypoint.TransactOpts, caller, isPermitted)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.TransferOwnership(&_ConsensusValidatorEntrypoint.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.TransferOwnership(&_ConsensusValidatorEntrypoint.TransactOpts, newOwner)
}

// Unjail is a paid mutator transaction binding the contract method 0x449ecfe6.
//
// Solidity: function unjail(address valAddr) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) Unjail(opts *bind.TransactOpts, valAddr common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "unjail", valAddr)
}

// Unjail is a paid mutator transaction binding the contract method 0x449ecfe6.
//
// Solidity: function unjail(address valAddr) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) Unjail(valAddr common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Unjail(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr)
}

// Unjail is a paid mutator transaction binding the contract method 0x449ecfe6.
//
// Solidity: function unjail(address valAddr) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) Unjail(valAddr common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Unjail(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr)
}

// UpdateExtraVotingPower is a paid mutator transaction binding the contract method 0xf09d3149.
//
// Solidity: function updateExtraVotingPower(address valAddr, uint256 extraVotingPower) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) UpdateExtraVotingPower(opts *bind.TransactOpts, valAddr common.Address, extraVotingPower *big.Int) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "updateExtraVotingPower", valAddr, extraVotingPower)
}

// UpdateExtraVotingPower is a paid mutator transaction binding the contract method 0xf09d3149.
//
// Solidity: function updateExtraVotingPower(address valAddr, uint256 extraVotingPower) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) UpdateExtraVotingPower(valAddr common.Address, extraVotingPower *big.Int) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.UpdateExtraVotingPower(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, extraVotingPower)
}

// UpdateExtraVotingPower is a paid mutator transaction binding the contract method 0xf09d3149.
//
// Solidity: function updateExtraVotingPower(address valAddr, uint256 extraVotingPower) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) UpdateExtraVotingPower(valAddr common.Address, extraVotingPower *big.Int) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.UpdateExtraVotingPower(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, extraVotingPower)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.UpgradeToAndCall(&_ConsensusValidatorEntrypoint.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.UpgradeToAndCall(&_ConsensusValidatorEntrypoint.TransactOpts, newImplementation, data)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0xbfbac82b.
//
// Solidity: function withdrawCollateral(address valAddr, uint256 amount, address receiver, uint48 maturesAt) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) WithdrawCollateral(opts *bind.TransactOpts, valAddr common.Address, amount *big.Int, receiver common.Address, maturesAt *big.Int) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "withdrawCollateral", valAddr, amount, receiver, maturesAt)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0xbfbac82b.
//
// Solidity: function withdrawCollateral(address valAddr, uint256 amount, address receiver, uint48 maturesAt) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) WithdrawCollateral(valAddr common.Address, amount *big.Int, receiver common.Address, maturesAt *big.Int) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.WithdrawCollateral(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, amount, receiver, maturesAt)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0xbfbac82b.
//
// Solidity: function withdrawCollateral(address valAddr, uint256 amount, address receiver, uint48 maturesAt) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) WithdrawCollateral(valAddr common.Address, amount *big.Int, receiver common.Address, maturesAt *big.Int) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.WithdrawCollateral(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, amount, receiver, maturesAt)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Fallback(&_ConsensusValidatorEntrypoint.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Fallback(&_ConsensusValidatorEntrypoint.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) Receive() (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Receive(&_ConsensusValidatorEntrypoint.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) Receive() (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Receive(&_ConsensusValidatorEntrypoint.TransactOpts)
}

// ConsensusValidatorEntrypointInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointInitializedIterator struct {
	Event *ConsensusValidatorEntrypointInitialized // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointInitialized)
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
		it.Event = new(ConsensusValidatorEntrypointInitialized)
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
func (it *ConsensusValidatorEntrypointInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointInitialized represents a Initialized event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterInitialized(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointInitializedIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointInitializedIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointInitialized) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointInitialized)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseInitialized(log types.Log) (*ConsensusValidatorEntrypointInitialized, error) {
	event := new(ConsensusValidatorEntrypointInitialized)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointMsgDepositCollateralIterator is returned from FilterMsgDepositCollateral and is used to iterate over the raw logs and unpacked data for MsgDepositCollateral events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgDepositCollateralIterator struct {
	Event *ConsensusValidatorEntrypointMsgDepositCollateral // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointMsgDepositCollateralIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointMsgDepositCollateral)
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
		it.Event = new(ConsensusValidatorEntrypointMsgDepositCollateral)
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
func (it *ConsensusValidatorEntrypointMsgDepositCollateralIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointMsgDepositCollateralIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointMsgDepositCollateral represents a MsgDepositCollateral event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgDepositCollateral struct {
	ValAddr    common.Address
	AmountGwei *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterMsgDepositCollateral is a free log retrieval operation binding the contract event 0x0ece843e2f22b53f4c2c333294e8ec7a628801a1a6238d0c7250651ec700f056.
//
// Solidity: event MsgDepositCollateral(address valAddr, uint256 amountGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterMsgDepositCollateral(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointMsgDepositCollateralIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "MsgDepositCollateral")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointMsgDepositCollateralIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "MsgDepositCollateral", logs: logs, sub: sub}, nil
}

// WatchMsgDepositCollateral is a free log subscription operation binding the contract event 0x0ece843e2f22b53f4c2c333294e8ec7a628801a1a6238d0c7250651ec700f056.
//
// Solidity: event MsgDepositCollateral(address valAddr, uint256 amountGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchMsgDepositCollateral(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointMsgDepositCollateral) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "MsgDepositCollateral")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointMsgDepositCollateral)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgDepositCollateral", log); err != nil {
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

// ParseMsgDepositCollateral is a log parse operation binding the contract event 0x0ece843e2f22b53f4c2c333294e8ec7a628801a1a6238d0c7250651ec700f056.
//
// Solidity: event MsgDepositCollateral(address valAddr, uint256 amountGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseMsgDepositCollateral(log types.Log) (*ConsensusValidatorEntrypointMsgDepositCollateral, error) {
	event := new(ConsensusValidatorEntrypointMsgDepositCollateral)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgDepositCollateral", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointMsgRegisterValidatorIterator is returned from FilterMsgRegisterValidator and is used to iterate over the raw logs and unpacked data for MsgRegisterValidator events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgRegisterValidatorIterator struct {
	Event *ConsensusValidatorEntrypointMsgRegisterValidator // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointMsgRegisterValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointMsgRegisterValidator)
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
		it.Event = new(ConsensusValidatorEntrypointMsgRegisterValidator)
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
func (it *ConsensusValidatorEntrypointMsgRegisterValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointMsgRegisterValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointMsgRegisterValidator represents a MsgRegisterValidator event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgRegisterValidator struct {
	ValAddr                     common.Address
	PubKey                      []byte
	InitialCollateralAmountGwei *big.Int
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterMsgRegisterValidator is a free log retrieval operation binding the contract event 0x4981a1131eb40102a34e93fcdbbba169712bf298728708d07b83a8543866bc9f.
//
// Solidity: event MsgRegisterValidator(address valAddr, bytes pubKey, uint256 initialCollateralAmountGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterMsgRegisterValidator(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointMsgRegisterValidatorIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "MsgRegisterValidator")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointMsgRegisterValidatorIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "MsgRegisterValidator", logs: logs, sub: sub}, nil
}

// WatchMsgRegisterValidator is a free log subscription operation binding the contract event 0x4981a1131eb40102a34e93fcdbbba169712bf298728708d07b83a8543866bc9f.
//
// Solidity: event MsgRegisterValidator(address valAddr, bytes pubKey, uint256 initialCollateralAmountGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchMsgRegisterValidator(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointMsgRegisterValidator) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "MsgRegisterValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointMsgRegisterValidator)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgRegisterValidator", log); err != nil {
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

// ParseMsgRegisterValidator is a log parse operation binding the contract event 0x4981a1131eb40102a34e93fcdbbba169712bf298728708d07b83a8543866bc9f.
//
// Solidity: event MsgRegisterValidator(address valAddr, bytes pubKey, uint256 initialCollateralAmountGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseMsgRegisterValidator(log types.Log) (*ConsensusValidatorEntrypointMsgRegisterValidator, error) {
	event := new(ConsensusValidatorEntrypointMsgRegisterValidator)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgRegisterValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointMsgUnjailIterator is returned from FilterMsgUnjail and is used to iterate over the raw logs and unpacked data for MsgUnjail events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgUnjailIterator struct {
	Event *ConsensusValidatorEntrypointMsgUnjail // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointMsgUnjailIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointMsgUnjail)
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
		it.Event = new(ConsensusValidatorEntrypointMsgUnjail)
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
func (it *ConsensusValidatorEntrypointMsgUnjailIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointMsgUnjailIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointMsgUnjail represents a MsgUnjail event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgUnjail struct {
	ValAddr common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMsgUnjail is a free log retrieval operation binding the contract event 0xc2c03e4fbe86816915120cf64410100921065083a63a94c0a510d190bb79a893.
//
// Solidity: event MsgUnjail(address valAddr)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterMsgUnjail(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointMsgUnjailIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "MsgUnjail")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointMsgUnjailIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "MsgUnjail", logs: logs, sub: sub}, nil
}

// WatchMsgUnjail is a free log subscription operation binding the contract event 0xc2c03e4fbe86816915120cf64410100921065083a63a94c0a510d190bb79a893.
//
// Solidity: event MsgUnjail(address valAddr)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchMsgUnjail(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointMsgUnjail) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "MsgUnjail")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointMsgUnjail)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgUnjail", log); err != nil {
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

// ParseMsgUnjail is a log parse operation binding the contract event 0xc2c03e4fbe86816915120cf64410100921065083a63a94c0a510d190bb79a893.
//
// Solidity: event MsgUnjail(address valAddr)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseMsgUnjail(log types.Log) (*ConsensusValidatorEntrypointMsgUnjail, error) {
	event := new(ConsensusValidatorEntrypointMsgUnjail)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgUnjail", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointMsgUpdateExtraVotingPowerIterator is returned from FilterMsgUpdateExtraVotingPower and is used to iterate over the raw logs and unpacked data for MsgUpdateExtraVotingPower events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgUpdateExtraVotingPowerIterator struct {
	Event *ConsensusValidatorEntrypointMsgUpdateExtraVotingPower // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointMsgUpdateExtraVotingPowerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointMsgUpdateExtraVotingPower)
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
		it.Event = new(ConsensusValidatorEntrypointMsgUpdateExtraVotingPower)
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
func (it *ConsensusValidatorEntrypointMsgUpdateExtraVotingPowerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointMsgUpdateExtraVotingPowerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointMsgUpdateExtraVotingPower represents a MsgUpdateExtraVotingPower event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgUpdateExtraVotingPower struct {
	ValAddr              common.Address
	ExtraVotingPowerGwei *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterMsgUpdateExtraVotingPower is a free log retrieval operation binding the contract event 0x38a463da3fce9952cb48e4bbe1b35ba1c3cfed5aec3043a4a0461a3168191d9c.
//
// Solidity: event MsgUpdateExtraVotingPower(address valAddr, uint256 extraVotingPowerGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterMsgUpdateExtraVotingPower(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointMsgUpdateExtraVotingPowerIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "MsgUpdateExtraVotingPower")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointMsgUpdateExtraVotingPowerIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "MsgUpdateExtraVotingPower", logs: logs, sub: sub}, nil
}

// WatchMsgUpdateExtraVotingPower is a free log subscription operation binding the contract event 0x38a463da3fce9952cb48e4bbe1b35ba1c3cfed5aec3043a4a0461a3168191d9c.
//
// Solidity: event MsgUpdateExtraVotingPower(address valAddr, uint256 extraVotingPowerGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchMsgUpdateExtraVotingPower(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointMsgUpdateExtraVotingPower) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "MsgUpdateExtraVotingPower")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointMsgUpdateExtraVotingPower)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgUpdateExtraVotingPower", log); err != nil {
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

// ParseMsgUpdateExtraVotingPower is a log parse operation binding the contract event 0x38a463da3fce9952cb48e4bbe1b35ba1c3cfed5aec3043a4a0461a3168191d9c.
//
// Solidity: event MsgUpdateExtraVotingPower(address valAddr, uint256 extraVotingPowerGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseMsgUpdateExtraVotingPower(log types.Log) (*ConsensusValidatorEntrypointMsgUpdateExtraVotingPower, error) {
	event := new(ConsensusValidatorEntrypointMsgUpdateExtraVotingPower)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgUpdateExtraVotingPower", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointMsgWithdrawCollateralIterator is returned from FilterMsgWithdrawCollateral and is used to iterate over the raw logs and unpacked data for MsgWithdrawCollateral events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgWithdrawCollateralIterator struct {
	Event *ConsensusValidatorEntrypointMsgWithdrawCollateral // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointMsgWithdrawCollateralIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointMsgWithdrawCollateral)
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
		it.Event = new(ConsensusValidatorEntrypointMsgWithdrawCollateral)
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
func (it *ConsensusValidatorEntrypointMsgWithdrawCollateralIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointMsgWithdrawCollateralIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointMsgWithdrawCollateral represents a MsgWithdrawCollateral event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgWithdrawCollateral struct {
	ValAddr    common.Address
	AmountGwei *big.Int
	Receiver   common.Address
	MaturesAt  *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterMsgWithdrawCollateral is a free log retrieval operation binding the contract event 0xa2f5c5687a98936a776591ca983c49720f98220a01257a1414160b7035b02e80.
//
// Solidity: event MsgWithdrawCollateral(address valAddr, uint256 amountGwei, address receiver, uint48 maturesAt)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterMsgWithdrawCollateral(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointMsgWithdrawCollateralIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "MsgWithdrawCollateral")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointMsgWithdrawCollateralIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "MsgWithdrawCollateral", logs: logs, sub: sub}, nil
}

// WatchMsgWithdrawCollateral is a free log subscription operation binding the contract event 0xa2f5c5687a98936a776591ca983c49720f98220a01257a1414160b7035b02e80.
//
// Solidity: event MsgWithdrawCollateral(address valAddr, uint256 amountGwei, address receiver, uint48 maturesAt)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchMsgWithdrawCollateral(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointMsgWithdrawCollateral) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "MsgWithdrawCollateral")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointMsgWithdrawCollateral)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgWithdrawCollateral", log); err != nil {
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

// ParseMsgWithdrawCollateral is a log parse operation binding the contract event 0xa2f5c5687a98936a776591ca983c49720f98220a01257a1414160b7035b02e80.
//
// Solidity: event MsgWithdrawCollateral(address valAddr, uint256 amountGwei, address receiver, uint48 maturesAt)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseMsgWithdrawCollateral(log types.Log) (*ConsensusValidatorEntrypointMsgWithdrawCollateral, error) {
	event := new(ConsensusValidatorEntrypointMsgWithdrawCollateral)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgWithdrawCollateral", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointOwnershipTransferStartedIterator struct {
	Event *ConsensusValidatorEntrypointOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointOwnershipTransferStarted)
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
		it.Event = new(ConsensusValidatorEntrypointOwnershipTransferStarted)
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
func (it *ConsensusValidatorEntrypointOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ConsensusValidatorEntrypointOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointOwnershipTransferStartedIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointOwnershipTransferStarted)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseOwnershipTransferStarted(log types.Log) (*ConsensusValidatorEntrypointOwnershipTransferStarted, error) {
	event := new(ConsensusValidatorEntrypointOwnershipTransferStarted)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointOwnershipTransferredIterator struct {
	Event *ConsensusValidatorEntrypointOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointOwnershipTransferred)
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
		it.Event = new(ConsensusValidatorEntrypointOwnershipTransferred)
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
func (it *ConsensusValidatorEntrypointOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointOwnershipTransferred represents a OwnershipTransferred event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ConsensusValidatorEntrypointOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointOwnershipTransferredIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointOwnershipTransferred)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseOwnershipTransferred(log types.Log) (*ConsensusValidatorEntrypointOwnershipTransferred, error) {
	event := new(ConsensusValidatorEntrypointOwnershipTransferred)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointPermittedCallerSetIterator is returned from FilterPermittedCallerSet and is used to iterate over the raw logs and unpacked data for PermittedCallerSet events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointPermittedCallerSetIterator struct {
	Event *ConsensusValidatorEntrypointPermittedCallerSet // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointPermittedCallerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointPermittedCallerSet)
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
		it.Event = new(ConsensusValidatorEntrypointPermittedCallerSet)
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
func (it *ConsensusValidatorEntrypointPermittedCallerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointPermittedCallerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointPermittedCallerSet represents a PermittedCallerSet event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointPermittedCallerSet struct {
	Caller      common.Address
	IsPermitted bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPermittedCallerSet is a free log retrieval operation binding the contract event 0x58b0246a79531a991271a8abe150f2c09805dec04338c87eca66ed423855d4c5.
//
// Solidity: event PermittedCallerSet(address caller, bool isPermitted)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterPermittedCallerSet(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointPermittedCallerSetIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "PermittedCallerSet")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointPermittedCallerSetIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "PermittedCallerSet", logs: logs, sub: sub}, nil
}

// WatchPermittedCallerSet is a free log subscription operation binding the contract event 0x58b0246a79531a991271a8abe150f2c09805dec04338c87eca66ed423855d4c5.
//
// Solidity: event PermittedCallerSet(address caller, bool isPermitted)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchPermittedCallerSet(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointPermittedCallerSet) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "PermittedCallerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointPermittedCallerSet)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "PermittedCallerSet", log); err != nil {
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

// ParsePermittedCallerSet is a log parse operation binding the contract event 0x58b0246a79531a991271a8abe150f2c09805dec04338c87eca66ed423855d4c5.
//
// Solidity: event PermittedCallerSet(address caller, bool isPermitted)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParsePermittedCallerSet(log types.Log) (*ConsensusValidatorEntrypointPermittedCallerSet, error) {
	event := new(ConsensusValidatorEntrypointPermittedCallerSet)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "PermittedCallerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointUpgradedIterator struct {
	Event *ConsensusValidatorEntrypointUpgraded // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointUpgraded)
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
		it.Event = new(ConsensusValidatorEntrypointUpgraded)
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
func (it *ConsensusValidatorEntrypointUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointUpgraded represents a Upgraded event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ConsensusValidatorEntrypointUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointUpgradedIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointUpgraded)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseUpgraded(log types.Log) (*ConsensusValidatorEntrypointUpgraded, error) {
	event := new(ConsensusValidatorEntrypointUpgraded)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
