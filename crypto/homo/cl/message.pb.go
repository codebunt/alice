// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/getamis/alice/crypto/homo/cl/message.proto

package cl

import (
	fmt "fmt"
	binaryquadraticform "github.com/getamis/alice/crypto/binaryquadraticform"
	ecpointgrouplaw "github.com/getamis/alice/crypto/ecpointgrouplaw"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type PubKeyMessage struct {
	P                    []byte                      `protobuf:"bytes,1,opt,name=p,proto3" json:"p,omitempty"`
	A                    []byte                      `protobuf:"bytes,2,opt,name=a,proto3" json:"a,omitempty"`
	Q                    []byte                      `protobuf:"bytes,3,opt,name=q,proto3" json:"q,omitempty"`
	G                    *binaryquadraticform.BQForm `protobuf:"bytes,4,opt,name=g,proto3" json:"g,omitempty"`
	F                    *binaryquadraticform.BQForm `protobuf:"bytes,5,opt,name=f,proto3" json:"f,omitempty"`
	H                    *binaryquadraticform.BQForm `protobuf:"bytes,6,opt,name=h,proto3" json:"h,omitempty"`
	C                    []byte                      `protobuf:"bytes,7,opt,name=c,proto3" json:"c,omitempty"`
	D                    uint32                      `protobuf:"varint,8,opt,name=d,proto3" json:"d,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *PubKeyMessage) Reset()         { *m = PubKeyMessage{} }
func (m *PubKeyMessage) String() string { return proto.CompactTextString(m) }
func (*PubKeyMessage) ProtoMessage()    {}
func (*PubKeyMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_69909e00b9236d45, []int{0}
}

func (m *PubKeyMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PubKeyMessage.Unmarshal(m, b)
}
func (m *PubKeyMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PubKeyMessage.Marshal(b, m, deterministic)
}
func (m *PubKeyMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PubKeyMessage.Merge(m, src)
}
func (m *PubKeyMessage) XXX_Size() int {
	return xxx_messageInfo_PubKeyMessage.Size(m)
}
func (m *PubKeyMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_PubKeyMessage.DiscardUnknown(m)
}

var xxx_messageInfo_PubKeyMessage proto.InternalMessageInfo

func (m *PubKeyMessage) GetP() []byte {
	if m != nil {
		return m.P
	}
	return nil
}

func (m *PubKeyMessage) GetA() []byte {
	if m != nil {
		return m.A
	}
	return nil
}

func (m *PubKeyMessage) GetQ() []byte {
	if m != nil {
		return m.Q
	}
	return nil
}

func (m *PubKeyMessage) GetG() *binaryquadraticform.BQForm {
	if m != nil {
		return m.G
	}
	return nil
}

func (m *PubKeyMessage) GetF() *binaryquadraticform.BQForm {
	if m != nil {
		return m.F
	}
	return nil
}

func (m *PubKeyMessage) GetH() *binaryquadraticform.BQForm {
	if m != nil {
		return m.H
	}
	return nil
}

func (m *PubKeyMessage) GetC() []byte {
	if m != nil {
		return m.C
	}
	return nil
}

func (m *PubKeyMessage) GetD() uint32 {
	if m != nil {
		return m.D
	}
	return 0
}

type EncryptedMessage struct {
	M1                   *binaryquadraticform.BQForm `protobuf:"bytes,1,opt,name=m1,proto3" json:"m1,omitempty"`
	M2                   *binaryquadraticform.BQForm `protobuf:"bytes,2,opt,name=m2,proto3" json:"m2,omitempty"`
	Proof                *ProofMessage               `protobuf:"bytes,3,opt,name=proof,proto3" json:"proof,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *EncryptedMessage) Reset()         { *m = EncryptedMessage{} }
func (m *EncryptedMessage) String() string { return proto.CompactTextString(m) }
func (*EncryptedMessage) ProtoMessage()    {}
func (*EncryptedMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_69909e00b9236d45, []int{1}
}

func (m *EncryptedMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EncryptedMessage.Unmarshal(m, b)
}
func (m *EncryptedMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EncryptedMessage.Marshal(b, m, deterministic)
}
func (m *EncryptedMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EncryptedMessage.Merge(m, src)
}
func (m *EncryptedMessage) XXX_Size() int {
	return xxx_messageInfo_EncryptedMessage.Size(m)
}
func (m *EncryptedMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_EncryptedMessage.DiscardUnknown(m)
}

var xxx_messageInfo_EncryptedMessage proto.InternalMessageInfo

func (m *EncryptedMessage) GetM1() *binaryquadraticform.BQForm {
	if m != nil {
		return m.M1
	}
	return nil
}

func (m *EncryptedMessage) GetM2() *binaryquadraticform.BQForm {
	if m != nil {
		return m.M2
	}
	return nil
}

func (m *EncryptedMessage) GetProof() *ProofMessage {
	if m != nil {
		return m.Proof
	}
	return nil
}

type ProofMessage struct {
	Salt                 []byte                      `protobuf:"bytes,1,opt,name=salt,proto3" json:"salt,omitempty"`
	U1                   []byte                      `protobuf:"bytes,2,opt,name=u1,proto3" json:"u1,omitempty"`
	U2                   []byte                      `protobuf:"bytes,3,opt,name=u2,proto3" json:"u2,omitempty"`
	T1                   *binaryquadraticform.BQForm `protobuf:"bytes,4,opt,name=t1,proto3" json:"t1,omitempty"`
	T2                   *binaryquadraticform.BQForm `protobuf:"bytes,5,opt,name=t2,proto3" json:"t2,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *ProofMessage) Reset()         { *m = ProofMessage{} }
func (m *ProofMessage) String() string { return proto.CompactTextString(m) }
func (*ProofMessage) ProtoMessage()    {}
func (*ProofMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_69909e00b9236d45, []int{2}
}

func (m *ProofMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProofMessage.Unmarshal(m, b)
}
func (m *ProofMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProofMessage.Marshal(b, m, deterministic)
}
func (m *ProofMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProofMessage.Merge(m, src)
}
func (m *ProofMessage) XXX_Size() int {
	return xxx_messageInfo_ProofMessage.Size(m)
}
func (m *ProofMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ProofMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ProofMessage proto.InternalMessageInfo

func (m *ProofMessage) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *ProofMessage) GetU1() []byte {
	if m != nil {
		return m.U1
	}
	return nil
}

func (m *ProofMessage) GetU2() []byte {
	if m != nil {
		return m.U2
	}
	return nil
}

func (m *ProofMessage) GetT1() *binaryquadraticform.BQForm {
	if m != nil {
		return m.T1
	}
	return nil
}

func (m *ProofMessage) GetT2() *binaryquadraticform.BQForm {
	if m != nil {
		return m.T2
	}
	return nil
}

type VerifyMtaMessage struct {
	BetaG                *ecpointgrouplaw.EcPointMessage `protobuf:"bytes,1,opt,name=betaG,proto3" json:"betaG,omitempty"`
	BG                   *ecpointgrouplaw.EcPointMessage `protobuf:"bytes,2,opt,name=bG,proto3" json:"bG,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *VerifyMtaMessage) Reset()         { *m = VerifyMtaMessage{} }
func (m *VerifyMtaMessage) String() string { return proto.CompactTextString(m) }
func (*VerifyMtaMessage) ProtoMessage()    {}
func (*VerifyMtaMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_69909e00b9236d45, []int{3}
}

func (m *VerifyMtaMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VerifyMtaMessage.Unmarshal(m, b)
}
func (m *VerifyMtaMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VerifyMtaMessage.Marshal(b, m, deterministic)
}
func (m *VerifyMtaMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VerifyMtaMessage.Merge(m, src)
}
func (m *VerifyMtaMessage) XXX_Size() int {
	return xxx_messageInfo_VerifyMtaMessage.Size(m)
}
func (m *VerifyMtaMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_VerifyMtaMessage.DiscardUnknown(m)
}

var xxx_messageInfo_VerifyMtaMessage proto.InternalMessageInfo

func (m *VerifyMtaMessage) GetBetaG() *ecpointgrouplaw.EcPointMessage {
	if m != nil {
		return m.BetaG
	}
	return nil
}

func (m *VerifyMtaMessage) GetBG() *ecpointgrouplaw.EcPointMessage {
	if m != nil {
		return m.BG
	}
	return nil
}

type Hash struct {
	T1                   *binaryquadraticform.BQForm `protobuf:"bytes,1,opt,name=t1,proto3" json:"t1,omitempty"`
	T2                   *binaryquadraticform.BQForm `protobuf:"bytes,2,opt,name=t2,proto3" json:"t2,omitempty"`
	G                    *binaryquadraticform.BQForm `protobuf:"bytes,3,opt,name=g,proto3" json:"g,omitempty"`
	F                    *binaryquadraticform.BQForm `protobuf:"bytes,4,opt,name=f,proto3" json:"f,omitempty"`
	H                    *binaryquadraticform.BQForm `protobuf:"bytes,5,opt,name=h,proto3" json:"h,omitempty"`
	P                    []byte                      `protobuf:"bytes,6,opt,name=p,proto3" json:"p,omitempty"`
	Q                    []byte                      `protobuf:"bytes,7,opt,name=q,proto3" json:"q,omitempty"`
	A                    []byte                      `protobuf:"bytes,8,opt,name=a,proto3" json:"a,omitempty"`
	C                    []byte                      `protobuf:"bytes,9,opt,name=c,proto3" json:"c,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *Hash) Reset()         { *m = Hash{} }
func (m *Hash) String() string { return proto.CompactTextString(m) }
func (*Hash) ProtoMessage()    {}
func (*Hash) Descriptor() ([]byte, []int) {
	return fileDescriptor_69909e00b9236d45, []int{4}
}

func (m *Hash) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Hash.Unmarshal(m, b)
}
func (m *Hash) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Hash.Marshal(b, m, deterministic)
}
func (m *Hash) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Hash.Merge(m, src)
}
func (m *Hash) XXX_Size() int {
	return xxx_messageInfo_Hash.Size(m)
}
func (m *Hash) XXX_DiscardUnknown() {
	xxx_messageInfo_Hash.DiscardUnknown(m)
}

var xxx_messageInfo_Hash proto.InternalMessageInfo

func (m *Hash) GetT1() *binaryquadraticform.BQForm {
	if m != nil {
		return m.T1
	}
	return nil
}

func (m *Hash) GetT2() *binaryquadraticform.BQForm {
	if m != nil {
		return m.T2
	}
	return nil
}

func (m *Hash) GetG() *binaryquadraticform.BQForm {
	if m != nil {
		return m.G
	}
	return nil
}

func (m *Hash) GetF() *binaryquadraticform.BQForm {
	if m != nil {
		return m.F
	}
	return nil
}

func (m *Hash) GetH() *binaryquadraticform.BQForm {
	if m != nil {
		return m.H
	}
	return nil
}

func (m *Hash) GetP() []byte {
	if m != nil {
		return m.P
	}
	return nil
}

func (m *Hash) GetQ() []byte {
	if m != nil {
		return m.Q
	}
	return nil
}

func (m *Hash) GetA() []byte {
	if m != nil {
		return m.A
	}
	return nil
}

func (m *Hash) GetC() []byte {
	if m != nil {
		return m.C
	}
	return nil
}

func init() {
	proto.RegisterType((*PubKeyMessage)(nil), "cl.PubKeyMessage")
	proto.RegisterType((*EncryptedMessage)(nil), "cl.EncryptedMessage")
	proto.RegisterType((*ProofMessage)(nil), "cl.ProofMessage")
	proto.RegisterType((*VerifyMtaMessage)(nil), "cl.VerifyMtaMessage")
	proto.RegisterType((*Hash)(nil), "cl.Hash")
}

func init() {
	proto.RegisterFile("github.com/getamis/alice/crypto/homo/cl/message.proto", fileDescriptor_69909e00b9236d45)
}

var fileDescriptor_69909e00b9236d45 = []byte{
	// 439 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0x4f, 0x8b, 0x13, 0x31,
	0x1c, 0x25, 0xd9, 0x76, 0x5d, 0xb3, 0x5d, 0x29, 0x73, 0x0a, 0xeb, 0xc1, 0x32, 0x07, 0xa9, 0x08,
	0x33, 0x4c, 0x64, 0x4f, 0x9e, 0x14, 0xd6, 0x0a, 0xb2, 0x50, 0x7b, 0xf0, 0x9e, 0x49, 0x33, 0x7f,
	0x60, 0xd2, 0x64, 0x32, 0x19, 0xa4, 0x7e, 0x0e, 0x2f, 0x82, 0x9f, 0xcb, 0xcf, 0x23, 0xf9, 0x53,
	0xa9, 0x22, 0x4c, 0xf6, 0xd6, 0x97, 0xbe, 0xdf, 0x6f, 0xde, 0x7b, 0xbc, 0x1f, 0xba, 0xab, 0x5b,
	0xd3, 0x8c, 0x65, 0xc6, 0xa4, 0xc8, 0x6b, 0x6e, 0xa8, 0x68, 0x87, 0x9c, 0x76, 0x2d, 0xe3, 0x39,
	0xd3, 0x47, 0x65, 0x64, 0xde, 0x48, 0x21, 0x73, 0xd6, 0xe5, 0x82, 0x0f, 0x03, 0xad, 0x79, 0xa6,
	0xb4, 0x34, 0x32, 0x81, 0xac, 0xbb, 0x7d, 0x37, 0x35, 0x5a, 0xb6, 0x07, 0xaa, 0x8f, 0xfd, 0x48,
	0xf7, 0x9a, 0x9a, 0x96, 0x55, 0x52, 0x8b, 0xbf, 0xd7, 0xdc, 0xbe, 0x9d, 0x5a, 0xc1, 0x99, 0x92,
	0xed, 0xc1, 0xd4, 0x5a, 0x8e, 0xaa, 0xa3, 0x5f, 0x73, 0x87, 0xfc, 0x70, 0xfa, 0x0b, 0xa0, 0x9b,
	0xed, 0x58, 0x7e, 0xe2, 0xc7, 0x07, 0xbf, 0x34, 0x59, 0x20, 0xa0, 0x30, 0x58, 0x81, 0xf5, 0x62,
	0x07, 0x94, 0x45, 0x14, 0x43, 0x8f, 0xa8, 0x45, 0x3d, 0xbe, 0xf0, 0xa8, 0x4f, 0x5e, 0x21, 0x50,
	0xe3, 0xd9, 0x0a, 0xac, 0xaf, 0xc9, 0xf3, 0xec, 0x3f, 0x3a, 0xb3, 0xf7, 0x9f, 0x3f, 0x48, 0x2d,
	0x76, 0xa0, 0xb6, 0xd4, 0x0a, 0xcf, 0x23, 0xa8, 0x95, 0xa5, 0x36, 0xf8, 0x32, 0x82, 0xda, 0x58,
	0x39, 0x0c, 0x3f, 0xf1, 0x72, 0x98, 0x45, 0x7b, 0x7c, 0xb5, 0x02, 0xeb, 0x9b, 0x1d, 0xd8, 0xa7,
	0xdf, 0x01, 0x5a, 0xde, 0x1f, 0x5c, 0x02, 0x7c, 0x7f, 0xf2, 0xf6, 0x1a, 0x41, 0x51, 0x38, 0x73,
	0x13, 0xcb, 0xa1, 0x28, 0x1c, 0x99, 0x38, 0xef, 0x93, 0x64, 0x92, 0xbc, 0x44, 0x73, 0xa5, 0xa5,
	0xac, 0x5c, 0x3a, 0xd7, 0x64, 0x99, 0xb1, 0x2e, 0xdb, 0xda, 0x87, 0xf0, 0xe9, 0x9d, 0xff, 0x3b,
	0xfd, 0x09, 0xd0, 0xe2, 0xfc, 0x3d, 0x49, 0xd0, 0x6c, 0xa0, 0x9d, 0x09, 0x89, 0xbb, 0xdf, 0xc9,
	0x33, 0x04, 0xc7, 0x22, 0xa4, 0x0e, 0xc7, 0xc2, 0x61, 0x12, 0x72, 0x87, 0x23, 0xb1, 0xca, 0x4c,
	0x11, 0x93, 0x3c, 0x34, 0xce, 0x86, 0x21, 0x31, 0xd9, 0x43, 0x43, 0xd2, 0x6f, 0x68, 0xf9, 0x85,
	0xeb, 0xb6, 0x3a, 0x3e, 0x18, 0x7a, 0x52, 0x78, 0x87, 0xe6, 0x25, 0x37, 0x74, 0x13, 0x72, 0x7b,
	0x91, 0xfd, 0xd3, 0xa7, 0xec, 0x9e, 0x6d, 0x2d, 0xfe, 0xe3, 0xd4, 0xb1, 0x93, 0x1c, 0xc1, 0x72,
	0x13, 0xe2, 0x9b, 0x9c, 0x81, 0xe5, 0x26, 0xfd, 0x01, 0xd1, 0xec, 0x23, 0x1d, 0x9a, 0x60, 0x0f,
	0x3c, 0xc6, 0x1e, 0x8c, 0xb2, 0xe7, 0x1b, 0x7b, 0x11, 0xdf, 0xd8, 0x59, 0x7c, 0x63, 0xe7, 0xb1,
	0x8d, 0x55, 0xae, 0xdc, 0xa7, 0xe3, 0xea, 0x4f, 0xfd, 0xed, 0xfd, 0xa9, 0x5d, 0x9d, 0x9d, 0x1a,
	0xc3, 0x4f, 0x43, 0xb7, 0xcb, 0x4b, 0x77, 0xad, 0x6f, 0x7e, 0x07, 0x00, 0x00, 0xff, 0xff, 0x54,
	0xa2, 0x39, 0x6e, 0x6a, 0x04, 0x00, 0x00,
}
