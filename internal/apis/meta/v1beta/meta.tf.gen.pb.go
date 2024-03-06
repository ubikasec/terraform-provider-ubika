// Code generated by protoc-gen-api-client. DO NOT EDIT.
package metav1beta

import (
	context "context"
	diag "github.com/hashicorp/terraform-plugin-framework/diag"
	types "github.com/hashicorp/terraform-plugin-framework/types"
	basetypes "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type TypeMetaResourceModel struct {
	ApiVersion string `tfsdk:"api_version"`
	Kind       string `tfsdk:"kind"`
}

// FromProto imports field values from protobuf message
func (m *TypeMetaResourceModel) FromProto(r *TypeMeta) (_ *TypeMetaResourceModel, err error) {
	if m == nil {
		m = new(TypeMetaResourceModel)
	}
	m.ApiVersion = r.GetApiVersion()
	m.Kind = r.GetKind()
	return m, nil
}

type TypeMetaResourceTFModel struct {
	ApiVersion types.String `tfsdk:"api_version"`
	Kind       types.String `tfsdk:"kind"`
}

// ToProto converts the model to the corresponding protobuf struct
func (m *TypeMetaResourceTFModel) ToProto(ctx context.Context) (*TypeMeta, diag.Diagnostics) {
	if m == nil {
		return nil, nil
	}
	r := &TypeMeta{}
	if !m.ApiVersion.IsNull() && !m.ApiVersion.IsUnknown() {
		r.ApiVersion = m.ApiVersion.ValueString()
	}
	if !m.Kind.IsNull() && !m.Kind.IsUnknown() {
		r.Kind = m.Kind.ValueString()
	}
	return r, nil
}

type ObjectMetaResourceModel struct {
	Name      string      `tfsdk:"name"`
	Namespace string      `tfsdk:"namespace"`
	Created   types.Int64 `tfsdk:"created"`
	Updated   types.Int64 `tfsdk:"updated"`
	Version   int64       `tfsdk:"version"`
}

// FromProto imports field values from protobuf message
func (m *ObjectMetaResourceModel) FromProto(r *ObjectMeta) (_ *ObjectMetaResourceModel, err error) {
	if m == nil {
		m = new(ObjectMetaResourceModel)
	}
	m.Name = r.GetName()
	m.Namespace = r.GetNamespace()
	if r.GetCreated() != nil {
	}
	if r.GetUpdated() != nil {
	}
	m.Version = r.GetVersion()
	return m, nil
}

type ObjectMetaResourceTFModel struct {
	Name      types.String `tfsdk:"name"`
	Namespace types.String `tfsdk:"namespace"`
	Created   types.Int64  `tfsdk:"created"`
	Updated   types.Int64  `tfsdk:"updated"`
	Version   types.Int64  `tfsdk:"version"`
}

// ToProto converts the model to the corresponding protobuf struct
func (m *ObjectMetaResourceTFModel) ToProto(ctx context.Context) (*ObjectMeta, diag.Diagnostics) {
	if m == nil {
		return nil, nil
	}
	r := &ObjectMeta{}
	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		r.Name = m.Name.ValueString()
	}
	if !m.Namespace.IsNull() && !m.Namespace.IsUnknown() {
		r.Namespace = m.Namespace.ValueString()
	}
	// TODO(generateTFSDKToProtoField): need to handle type for field Created of external source ("google.golang.org/protobuf/types/known/timestamppb")
	// TODO(generateTFSDKToProtoField): need to handle type for field Updated of external source ("google.golang.org/protobuf/types/known/timestamppb")
	if !m.Version.IsNull() && !m.Version.IsUnknown() {
		r.Version = m.Version.ValueInt64()
	}
	return r, nil
}

type ListMetaResourceModel struct {
}

// FromProto imports field values from protobuf message
func (m *ListMetaResourceModel) FromProto(r *ListMeta) (_ *ListMetaResourceModel, err error) {
	if m == nil {
		m = new(ListMetaResourceModel)
	}
	return m, nil
}

type ListMetaResourceTFModel struct {
}

// ToProto converts the model to the corresponding protobuf struct
func (m *ListMetaResourceTFModel) ToProto(ctx context.Context) (*ListMeta, diag.Diagnostics) {
	if m == nil {
		return nil, nil
	}
	r := &ListMeta{}
	return r, nil
}

type UnstructuredObjectResourceModel struct {
	ApiVersion string                   `tfsdk:"api_version"`
	Kind       string                   `tfsdk:"kind"`
	Metadata   *ObjectMetaResourceModel `tfsdk:"metadata"`
}

