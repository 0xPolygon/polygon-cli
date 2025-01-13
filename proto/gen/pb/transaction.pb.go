// If you make changes, recompile protos with `make generate`

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.29.2
// source: transaction.proto

package pb

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
	mi := &file_transaction_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[0]
	if x != nil {
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
	return file_transaction_proto_rawDescGZIP(), []int{0}
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

var File_transaction_proto protoreflect.FileDescriptor

var file_transaction_proto_rawDesc = []byte{
	0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xff, 0x02, 0x0a, 0x0b, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61,
	0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x14,
	0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6e,
	0x6f, 0x6e, 0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73,
	0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61,
	0x73, 0x68, 0x12, 0x20, 0x0a, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x12, 0x2a, 0x0a, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x66, 0x72, 0x6f, 0x6d, 0x12, 0x13, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x02, 0x74, 0x6f, 0x88, 0x01, 0x01, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x67, 0x61, 0x73, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x67, 0x61, 0x73, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x67,
	0x61, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x67, 0x61, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x76,
	0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x76, 0x12, 0x0c, 0x0a, 0x01, 0x73, 0x18, 0x0f,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x73, 0x12, 0x0c, 0x0a, 0x01, 0x72, 0x18, 0x10, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x01, 0x72, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x74, 0x6f, 0x42, 0x32, 0x5a, 0x30,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x30, 0x78, 0x50, 0x6f, 0x6c,
	0x79, 0x67, 0x6f, 0x6e, 0x2f, 0x70, 0x6f, 0x6c, 0x79, 0x67, 0x6f, 0x6e, 0x2d, 0x63, 0x6c, 0x69,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_transaction_proto_rawDescOnce sync.Once
	file_transaction_proto_rawDescData = file_transaction_proto_rawDesc
)

func file_transaction_proto_rawDescGZIP() []byte {
	file_transaction_proto_rawDescOnce.Do(func() {
		file_transaction_proto_rawDescData = protoimpl.X.CompressGZIP(file_transaction_proto_rawDescData)
	})
	return file_transaction_proto_rawDescData
}

var file_transaction_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_transaction_proto_goTypes = []any{
	(*Transaction)(nil), // 0: proto.Transaction
}
var file_transaction_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_transaction_proto_init() }
func file_transaction_proto_init() {
	if File_transaction_proto != nil {
		return
	}
	file_transaction_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_transaction_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_transaction_proto_goTypes,
		DependencyIndexes: file_transaction_proto_depIdxs,
		MessageInfos:      file_transaction_proto_msgTypes,
	}.Build()
	File_transaction_proto = out.File
	file_transaction_proto_rawDesc = nil
	file_transaction_proto_goTypes = nil
	file_transaction_proto_depIdxs = nil
}
