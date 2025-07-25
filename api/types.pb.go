// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: types.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// RelationTuple represents a relationship between a user and an object
// Following Zanzibar's tuple format: <object>#<relation>@<user>
type RelationTuple struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The namespace that defines the schema for this tuple
	Namespace string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// The ID of the object in the relationship
	ObjectId string `protobuf:"bytes,2,opt,name=object_id,json=objectId,proto3" json:"object_id,omitempty"`
	// The relation name (e.g., "viewer", "editor", "owner")
	Relation string `protobuf:"bytes,3,opt,name=relation,proto3" json:"relation,omitempty"`
	// The user ID in the relationship
	UserId string `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// Optional userset for indirect relationships (e.g., "group:eng#member")
	Userset string `protobuf:"bytes,5,opt,name=userset,proto3" json:"userset,omitempty"`
	// When this tuple was created
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// When this tuple was last updated
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RelationTuple) Reset() {
	*x = RelationTuple{}
	mi := &file_types_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RelationTuple) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelationTuple) ProtoMessage() {}

func (x *RelationTuple) ProtoReflect() protoreflect.Message {
	mi := &file_types_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelationTuple.ProtoReflect.Descriptor instead.
func (*RelationTuple) Descriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{0}
}

func (x *RelationTuple) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *RelationTuple) GetObjectId() string {
	if x != nil {
		return x.ObjectId
	}
	return ""
}

func (x *RelationTuple) GetRelation() string {
	if x != nil {
		return x.Relation
	}
	return ""
}

func (x *RelationTuple) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *RelationTuple) GetUserset() string {
	if x != nil {
		return x.Userset
	}
	return ""
}

func (x *RelationTuple) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *RelationTuple) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

// NamespaceConfig defines the schema and rules for a namespace
type NamespaceConfig struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The namespace name (e.g., "documents", "folders")
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Relations defined for this namespace
	Relations []*RelationConfig `protobuf:"bytes,2,rep,name=relations,proto3" json:"relations,omitempty"`
	// When this namespace was created
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// When this namespace was last updated
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NamespaceConfig) Reset() {
	*x = NamespaceConfig{}
	mi := &file_types_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NamespaceConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NamespaceConfig) ProtoMessage() {}

func (x *NamespaceConfig) ProtoReflect() protoreflect.Message {
	mi := &file_types_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NamespaceConfig.ProtoReflect.Descriptor instead.
func (*NamespaceConfig) Descriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{1}
}

func (x *NamespaceConfig) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *NamespaceConfig) GetRelations() []*RelationConfig {
	if x != nil {
		return x.Relations
	}
	return nil
}

func (x *NamespaceConfig) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *NamespaceConfig) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

// RelationConfig defines a single relation within a namespace
type RelationConfig struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The relation name (e.g., "viewer", "editor")
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Userset rewrite rules in JSON format
	// Defines how this relation can be computed from other relations
	RewriteRules string `protobuf:"bytes,2,opt,name=rewrite_rules,json=rewriteRules,proto3" json:"rewrite_rules,omitempty"`
	// Optional description of this relation
	Description   string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RelationConfig) Reset() {
	*x = RelationConfig{}
	mi := &file_types_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RelationConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelationConfig) ProtoMessage() {}

func (x *RelationConfig) ProtoReflect() protoreflect.Message {
	mi := &file_types_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelationConfig.ProtoReflect.Descriptor instead.
func (*RelationConfig) Descriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{2}
}

func (x *RelationConfig) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RelationConfig) GetRewriteRules() string {
	if x != nil {
		return x.RewriteRules
	}
	return ""
}

func (x *RelationConfig) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// UserSet represents a set of users that can be computed
type UserSet struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Userset:
	//
	//	*UserSet_UserId
	//	*UserSet_ObjectRelation
	//	*UserSet_Union
	//	*UserSet_Intersection
	//	*UserSet_Exclusion
	Userset       isUserSet_Userset `protobuf_oneof:"userset"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserSet) Reset() {
	*x = UserSet{}
	mi := &file_types_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserSet) ProtoMessage() {}

