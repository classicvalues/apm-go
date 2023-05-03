// © 2023 SolarWinds Worldwide, LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: collector.proto

package collector

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ResultCode int32

const (
	ResultCode_OK              ResultCode = 0
	ResultCode_TRY_LATER       ResultCode = 1
	ResultCode_INVALID_API_KEY ResultCode = 2
	ResultCode_LIMIT_EXCEEDED  ResultCode = 3
	ResultCode_REDIRECT        ResultCode = 4
)

var ResultCode_name = map[int32]string{
	0: "OK",
	1: "TRY_LATER",
	2: "INVALID_API_KEY",
	3: "LIMIT_EXCEEDED",
	4: "REDIRECT",
}

var ResultCode_value = map[string]int32{
	"OK":              0,
	"TRY_LATER":       1,
	"INVALID_API_KEY": 2,
	"LIMIT_EXCEEDED":  3,
	"REDIRECT":        4,
}

func (x ResultCode) String() string {
	return proto.EnumName(ResultCode_name, int32(x))
}

func (ResultCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_9305884a292fdf82, []int{0}
}

type EncodingType int32

const (
	EncodingType_BSON     EncodingType = 0
	EncodingType_PROTOBUF EncodingType = 1
)

var EncodingType_name = map[int32]string{
	0: "BSON",
	1: "PROTOBUF",
}

var EncodingType_value = map[string]int32{
	"BSON":     0,
	"PROTOBUF": 1,
}

func (x EncodingType) String() string {
	return proto.EnumName(EncodingType_name, int32(x))
}

func (EncodingType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_9305884a292fdf82, []int{1}
}

type HostType int32

const (
	HostType_PERSISTENT HostType = 0
	HostType_AWS_LAMBDA HostType = 1
)

var HostType_name = map[int32]string{
	0: "PERSISTENT",
	1: "AWS_LAMBDA",
}

var HostType_value = map[string]int32{
	"PERSISTENT": 0,
	"AWS_LAMBDA": 1,
}

func (x HostType) String() string {
	return proto.EnumName(HostType_name, int32(x))
}

func (HostType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_9305884a292fdf82, []int{2}
}

type OboeSettingType int32

const (
	OboeSettingType_DEFAULT_SAMPLE_RATE        OboeSettingType = 0
	OboeSettingType_LAYER_SAMPLE_RATE          OboeSettingType = 1
	OboeSettingType_LAYER_APP_SAMPLE_RATE      OboeSettingType = 2
	OboeSettingType_LAYER_HTTPHOST_SAMPLE_RATE OboeSettingType = 3
	OboeSettingType_CONFIG_STRING              OboeSettingType = 4
	OboeSettingType_CONFIG_INT                 OboeSettingType = 5
)

var OboeSettingType_name = map[int32]string{
	0: "DEFAULT_SAMPLE_RATE",
	1: "LAYER_SAMPLE_RATE",
	2: "LAYER_APP_SAMPLE_RATE",
	3: "LAYER_HTTPHOST_SAMPLE_RATE",
	4: "CONFIG_STRING",
	5: "CONFIG_INT",
}

var OboeSettingType_value = map[string]int32{
	"DEFAULT_SAMPLE_RATE":        0,
	"LAYER_SAMPLE_RATE":          1,
	"LAYER_APP_SAMPLE_RATE":      2,
	"LAYER_HTTPHOST_SAMPLE_RATE": 3,
	"CONFIG_STRING":              4,
	"CONFIG_INT":                 5,
}

func (x OboeSettingType) String() string {
	return proto.EnumName(OboeSettingType_name, int32(x))
}

func (OboeSettingType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_9305884a292fdf82, []int{3}
}

type HostID struct {
	Hostname               string   `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	IpAddresses            []string `protobuf:"bytes,2,rep,name=ip_addresses,json=ipAddresses,proto3" json:"ip_addresses,omitempty"`
	Uuid                   string   `protobuf:"bytes,3,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Pid                    int32    `protobuf:"varint,4,opt,name=pid,proto3" json:"pid,omitempty"`
	Ec2InstanceID          string   `protobuf:"bytes,5,opt,name=ec2InstanceID,proto3" json:"ec2InstanceID,omitempty"`
	Ec2AvailabilityZone    string   `protobuf:"bytes,6,opt,name=ec2AvailabilityZone,proto3" json:"ec2AvailabilityZone,omitempty"`
	DockerContainerID      string   `protobuf:"bytes,7,opt,name=dockerContainerID,proto3" json:"dockerContainerID,omitempty"`
	MacAddresses           []string `protobuf:"bytes,8,rep,name=macAddresses,proto3" json:"macAddresses,omitempty"`
	HerokuDynoID           string   `protobuf:"bytes,9,opt,name=herokuDynoID,proto3" json:"herokuDynoID,omitempty"`
	AzAppServiceInstanceID string   `protobuf:"bytes,10,opt,name=azAppServiceInstanceID,proto3" json:"azAppServiceInstanceID,omitempty"`
	HostType               HostType `protobuf:"varint,11,opt,name=hostType,proto3,enum=collectorpb.HostType" json:"hostType,omitempty"`
	XXX_NoUnkeyedLiteral   struct{} `json:"-"`
	XXX_unrecognized       []byte   `json:"-"`
	XXX_sizecache          int32    `json:"-"`
}

