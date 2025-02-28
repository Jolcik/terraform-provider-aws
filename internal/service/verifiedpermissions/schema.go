// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package verifiedpermissions

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/verifiedpermissions"
	awstypes "github.com/aws/aws-sdk-go-v2/service/verifiedpermissions/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
	fwvalidators "github.com/hashicorp/terraform-provider-aws/internal/framework/validators"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkResource(name="Schema")
func newResourceSchema(context.Context) (resource.ResourceWithConfigure, error) {
	r := &resourceSchema{}

	return r, nil
}

const (
	ResNamePolicyStoreSchema = "Schema"
)

type resourceSchema struct {
	framework.ResourceWithConfigure
}

func (r *resourceSchema) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "aws_verifiedpermissions_schema"
}

func (r *resourceSchema) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"namespaces": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"policy_store_id": schema.StringAttribute{
				Required: true,
			},
		},
		Blocks: map[string]schema.Block{
			"definition": schema.SingleNestedBlock{
				Validators: []validator.Object{
					objectvalidator.IsRequired(),
				},
				Attributes: map[string]schema.Attribute{
					"value": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							fwvalidators.JSON(),
						},
					},
				},
			},
		},
	}

	response.Schema = s
}

func (r *resourceSchema) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	conn := r.Meta().VerifiedPermissionsClient(ctx)
	var plan resourceSchemaData

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)

	if response.Diagnostics.HasError() {
		return
	}

	input := &verifiedpermissions.PutSchemaInput{
		PolicyStoreId: flex.StringFromFramework(ctx, plan.PolicyStoreID),
		Definition:    expandDefinition(ctx, plan.Definition, &response.Diagnostics),
	}

	if response.Diagnostics.HasError() {
		return
	}

	output, err := conn.PutSchema(ctx, input)

	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.VerifiedPermissions, create.ErrActionCreating, ResNamePolicyStoreSchema, plan.PolicyStoreID.ValueString(), err),
			err.Error(),
		)
		return
	}

	state := plan
	state.ID = flex.StringToFramework(ctx, output.PolicyStoreId)

	state.Namespaces = flex.FlattenFrameworkStringValueSet(ctx, output.Namespaces)

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (r *resourceSchema) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	conn := r.Meta().VerifiedPermissionsClient(ctx)
	var state resourceSchemaData

	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if response.Diagnostics.HasError() {
		return
	}

	output, err := findSchemaByPolicyStoreID(ctx, conn, state.ID.ValueString())

	if tfresource.NotFound(err) {
		response.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.VerifiedPermissions, create.ErrActionReading, ResNamePolicyStoreSchema, state.PolicyStoreID.ValueString(), err),
			err.Error(),
		)
		return
	}

	state.PolicyStoreID = flex.StringToFramework(ctx, output.PolicyStoreId)
	state.Namespaces = flex.FlattenFrameworkStringValueSet(ctx, output.Namespaces)
	state.Definition = flattenDefinition(ctx, output, &response.Diagnostics)

	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (r *resourceSchema) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	conn := r.Meta().VerifiedPermissionsClient(ctx)
	var state, plan resourceSchemaData

	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)

	if response.Diagnostics.HasError() {
		return
	}

	if !plan.Definition.Equal(state.Definition) {
		input := &verifiedpermissions.PutSchemaInput{
			PolicyStoreId: flex.StringFromFramework(ctx, state.ID),
			Definition:    expandDefinition(ctx, plan.Definition, &response.Diagnostics),
		}

		if response.Diagnostics.HasError() {
			return
		}

		_, err := conn.PutSchema(ctx, input)

		if err != nil {
			response.Diagnostics.AddError(
				create.ProblemStandardMessage(names.VerifiedPermissions, create.ErrActionUpdating, ResNamePolicyStoreSchema, state.PolicyStoreID.ValueString(), err),
				err.Error(),
			)
			return
		}

		out, err := findSchemaByPolicyStoreID(ctx, conn, state.ID.ValueString())

		if err != nil {
			response.Diagnostics.AddError(
				create.ProblemStandardMessage(names.VerifiedPermissions, create.ErrActionUpdating, ResNamePolicyStoreSchema, state.PolicyStoreID.ValueString(), err),
				err.Error(),
			)
			return
		}

		plan.Namespaces = flex.FlattenFrameworkStringValueSet(ctx, out.Namespaces)
	}

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *resourceSchema) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	conn := r.Meta().VerifiedPermissionsClient(ctx)
	var state resourceSchemaData

	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "deleting Verified Permissions Policy Store Schema", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	input := &verifiedpermissions.PutSchemaInput{
		PolicyStoreId: flex.StringFromFramework(ctx, state.ID),
		Definition: &awstypes.SchemaDefinitionMemberCedarJson{
			Value: "{}",
		},
	}

	_, err := conn.PutSchema(ctx, input)

	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.VerifiedPermissions, create.ErrActionDeleting, ResNamePolicyStoreSchema, state.PolicyStoreID.ValueString(), err),
			err.Error(),
		)
		return
	}
}

func (r *resourceSchema) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

type resourceSchemaData struct {
	ID            types.String `tfsdk:"id"`
	Definition    types.Object `tfsdk:"definition"`
	Namespaces    types.Set    `tfsdk:"namespaces"`
	PolicyStoreID types.String `tfsdk:"policy_store_id"`
}

type definition struct {
	Value types.String `tfsdk:"value"`
}

func findSchemaByPolicyStoreID(ctx context.Context, conn *verifiedpermissions.Client, id string) (*verifiedpermissions.GetSchemaOutput, error) {
	in := &verifiedpermissions.GetSchemaInput{
		PolicyStoreId: aws.String(id),
	}

	out, err := conn.GetSchema(ctx, in)
	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: in,
		}
	}
	if err != nil {
		return nil, err
	}

	if out == nil || out.Schema == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}

func expandDefinition(ctx context.Context, object types.Object, diags *diag.Diagnostics) *awstypes.SchemaDefinitionMemberCedarJson {
	var de definition
	diags.Append(object.As(ctx, &de, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	out := &awstypes.SchemaDefinitionMemberCedarJson{
		Value: de.Value.ValueString(),
	}

	return out
}

func flattenDefinition(ctx context.Context, input *verifiedpermissions.GetSchemaOutput, diags *diag.Diagnostics) types.Object {
	if input == nil {
		return fwtypes.NewObjectValueOfNull[definition](ctx).ObjectValue
	}

	var data any
	err := json.Unmarshal([]byte(aws.ToString(input.Schema)), &data)
	if err != nil {
		diags.AddError(
			"unable to unmarshal schema",
			err.Error(),
		)
		return fwtypes.NewObjectValueOfNull[definition](ctx).ObjectValue
	}

	val, err := json.Marshal(data)
	if err != nil {
		diags.AddError(
			"unable to marshal schema",
			err.Error(),
		)
		return fwtypes.NewObjectValueOfNull[definition](ctx).ObjectValue
	}

	attributeTypes := fwtypes.AttributeTypesMust[definition](ctx)
	attrs := map[string]attr.Value{}
	attrs["value"] = flex.StringValueToFramework(ctx, string(val))

	return types.ObjectValueMust(attributeTypes, attrs)
}
