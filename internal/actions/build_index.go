// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	internalerrors "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ action.Action = &BuildIndexAction{}
var _ action.ActionWithConfigure = &BuildIndexAction{}

func NewBuildIndexAction() action.Action {
	return &BuildIndexAction{}
}

// BuildIndexAction defines the action implementation.
type BuildIndexAction struct {
	client *providerschema.Data
}

// BuildIndexActionModel describes the action data model.
type BuildIndexActionModel struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	ClusterId      types.String `tfsdk:"cluster_id"`
	BucketName     types.String `tfsdk:"bucket_name"`
	IndexName      types.String `tfsdk:"index_name"`
	ScopeName      types.String `tfsdk:"scope_name"`
	CollectionName types.String `tfsdk:"collection_name"`
}

func (bi *BuildIndexAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_build_index"
}

func (bi *BuildIndexAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Triggers the building of an index when its in the created state.",

		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "The organization id where the index is located.",
				Required:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The project id where the index is located.",
				Required:            true,
			},
			"cluster_id": schema.StringAttribute{
				MarkdownDescription: "The cluster id where the index is located.",
				Required:            true,
			},
			"bucket_name": schema.StringAttribute{
				MarkdownDescription: "The bucket name where the index is located.",
				Required:            true,
			},
			"index_name": schema.StringAttribute{
				MarkdownDescription: "The name of the index to build.",
				Required:            true,
			},
			"collection_name": schema.StringAttribute{
				MarkdownDescription: "The name of the collection where the index is located.",
				Optional:            true,
			},
			"scope_name": schema.StringAttribute{
				MarkdownDescription: "The name of the scope where the index is located.",
				Optional:            true,
			},
		},
	}
}

func (bi *BuildIndexAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*providerschema.Data)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	bi.client = client
}

func (bi *BuildIndexAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	// Send a progress message back to Terraform
	resp.SendProgress(action.InvokeProgressEvent{
		Message: "starting action invocation",
	})

	var data BuildIndexActionModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var scope, collection string
	if data.ScopeName.IsNull() {
		scope = "_default"
	} else {
		scope = data.ScopeName.ValueString()
	}
	if data.CollectionName.IsNull() {
		collection = "_default"
	} else {
		collection = data.CollectionName.ValueString()
	}

	ddl := fmt.Sprintf(
		"BUILD INDEX ON `%s`.`%s`.`%s`(%s)",
		data.BucketName.ValueString(),
		scope,
		collection,
		data.IndexName.ValueString(),
	)

	monitor := func(cfg api.EndpointCfg) (response *api.Response, err error) {
		if err := api.Limiter.Wait(ctx); err != nil {
			// do not block if rate limiter fails
			tflog.Error(ctx, "rate limiter error: "+err.Error())
		}

		res, err := bi.client.ClientV1.ExecuteWithRetry(
			ctx,
			cfg,
			nil,
			bi.client.Token,
			nil,
		)

		status := api.IndexBuildStatusResponse{}
		if err = json.Unmarshal(response.Body, &status); err != nil {
			return res, err
		}
		resp.SendProgress(action.InvokeProgressEvent{
			Message: fmt.Sprintf("current status: %s", status.Status),
		})

		return res, err
	}

	err := api.WatchIndexes(
		ctx,
		"Created",
		[]string{data.IndexName.ValueString()},
		monitor,
		api.Options{
			Host:       bi.client.HostURL,
			OrgId:      data.OrganizationId.ValueString(),
			ProjectId:  data.ProjectId.ValueString(),
			ClusterId:  data.ClusterId.ValueString(),
			Bucket:     data.BucketName.ValueString(),
			Scope:      scope,
			Collection: collection,
		},
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Ensure Indexes Are Created Failed",
			fmt.Sprintf("Cannot do build index as all indexes are not in created state.  Error: %v\n", err.Error()),
		)

		return
	}

	state := data

	err = bi.executeGsiDdl(ctx, &state, ddl)

	tflog.Trace(ctx, "invoke an action")

	// Send a progress message back to Terraform
	resp.SendProgress(action.InvokeProgressEvent{
		Message: "finished action invocation",
	})
}

func (bi *BuildIndexAction) executeGsiDdl(ctx context.Context, plan *BuildIndexActionModel, ddl string) error {
	uri := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/queryService/indexes",
		bi.client.HostURL,
		plan.OrganizationId.ValueString(),
		plan.ProjectId.ValueString(),
		plan.ClusterId.ValueString(),
	)

	cfg := api.EndpointCfg{Url: uri, Method: http.MethodPost, SuccessStatus: http.StatusOK}
	ddlRequest := api.IndexDDLRequest{Definition: ddl}

	if err := api.Limiter.Wait(ctx); err != nil {
		// do not block if rate limiter fails
		tflog.Error(ctx, "rate limiter error: "+err.Error())
	}
	response, err := bi.client.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		ddlRequest,
		bi.client.Token,
		nil,
	)
	switch {
	case err == nil:
	// for large non-deferred indexes API server will timeout.
	case errors.Is(err, internalerrors.ErrGatewayTimeoutForIndexDDL):
		return internalerrors.ErrGatewayTimeoutForIndexDDL
	case errors.Is(err, &api.Error{}):
		// Indexer doesn't allow concurrent index builds.
		// Index build is resource intensive operation from indexer and KV perspective
		// as indexer will request data for the keyspace from the beginning.
		//
		// Indexer will automatically retry in the background.
		apiError, _ := err.(*api.Error)
		if strings.Contains(strings.ToLower(apiError.Message), "build already in progress") ||
			strings.Contains(strings.ToLower(apiError.Message), "concurrent create index request") {

			return internalerrors.ErrConcurrentIndexCreation
		}

		// Some other error from the API server.
		return err
	default:
		return err
	}

	ddlResponse := api.IndexDDLResponse{}
	err = json.Unmarshal(response.Body, &ddlResponse)
	if err != nil {
		return err
	}

	//  There are some cases where an operation fails yet query service returns 200 OK.
	//	For example, when an index is not found or already exists.
	//  In this case, query service returns errors attribute.
	//  See MB-62943 for more details.
	if ddlResponse.Errors != nil {
		var message string
		if len(ddlResponse.Errors) > 0 {
			message = ddlResponse.Errors[0].Msg
		}

		return errors.New(message)
	}

	return nil
}
