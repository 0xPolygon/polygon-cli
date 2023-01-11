// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: block.proto

package gen

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash             string  `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Nonce            string  `protobuf:"bytes,2,opt,name=nonce,proto3" json:"nonce,omitempty"`
	BlockHash        string  `protobuf:"bytes,3,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	BlockNumber      string  `protobuf:"bytes,4,opt,name=blockNumber,proto3" json:"blockNumber,omitempty"`
	TransactionIndex string  `protobuf:"bytes,5,opt,name=transactionIndex,proto3" json:"transactionIndex,omitempty"`
	From             string  `protobuf:"bytes,6,opt,name=from,proto3" json:"from,omitempty"`
	To               *string `protobuf:"bytes,7,opt,name=to,proto3,oneof" json:"to,omitempty"`
	Value            string  `protobuf:"bytes,8,opt,name=value,proto3" json:"value,omitempty"`
	GasPrice         string  `protobuf:"bytes,9,opt,name=gasPrice,proto3" json:"gasPrice,omitempty"`
	Gas              string  `protobuf:"bytes,10,opt,name=gas,proto3" json:"gas,omitempty"`
	Data             string  `protobuf:"bytes,11,opt,name=data,proto3" json:"data,omitempty"`
	Input            string  `protobuf:"bytes,12,opt,name=input,proto3" json:"input,omitempty"`
	Type             string  `protobuf:"bytes,13,opt,name=type,proto3" json:"type,omitempty"`
	V                string  `protobuf:"bytes,14,opt,name=v,proto3" json:"v,omitempty"`
	S                string  `protobuf:"bytes,15,opt,name=s,proto3" json:"s,omitempty"`
	R                string  `protobuf:"bytes,16,opt,name=r,proto3" json:"r,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{0}
}

func (x *Transaction) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *Transaction) GetNonce() string {
	if x != nil {
		return x.Nonce
	}
	return ""
}

func (x *Transaction) GetBlockHash() string {
	if x != nil {
		return x.BlockHash
	}
	return ""
}

func (x *Transaction) GetBlockNumber() string {
	if x != nil {
		return x.BlockNumber
	}
	return ""
}

func (x *Transaction) GetTransactionIndex() string {
	if x != nil {
		return x.TransactionIndex
	}
	return ""
}

func (x *Transaction) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *Transaction) GetTo() string {
	if x != nil && x.To != nil {
		return *x.To
	}
	return ""
}

func (x *Transaction) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Transaction) GetGasPrice() string {
	if x != nil {
		return x.GasPrice
	}
	return ""
}

func (x *Transaction) GetGas() string {
	if x != nil {
		return x.Gas
	}
	return ""
}

func (x *Transaction) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

func (x *Transaction) GetInput() string {
	if x != nil {
		return x.Input
	}
	return ""
}

func (x *Transaction) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Transaction) GetV() string {
	if x != nil {
		return x.V
	}
	return ""
}

func (x *Transaction) GetS() string {
	if x != nil {
		return x.S
	}
	return ""
}

func (x *Transaction) GetR() string {
	if x != nil {
		return x.R
	}
	return ""
}

type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Author           string         `protobuf:"bytes,1,opt,name=author,proto3" json:"author,omitempty"`
	Difficulty       string         `protobuf:"bytes,2,opt,name=difficulty,proto3" json:"difficulty,omitempty"`
	ExtraData        string         `protobuf:"bytes,3,opt,name=extraData,proto3" json:"extraData,omitempty"`
	GasLimit         string         `protobuf:"bytes,4,opt,name=gasLimit,proto3" json:"gasLimit,omitempty"`
	GasUsed          string         `protobuf:"bytes,5,opt,name=gasUsed,proto3" json:"gasUsed,omitempty"`
	Hash             string         `protobuf:"bytes,6,opt,name=hash,proto3" json:"hash,omitempty"`
	LogsBloom        string         `protobuf:"bytes,7,opt,name=logsBloom,proto3" json:"logsBloom,omitempty"`
	Miner            string         `protobuf:"bytes,8,opt,name=miner,proto3" json:"miner,omitempty"`
	Number           string         `protobuf:"bytes,9,opt,name=number,proto3" json:"number,omitempty"`
	ParentHash       string         `protobuf:"bytes,10,opt,name=parentHash,proto3" json:"parentHash,omitempty"`
	ReceiptsRoot     string         `protobuf:"bytes,11,opt,name=receiptsRoot,proto3" json:"receiptsRoot,omitempty"`
	Sha3Uncles       string         `protobuf:"bytes,12,opt,name=sha3Uncles,proto3" json:"sha3Uncles,omitempty"`
	Signature        string         `protobuf:"bytes,13,opt,name=signature,proto3" json:"signature,omitempty"`
	Size             string         `protobuf:"bytes,14,opt,name=size,proto3" json:"size,omitempty"`
	StateRoot        string         `protobuf:"bytes,15,opt,name=stateRoot,proto3" json:"stateRoot,omitempty"`
	Step             uint32         `protobuf:"varint,16,opt,name=step,proto3" json:"step,omitempty"`
	TotalDifficulty  string         `protobuf:"bytes,17,opt,name=totalDifficulty,proto3" json:"totalDifficulty,omitempty"`
	Timestamp        string         `protobuf:"bytes,18,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Transactions     []*Transaction `protobuf:"bytes,19,rep,name=transactions,proto3" json:"transactions,omitempty"`
	TransactionsRoot string         `protobuf:"bytes,20,opt,name=transactionsRoot,proto3" json:"transactionsRoot,omitempty"`
	Uncles           []string       `protobuf:"bytes,21,rep,name=uncles,proto3" json:"uncles,omitempty"`
	BaseFeePerGas    string         `protobuf:"bytes,22,opt,name=baseFeePerGas,proto3" json:"baseFeePerGas,omitempty"`
	MixHash          string         `protobuf:"bytes,23,opt,name=mixHash,proto3" json:"mixHash,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{1}
}

func (x *Block) GetAuthor() string {
	if x != nil {
		return x.Author
	}
	return ""
}

func (x *Block) GetDifficulty() string {
	if x != nil {
		return x.Difficulty
	}
	return ""
}

func (x *Block) GetExtraData() string {
	if x != nil {
		return x.ExtraData
	}
	return ""
}

func (x *Block) GetGasLimit() string {
	if x != nil {
		return x.GasLimit
	}
	return ""
}

func (x *Block) GetGasUsed() string {
	if x != nil {
		return x.GasUsed
	}
	return ""
}

func (x *Block) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *Block) GetLogsBloom() string {
	if x != nil {
		return x.LogsBloom
	}
	return ""
}

func (x *Block) GetMiner() string {
	if x != nil {
		return x.Miner
	}
	return ""
}

func (x *Block) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *Block) GetParentHash() string {
	if x != nil {
		return x.ParentHash
	}
	return ""
}

func (x *Block) GetReceiptsRoot() string {
	if x != nil {
		return x.ReceiptsRoot
	}
	return ""
}

func (x *Block) GetSha3Uncles() string {
	if x != nil {
		return x.Sha3Uncles
	}
	return ""
}

func (x *Block) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

func (x *Block) GetSize() string {
	if x != nil {
		return x.Size
	}
	return ""
}

func (x *Block) GetStateRoot() string {
	if x != nil {
		return x.StateRoot
	}
	return ""
}

func (x *Block) GetStep() uint32 {
	if x != nil {
		return x.Step
	}
	return 0
}

func (x *Block) GetTotalDifficulty() string {
	if x != nil {
		return x.TotalDifficulty
	}
	return ""
}

func (x *Block) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *Block) GetTransactions() []*Transaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

func (x *Block) GetTransactionsRoot() string {
	if x != nil {
		return x.TransactionsRoot
	}
	return ""
}

func (x *Block) GetUncles() []string {
	if x != nil {
		return x.Uncles
	}
	return nil
}

func (x *Block) GetBaseFeePerGas() string {
	if x != nil {
		return x.BaseFeePerGas
	}
	return ""
}

func (x *Block) GetMixHash() string {
	if x != nil {
		return x.MixHash
	}
	return ""
}

var File_block_proto protoreflect.FileDescriptor

var file_block_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x75,
	0x74, 0x69, 0x6c, 0x22, 0xff, 0x02, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x20, 0x0a, 0x0b, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x2a, 0x0a,
	0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f,
	0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x13, 0x0a,
	0x02, 0x74, 0x6f, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x02, 0x74, 0x6f, 0x88,
	0x01, 0x01, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x67, 0x61, 0x73, 0x50,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x67, 0x61, 0x73, 0x50,
	0x72, 0x69, 0x63, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x61, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x67, 0x61, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e,
	0x70, 0x75, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x76, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x01, 0x76, 0x12, 0x0c, 0x0a, 0x01, 0x73, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x73,
	0x12, 0x0c, 0x0a, 0x01, 0x72, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x72, 0x42, 0x05,
	0x0a, 0x03, 0x5f, 0x74, 0x6f, 0x22, 0xbe, 0x05, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12,
	0x16, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x69, 0x66, 0x66, 0x69,
	0x63, 0x75, 0x6c, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x69, 0x66,
	0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x78, 0x74, 0x72, 0x61,
	0x44, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x78, 0x74, 0x72,
	0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x67, 0x61, 0x73, 0x4c, 0x69, 0x6d, 0x69,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x67, 0x61, 0x73, 0x4c, 0x69, 0x6d, 0x69,
	0x74, 0x12, 0x18, 0x0a, 0x07, 0x67, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x67, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12,
	0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x67, 0x73, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x6c, 0x6f, 0x67, 0x73, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x12, 0x14, 0x0a,
	0x05, 0x6d, 0x69, 0x6e, 0x65, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x69,
	0x6e, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x70,
	0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x22, 0x0a, 0x0c, 0x72,
	0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x12,
	0x1e, 0x0a, 0x0a, 0x73, 0x68, 0x61, 0x33, 0x55, 0x6e, 0x63, 0x6c, 0x65, 0x73, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x68, 0x61, 0x33, 0x55, 0x6e, 0x63, 0x6c, 0x65, 0x73, 0x12,
	0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x0d, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x0f,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x73, 0x74, 0x65, 0x70, 0x18, 0x10, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x73,
	0x74, 0x65, 0x70, 0x12, 0x28, 0x0a, 0x0f, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x44, 0x69, 0x66, 0x66,
	0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x44, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x12, 0x1c, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x12, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x35, 0x0a, 0x0c, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x13, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x2a, 0x0a, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x75, 0x6e, 0x63, 0x6c, 0x65, 0x73, 0x18, 0x15, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06,
	0x75, 0x6e, 0x63, 0x6c, 0x65, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x62, 0x61, 0x73, 0x65, 0x46, 0x65,
	0x65, 0x50, 0x65, 0x72, 0x47, 0x61, 0x73, 0x18, 0x16, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x62,
	0x61, 0x73, 0x65, 0x46, 0x65, 0x65, 0x50, 0x65, 0x72, 0x47, 0x61, 0x73, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x69, 0x78, 0x48, 0x61, 0x73, 0x68, 0x18, 0x17, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x69, 0x78, 0x48, 0x61, 0x73, 0x68, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x61, 0x74, 0x69, 0x63, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x2f, 0x70, 0x6f, 0x6c, 0x79, 0x67, 0x6f, 0x6e, 0x2d, 0x63, 0x6c, 0x69, 0x2f, 0x75, 0x74,
	0x69, 0x6c, 0x2f, 0x67, 0x65, 0x6e, 0x3b, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_block_proto_rawDescOnce sync.Once
	file_block_proto_rawDescData = file_block_proto_rawDesc
)

func file_block_proto_rawDescGZIP() []byte {
	file_block_proto_rawDescOnce.Do(func() {
		file_block_proto_rawDescData = protoimpl.X.CompressGZIP(file_block_proto_rawDescData)
	})
	return file_block_proto_rawDescData
}

var file_block_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_block_proto_goTypes = []interface{}{
	(*Transaction)(nil), // 0: util.Transaction
	(*Block)(nil),       // 1: util.Block
}
var file_block_proto_depIdxs = []int32{
	0, // 0: util.Block.transactions:type_name -> util.Transaction
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_block_proto_init() }
func file_block_proto_init() {
	if File_block_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_block_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_block_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_block_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_block_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_block_proto_goTypes,
		DependencyIndexes: file_block_proto_depIdxs,
		MessageInfos:      file_block_proto_msgTypes,
	}.Build()
	File_block_proto = out.File
	file_block_proto_rawDesc = nil
	file_block_proto_goTypes = nil
	file_block_proto_depIdxs = nil
}