func (x *UserSet) ProtoReflect() protoreflect.Message {
	mi := &file_types_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserSet.ProtoReflect.Descriptor instead.
func (*UserSet) Descriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{3}
}

func (x *UserSet) GetUserset() isUserSet_Userset {
	if x != nil {
		return x.Userset
	}
	return nil
}

func (x *UserSet) GetUserId() string {
	if x != nil {
		if x, ok := x.Userset.(*UserSet_UserId); ok {
			return x.UserId
		}
	}
	return ""
}

func (x *UserSet) GetObjectRelation() *ObjectRelation {
	if x != nil {
		if x, ok := x.Userset.(*UserSet_ObjectRelation); ok {
			return x.ObjectRelation
		}
	}
	return nil
}

func (x *UserSet) GetUnion() *UserSetUnion {
	if x != nil {
		if x, ok := x.Userset.(*UserSet_Union); ok {
			return x.Union
		}
	}
	return nil
}

func (x *UserSet) GetIntersection() *UserSetIntersection {
	if x != nil {
		if x, ok := x.Userset.(*UserSet_Intersection); ok {
			return x.Intersection
		}
	}
	return nil
}

func (x *UserSet) GetExclusion() *UserSetExclusion {
	if x != nil {
		if x, ok := x.Userset.(*UserSet_Exclusion); ok {
			return x.Exclusion
		}
	}
	return nil
}

type isUserSet_Userset interface {
	isUserSet_Userset()
}

type UserSet_UserId struct {
	// Direct user reference
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3,oneof"`
}

type UserSet_ObjectRelation struct {
	// Reference to users with a specific relation to an object
	ObjectRelation *ObjectRelation `protobuf:"bytes,2,opt,name=object_relation,json=objectRelation,proto3,oneof"`
}

type UserSet_Union struct {
	// Union of multiple usersets
	Union *UserSetUnion `protobuf:"bytes,3,opt,name=union,proto3,oneof"`
}

type UserSet_Intersection struct {
	// Intersection of multiple usersets
	Intersection *UserSetIntersection `protobuf:"bytes,4,opt,name=intersection,proto3,oneof"`
}

type UserSet_Exclusion struct {
	// Exclusion of one userset from another
	Exclusion *UserSetExclusion `protobuf:"bytes,5,opt,name=exclusion,proto3,oneof"`
}

func (*UserSet_UserId) isUserSet_Userset() {}

func (*UserSet_ObjectRelation) isUserSet_Userset() {}

func (*UserSet_Union) isUserSet_Userset() {}

func (*UserSet_Intersection) isUserSet_Userset() {}

func (*UserSet_Exclusion) isUserSet_Userset() {}

// ObjectRelation represents users with a relation to an object
type ObjectRelation struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The namespace of the object
	Namespace string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// The object ID
	ObjectId string `protobuf:"bytes,2,opt,name=object_id,json=objectId,proto3" json:"object_id,omitempty"`
	// The relation name
	Relation      string `protobuf:"bytes,3,opt,name=relation,proto3" json:"relation,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ObjectRelation) Reset() {
	*x = ObjectRelation{}
	mi := &file_types_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ObjectRelation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ObjectRelation) ProtoMessage() {}

func (x *ObjectRelation) ProtoReflect() protoreflect.Message {
	mi := &file_types_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ObjectRelation.ProtoReflect.Descriptor instead.
func (*ObjectRelation) Descriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{4}
}

func (x *ObjectRelation) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *ObjectRelation) GetObjectId() string {
	if x != nil {
		return x.ObjectId
	}
	return ""
}

func (x *ObjectRelation) GetRelation() string {
	if x != nil {
		return x.Relation
	}
	return ""
}

// UserSetUnion represents the union of multiple usersets
type UserSetUnion struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Children      []*UserSet             `protobuf:"bytes,1,rep,name=children,proto3" json:"children,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserSetUnion) Reset() {
	*x = UserSetUnion{}
	mi := &file_types_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserSetUnion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserSetUnion) ProtoMessage() {}