// FromProto imports field values from protobuf message
func (m *UnstructuredObjectResourceModel) FromProto(r *UnstructuredObject) (_ *UnstructuredObjectResourceModel, err error) {
	if m == nil {
		m = new(UnstructuredObjectResourceModel)
	}
	m.ApiVersion = r.GetApiVersion()
	m.Kind = r.GetKind()
	if r.GetMetadata() != nil {
		if m.Metadata, err = m.Metadata.FromProto(r.GetMetadata()); err != nil {
			return m, err
		}
	}
	return m, nil
}

type UnstructuredObjectResourceTFModel struct {
	ApiVersion types.String `tfsdk:"api_version"`
	Kind       types.String `tfsdk:"kind"`
	Metadata   types.Object `tfsdk:"metadata"`
}

// ToProto converts the model to the corresponding protobuf struct
func (m *UnstructuredObjectResourceTFModel) ToProto(ctx context.Context) (*UnstructuredObject, diag.Diagnostics) {
	if m == nil {
		return nil, nil
	}
	r := &UnstructuredObject{}
	if !m.ApiVersion.IsNull() && !m.ApiVersion.IsUnknown() {
		r.ApiVersion = m.ApiVersion.ValueString()
	}
	if !m.Kind.IsNull() && !m.Kind.IsUnknown() {
		r.Kind = m.Kind.ValueString()
	}
	var metadata *ObjectMetaResourceTFModel
	if diags := m.Metadata.As(ctx, &metadata, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true}); diags.HasError() {
		return r, diags
	}
	if MetadataTmp, diags := metadata.ToProto(ctx); diags.HasError() {
		return r, diags
	} else {
		r.Metadata = MetadataTmp
	}
	return r, nil
}

type WatchEventResourceModel struct {
	Type string `tfsdk:"type"`
}

// FromProto imports field values from protobuf message
func (m *WatchEventResourceModel) FromProto(r *WatchEvent) (_ *WatchEventResourceModel, err error) {
	if m == nil {
		m = new(WatchEventResourceModel)
	}
	m.Type = r.GetType().String()
	// TODO(generateTerraformFromProtoField): need to handle type for field Object of Kind() bytes
	// TODO(generateTerraformFromProtoField): need to handle type for field Prev of Kind() bytes
	return m, nil
}

type WatchEventResourceTFModel struct {
	Type types.String `tfsdk:"type"`
}

// ToProto converts the model to the corresponding protobuf struct
func (m *WatchEventResourceTFModel) ToProto(ctx context.Context) (*WatchEvent, diag.Diagnostics) {
	if m == nil {
		return nil, nil
	}
	r := &WatchEvent{}
	if !m.Type.IsNull() && !m.Type.IsUnknown() {
		r.Type.UnmarshalText([]byte(m.Type.ValueString()))
	}
	// TODO(generateTFSDKToProtoField): need to handle type for field Object of Kind() bytes
	// TODO(generateTFSDKToProtoField): need to handle type for field Prev of Kind() bytes
	return r, nil
}

type ListOptionsResourceModel struct {
	Namespace string `tfsdk:"namespace"`
}

// FromProto imports field values from protobuf message
func (m *ListOptionsResourceModel) FromProto(r *ListOptions) (_ *ListOptionsResourceModel, err error) {
	if m == nil {
		m = new(ListOptionsResourceModel)
	}
	m.Namespace = r.GetNamespace()
	return m, nil
}

type ListOptionsResourceTFModel struct {
	Namespace types.String `tfsdk:"namespace"`
}

// ToProto converts the model to the corresponding protobuf struct
func (m *ListOptionsResourceTFModel) ToProto(ctx context.Context) (*ListOptions, diag.Diagnostics) {
	if m == nil {
		return nil, nil
	}
	r := &ListOptions{}
	if !m.Namespace.IsNull() && !m.Namespace.IsUnknown() {
		r.Namespace = m.Namespace.ValueString()
	}
	return r, nil
}

type GetOptionsResourceModel struct {
	Namespace string `tfsdk:"namespace"`
	Name      string `tfsdk:"name"`
}

// FromProto imports field values from protobuf message
func (m *GetOptionsResourceModel) FromProto(r *GetOptions) (_ *GetOptionsResourceModel, err error) {
	if m == nil {
		m = new(GetOptionsResourceModel)
	}
	m.Namespace = r.GetNamespace()
	m.Name = r.GetName()
	return m, nil
}

type GetOptionsResourceTFModel struct {
	Namespace types.String `tfsdk:"namespace"`
	Name      types.String `tfsdk:"name"`
}

