// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: twitter.proto

package twitter

import (
	models "github.com/twitter/models"
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

var File_twitter_proto protoreflect.FileDescriptor

var file_twitter_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x74, 0x77, 0x69, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x74, 0x77, 0x69, 0x74, 0x74, 0x65, 0x72, 0x1a, 0x0c, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xc8, 0x04, 0x0a, 0x07, 0x54, 0x77, 0x69, 0x74, 0x74,
	0x65, 0x72, 0x12, 0x2b, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x12, 0x0d, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x0d, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12,
	0x2a, 0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x12,
	0x0c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x0c, 0x2e,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x27, 0x0a, 0x09, 0x4c,
	0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x0c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x12, 0x29, 0x0a, 0x0a, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x55, 0x73,
	0x65, 0x72, 0x12, 0x0c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x1a, 0x0d, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12,
	0x2b, 0x0a, 0x0c, 0x55, 0x6e, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x55, 0x73, 0x65, 0x72, 0x12,
	0x0c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x0d, 0x2e,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x28, 0x0a, 0x0a,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x0c, 0x2e, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x73, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x46, 0x65, 0x65,
	0x64, 0x12, 0x0d, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x15, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70,
	0x6c, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x12, 0x29, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x0c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x50,
	0x6f, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x25, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0c, 0x2e,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x0c, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x33, 0x0a, 0x0e, 0x47, 0x65, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x0c, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x13, 0x2e, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x26,
	0x0a, 0x07, 0x47, 0x65, 0x74, 0x53, 0x65, 0x6c, 0x66, 0x12, 0x0d, 0x2e, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x32, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x4d, 0x79, 0x50,
	0x6f, 0x73, 0x74, 0x73, 0x12, 0x0d, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x15, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x4d, 0x75, 0x6c,
	0x74, 0x69, 0x70, 0x6c, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x12, 0x25, 0x0a, 0x07, 0x47, 0x65,
	0x74, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x0c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x50,
	0x6f, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x50, 0x6f, 0x73,
	0x74, 0x42, 0x1c, 0x5a, 0x1a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x74, 0x77, 0x69, 0x74, 0x74, 0x65, 0x72, 0x2f, 0x74, 0x77, 0x69, 0x74, 0x74, 0x65, 0x72, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_twitter_proto_goTypes = []interface{}{
	(*models.Empty)(nil),         // 0: models.Empty
	(*models.User)(nil),          // 1: models.User
	(*models.Post)(nil),          // 2: models.Post
	(*models.MultiplePosts)(nil), // 3: models.MultiplePosts
	(*models.UserProfile)(nil),   // 4: models.UserProfile
}
var file_twitter_proto_depIdxs = []int32{
	0,  // 0: twitter.Twitter.HealthCheck:input_type -> models.Empty
	1,  // 1: twitter.Twitter.RegisterUser:input_type -> models.User
	1,  // 2: twitter.Twitter.LoginUser:input_type -> models.User
	1,  // 3: twitter.Twitter.FollowUser:input_type -> models.User
	1,  // 4: twitter.Twitter.UnFollowUser:input_type -> models.User
	2,  // 5: twitter.Twitter.CreatePost:input_type -> models.Post
	0,  // 6: twitter.Twitter.GetFeed:input_type -> models.Empty
	2,  // 7: twitter.Twitter.DeletePost:input_type -> models.Post
	1,  // 8: twitter.Twitter.GetUser:input_type -> models.User
	1,  // 9: twitter.Twitter.GetUserProfile:input_type -> models.User
	0,  // 10: twitter.Twitter.GetSelf:input_type -> models.Empty
	0,  // 11: twitter.Twitter.GetMyPosts:input_type -> models.Empty
	2,  // 12: twitter.Twitter.GetPost:input_type -> models.Post
	0,  // 13: twitter.Twitter.HealthCheck:output_type -> models.Empty
	1,  // 14: twitter.Twitter.RegisterUser:output_type -> models.User
	1,  // 15: twitter.Twitter.LoginUser:output_type -> models.User
	0,  // 16: twitter.Twitter.FollowUser:output_type -> models.Empty
	0,  // 17: twitter.Twitter.UnFollowUser:output_type -> models.Empty
	2,  // 18: twitter.Twitter.CreatePost:output_type -> models.Post
	3,  // 19: twitter.Twitter.GetFeed:output_type -> models.MultiplePosts
	0,  // 20: twitter.Twitter.DeletePost:output_type -> models.Empty
	1,  // 21: twitter.Twitter.GetUser:output_type -> models.User
	4,  // 22: twitter.Twitter.GetUserProfile:output_type -> models.UserProfile
	1,  // 23: twitter.Twitter.GetSelf:output_type -> models.User
	3,  // 24: twitter.Twitter.GetMyPosts:output_type -> models.MultiplePosts
	2,  // 25: twitter.Twitter.GetPost:output_type -> models.Post
	13, // [13:26] is the sub-list for method output_type
	0,  // [0:13] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_twitter_proto_init() }
func file_twitter_proto_init() {
	if File_twitter_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_twitter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_twitter_proto_goTypes,
		DependencyIndexes: file_twitter_proto_depIdxs,
	}.Build()
	File_twitter_proto = out.File
	file_twitter_proto_rawDesc = nil
	file_twitter_proto_goTypes = nil
	file_twitter_proto_depIdxs = nil
}
