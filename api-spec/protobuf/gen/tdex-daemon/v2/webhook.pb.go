// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: tdex-daemon/v2/webhook.proto

package tdex_daemonv2

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

type AddWebhookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The endpoint of the webhook to call whenever the target event occurs.
	Endpoint string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	// The event for which the webhook endpoint should be called.
	Event WebhookEvent `protobuf:"varint,2,opt,name=event,proto3,enum=tdex_daemon.v2.WebhookEvent" json:"event,omitempty"`
	// The secret to use to generate an OAuth token for making authenticated
	// requests to the webhook endpoint.
	Secret string `protobuf:"bytes,3,opt,name=secret,proto3" json:"secret,omitempty"`
}

func (x *AddWebhookRequest) Reset() {
	*x = AddWebhookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tdex_daemon_v2_webhook_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddWebhookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddWebhookRequest) ProtoMessage() {}

func (x *AddWebhookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tdex_daemon_v2_webhook_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddWebhookRequest.ProtoReflect.Descriptor instead.
func (*AddWebhookRequest) Descriptor() ([]byte, []int) {
	return file_tdex_daemon_v2_webhook_proto_rawDescGZIP(), []int{0}
}

func (x *AddWebhookRequest) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *AddWebhookRequest) GetEvent() WebhookEvent {
	if x != nil {
		return x.Event
	}
	return WebhookEvent_WEBHOOK_EVENT_UNSPECIFIED
}

func (x *AddWebhookRequest) GetSecret() string {
	if x != nil {
		return x.Secret
	}
	return ""
}

type AddWebhookResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The id of the new webhook.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *AddWebhookResponse) Reset() {
	*x = AddWebhookResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tdex_daemon_v2_webhook_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddWebhookResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddWebhookResponse) ProtoMessage() {}

