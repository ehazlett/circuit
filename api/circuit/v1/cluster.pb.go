// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/ehazlett/circuit/api/circuit/v1/cluster.proto

package v1

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// OpComplete is used to signal completed response
type OpComplete struct {
	Node                 string   `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OpComplete) Reset()         { *m = OpComplete{} }
func (m *OpComplete) String() string { return proto.CompactTextString(m) }
func (*OpComplete) ProtoMessage()    {}
func (*OpComplete) Descriptor() ([]byte, []int) {
	return fileDescriptor_b7bb2fc0c6485ec0, []int{0}
}
func (m *OpComplete) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OpComplete.Unmarshal(m, b)
}
func (m *OpComplete) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OpComplete.Marshal(b, m, deterministic)
}
func (m *OpComplete) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OpComplete.Merge(m, src)
}
func (m *OpComplete) XXX_Size() int {
	return xxx_messageInfo_OpComplete.Size(m)
}
func (m *OpComplete) XXX_DiscardUnknown() {
	xxx_messageInfo_OpComplete.DiscardUnknown(m)
}

var xxx_messageInfo_OpComplete proto.InternalMessageInfo

func (m *OpComplete) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

type NodeInfo struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version              string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeInfo) Reset()         { *m = NodeInfo{} }
func (m *NodeInfo) String() string { return proto.CompactTextString(m) }
func (*NodeInfo) ProtoMessage()    {}
func (*NodeInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_b7bb2fc0c6485ec0, []int{1}
}
func (m *NodeInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeInfo.Unmarshal(m, b)
}
func (m *NodeInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeInfo.Marshal(b, m, deterministic)
}
func (m *NodeInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeInfo.Merge(m, src)
}
func (m *NodeInfo) XXX_Size() int {
	return xxx_messageInfo_NodeInfo.Size(m)
}
func (m *NodeInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeInfo.DiscardUnknown(m)
}

var xxx_messageInfo_NodeInfo proto.InternalMessageInfo

func (m *NodeInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *NodeInfo) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

type NodesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodesRequest) Reset()         { *m = NodesRequest{} }
func (m *NodesRequest) String() string { return proto.CompactTextString(m) }
func (*NodesRequest) ProtoMessage()    {}
func (*NodesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b7bb2fc0c6485ec0, []int{2}
}
func (m *NodesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodesRequest.Unmarshal(m, b)
}
func (m *NodesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodesRequest.Marshal(b, m, deterministic)
}
func (m *NodesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodesRequest.Merge(m, src)
}
func (m *NodesRequest) XXX_Size() int {
	return xxx_messageInfo_NodesRequest.Size(m)
}
func (m *NodesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NodesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NodesRequest proto.InternalMessageInfo

type NodesResponse struct {
	Nodes                []*NodeInfo `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *NodesResponse) Reset()         { *m = NodesResponse{} }
