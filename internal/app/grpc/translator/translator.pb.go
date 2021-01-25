// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: translator.proto

package translator

import (
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type LangType int32

const (
	LangType_UNKOWN LangType = 0
	LangType_JP     LangType = 1
	LangType_EN     LangType = 2
)

// Enum value maps for LangType.
var (
	LangType_name = map[int32]string{
		0: "UNKOWN",
		1: "JP",
		2: "EN",
	}
	LangType_value = map[string]int32{
		"UNKOWN": 0,
		"JP":     1,
		"EN":     2,
	}
)

func (x LangType) Enum() *LangType {
	p := new(LangType)
	*p = x
	return p
}

func (x LangType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LangType) Descriptor() protoreflect.EnumDescriptor {
	return file_translator_proto_enumTypes[0].Descriptor()
}

func (LangType) Type() protoreflect.EnumType {
	return &file_translator_proto_enumTypes[0]
}

func (x LangType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LangType.Descriptor instead.
func (LangType) EnumDescriptor() ([]byte, []int) {
	return file_translator_proto_rawDescGZIP(), []int{0}
}

// The request message containing the user's name.
type TranslateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text       string   `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	SrcLang    LangType `protobuf:"varint,2,opt,name=srcLang,proto3,enum=translator.LangType" json:"srcLang,omitempty"`
	TargetLang LangType `protobuf:"varint,3,opt,name=targetLang,proto3,enum=translator.LangType" json:"targetLang,omitempty"`
}

func (x *TranslateRequest) Reset() {
	*x = TranslateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_translator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TranslateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TranslateRequest) ProtoMessage() {}

func (x *TranslateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_translator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TranslateRequest.ProtoReflect.Descriptor instead.
func (*TranslateRequest) Descriptor() ([]byte, []int) {
	return file_translator_proto_rawDescGZIP(), []int{0}
}

func (x *TranslateRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *TranslateRequest) GetSrcLang() LangType {
	if x != nil {
		return x.SrcLang
	}
	return LangType_UNKOWN
}

func (x *TranslateRequest) GetTargetLang() LangType {
	if x != nil {
		return x.TargetLang
	}
	return LangType_UNKOWN
}

// The response message containing the greetings
type TranslateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text               string                     `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	SrcLang            LangType                   `protobuf:"varint,2,opt,name=srcLang,proto3,enum=translator.LangType" json:"srcLang,omitempty"`
	TranslatedTextList map[string]*TranslatedText `protobuf:"bytes,3,rep,name=translatedTextList,proto3" json:"translatedTextList,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *TranslateResponse) Reset() {
	*x = TranslateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_translator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TranslateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TranslateResponse) ProtoMessage() {}

func (x *TranslateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_translator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TranslateResponse.ProtoReflect.Descriptor instead.
func (*TranslateResponse) Descriptor() ([]byte, []int) {
	return file_translator_proto_rawDescGZIP(), []int{1}
}

func (x *TranslateResponse) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *TranslateResponse) GetSrcLang() LangType {
	if x != nil {
		return x.SrcLang
	}
	return LangType_UNKOWN
}

func (x *TranslateResponse) GetTranslatedTextList() map[string]*TranslatedText {
	if x != nil {
		return x.TranslatedTextList
	}
	return nil
}

// https://developers.google.com/protocol-buffers/docs/proto3#maps
type TranslatedText struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text string   `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Lang LangType `protobuf:"varint,2,opt,name=lang,proto3,enum=translator.LangType" json:"lang,omitempty"`
}

func (x *TranslatedText) Reset() {
	*x = TranslatedText{}
	if protoimpl.UnsafeEnabled {
		mi := &file_translator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TranslatedText) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TranslatedText) ProtoMessage() {}

func (x *TranslatedText) ProtoReflect() protoreflect.Message {
	mi := &file_translator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TranslatedText.ProtoReflect.Descriptor instead.
func (*TranslatedText) Descriptor() ([]byte, []int) {
	return file_translator_proto_rawDescGZIP(), []int{2}
}

func (x *TranslatedText) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *TranslatedText) GetLang() LangType {
	if x != nil {
		return x.Lang
	}
	return LangType_UNKOWN
}

var File_translator_proto protoreflect.FileDescriptor

var file_translator_proto_rawDesc = []byte{
	0x0a, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x22, 0x8c,
	0x01, 0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x73, 0x72, 0x63, 0x4c, 0x61,
	0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x4c, 0x61, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x07,
	0x73, 0x72, 0x63, 0x4c, 0x61, 0x6e, 0x67, 0x12, 0x34, 0x0a, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x4c, 0x61, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x4c, 0x61, 0x6e, 0x67, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4c, 0x61, 0x6e, 0x67, 0x22, 0xa1, 0x02,
	0x0a, 0x11, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x73, 0x72, 0x63, 0x4c, 0x61,
	0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x4c, 0x61, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x07,
	0x73, 0x72, 0x63, 0x4c, 0x61, 0x6e, 0x67, 0x12, 0x65, 0x0a, 0x12, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x6c, 0x61, 0x74, 0x65, 0x64, 0x54, 0x65, 0x78, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f, 0x72,
	0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x54, 0x65, 0x78,
	0x74, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x12, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x54, 0x65, 0x78, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x1a, 0x61,
	0x0a, 0x17, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x54, 0x65, 0x78, 0x74,
	0x4c, 0x69, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x30, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74,
	0x65, 0x64, 0x54, 0x65, 0x78, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0x4e, 0x0a, 0x0e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x54,
	0x65, 0x78, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x28, 0x0a, 0x04, 0x6c, 0x61, 0x6e, 0x67, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74,
	0x6f, 0x72, 0x2e, 0x4c, 0x61, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x6c, 0x61, 0x6e,
	0x67, 0x2a, 0x26, 0x0a, 0x08, 0x4c, 0x61, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a,
	0x06, 0x55, 0x4e, 0x4b, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x06, 0x0a, 0x02, 0x4a, 0x50, 0x10,
	0x01, 0x12, 0x06, 0x0a, 0x02, 0x45, 0x4e, 0x10, 0x02, 0x32, 0x58, 0x0a, 0x0a, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x4a, 0x0a, 0x09, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x6c, 0x61, 0x74, 0x65, 0x12, 0x1c, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f,
	0x72, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x22, 0x5a, 0x20, 0x61, 0x70, 0x69, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_translator_proto_rawDescOnce sync.Once
	file_translator_proto_rawDescData = file_translator_proto_rawDesc
)

func file_translator_proto_rawDescGZIP() []byte {
	file_translator_proto_rawDescOnce.Do(func() {
		file_translator_proto_rawDescData = protoimpl.X.CompressGZIP(file_translator_proto_rawDescData)
	})
	return file_translator_proto_rawDescData
}

var file_translator_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_translator_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_translator_proto_goTypes = []interface{}{
	(LangType)(0),             // 0: translator.LangType
	(*TranslateRequest)(nil),  // 1: translator.TranslateRequest
	(*TranslateResponse)(nil), // 2: translator.TranslateResponse
	(*TranslatedText)(nil),    // 3: translator.TranslatedText
	nil,                       // 4: translator.TranslateResponse.TranslatedTextListEntry
}
var file_translator_proto_depIdxs = []int32{
	0, // 0: translator.TranslateRequest.srcLang:type_name -> translator.LangType
	0, // 1: translator.TranslateRequest.targetLang:type_name -> translator.LangType
	0, // 2: translator.TranslateResponse.srcLang:type_name -> translator.LangType
	4, // 3: translator.TranslateResponse.translatedTextList:type_name -> translator.TranslateResponse.TranslatedTextListEntry
	0, // 4: translator.TranslatedText.lang:type_name -> translator.LangType
	3, // 5: translator.TranslateResponse.TranslatedTextListEntry.value:type_name -> translator.TranslatedText
	1, // 6: translator.Translator.Translate:input_type -> translator.TranslateRequest
	2, // 7: translator.Translator.Translate:output_type -> translator.TranslateResponse
	7, // [7:8] is the sub-list for method output_type
	6, // [6:7] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_translator_proto_init() }
func file_translator_proto_init() {
	if File_translator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_translator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TranslateRequest); i {
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
		file_translator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TranslateResponse); i {
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
		file_translator_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TranslatedText); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_translator_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_translator_proto_goTypes,
		DependencyIndexes: file_translator_proto_depIdxs,
		EnumInfos:         file_translator_proto_enumTypes,
		MessageInfos:      file_translator_proto_msgTypes,
	}.Build()
	File_translator_proto = out.File
	file_translator_proto_rawDesc = nil
	file_translator_proto_goTypes = nil
	file_translator_proto_depIdxs = nil
}