func (m *HostID) Reset()         { *m = HostID{} }
func (m *HostID) String() string { return proto.CompactTextString(m) }
func (*HostID) ProtoMessage()    {}
func (*HostID) Descriptor() ([]byte, []int) {
	return fileDescriptor_9305884a292fdf82, []int{0}
}

func (m *HostID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HostID.Unmarshal(m, b)
}
func (m *HostID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HostID.Marshal(b, m, deterministic)
}
func (m *HostID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HostID.Merge(m, src)
}
func (m *HostID) XXX_Size() int {
	return xxx_messageInfo_HostID.Size(m)
}
func (m *HostID) XXX_DiscardUnknown() {
	xxx_messageInfo_HostID.DiscardUnknown(m)
}

var xxx_messageInfo_HostID proto.InternalMessageInfo

func (m *HostID) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *HostID) GetIpAddresses() []string {
	if m != nil {
		return m.IpAddresses
	}
	return nil
}

func (m *HostID) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *HostID) GetPid() int32 {
	if m != nil {
		return m.Pid
	}
	return 0
}

func (m *HostID) GetEc2InstanceID() string {
	if m != nil {
		return m.Ec2InstanceID
	}
	return ""
}

func (m *HostID) GetEc2AvailabilityZone() string {
	if m != nil {
		return m.Ec2AvailabilityZone
	}
	return ""
}

func (m *HostID) GetDockerContainerID() string {
	if m != nil {
		return m.DockerContainerID
	}
	return ""
}

func (m *HostID) GetMacAddresses() []string {
	if m != nil {
		return m.MacAddresses
	}
	return nil
}

func (m *HostID) GetHerokuDynoID() string {
	if m != nil {
		return m.HerokuDynoID
	}
	return ""
}

func (m *HostID) GetAzAppServiceInstanceID() string {
	if m != nil {
		return m.AzAppServiceInstanceID
	}
	return ""
}

func (m *HostID) GetHostType() HostType {
	if m != nil {
		return m.HostType
	}
	return HostType_PERSISTENT
}

type OboeSetting struct {
	Type                 OboeSettingType   `protobuf:"varint,1,opt,name=type,proto3,enum=collectorpb.OboeSettingType" json:"type,omitempty"`
	Flags                []byte            `protobuf:"bytes,2,opt,name=flags,proto3" json:"flags,omitempty"`
	Timestamp            int64             `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Value                int64             `protobuf:"varint,4,opt,name=value,proto3" json:"value,omitempty"`
	Layer                []byte            `protobuf:"bytes,5,opt,name=layer,proto3" json:"layer,omitempty"`
	Arguments            map[string][]byte `protobuf:"bytes,7,rep,name=arguments,proto3" json:"arguments,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Ttl                  int64             `protobuf:"varint,8,opt,name=ttl,proto3" json:"ttl,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *OboeSetting) Reset()         { *m = OboeSetting{} }
func (m *OboeSetting) String() string { return proto.CompactTextString(m) }
func (*OboeSetting) ProtoMessage()    {}
func (*OboeSetting) Descriptor() ([]byte, []int) {
	return fileDescriptor_9305884a292fdf82, []int{1}
}

func (m *OboeSetting) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OboeSetting.Unmarshal(m, b)
}
func (m *OboeSetting) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OboeSetting.Marshal(b, m, deterministic)
}
func (m *OboeSetting) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OboeSetting.Merge(m, src)
}
func (m *OboeSetting) XXX_Size() int {
	return xxx_messageInfo_OboeSetting.Size(m)
}
func (m *OboeSetting) XXX_DiscardUnknown() {
	xxx_messageInfo_OboeSetting.DiscardUnknown(m)
}

var xxx_messageInfo_OboeSetting proto.InternalMessageInfo

