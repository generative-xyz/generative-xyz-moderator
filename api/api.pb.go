// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.7
// source: api.proto

package api

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_api_proto protoreflect.FileDescriptor

var file_api_proto_rawDesc = []byte{
	0x0a, 0x09, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x61, 0x70, 0x69,
	0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c,
	0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0e, 0x74, 0x65,
	0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x72, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x65, 0x64, 0x2d, 0x6e, 0x66, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x32, 0x82, 0x0f, 0x0a, 0x0a, 0x41, 0x70, 0x69, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x64, 0x0a, 0x04, 0x4c, 0x69, 0x76, 0x65, 0x12, 0x20, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x4c, 0x69,
	0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e,
	0x4c, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x11, 0x12, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68,
	0x2f, 0x6c, 0x69, 0x76, 0x65, 0x12, 0x5d, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x20, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x68, 0x75, 0x62,
	0x2e, 0x69, 0x6f, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x21, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x68,
	0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x10, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0a, 0x12, 0x08, 0x2f, 0x76, 0x31, 0x2f,
	0x70, 0x69, 0x6e, 0x67, 0x12, 0x76, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x54, 0x65, 0x6d, 0x70, 0x6c,
	0x61, 0x74, 0x65, 0x12, 0x27, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x69, 0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x65, 0x6d,
	0x70, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e,
	0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x12, 0x0c,
	0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x8f, 0x01, 0x0a,
	0x11, 0x47, 0x65, 0x74, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x44, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x12, 0x2d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69,
	0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x65, 0x6d, 0x70,
	0x6c, 0x61, 0x74, 0x65, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e,
	0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x65, 0x6d, 0x70, 0x6c,
	0x61, 0x74, 0x65, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x12, 0x13, 0x2f, 0x76, 0x31, 0x2f, 0x74,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2d, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x92,
	0x01, 0x0a, 0x11, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x69, 0x6e, 0x67, 0x12, 0x2d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x69, 0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x54, 0x65, 0x6d, 0x70, 0x6c,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x69, 0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x22, 0x13, 0x2f, 0x76, 0x31,
	0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2d, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x3a, 0x01, 0x2a, 0x12, 0xb9, 0x01, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x65, 0x64, 0x4e, 0x66, 0x74, 0x12, 0x2a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74,
	0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x64, 0x4e, 0x66, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69,
	0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x65, 0x64, 0x4e, 0x66, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x4e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x48, 0x12, 0x46, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x65, 0x64, 0x2d, 0x6e, 0x66, 0x74, 0x2f, 0x7b, 0x63, 0x68, 0x61, 0x69, 0x6e,
	0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x7b, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x7d, 0x2f, 0x7b, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x7b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x7d, 0x12,
	0xa8, 0x01, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x43, 0x61, 0x6e, 0x64, 0x79, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x2c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x69, 0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x61,
	0x6e, 0x64, 0x79, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69,
	0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x6e, 0x64,
	0x79, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x37, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x31, 0x12, 0x2f, 0x2f, 0x76, 0x31, 0x2f, 0x72,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x64, 0x2d, 0x6e, 0x66, 0x74, 0x2f, 0x7b, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2f,
	0x7b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x7d, 0x12, 0xaa, 0x01, 0x0a, 0x11, 0x47,
	0x65, 0x74, 0x43, 0x61, 0x6e, 0x64, 0x79, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x73,
	0x12, 0x2d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67,
	0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x6e, 0x64, 0x79, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x68,
	0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x6e, 0x64, 0x79, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x36, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x30, 0x12, 0x2e, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x77, 0x65,
	0x65, 0x74, 0x2d, 0x6e, 0x66, 0x74, 0x2f, 0x7b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64,
	0x7d, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x73, 0x2f, 0x7b, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x5f, 0x69, 0x64, 0x73, 0x7d, 0x12, 0xaf, 0x01, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x43,
	0x61, 0x6e, 0x64, 0x79, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x73, 0x74,
	0x12, 0x2c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67,
	0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x6e, 0x64, 0x79, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x68, 0x75,
	0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x6e, 0x64, 0x79, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3a, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x34, 0x22, 0x2f, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x65, 0x64, 0x2d, 0x6e, 0x66, 0x74, 0x2f, 0x7b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69,
	0x64, 0x7d, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x7b, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x5f, 0x69, 0x64, 0x7d, 0x3a, 0x01, 0x2a, 0x12, 0xc0, 0x01, 0x0a, 0x12, 0x47, 0x65,
	0x74, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x64, 0x4e, 0x66, 0x74, 0x50, 0x6f, 0x73, 0x74,
	0x12, 0x2a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67,
	0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x65, 0x64, 0x4e, 0x66, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e,
	0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x64, 0x4e, 0x66,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x51, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x4b, 0x22, 0x46, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x64, 0x2d,
	0x6e, 0x66, 0x74, 0x2f, 0x7b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x7b,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x7d, 0x2f, 0x7b, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x7b,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x7d, 0x3a, 0x01, 0x2a, 0x12, 0xbe, 0x01, 0x0a,
	0x11, 0x47, 0x65, 0x74, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x2d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69,
	0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x76, 0x61, 0x74,
	0x61, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e,
	0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x76, 0x61, 0x74, 0x61,
	0x72, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x4a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x44, 0x12, 0x42, 0x2f, 0x76, 0x31, 0x2f, 0x72,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x64, 0x2d, 0x6e, 0x66, 0x74, 0x2f, 0x7b, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x7b, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x7d, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x2f, 0x7b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x7d, 0x12, 0xc5, 0x01,
	0x0a, 0x15, 0x47, 0x65, 0x74, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x2d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65,
	0x74, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74,
	0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x4d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x47, 0x22, 0x42,
	0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x64, 0x2d, 0x6e, 0x66, 0x74,
	0x2f, 0x7b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x7b, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x7d, 0x2f, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x7b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69,
	0x64, 0x7d, 0x3a, 0x01, 0x2a, 0x42, 0x19, 0x5a, 0x17, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69,
	0x6e, 0x67, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x3b, 0x61, 0x70, 0x69,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_api_proto_goTypes = []interface{}{
	(*LiveRequest)(nil),               // 0: api.renderinghub.io.LiveRequest
	(*PingRequest)(nil),               // 1: api.renderinghub.io.PingRequest
	(*GetTemplateRequest)(nil),        // 2: api.renderinghub.io.GetTemplateRequest
	(*GetTemplateDetailRequest)(nil),  // 3: api.renderinghub.io.GetTemplateDetailRequest
	(*TemplateRenderingRequest)(nil),  // 4: api.renderinghub.io.TemplateRenderingRequest
	(*GetRenderedNftRequest)(nil),     // 5: api.renderinghub.io.GetRenderedNftRequest
	(*GetCandyMetadataRequest)(nil),   // 6: api.renderinghub.io.GetCandyMetadataRequest
	(*GetCandyMetadatasRequest)(nil),  // 7: api.renderinghub.io.GetCandyMetadatasRequest
	(*GetAvatarMetadataRequest)(nil),  // 8: api.renderinghub.io.GetAvatarMetadataRequest
	(*LiveResponse)(nil),              // 9: api.renderinghub.io.LiveResponse
	(*PingResponse)(nil),              // 10: api.renderinghub.io.PingResponse
	(*GetTemplateResponse)(nil),       // 11: api.renderinghub.io.GetTemplateResponse
	(*GetTemplateDetailResponse)(nil), // 12: api.renderinghub.io.GetTemplateDetailResponse
	(*TemplateRenderingResponse)(nil), // 13: api.renderinghub.io.TemplateRenderingResponse
	(*GetRenderedNftResponse)(nil),    // 14: api.renderinghub.io.GetRenderedNftResponse
	(*GetCandyMetadataResponse)(nil),  // 15: api.renderinghub.io.GetCandyMetadataResponse
	(*GetCandyMetadatasResponse)(nil), // 16: api.renderinghub.io.GetCandyMetadatasResponse
	(*GetAvatarMetadataResponse)(nil), // 17: api.renderinghub.io.GetAvatarMetadataResponse
}
var file_api_proto_depIdxs = []int32{
	0,  // 0: api.renderinghub.io.ApiService.Live:input_type -> api.renderinghub.io.LiveRequest
	1,  // 1: api.renderinghub.io.ApiService.Ping:input_type -> api.renderinghub.io.PingRequest
	2,  // 2: api.renderinghub.io.ApiService.GetTemplate:input_type -> api.renderinghub.io.GetTemplateRequest
	3,  // 3: api.renderinghub.io.ApiService.GetTemplateDetail:input_type -> api.renderinghub.io.GetTemplateDetailRequest
	4,  // 4: api.renderinghub.io.ApiService.TemplateRendering:input_type -> api.renderinghub.io.TemplateRenderingRequest
	5,  // 5: api.renderinghub.io.ApiService.GetRenderedNft:input_type -> api.renderinghub.io.GetRenderedNftRequest
	6,  // 6: api.renderinghub.io.ApiService.GetCandyMetadata:input_type -> api.renderinghub.io.GetCandyMetadataRequest
	7,  // 7: api.renderinghub.io.ApiService.GetCandyMetadatas:input_type -> api.renderinghub.io.GetCandyMetadatasRequest
	6,  // 8: api.renderinghub.io.ApiService.GetCandyMetadataPost:input_type -> api.renderinghub.io.GetCandyMetadataRequest
	5,  // 9: api.renderinghub.io.ApiService.GetRenderedNftPost:input_type -> api.renderinghub.io.GetRenderedNftRequest
	8,  // 10: api.renderinghub.io.ApiService.GetAvatarMetadata:input_type -> api.renderinghub.io.GetAvatarMetadataRequest
	8,  // 11: api.renderinghub.io.ApiService.GetAvatarMetadataPost:input_type -> api.renderinghub.io.GetAvatarMetadataRequest
	9,  // 12: api.renderinghub.io.ApiService.Live:output_type -> api.renderinghub.io.LiveResponse
	10, // 13: api.renderinghub.io.ApiService.Ping:output_type -> api.renderinghub.io.PingResponse
	11, // 14: api.renderinghub.io.ApiService.GetTemplate:output_type -> api.renderinghub.io.GetTemplateResponse
	12, // 15: api.renderinghub.io.ApiService.GetTemplateDetail:output_type -> api.renderinghub.io.GetTemplateDetailResponse
	13, // 16: api.renderinghub.io.ApiService.TemplateRendering:output_type -> api.renderinghub.io.TemplateRenderingResponse
	14, // 17: api.renderinghub.io.ApiService.GetRenderedNft:output_type -> api.renderinghub.io.GetRenderedNftResponse
	15, // 18: api.renderinghub.io.ApiService.GetCandyMetadata:output_type -> api.renderinghub.io.GetCandyMetadataResponse
	16, // 19: api.renderinghub.io.ApiService.GetCandyMetadatas:output_type -> api.renderinghub.io.GetCandyMetadatasResponse
	15, // 20: api.renderinghub.io.ApiService.GetCandyMetadataPost:output_type -> api.renderinghub.io.GetCandyMetadataResponse
	14, // 21: api.renderinghub.io.ApiService.GetRenderedNftPost:output_type -> api.renderinghub.io.GetRenderedNftResponse
	17, // 22: api.renderinghub.io.ApiService.GetAvatarMetadata:output_type -> api.renderinghub.io.GetAvatarMetadataResponse
	17, // 23: api.renderinghub.io.ApiService.GetAvatarMetadataPost:output_type -> api.renderinghub.io.GetAvatarMetadataResponse
	12, // [12:24] is the sub-list for method output_type
	0,  // [0:12] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_api_proto_init() }
func file_api_proto_init() {
	if File_api_proto != nil {
		return
	}
	file_health_proto_init()
	file_template_proto_init()
	file_rendered_nft_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_goTypes,
		DependencyIndexes: file_api_proto_depIdxs,
	}.Build()
	File_api_proto = out.File
	file_api_proto_rawDesc = nil
	file_api_proto_goTypes = nil
	file_api_proto_depIdxs = nil
}