func (x *AddWebhookResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tdex_daemon_v2_webhook_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddWebhookResponse.ProtoReflect.Descriptor instead.
func (*AddWebhookResponse) Descriptor() ([]byte, []int) {
	return file_tdex_daemon_v2_webhook_proto_rawDescGZIP(), []int{1}
}

func (x *AddWebhookResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type RemoveWebhookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The id of the webhook to remove.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RemoveWebhookRequest) Reset() {
	*x = RemoveWebhookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tdex_daemon_v2_webhook_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveWebhookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveWebhookRequest) ProtoMessage() {}

func (x *RemoveWebhookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tdex_daemon_v2_webhook_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveWebhookRequest.ProtoReflect.Descriptor instead.
func (*RemoveWebhookRequest) Descriptor() ([]byte, []int) {
	return file_tdex_daemon_v2_webhook_proto_rawDescGZIP(), []int{2}
}

func (x *RemoveWebhookRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type RemoveWebhookResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RemoveWebhookResponse) Reset() {
	*x = RemoveWebhookResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tdex_daemon_v2_webhook_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveWebhookResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveWebhookResponse) ProtoMessage() {}

func (x *RemoveWebhookResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tdex_daemon_v2_webhook_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveWebhookResponse.ProtoReflect.Descriptor instead.
func (*RemoveWebhookResponse) Descriptor() ([]byte, []int) {
	return file_tdex_daemon_v2_webhook_proto_rawDescGZIP(), []int{3}
}

type ListWebhooksRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Filter the list of webhooks by event.
	Event WebhookEvent `protobuf:"varint,1,opt,name=event,proto3,enum=tdex_daemon.v2.WebhookEvent" json:"event,omitempty"`
}

func (x *ListWebhooksRequest) Reset() {
	*x = ListWebhooksRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tdex_daemon_v2_webhook_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListWebhooksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListWebhooksRequest) ProtoMessage() {}

func (x *ListWebhooksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tdex_daemon_v2_webhook_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListWebhooksRequest.ProtoReflect.Descriptor instead.
func (*ListWebhooksRequest) Descriptor() ([]byte, []int) {
	return file_tdex_daemon_v2_webhook_proto_rawDescGZIP(), []int{4}
}

func (x *ListWebhooksRequest) GetEvent() WebhookEvent {
	if x != nil {
		return x.Event
	}
	return WebhookEvent_WEBHOOK_EVENT_UNSPECIFIED
}

type ListWebhooksResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The list of info about the webhooks.
	WebhookInfo []*WebhookInfo `protobuf:"bytes,1,rep,name=webhook_info,json=webhookInfo,proto3" json:"webhook_info,omitempty"`
}

func (x *ListWebhooksResponse) Reset() {
	*x = ListWebhooksResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tdex_daemon_v2_webhook_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListWebhooksResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListWebhooksResponse) ProtoMessage() {}

func (x *ListWebhooksResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tdex_daemon_v2_webhook_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListWebhooksResponse.ProtoReflect.Descriptor instead.
func (*ListWebhooksResponse) Descriptor() ([]byte, []int) {
	return file_tdex_daemon_v2_webhook_proto_rawDescGZIP(), []int{5}
}

func (x *ListWebhooksResponse) GetWebhookInfo() []*WebhookInfo {
	if x != nil {
		return x.WebhookInfo
	}
	return nil
}

var File_tdex_daemon_v2_webhook_proto protoreflect.FileDescriptor

var file_tdex_daemon_v2_webhook_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x74, 0x64, 0x65, 0x78, 0x2d, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32,
	0x2f, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e,
	0x74, 0x64, 0x65, 0x78, 0x5f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x32, 0x1a, 0x1a,
	0x74, 0x64, 0x65, 0x78, 0x2d, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7b, 0x0a, 0x11, 0x41, 0x64,
	0x64, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x32, 0x0a, 0x05, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x74, 0x64, 0x65,
	0x78, 0x5f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x32, 0x2e, 0x57, 0x65, 0x62, 0x68,
	0x6f, 0x6f, 0x6b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x22, 0x24, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x57, 0x65,
	0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x26, 0x0a,
	0x14, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x17, 0x0a, 0x15, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x57,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x49,
	0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x32, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x74, 0x64, 0x65, 0x78, 0x5f, 0x64, 0x61, 0x65, 0x6d,
	0x6f, 0x6e, 0x2e, 0x76, 0x32, 0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x56, 0x0a, 0x14, 0x4c, 0x69, 0x73,
	0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x3e, 0x0a, 0x0c, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x69, 0x6e, 0x66,
	0x6f, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x64, 0x65, 0x78, 0x5f, 0x64,
	0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x32, 0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0b, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x49, 0x6e, 0x66,
	0x6f, 0x32, 0xa4, 0x02, 0x0a, 0x0e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x55, 0x0a, 0x0a, 0x41, 0x64, 0x64, 0x57, 0x65, 0x62, 0x68, 0x6f,
	0x6f, 0x6b, 0x12, 0x21, 0x2e, 0x74, 0x64, 0x65, 0x78, 0x5f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e,
	0x2e, 0x76, 0x32, 0x2e, 0x41, 0x64, 0x64, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x74, 0x64, 0x65, 0x78, 0x5f, 0x64, 0x61, 0x65,
	0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x32, 0x2e, 0x41, 0x64, 0x64, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f,
	0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5e, 0x0a, 0x0d, 0x52,
	0x65, 0x6d, 0x6f, 0x76, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x12, 0x24, 0x2e, 0x74,
	0x64, 0x65, 0x78, 0x5f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x32, 0x2e, 0x52, 0x65,
	0x6d, 0x6f, 0x76, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x25, 0x2e, 0x74, 0x64, 0x65, 0x78, 0x5f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e,
	0x2e, 0x76, 0x32, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f,
	0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5b, 0x0a, 0x0c, 0x4c,
	0x69, 0x73, 0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x12, 0x23, 0x2e, 0x74, 0x64,
	0x65, 0x78, 0x5f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x32, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x24, 0x2e, 0x74, 0x64, 0x65, 0x78, 0x5f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x76,
	0x32, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0xcf, 0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d,
	0x2e, 0x74, 0x64, 0x65, 0x78, 0x5f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x32, 0x42,
	0x0c, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x56, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x64, 0x65, 0x78,
	0x2d, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x74, 0x64, 0x65, 0x78, 0x2d, 0x64, 0x61,
	0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2d, 0x73, 0x70, 0x65, 0x63, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x74, 0x64, 0x65, 0x78, 0x2d,
	0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x3b, 0x74, 0x64, 0x65, 0x78, 0x5f, 0x64,
	0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x76, 0x32, 0xa2, 0x02, 0x03, 0x54, 0x58, 0x58, 0xaa, 0x02, 0x0d,
	0x54, 0x64, 0x65, 0x78, 0x44, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x56, 0x32, 0xca, 0x02, 0x0d,
	0x54, 0x64, 0x65, 0x78, 0x44, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x5c, 0x56, 0x32, 0xe2, 0x02, 0x19,
	0x54, 0x64, 0x65, 0x78, 0x44, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x5c, 0x56, 0x32, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0e, 0x54, 0x64, 0x65, 0x78,
	0x44, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x3a, 0x3a, 0x56, 0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_tdex_daemon_v2_webhook_proto_rawDescOnce sync.Once
	file_tdex_daemon_v2_webhook_proto_rawDescData = file_tdex_daemon_v2_webhook_proto_rawDesc
)

func file_tdex_daemon_v2_webhook_proto_rawDescGZIP() []byte {
	file_tdex_daemon_v2_webhook_proto_rawDescOnce.Do(func() {
		file_tdex_daemon_v2_webhook_proto_rawDescData = protoimpl.X.CompressGZIP(file_tdex_daemon_v2_webhook_proto_rawDescData)
	})
	return file_tdex_daemon_v2_webhook_proto_rawDescData
}

var file_tdex_daemon_v2_webhook_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_tdex_daemon_v2_webhook_proto_goTypes = []interface{}{
	(*AddWebhookRequest)(nil),     // 0: tdex_daemon.v2.AddWebhookRequest
	(*AddWebhookResponse)(nil),    // 1: tdex_daemon.v2.AddWebhookResponse
	(*RemoveWebhookRequest)(nil),  // 2: tdex_daemon.v2.RemoveWebhookRequest
	(*RemoveWebhookResponse)(nil), // 3: tdex_daemon.v2.RemoveWebhookResponse
	(*ListWebhooksRequest)(nil),   // 4: tdex_daemon.v2.ListWebhooksRequest
	(*ListWebhooksResponse)(nil),  // 5: tdex_daemon.v2.ListWebhooksResponse
	(WebhookEvent)(0),             // 6: tdex_daemon.v2.WebhookEvent
	(*WebhookInfo)(nil),           // 7: tdex_daemon.v2.WebhookInfo
}
var file_tdex_daemon_v2_webhook_proto_depIdxs = []int32{
	6, // 0: tdex_daemon.v2.AddWebhookRequest.event:type_name -> tdex_daemon.v2.WebhookEvent
	6, // 1: tdex_daemon.v2.ListWebhooksRequest.event:type_name -> tdex_daemon.v2.WebhookEvent
	7, // 2: tdex_daemon.v2.ListWebhooksResponse.webhook_info:type_name -> tdex_daemon.v2.WebhookInfo
	0, // 3: tdex_daemon.v2.WebhookService.AddWebhook:input_type -> tdex_daemon.v2.AddWebhookRequest
	2, // 4: tdex_daemon.v2.WebhookService.RemoveWebhook:input_type -> tdex_daemon.v2.RemoveWebhookRequest
	4, // 5: tdex_daemon.v2.WebhookService.ListWebhooks:input_type -> tdex_daemon.v2.ListWebhooksRequest
	1, // 6: tdex_daemon.v2.WebhookService.AddWebhook:output_type -> tdex_daemon.v2.AddWebhookResponse
	3, // 7: tdex_daemon.v2.WebhookService.RemoveWebhook:output_type -> tdex_daemon.v2.RemoveWebhookResponse
	5, // 8: tdex_daemon.v2.WebhookService.ListWebhooks:output_type -> tdex_daemon.v2.ListWebhooksResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_tdex_daemon_v2_webhook_proto_init() }
func file_tdex_daemon_v2_webhook_proto_init() {
	if File_tdex_daemon_v2_webhook_proto != nil {
		return
	}
	file_tdex_daemon_v2_types_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_tdex_daemon_v2_webhook_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddWebhookRequest); i {
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
		file_tdex_daemon_v2_webhook_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddWebhookResponse); i {
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
		file_tdex_daemon_v2_webhook_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveWebhookRequest); i {
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
		file_tdex_daemon_v2_webhook_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveWebhookResponse); i {
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
		file_tdex_daemon_v2_webhook_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListWebhooksRequest); i {
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
		file_tdex_daemon_v2_webhook_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListWebhooksResponse); i {
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
			RawDescriptor: file_tdex_daemon_v2_webhook_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tdex_daemon_v2_webhook_proto_goTypes,
		DependencyIndexes: file_tdex_daemon_v2_webhook_proto_depIdxs,
		MessageInfos:      file_tdex_daemon_v2_webhook_proto_msgTypes,
	}.Build()
	File_tdex_daemon_v2_webhook_proto = out.File
	file_tdex_daemon_v2_webhook_proto_rawDesc = nil
	file_tdex_daemon_v2_webhook_proto_goTypes = nil
	file_tdex_daemon_v2_webhook_proto_depIdxs = nil
}
