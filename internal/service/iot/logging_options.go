// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iot

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	awstypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKResource("aws_iot_logging_options")
func ResourceLoggingOptions() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceLoggingOptionsPut,
		ReadWithoutTimeout:   resourceLoggingOptionsRead,
		UpdateWithoutTimeout: resourceLoggingOptionsPut,
		DeleteWithoutTimeout: schema.NoopContext,

		Schema: map[string]*schema.Schema{
			"default_log_level": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: enum.Validate[awstypes.LogLevel](),
			},
			"disable_all_logs": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"role_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},
		},
	}
}

func resourceLoggingOptionsPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).IoTClient(ctx)

	input := &iot.SetV2LoggingOptionsInput{}

	if v, ok := d.GetOk("default_log_level"); ok {
		input.DefaultLogLevel = awstypes.LogLevel(v.(string))
	}

	if v, ok := d.GetOk("disable_all_logs"); ok {
		input.DisableAllLogs = v.(bool)
	}

	if v, ok := d.GetOk("role_arn"); ok {
		input.RoleArn = aws.String(v.(string))
	}

	err := tfresource.Retry(ctx, propagationTimeout, func() *retry.RetryError {
		var err error
		_, err = conn.SetV2LoggingOptions(ctx, input)
		if err != nil {
			var ire *awstypes.InvalidRequestException
			if errors.As(err, &ire) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "setting IoT logging options: %s", err)
	}

	d.SetId(meta.(*conns.AWSClient).Region)

	return append(diags, resourceLoggingOptionsRead(ctx, d, meta)...)
}

func resourceLoggingOptionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).IoTClient(ctx)

	output, err := conn.GetV2LoggingOptions(ctx, &iot.GetV2LoggingOptionsInput{})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading IoT logging options: %s", err)
	}

	d.Set("default_log_level", output.DefaultLogLevel)
	d.Set("disable_all_logs", output.DisableAllLogs)
	d.Set("role_arn", output.RoleArn)

	return diags
}
