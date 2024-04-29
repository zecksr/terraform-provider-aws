// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iot

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	awstypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_iot_thing_type", name="Thing Type")
// @Tags(identifierAttribute="arn")
func ResourceThingType() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceThingTypeCreate,
		ReadWithoutTimeout:   resourceThingTypeRead,
		UpdateWithoutTimeout: resourceThingTypeUpdate,
		DeleteWithoutTimeout: resourceThingTypeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("name", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deprecated": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validThingTypeName,
			},
			"properties": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validThingTypeDescription,
						},
						"searchable_attributes": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							ForceNew: true,
							MaxItems: 3,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validThingTypeSearchableAttribute,
							},
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceThingTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).IoTClient(ctx)

	name := d.Get("name").(string)
	input := &iot.CreateThingTypeInput{
		Tags:          getTagsIn(ctx),
		ThingTypeName: aws.String(name),
	}

	if v, ok := d.GetOk("properties"); ok {
		configs := v.([]interface{})
		if config, ok := configs[0].(map[string]interface{}); ok && config != nil {
			input.ThingTypeProperties = expandThingTypeProperties(config)
		}
	}

	out, err := conn.CreateThingType(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating IoT Thing Type (%s): %s", name, err)
	}

	d.SetId(aws.ToString(out.ThingTypeName))

	if v := d.Get("deprecated").(bool); v {
		input := &iot.DeprecateThingTypeInput{
			ThingTypeName: aws.String(d.Id()),
			UndoDeprecate: false,
		}

		_, err := conn.DeprecateThingType(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "deprecating IoT Thing Type (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceThingTypeRead(ctx, d, meta)...)
}

func resourceThingTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).IoTClient(ctx)

	output, err := findThingTypeByName(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] IoT Thing Type (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading IoT Thing Type (%s): %s", d.Id(), err)
	}

	d.Set("arn", output.ThingTypeArn)
	if output.ThingTypeMetadata != nil {
		d.Set("deprecated", output.ThingTypeMetadata.Deprecated)
	}
	if err := d.Set("properties", flattenThingTypeProperties(output.ThingTypeProperties)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting properties: %s", err)
	}

	return diags
}

func resourceThingTypeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).IoTClient(ctx)

	if d.HasChange("deprecated") {
		input := &iot.DeprecateThingTypeInput{
			ThingTypeName: aws.String(d.Id()),
			UndoDeprecate: !d.Get("deprecated").(bool),
		}

		_, err := conn.DeprecateThingType(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "deprecating IoT Thing Type (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceThingTypeRead(ctx, d, meta)...)
}

func resourceThingTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).IoTClient(ctx)

	// In order to delete an IoT Thing Type, you must deprecate it first and wait at least 5 minutes.
	_, err := conn.DeprecateThingType(ctx, &iot.DeprecateThingTypeInput{
		ThingTypeName: aws.String(d.Id()),
	})

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deprecating IoT Thing Type (%s): %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Deleting IoT Thing Type: %s", d.Id())
	_, err = tfresource.RetryWhenIsA[*awstypes.InvalidRequestException](ctx, 10*time.Minute,
		func() (interface{}, error) {
			return conn.DeleteThingType(ctx, &iot.DeleteThingTypeInput{
				ThingTypeName: aws.String(d.Id()),
			})
		})

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting IoT Thing Type (%s): %s", d.Id(), err)
	}

	return diags
}

func findThingTypeByName(ctx context.Context, conn *iot.Client, name string) (*iot.DescribeThingTypeOutput, error) {
	input := &iot.DescribeThingTypeInput{
		ThingTypeName: aws.String(name),
	}

	output, err := conn.DescribeThingType(ctx, input)

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}