// ToProto converts the model to the corresponding protobuf struct
func (m *GetOptionsResourceTFModel) ToProto(ctx context.Context) (*GetOptions, diag.Diagnostics) {
	if m == nil {
		return nil, nil
	}
	r := &GetOptions{}
	if !m.Namespace.IsNull() && !m.Namespace.IsUnknown() {
		r.Namespace = m.Namespace.ValueString()
	}
	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		r.Name = m.Name.ValueString()
	}
	return r, nil
}

type PatchOptionsResourceModel struct {
	Namespace string `tfsdk:"namespace"`
	Name      string `tfsdk:"name"`
	// skipping f Item is a wellknown
	// skipping f FieldMask is a wellknown
}

// FromProto imports field values from protobuf message
func (m *PatchOptionsResourceModel) FromProto(r *PatchOptions) (_ *PatchOptionsResourceModel, err error) {
	if m == nil {
		m = new(PatchOptionsResourceModel)
	}
	m.Namespace = r.GetNamespace()
	m.Name = r.GetName()
	if r.GetItem() != nil {
		// TODO(generateTerraformFromProtoField): handle wellknown field for Item
	}
	if r.GetFieldMask() != nil {
		// TODO(generateTerraformFromProtoField): handle wellknown field for FieldMask
	}
	return m, nil
}

type PatchOptionsResourceTFModel struct {
	Namespace types.String `tfsdk:"namespace"`
	Name      types.String `tfsdk:"name"`
	// skipping f Item is a wellknown
	// skipping f FieldMask is a wellknown
}

// ToProto converts the model to the corresponding protobuf struct
func (m *PatchOptionsResourceTFModel) ToProto(ctx context.Context) (*PatchOptions, diag.Diagnostics) {
	if m == nil {
		return nil, nil
	}
	r := &PatchOptions{}
	if !m.Namespace.IsNull() && !m.Namespace.IsUnknown() {
		r.Namespace = m.Namespace.ValueString()
	}
	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		r.Name = m.Name.ValueString()
	}
	// TODO(generateTFSDKToProtoField): need to handle type for field Item of external source ("google.golang.org/protobuf/types/known/anypb")
	// TODO(generateTFSDKToProtoField): need to handle type for field FieldMask of external source ("google.golang.org/protobuf/types/known/fieldmaskpb")
	return r, nil
}

type DeleteOptionsResourceModel struct {
	Namespace string `tfsdk:"namespace"`
	Name      string `tfsdk:"name"`
}

// FromProto imports field values from protobuf message
func (m *DeleteOptionsResourceModel) FromProto(r *DeleteOptions) (_ *DeleteOptionsResourceModel, err error) {
	if m == nil {
		m = new(DeleteOptionsResourceModel)
	}
	m.Namespace = r.GetNamespace()
	m.Name = r.GetName()
	return m, nil
}

type DeleteOptionsResourceTFModel struct {
	Namespace types.String `tfsdk:"namespace"`
	Name      types.String `tfsdk:"name"`
}

// ToProto converts the model to the corresponding protobuf struct
func (m *DeleteOptionsResourceTFModel) ToProto(ctx context.Context) (*DeleteOptions, diag.Diagnostics) {
	if m == nil {
		return nil, nil
	}
	r := &DeleteOptions{}
	if !m.Namespace.IsNull() && !m.Namespace.IsUnknown() {
		r.Namespace = m.Namespace.ValueString()
	}
	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		r.Name = m.Name.ValueString()
	}
	return r, nil
}

type WatchOptionsResourceModel struct {
	Namespace string `tfsdk:"namespace"`
	Name      string `tfsdk:"name"`
}

// FromProto imports field values from protobuf message
func (m *WatchOptionsResourceModel) FromProto(r *WatchOptions) (_ *WatchOptionsResourceModel, err error) {
	if m == nil {
		m = new(WatchOptionsResourceModel)
	}
	m.Namespace = r.GetNamespace()
	m.Name = r.GetName()
	return m, nil
}

type WatchOptionsResourceTFModel struct {
	Namespace types.String `tfsdk:"namespace"`
	Name      types.String `tfsdk:"name"`
}

// ToProto converts the model to the corresponding protobuf struct
func (m *WatchOptionsResourceTFModel) ToProto(ctx context.Context) (*WatchOptions, diag.Diagnostics) {
	if m == nil {
		return nil, nil
	}
	r := &WatchOptions{}
	if !m.Namespace.IsNull() && !m.Namespace.IsUnknown() {
		r.Namespace = m.Namespace.ValueString()
	}
	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		r.Name = m.Name.ValueString()
	}
	return r, nil
}