func (m *OboeSetting) GetType() OboeSettingType {
	if m != nil {
		return m.Type
	}
	return OboeSettingType_DEFAULT_SAMPLE_RATE
}

func (m *OboeSetting) GetFlags() []byte {
	if m != nil {
		return m.Flags
	}
	return nil
}

func (m *OboeSetting) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *OboeSetting) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *OboeSetting) GetLayer() []byte {
	if m != nil {
		return m.Layer
	}
	return nil
}

func (m *OboeSetting) GetArguments() map[string][]byte {
	if m != nil {
		return m.Arguments
	}
	return nil
}

func (m *OboeSetting) GetTtl() int64 {
	if m != nil {
		return m.Ttl
	}
	return 0
}

type MessageRequest struct {
	ApiKey               string       `protobuf:"bytes,1,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
	Messages             [][]byte     `protobuf:"bytes,2,rep,name=messages,proto3" json:"messages,omitempty"`
	Encoding             EncodingType `protobuf:"varint,3,opt,name=encoding,proto3,enum=collectorpb.EncodingType" json:"encoding,omitempty"`
	Identity             *HostID      `protobuf:"bytes,4,opt,name=identity,proto3" json:"identity,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *MessageRequest) Reset()         { *m = MessageRequest{} }
func (m *MessageRequest) String() string { return proto.CompactTextString(m) }
func (*MessageRequest) ProtoMessage()    {}
func (*MessageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9305884a292fdf82, []int{2}
}

func (m *MessageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageRequest.Unmarshal(m, b)
}
func (m *MessageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageRequest.Marshal(b, m, deterministic)
}
func (m *MessageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageRequest.Merge(m, src)
}
func (m *MessageRequest) XXX_Size() int {
	return xxx_messageInfo_MessageRequest.Size(m)
}
func (m *MessageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MessageRequest proto.InternalMessageInfo

func (m *MessageRequest) GetApiKey() string {
	if m != nil {
		return m.ApiKey
	}
	return ""
}

func (m *MessageRequest) GetMessages() [][]byte {
	if m != nil {
		return m.Messages
	}
	return nil
}

func (m *MessageRequest) GetEncoding() EncodingType {
	if m != nil {
		return m.Encoding
	}
	return EncodingType_BSON
}

func (m *MessageRequest) GetIdentity() *HostID {
	if m != nil {
		return m.Identity
	}
	return nil
}

type MessageResult struct {
	Result ResultCode `protobuf:"varint,1,opt,name=result,proto3,enum=collectorpb.ResultCode" json:"result,omitempty"`
	Arg    string     `protobuf:"bytes,2,opt,name=arg,proto3" json:"arg,omitempty"`
	// warning specifies a user-facing warning message; agents attempt to squelch repeated warnings,
	// so care should be taken to ensure that warning messages are consistent across all RPCs.
	Warning              string   `protobuf:"bytes,4,opt,name=warning,proto3" json:"warning,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MessageResult) Reset()         { *m = MessageResult{} }
func (m *MessageResult) String() string { return proto.CompactTextString(m) }
func (*MessageResult) ProtoMessage()    {}
func (*MessageResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_9305884a292fdf82, []int{3}
}

func (m *MessageResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageResult.Unmarshal(m, b)
}
func (m *MessageResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageResult.Marshal(b, m, deterministic)
}
func (m *MessageResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageResult.Merge(m, src)
}
func (m *MessageResult) XXX_Size() int {
	return xxx_messageInfo_MessageResult.Size(m)
}
func (m *MessageResult) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageResult.DiscardUnknown(m)
}

var xxx_messageInfo_MessageResult proto.InternalMessageInfo

func (m *MessageResult) GetResult() ResultCode {
	if m != nil {
		return m.Result
	}
	return ResultCode_OK
}

func (m *MessageResult) GetArg() string {
	if m != nil {
		return m.Arg
	}
	return ""
}

func (m *MessageResult) GetWarning() string {
	if m != nil {
		return m.Warning
	}
	return ""
}

type SettingsRequest struct {
	ApiKey               string   `protobuf:"bytes,1,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
	Identity             *HostID  `protobuf:"bytes,2,opt,name=identity,proto3" json:"identity,omitempty"`
	ClientVersion        string   `protobuf:"bytes,3,opt,name=clientVersion,proto3" json:"clientVersion,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SettingsRequest) Reset()         { *m = SettingsRequest{} }
func (m *SettingsRequest) String() string { return proto.CompactTextString(m) }
func (*SettingsRequest) ProtoMessage()    {}
func (*SettingsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9305884a292fdf82, []int{4}
}

func (m *SettingsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SettingsRequest.Unmarshal(m, b)
}
func (m *SettingsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SettingsRequest.Marshal(b, m, deterministic)
}
func (m *SettingsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SettingsRequest.Merge(m, src)
}
func (m *SettingsRequest) XXX_Size() int {
	return xxx_messageInfo_SettingsRequest.Size(m)
}
func (m *SettingsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SettingsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SettingsRequest proto.InternalMessageInfo

func (m *SettingsRequest) GetApiKey() string {
	if m != nil {
		return m.ApiKey
	}
	return ""
}

func (m *SettingsRequest) GetIdentity() *HostID {
	if m != nil {
		return m.Identity
	}
	return nil
}

func (m *SettingsRequest) GetClientVersion() string {
	if m != nil {
		return m.ClientVersion
	}
	return ""
}

type SettingsResult struct {
	Result   ResultCode     `protobuf:"varint,1,opt,name=result,proto3,enum=collectorpb.ResultCode" json:"result,omitempty"`
	Arg      string         `protobuf:"bytes,2,opt,name=arg,proto3" json:"arg,omitempty"`
	Settings []*OboeSetting `protobuf:"bytes,3,rep,name=settings,proto3" json:"settings,omitempty"`
	// warning specifies a user-facing warning message; see note on MessageResult.warning for details.
	Warning              string   `protobuf:"bytes,4,opt,name=warning,proto3" json:"warning,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SettingsResult) Reset()         { *m = SettingsResult{} }
func (m *SettingsResult) String() string { return proto.CompactTextString(m) }
func (*SettingsResult) ProtoMessage()    {}
func (*SettingsResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_9305884a292fdf82, []int{5}
}

func (m *SettingsResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SettingsResult.Unmarshal(m, b)
}
func (m *SettingsResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SettingsResult.Marshal(b, m, deterministic)
}
func (m *SettingsResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SettingsResult.Merge(m, src)
}
func (m *SettingsResult) XXX_Size() int {
	return xxx_messageInfo_SettingsResult.Size(m)
}
func (m *SettingsResult) XXX_DiscardUnknown() {
	xxx_messageInfo_SettingsResult.DiscardUnknown(m)
}

var xxx_messageInfo_SettingsResult proto.InternalMessageInfo

func (m *SettingsResult) GetResult() ResultCode {
	if m != nil {
		return m.Result
	}
	return ResultCode_OK
}

func (m *SettingsResult) GetArg() string {
	if m != nil {
		return m.Arg
	}
	return ""
}

func (m *SettingsResult) GetSettings() []*OboeSetting {
	if m != nil {
		return m.Settings
	}
	return nil
}

func (m *SettingsResult) GetWarning() string {
	if m != nil {
		return m.Warning
	}
	return ""
}

type PingRequest struct {
	ApiKey               string   `protobuf:"bytes,1,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingRequest) Reset()         { *m = PingRequest{} }
func (m *PingRequest) String() string { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()    {}
func (*PingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9305884a292fdf82, []int{6}
}

func (m *PingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingRequest.Unmarshal(m, b)
}
func (m *PingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingRequest.Marshal(b, m, deterministic)
}
func (m *PingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingRequest.Merge(m, src)
}
func (m *PingRequest) XXX_Size() int {
	return xxx_messageInfo_PingRequest.Size(m)
}
func (m *PingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PingRequest proto.InternalMessageInfo

func (m *PingRequest) GetApiKey() string {
	if m != nil {
		return m.ApiKey
	}
	return ""
}

func init() {
	proto.RegisterEnum("collectorpb.ResultCode", ResultCode_name, ResultCode_value)
	proto.RegisterEnum("collectorpb.EncodingType", EncodingType_name, EncodingType_value)
	proto.RegisterEnum("collectorpb.HostType", HostType_name, HostType_value)
	proto.RegisterEnum("collectorpb.OboeSettingType", OboeSettingType_name, OboeSettingType_value)
	proto.RegisterType((*HostID)(nil), "collectorpb.HostID")
	proto.RegisterType((*OboeSetting)(nil), "collectorpb.OboeSetting")
	proto.RegisterMapType((map[string][]byte)(nil), "collectorpb.OboeSetting.ArgumentsEntry")
	proto.RegisterType((*MessageRequest)(nil), "collectorpb.MessageRequest")
	proto.RegisterType((*MessageResult)(nil), "collectorpb.MessageResult")
	proto.RegisterType((*SettingsRequest)(nil), "collectorpb.SettingsRequest")
	proto.RegisterType((*SettingsResult)(nil), "collectorpb.SettingsResult")
	proto.RegisterType((*PingRequest)(nil), "collectorpb.PingRequest")
}

func init() { proto.RegisterFile("collector.proto", fileDescriptor_9305884a292fdf82) }

var fileDescriptor_9305884a292fdf82 = []byte{
	// 954 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0x4d, 0x6f, 0xdb, 0x46,
	0x13, 0x16, 0x45, 0x59, 0x96, 0x46, 0xb2, 0x4c, 0xaf, 0x5f, 0xbf, 0x66, 0x94, 0xa0, 0x50, 0x85,
	0x22, 0x15, 0x84, 0xc2, 0x49, 0xd5, 0x0f, 0x14, 0x45, 0x51, 0x80, 0x16, 0x69, 0x8b, 0xb5, 0xbe,
	0xb0, 0x64, 0xd2, 0x3a, 0x17, 0x62, 0x4d, 0x6d, 0x15, 0xc2, 0x14, 0xc9, 0x92, 0x2b, 0x17, 0xea,
	0xad, 0xbf, 0xa4, 0x87, 0x02, 0xbd, 0xf5, 0xda, 0xdf, 0xd6, 0x63, 0xb1, 0x4b, 0xea, 0x83, 0x49,
	0x9c, 0x14, 0x48, 0x6f, 0x3b, 0xcf, 0x3c, 0x33, 0x9c, 0x99, 0x67, 0x76, 0x41, 0x38, 0x74, 0x43,
	0xdf, 0xa7, 0x2e, 0x0b, 0xe3, 0xb3, 0x28, 0x0e, 0x59, 0x88, 0x6a, 0x1b, 0x20, 0xba, 0x69, 0xff,
	0x21, 0x43, 0x79, 0x10, 0x26, 0xcc, 0xd4, 0x51, 0x13, 0x2a, 0x2f, 0xc3, 0x84, 0x05, 0x64, 0x41,
	0x55, 0xa9, 0x25, 0x75, 0xaa, 0x78, 0x63, 0xa3, 0x0f, 0xa1, 0xee, 0x45, 0x0e, 0x99, 0xcd, 0x62,
	0x9a, 0x24, 0x34, 0x51, 0x8b, 0x2d, 0xb9, 0x53, 0xc5, 0x35, 0x2f, 0xd2, 0xd6, 0x10, 0x42, 0x50,
	0x5a, 0x2e, 0xbd, 0x99, 0x2a, 0x8b, 0x50, 0x71, 0x46, 0x0a, 0xc8, 0x91, 0x37, 0x53, 0x4b, 0x2d,
	0xa9, 0xb3, 0x87, 0xf9, 0x11, 0x7d, 0x04, 0x07, 0xd4, 0xed, 0x99, 0x41, 0xc2, 0x48, 0xe0, 0x52,
	0x53, 0x57, 0xf7, 0x04, 0x3d, 0x0f, 0xa2, 0xa7, 0x70, 0x4c, 0xdd, 0x9e, 0x76, 0x47, 0x3c, 0x9f,
	0xdc, 0x78, 0xbe, 0xc7, 0x56, 0x2f, 0xc2, 0x80, 0xaa, 0x65, 0xc1, 0x7d, 0x93, 0x0b, 0x7d, 0x02,
	0x47, 0xb3, 0xd0, 0xbd, 0xa5, 0x71, 0x3f, 0x0c, 0x18, 0xf1, 0x02, 0x1a, 0x9b, 0xba, 0xba, 0x2f,
	0xf8, 0xaf, 0x3b, 0x50, 0x1b, 0xea, 0x0b, 0xe2, 0x6e, 0x6a, 0x57, 0x2b, 0xa2, 0x9d, 0x1c, 0xc6,
	0x39, 0x2f, 0x69, 0x1c, 0xde, 0x2e, 0xf5, 0x55, 0x10, 0x9a, 0xba, 0x5a, 0x15, 0xc9, 0x72, 0x18,
	0xfa, 0x12, 0xfe, 0x4f, 0x7e, 0xd1, 0xa2, 0xc8, 0xa2, 0xf1, 0x9d, 0xe7, 0xd2, 0x9d, 0xb6, 0x40,
	0xb0, 0xef, 0xf1, 0xa2, 0x4f, 0xd3, 0x51, 0xdb, 0xab, 0x88, 0xaa, 0xb5, 0x96, 0xd4, 0x69, 0xf4,
	0x4e, 0xce, 0x76, 0x54, 0x39, 0x1b, 0x64, 0x4e, 0xbc, 0xa1, 0xb5, 0xff, 0x2a, 0x42, 0x6d, 0x72,
	0x13, 0x52, 0x8b, 0x32, 0xe6, 0x05, 0x73, 0xf4, 0x14, 0x4a, 0x8c, 0x87, 0x4b, 0x22, 0xfc, 0x51,
	0x2e, 0x7c, 0x87, 0x27, 0xb2, 0x08, 0x26, 0xfa, 0x1f, 0xec, 0xfd, 0xe8, 0x93, 0x39, 0x17, 0x4f,
	0xea, 0xd4, 0x71, 0x6a, 0xa0, 0x47, 0x50, 0x65, 0xde, 0x82, 0x26, 0x8c, 0x2c, 0x22, 0xa1, 0x9d,
	0x8c, 0xb7, 0x00, 0x8f, 0xb9, 0x23, 0xfe, 0x92, 0x0a, 0x09, 0x65, 0x9c, 0x1a, 0x1c, 0xf5, 0xc9,
	0x8a, 0xc6, 0x42, 0xbc, 0x3a, 0x4e, 0x0d, 0x64, 0x40, 0x95, 0xc4, 0xf3, 0xe5, 0x82, 0x06, 0x2c,
	0x51, 0xf7, 0x5b, 0x72, 0xa7, 0xd6, 0xfb, 0xf8, 0xbe, 0xb2, 0xce, 0xb4, 0x35, 0xd3, 0x08, 0x58,
	0xbc, 0xc2, 0xdb, 0x48, 0xbe, 0x33, 0x8c, 0xf9, 0x6a, 0x45, 0x7c, 0x90, 0x1f, 0x9b, 0xdf, 0x40,
	0x23, 0x4f, 0xe7, 0x9c, 0x5b, 0xba, 0xca, 0xb6, 0x94, 0x1f, 0xb7, 0x85, 0x66, 0xcd, 0x09, 0xe3,
	0xeb, 0xe2, 0x57, 0x52, 0xfb, 0x4f, 0x09, 0x1a, 0x23, 0x9a, 0x24, 0x64, 0x4e, 0x31, 0xfd, 0x69,
	0x49, 0x13, 0x86, 0x4e, 0x61, 0x9f, 0x44, 0x9e, 0xb3, 0x4d, 0x51, 0x26, 0x91, 0x77, 0x45, 0x57,
	0xfc, 0x0a, 0x2c, 0x52, 0x6a, 0xba, 0xe2, 0x75, 0xbc, 0xb1, 0xd1, 0x17, 0x50, 0xa1, 0x81, 0x1b,
	0xce, 0xbc, 0x60, 0x2e, 0xe6, 0xd4, 0xe8, 0x3d, 0xc8, 0x75, 0x67, 0x64, 0xce, 0x54, 0xb7, 0x35,
	0x15, 0x3d, 0x81, 0x8a, 0x37, 0xa3, 0x01, 0xf3, 0xd8, 0x4a, 0x0c, 0xb1, 0xd6, 0x3b, 0x7e, 0x4d,
	0x6a, 0x53, 0xc7, 0x1b, 0x52, 0xdb, 0x87, 0x83, 0x4d, 0xb9, 0xc9, 0xd2, 0x67, 0xe8, 0x09, 0x94,
	0x63, 0x71, 0xca, 0xb4, 0x3e, 0xcd, 0xc5, 0xa7, 0xa4, 0x7e, 0x38, 0xa3, 0x38, 0xa3, 0xf1, 0xe9,
	0x90, 0x78, 0x2e, 0x26, 0x51, 0xc5, 0xfc, 0x88, 0x54, 0xd8, 0xff, 0x99, 0xc4, 0x01, 0x2f, 0xbd,
	0x24, 0xd0, 0xb5, 0xd9, 0xfe, 0x55, 0x82, 0xc3, 0x4c, 0x93, 0xe4, 0x9d, 0xe3, 0xd9, 0xed, 0xa5,
	0xf8, 0x2f, 0x7a, 0xe1, 0xb7, 0xdd, 0xf5, 0x3d, 0x1a, 0xb0, 0xe7, 0x34, 0x4e, 0xbc, 0x30, 0xc8,
	0x1e, 0x87, 0x3c, 0xd8, 0xfe, 0x5d, 0x82, 0xc6, 0xb6, 0x86, 0xff, 0xaa, 0xe7, 0xcf, 0xa1, 0x92,
	0x64, 0x49, 0x55, 0x59, 0x6c, 0xa3, 0x7a, 0xdf, 0x36, 0xe2, 0x0d, 0xf3, 0x2d, 0x93, 0x7a, 0x0c,
	0xb5, 0x29, 0xe7, 0xbe, 0x63, 0x48, 0xdd, 0x17, 0x00, 0xdb, 0xfa, 0x50, 0x19, 0x8a, 0x93, 0x2b,
	0xa5, 0x80, 0x0e, 0xa0, 0x6a, 0xe3, 0x6b, 0x67, 0xa8, 0xd9, 0x06, 0x56, 0x24, 0x74, 0x0c, 0x87,
	0xe6, 0xf8, 0xb9, 0x36, 0x34, 0x75, 0x47, 0x9b, 0x9a, 0xce, 0x95, 0x71, 0xad, 0x14, 0x11, 0x82,
	0xc6, 0xd0, 0x1c, 0x99, 0xb6, 0x63, 0xfc, 0xd0, 0x37, 0x0c, 0xdd, 0xd0, 0x15, 0x19, 0xd5, 0xa1,
	0x82, 0x0d, 0xdd, 0xc4, 0x46, 0xdf, 0x56, 0x4a, 0xdd, 0xc7, 0x50, 0xdf, 0x5d, 0x33, 0x54, 0x81,
	0xd2, 0xb9, 0x35, 0x19, 0x2b, 0x05, 0xce, 0x9b, 0xe2, 0x89, 0x3d, 0x39, 0x7f, 0x76, 0xa1, 0x48,
	0xdd, 0x2e, 0x54, 0xd6, 0x4f, 0x08, 0x6a, 0x00, 0x4c, 0x0d, 0x6c, 0x99, 0x96, 0x6d, 0x8c, 0x6d,
	0xa5, 0xc0, 0x6d, 0xed, 0x7b, 0xcb, 0x19, 0x6a, 0xa3, 0x73, 0x5d, 0x53, 0xa4, 0xee, 0x6f, 0x12,
	0x1c, 0xbe, 0xf2, 0x60, 0xa0, 0x53, 0x38, 0xd6, 0x8d, 0x0b, 0xed, 0xd9, 0xd0, 0x76, 0x2c, 0x6d,
	0x34, 0x1d, 0x1a, 0x0e, 0xd6, 0x6c, 0x43, 0x29, 0xa0, 0x13, 0x38, 0x1a, 0x6a, 0xd7, 0x06, 0xce,
	0xc1, 0x12, 0x7a, 0x00, 0x27, 0x29, 0xac, 0x4d, 0xa7, 0x39, 0x57, 0x11, 0x7d, 0x00, 0xcd, 0xd4,
	0x35, 0xb0, 0xed, 0xe9, 0x60, 0x62, 0xe5, 0x33, 0xca, 0xe8, 0x08, 0x0e, 0xfa, 0x93, 0xf1, 0x85,
	0x79, 0xe9, 0x58, 0x36, 0x36, 0xc7, 0x97, 0x4a, 0x89, 0x57, 0x98, 0x41, 0xe6, 0xd8, 0x56, 0xf6,
	0x7a, 0x7f, 0x17, 0xa1, 0x61, 0xc7, 0xc4, 0xa5, 0xfd, 0xb5, 0x7c, 0xe8, 0x12, 0x20, 0x0a, 0x13,
	0x66, 0xdc, 0x89, 0x27, 0xe3, 0x61, 0x4e, 0xd8, 0xfc, 0x65, 0x6f, 0x36, 0xdf, 0xec, 0xe4, 0x0a,
	0xb5, 0x0b, 0x68, 0x00, 0x35, 0x9e, 0x68, 0x44, 0x59, 0xec, 0xb9, 0xef, 0x95, 0x29, 0x2b, 0xc9,
	0x62, 0x84, 0x2d, 0xdf, 0x2b, 0xd1, 0x77, 0x50, 0x9b, 0x53, 0xb6, 0xbe, 0x10, 0x28, 0xff, 0xb4,
	0xbf, 0x72, 0x57, 0x9b, 0x0f, 0xef, 0xf1, 0x66, 0xb9, 0xbe, 0x85, 0x52, 0xc4, 0x5f, 0xa1, 0xfc,
	0xea, 0xef, 0xec, 0xf1, 0xdb, 0x6b, 0xb9, 0x29, 0x8b, 0x5f, 0x86, 0xcf, 0xfe, 0x09, 0x00, 0x00,
	0xff, 0xff, 0x1a, 0x5e, 0xa4, 0x51, 0x45, 0x08, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TraceCollectorClient is the client API for TraceCollector service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TraceCollectorClient interface {
	PostEvents(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResult, error)
	PostMetrics(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResult, error)
	PostStatus(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResult, error)
	GetSettings(ctx context.Context, in *SettingsRequest, opts ...grpc.CallOption) (*SettingsResult, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*MessageResult, error)
}

type traceCollectorClient struct {
	cc *grpc.ClientConn
}

func NewTraceCollectorClient(cc *grpc.ClientConn) TraceCollectorClient {
	return &traceCollectorClient{cc}
}

func (c *traceCollectorClient) PostEvents(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResult, error) {
	out := new(MessageResult)
	err := c.cc.Invoke(ctx, "/collector.TraceCollector/postEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traceCollectorClient) PostMetrics(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResult, error) {
	out := new(MessageResult)
	err := c.cc.Invoke(ctx, "/collector.TraceCollector/postMetrics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traceCollectorClient) PostStatus(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResult, error) {
	out := new(MessageResult)
	err := c.cc.Invoke(ctx, "/collector.TraceCollector/postStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traceCollectorClient) GetSettings(ctx context.Context, in *SettingsRequest, opts ...grpc.CallOption) (*SettingsResult, error) {
	out := new(SettingsResult)
	err := c.cc.Invoke(ctx, "/collector.TraceCollector/getSettings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traceCollectorClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*MessageResult, error) {
	out := new(MessageResult)
	err := c.cc.Invoke(ctx, "/collector.TraceCollector/ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TraceCollectorServer is the server API for TraceCollector service.
type TraceCollectorServer interface {
	PostEvents(context.Context, *MessageRequest) (*MessageResult, error)
	PostMetrics(context.Context, *MessageRequest) (*MessageResult, error)
	PostStatus(context.Context, *MessageRequest) (*MessageResult, error)
	GetSettings(context.Context, *SettingsRequest) (*SettingsResult, error)
	Ping(context.Context, *PingRequest) (*MessageResult, error)
}

// UnimplementedTraceCollectorServer can be embedded to have forward compatible implementations.
type UnimplementedTraceCollectorServer struct {
}

func (*UnimplementedTraceCollectorServer) PostEvents(ctx context.Context, req *MessageRequest) (*MessageResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostEvents not implemented")
}
func (*UnimplementedTraceCollectorServer) PostMetrics(ctx context.Context, req *MessageRequest) (*MessageResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostMetrics not implemented")
}
func (*UnimplementedTraceCollectorServer) PostStatus(ctx context.Context, req *MessageRequest) (*MessageResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostStatus not implemented")
}
func (*UnimplementedTraceCollectorServer) GetSettings(ctx context.Context, req *SettingsRequest) (*SettingsResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSettings not implemented")
}
func (*UnimplementedTraceCollectorServer) Ping(ctx context.Context, req *PingRequest) (*MessageResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}

func RegisterTraceCollectorServer(s *grpc.Server, srv TraceCollectorServer) {
	s.RegisterService(&_TraceCollector_serviceDesc, srv)
}

func _TraceCollector_PostEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceCollectorServer).PostEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/collector.TraceCollector/PostEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceCollectorServer).PostEvents(ctx, req.(*MessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TraceCollector_PostMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceCollectorServer).PostMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/collector.TraceCollector/PostMetrics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceCollectorServer).PostMetrics(ctx, req.(*MessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TraceCollector_PostStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceCollectorServer).PostStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/collector.TraceCollector/PostStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceCollectorServer).PostStatus(ctx, req.(*MessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TraceCollector_GetSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SettingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceCollectorServer).GetSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/collector.TraceCollector/GetSettings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceCollectorServer).GetSettings(ctx, req.(*SettingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TraceCollector_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraceCollectorServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/collector.TraceCollector/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraceCollectorServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TraceCollector_serviceDesc = grpc.ServiceDesc{
	ServiceName: "collector.TraceCollector",
	HandlerType: (*TraceCollectorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "postEvents",
			Handler:    _TraceCollector_PostEvents_Handler,
		},
		{
			MethodName: "postMetrics",
			Handler:    _TraceCollector_PostMetrics_Handler,
		},
		{
			MethodName: "postStatus",
			Handler:    _TraceCollector_PostStatus_Handler,
		},
		{
			MethodName: "getSettings",
			Handler:    _TraceCollector_GetSettings_Handler,
		},
		{
			MethodName: "ping",
			Handler:    _TraceCollector_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "collector.proto",
}