func (x *UserSetUnion) ProtoReflect() protoreflect.Message {
	mi := &file_types_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserSetUnion.ProtoReflect.Descriptor instead.
func (*UserSetUnion) Descriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{5}
}

func (x *UserSetUnion) GetChildren() []*UserSet {
	if x != nil {
		return x.Children
	}
	return nil
}

// UserSetIntersection represents the intersection of multiple usersets
type UserSetIntersection struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Children      []*UserSet             `protobuf:"bytes,1,rep,name=children,proto3" json:"children,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserSetIntersection) Reset() {
	*x = UserSetIntersection{}
	mi := &file_types_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserSetIntersection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserSetIntersection) ProtoMessage() {}

func (x *UserSetIntersection) ProtoReflect() protoreflect.Message {
	mi := &file_types_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserSetIntersection.ProtoReflect.Descriptor instead.
func (*UserSetIntersection) Descriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{6}
}

func (x *UserSetIntersection) GetChildren() []*UserSet {
	if x != nil {
		return x.Children
	}
	return nil
}

// UserSetExclusion represents exclusion of one userset from another
type UserSetExclusion struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Base          *UserSet               `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	Exclude       *UserSet               `protobuf:"bytes,2,opt,name=exclude,proto3" json:"exclude,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserSetExclusion) Reset() {
	*x = UserSetExclusion{}
	mi := &file_types_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserSetExclusion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserSetExclusion) ProtoMessage() {}

func (x *UserSetExclusion) ProtoReflect() protoreflect.Message {
	mi := &file_types_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserSetExclusion.ProtoReflect.Descriptor instead.
func (*UserSetExclusion) Descriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{7}
}

func (x *UserSetExclusion) GetBase() *UserSet {
	if x != nil {
		return x.Base
	}
	return nil
}

func (x *UserSetExclusion) GetExclude() *UserSet {
	if x != nil {
		return x.Exclude
	}
	return nil
}

// ConsistencyToken represents a point-in-time consistency marker
// Similar to Zanzibar's "zookie" concept
type ConsistencyToken struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Opaque token value
	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	// Timestamp when this token was issued
	IssuedAt      *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=issued_at,json=issuedAt,proto3" json:"issued_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ConsistencyToken) Reset() {
	*x = ConsistencyToken{}
	mi := &file_types_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConsistencyToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsistencyToken) ProtoMessage() {}

func (x *ConsistencyToken) ProtoReflect() protoreflect.Message {
	mi := &file_types_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsistencyToken.ProtoReflect.Descriptor instead.
func (*ConsistencyToken) Descriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{8}
}

func (x *ConsistencyToken) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *ConsistencyToken) GetIssuedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.IssuedAt
	}
	return nil
}

// Permission represents an action that can be performed
type Permission struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The namespace this permission applies to
	Namespace string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// The object ID this permission applies to
	ObjectId string `protobuf:"bytes,2,opt,name=object_id,json=objectId,proto3" json:"object_id,omitempty"`
	// The relation/permission name
	Relation string `protobuf:"bytes,3,opt,name=relation,proto3" json:"relation,omitempty"`
	// Whether this permission is allowed
	Allowed       bool `protobuf:"varint,4,opt,name=allowed,proto3" json:"allowed,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Permission) Reset() {
	*x = Permission{}
	mi := &file_types_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Permission) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Permission) ProtoMessage() {}

func (x *Permission) ProtoReflect() protoreflect.Message {
	mi := &file_types_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Permission.ProtoReflect.Descriptor instead.
func (*Permission) Descriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{9}
}

func (x *Permission) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *Permission) GetObjectId() string {
	if x != nil {
		return x.ObjectId
	}
	return ""
}

func (x *Permission) GetRelation() string {
	if x != nil {
		return x.Relation
	}
	return ""
}