func (m *NodesResponse) String() string { return proto.CompactTextString(m) }
func (*NodesResponse) ProtoMessage()    {}
func (*NodesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b7bb2fc0c6485ec0, []int{3}
}
func (m *NodesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodesResponse.Unmarshal(m, b)
}
func (m *NodesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodesResponse.Marshal(b, m, deterministic)
}
func (m *NodesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodesResponse.Merge(m, src)
}
func (m *NodesResponse) XXX_Size() int {
	return xxx_messageInfo_NodesResponse.Size(m)
}
func (m *NodesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NodesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NodesResponse proto.InternalMessageInfo

func (m *NodesResponse) GetNodes() []*NodeInfo {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type ContainerIPQuery struct {
	Container            string   `protobuf:"bytes,1,opt,name=container,proto3" json:"container,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ContainerIPQuery) Reset()         { *m = ContainerIPQuery{} }
func (m *ContainerIPQuery) String() string { return proto.CompactTextString(m) }
func (*ContainerIPQuery) ProtoMessage()    {}
func (*ContainerIPQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_b7bb2fc0c6485ec0, []int{4}
}
func (m *ContainerIPQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ContainerIPQuery.Unmarshal(m, b)
}
func (m *ContainerIPQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ContainerIPQuery.Marshal(b, m, deterministic)
}
func (m *ContainerIPQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContainerIPQuery.Merge(m, src)
}
func (m *ContainerIPQuery) XXX_Size() int {
	return xxx_messageInfo_ContainerIPQuery.Size(m)
}
func (m *ContainerIPQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_ContainerIPQuery.DiscardUnknown(m)
}

var xxx_messageInfo_ContainerIPQuery proto.InternalMessageInfo

func (m *ContainerIPQuery) GetContainer() string {
	if m != nil {
		return m.Container
	}
	return ""
}

func init() {
	proto.RegisterType((*OpComplete)(nil), "io.circuit.v1.OpComplete")
	proto.RegisterType((*NodeInfo)(nil), "io.circuit.v1.NodeInfo")
	proto.RegisterType((*NodesRequest)(nil), "io.circuit.v1.NodesRequest")
	proto.RegisterType((*NodesResponse)(nil), "io.circuit.v1.NodesResponse")
	proto.RegisterType((*ContainerIPQuery)(nil), "io.circuit.v1.ContainerIPQuery")
}

func init() {
	proto.RegisterFile("github.com/ehazlett/circuit/api/circuit/v1/cluster.proto", fileDescriptor_b7bb2fc0c6485ec0)
}

var fileDescriptor_b7bb2fc0c6485ec0 = []byte{
	// 269 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xbb, 0x4f, 0xf3, 0x40,
	0x10, 0xc4, 0xe5, 0xef, 0x23, 0x84, 0x2c, 0x04, 0xa1, 0x6b, 0xb0, 0x20, 0x85, 0xe5, 0x2a, 0x4d,
	0xce, 0x38, 0x34, 0x91, 0x90, 0x28, 0xe2, 0x2a, 0x05, 0x2f, 0x97, 0x74, 0x8e, 0xb3, 0x90, 0x93,
	0xec, 0xdb, 0xe3, 0x1e, 0x96, 0xe0, 0xaf, 0x47, 0x39, 0xdb, 0x44, 0x48, 0x29, 0xe8, 0x76, 0x67,
	0xee, 0x34, 0xbf, 0x1d, 0x58, 0xbc, 0x0b, 0xbb, 0x75, 0x6b, 0x5e, 0x52, 0x9d, 0xe0, 0xb6, 0xf8,
	0xaa, 0xd0, 0xda, 0xa4, 0x14, 0xba, 0x74, 0xc2, 0x26, 0x85, 0x12, 0x3f, 0x73, 0x93, 0x26, 0x65,
	0xe5, 0x8c, 0x45, 0xcd, 0x95, 0x26, 0x4b, 0x6c, 0x2c, 0x88, 0x77, 0x26, 0x6f, 0xd2, 0x38, 0x02,
	0x78, 0x52, 0x19, 0xd5, 0xaa, 0x42, 0x8b, 0x8c, 0xc1, 0x91, 0xa4, 0x0d, 0x86, 0x41, 0x14, 0x4c,
	0x47, 0xb9, 0x9f, 0xe3, 0x05, 0x9c, 0x3c, 0xd2, 0x06, 0x57, 0xf2, 0x8d, 0xbc, 0x5f, 0xd4, 0x7b,
	0xbf, 0xa8, 0x91, 0x85, 0x30, 0x6c, 0x50, 0x1b, 0x41, 0x32, 0xfc, 0xe7, 0xe5, 0x7e, 0x8d, 0xcf,
	0xe1, 0x6c, 0xf7, 0xd3, 0xe4, 0xf8, 0xe1, 0xd0, 0xd8, 0xf8, 0x1e, 0xc6, 0xdd, 0x6e, 0x14, 0x49,
	0x83, 0x6c, 0x06, 0x83, 0x5d, 0x84, 0x09, 0x83, 0xe8, 0xff, 0xf4, 0x74, 0x7e, 0xc9, 0x7f, 0xb1,
	0xf1, 0x3e, 0x36, 0x6f, 0x5f, 0xc5, 0x37, 0x70, 0x91, 0x91, 0xb4, 0x85, 0x90, 0xa8, 0x57, 0xcf,
	0x2f, 0x0e, 0xf5, 0x27, 0x9b, 0xc0, 0xa8, 0xec, 0xb5, 0x0e, 0x6b, 0x2f, 0xcc, 0x1f, 0x60, 0x98,
	0xb5, 0xd7, 0xb3, 0x25, 0x0c, 0x7c, 0x38, 0xbb, 0x3e, 0x90, 0xd2, 0x23, 0x5e, 0x4d, 0x0e, 0x9b,
	0x2d, 0xef, 0x32, 0x79, 0x9d, 0xfd, 0xbd, 0xf7, 0xbb, 0x26, 0x5d, 0x1f, 0xfb, 0xce, 0x6f, 0xbf,
	0x03, 0x00, 0x00, 0xff, 0xff, 0xe7, 0x51, 0x1c, 0xa6, 0xaf, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ClusterClient is the client API for Cluster service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ClusterClient interface {
	Nodes(ctx context.Context, in *NodesRequest, opts ...grpc.CallOption) (*NodesResponse, error)
}

type clusterClient struct {
	cc *grpc.ClientConn
}

func NewClusterClient(cc *grpc.ClientConn) ClusterClient {
	return &clusterClient{cc}
}

func (c *clusterClient) Nodes(ctx context.Context, in *NodesRequest, opts ...grpc.CallOption) (*NodesResponse, error) {
	out := new(NodesResponse)
	err := c.cc.Invoke(ctx, "/io.circuit.v1.Cluster/Nodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterServer is the server API for Cluster service.
type ClusterServer interface {
	Nodes(context.Context, *NodesRequest) (*NodesResponse, error)
}

// UnimplementedClusterServer can be embedded to have forward compatible implementations.
type UnimplementedClusterServer struct {
}

func (*UnimplementedClusterServer) Nodes(ctx context.Context, req *NodesRequest) (*NodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Nodes not implemented")
}

func RegisterClusterServer(s *grpc.Server, srv ClusterServer) {
	s.RegisterService(&_Cluster_serviceDesc, srv)
}

func _Cluster_Nodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).Nodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/io.circuit.v1.Cluster/Nodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).Nodes(ctx, req.(*NodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Cluster_serviceDesc = grpc.ServiceDesc{
	ServiceName: "io.circuit.v1.Cluster",
	HandlerType: (*ClusterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Nodes",
			Handler:    _Cluster_Nodes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/ehazlett/circuit/api/circuit/v1/cluster.proto",
}