func (x *Permission) GetAllowed() bool {
	if x != nil {
		return x.Allowed
	}
	return false
}

var File_types_proto protoreflect.FileDescriptor

const file_types_proto_rawDesc = "" +
	"\n" +
	"\vtypes.proto\x12\bgoacl.v1\x1a\x1fgoogle/protobuf/timestamp.proto\"\x8f\x02\n" +
	"\rRelationTuple\x12\x1c\n" +
	"\tnamespace\x18\x01 \x01(\tR\tnamespace\x12\x1b\n" +
	"\tobject_id\x18\x02 \x01(\tR\bobjectId\x12\x1a\n" +
	"\brelation\x18\x03 \x01(\tR\brelation\x12\x17\n" +
	"\auser_id\x18\x04 \x01(\tR\x06userId\x12\x18\n" +
	"\auserset\x18\x05 \x01(\tR\auserset\x129\n" +
	"\n" +
	"created_at\x18\x06 \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x129\n" +
	"\n" +
	"updated_at\x18\a \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAt\"\xd3\x01\n" +
	"\x0fNamespaceConfig\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x126\n" +
	"\trelations\x18\x02 \x03(\v2\x18.goacl.v1.RelationConfigR\trelations\x129\n" +
	"\n" +
	"created_at\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x129\n" +
	"\n" +
	"updated_at\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAt\"k\n" +
	"\x0eRelationConfig\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12#\n" +
	"\rrewrite_rules\x18\x02 \x01(\tR\frewriteRules\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\"\xa5\x02\n" +
	"\aUserSet\x12\x19\n" +
	"\auser_id\x18\x01 \x01(\tH\x00R\x06userId\x12C\n" +
	"\x0fobject_relation\x18\x02 \x01(\v2\x18.goacl.v1.ObjectRelationH\x00R\x0eobjectRelation\x12.\n" +
	"\x05union\x18\x03 \x01(\v2\x16.goacl.v1.UserSetUnionH\x00R\x05union\x12C\n" +
	"\fintersection\x18\x04 \x01(\v2\x1d.goacl.v1.UserSetIntersectionH\x00R\fintersection\x12:\n" +
	"\texclusion\x18\x05 \x01(\v2\x1a.goacl.v1.UserSetExclusionH\x00R\texclusionB\t\n" +
	"\auserset\"g\n" +
	"\x0eObjectRelation\x12\x1c\n" +
	"\tnamespace\x18\x01 \x01(\tR\tnamespace\x12\x1b\n" +
	"\tobject_id\x18\x02 \x01(\tR\bobjectId\x12\x1a\n" +
	"\brelation\x18\x03 \x01(\tR\brelation\"=\n" +
	"\fUserSetUnion\x12-\n" +
	"\bchildren\x18\x01 \x03(\v2\x11.goacl.v1.UserSetR\bchildren\"D\n" +
	"\x13UserSetIntersection\x12-\n" +
	"\bchildren\x18\x01 \x03(\v2\x11.goacl.v1.UserSetR\bchildren\"f\n" +
	"\x10UserSetExclusion\x12%\n" +
	"\x04base\x18\x01 \x01(\v2\x11.goacl.v1.UserSetR\x04base\x12+\n" +
	"\aexclude\x18\x02 \x01(\v2\x11.goacl.v1.UserSetR\aexclude\"a\n" +
	"\x10ConsistencyToken\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\x127\n" +
	"\tissued_at\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\bissuedAt\"}\n" +
	"\n" +
	"Permission\x12\x1c\n" +
	"\tnamespace\x18\x01 \x01(\tR\tnamespace\x12\x1b\n" +
	"\tobject_id\x18\x02 \x01(\tR\bobjectId\x12\x1a\n" +
	"\brelation\x18\x03 \x01(\tR\brelation\x12\x18\n" +
	"\aallowed\x18\x04 \x01(\bR\aallowedB|\n" +
	"\fcom.goacl.v1B\n" +
	"TypesProtoP\x01Z\x1fgithub.com/DangVTNhan/goacl/api\xa2\x02\x03GXX\xaa\x02\bGoacl.V1\xca\x02\bGoacl\\V1\xe2\x02\x14Goacl\\V1\\GPBMetadata\xea\x02\tGoacl::V1b\x06proto3"

var (
	file_types_proto_rawDescOnce sync.Once
	file_types_proto_rawDescData []byte
)

func file_types_proto_rawDescGZIP() []byte {
	file_types_proto_rawDescOnce.Do(func() {
		file_types_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_types_proto_rawDesc), len(file_types_proto_rawDesc)))
	})
	return file_types_proto_rawDescData
}

var file_types_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_types_proto_goTypes = []any{
	(*RelationTuple)(nil),         // 0: goacl.v1.RelationTuple
	(*NamespaceConfig)(nil),       // 1: goacl.v1.NamespaceConfig
	(*RelationConfig)(nil),        // 2: goacl.v1.RelationConfig
	(*UserSet)(nil),               // 3: goacl.v1.UserSet
	(*ObjectRelation)(nil),        // 4: goacl.v1.ObjectRelation
	(*UserSetUnion)(nil),          // 5: goacl.v1.UserSetUnion
	(*UserSetIntersection)(nil),   // 6: goacl.v1.UserSetIntersection
	(*UserSetExclusion)(nil),      // 7: goacl.v1.UserSetExclusion
	(*ConsistencyToken)(nil),      // 8: goacl.v1.ConsistencyToken
	(*Permission)(nil),            // 9: goacl.v1.Permission
	(*timestamppb.Timestamp)(nil), // 10: google.protobuf.Timestamp
}
var file_types_proto_depIdxs = []int32{
	10, // 0: goacl.v1.RelationTuple.created_at:type_name -> google.protobuf.Timestamp
	10, // 1: goacl.v1.RelationTuple.updated_at:type_name -> google.protobuf.Timestamp
	2,  // 2: goacl.v1.NamespaceConfig.relations:type_name -> goacl.v1.RelationConfig
	10, // 3: goacl.v1.NamespaceConfig.created_at:type_name -> google.protobuf.Timestamp
	10, // 4: goacl.v1.NamespaceConfig.updated_at:type_name -> google.protobuf.Timestamp
	4,  // 5: goacl.v1.UserSet.object_relation:type_name -> goacl.v1.ObjectRelation
	5,  // 6: goacl.v1.UserSet.union:type_name -> goacl.v1.UserSetUnion
	6,  // 7: goacl.v1.UserSet.intersection:type_name -> goacl.v1.UserSetIntersection
	7,  // 8: goacl.v1.UserSet.exclusion:type_name -> goacl.v1.UserSetExclusion
	3,  // 9: goacl.v1.UserSetUnion.children:type_name -> goacl.v1.UserSet
	3,  // 10: goacl.v1.UserSetIntersection.children:type_name -> goacl.v1.UserSet
	3,  // 11: goacl.v1.UserSetExclusion.base:type_name -> goacl.v1.UserSet
	3,  // 12: goacl.v1.UserSetExclusion.exclude:type_name -> goacl.v1.UserSet
	10, // 13: goacl.v1.ConsistencyToken.issued_at:type_name -> google.protobuf.Timestamp
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_types_proto_init() }
func file_types_proto_init() {
	if File_types_proto != nil {
		return
	}
	file_types_proto_msgTypes[3].OneofWrappers = []any{
		(*UserSet_UserId)(nil),
		(*UserSet_ObjectRelation)(nil),
		(*UserSet_Union)(nil),
		(*UserSet_Intersection)(nil),
		(*UserSet_Exclusion)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_types_proto_rawDesc), len(file_types_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_types_proto_goTypes,
		DependencyIndexes: file_types_proto_depIdxs,
		MessageInfos:      file_types_proto_msgTypes,
	}.Build()
	File_types_proto = out.File
	file_types_proto_goTypes = nil
	file_types_proto_depIdxs = nil
}
