---
subcategory: "OpsWorks"
layout: "aws"
page_title: "AWS: aws_opsworks_stack"
description: |-
  Provides an OpsWorks stack resource.
---


<!-- Please do not edit this file, it is generated. -->
# Resource: aws_opsworks_stack

Provides an OpsWorks stack resource.

!> **ALERT:** AWS no longer supports OpsWorks Stacks. All related resources will be removed from the Terraform AWS Provider in the next major version.

## Example Usage

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { Token, TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { OpsworksStack } from "./.gen/providers/aws/opsworks-stack";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    new OpsworksStack(this, "main", {
      customJson: '{\n "foobar": {\n    "version": "1.0.0"\n  }\n}\n\n',
      defaultInstanceProfileArn: opsworks.arn,
      name: "awesome-stack",
      region: "us-west-1",
      serviceRoleArn: Token.asString(awsIamRoleOpsworks.arn),
      tags: {
        Name: "foobar-terraform-stack",
      },
    });
  }
}

```

## Argument Reference

This resource supports the following arguments:

* `name` - (Required) The name of the stack.
* `region` - (Required) The name of the region where the stack will exist.
* `serviceRoleArn` - (Required) The ARN of an IAM role that the OpsWorks service will act as.
* `defaultInstanceProfileArn` - (Required) The ARN of an IAM Instance Profile that created instances will have by default.
* `agentVersion` - (Optional) If set to `"LATEST"`, OpsWorks will automatically install the latest version.
* `berkshelfVersion` - (Optional) If `manageBerkshelf` is enabled, the version of Berkshelf to use.
* `color` - (Optional) Color to paint next to the stack's resources in the OpsWorks console.
* `configurationManagerName` - (Optional) Name of the configuration manager to use. Defaults to "Chef".
* `configurationManagerVersion` - (Optional) Version of the configuration manager to use. Defaults to "11.4".
* `customCookbooksSource` - (Optional) When `useCustomCookbooks` is set, provide this sub-object as described below.
* `customJson` - (Optional) User defined JSON passed to "Chef". Use a "here doc" for multiline JSON.
* `defaultAvailabilityZone` - (Optional) Name of the availability zone where instances will be created by default.
  Cannot be set when `vpcId` is set.
* `defaultOs` - (Optional) Name of OS that will be installed on instances by default.
* `defaultRootDeviceType` - (Optional) Name of the type of root device instances will have by default.
* `defaultSshKeyName` - (Optional) Name of the SSH keypair that instances will have by default.
* `defaultSubnetId` - (Optional) ID of the subnet in which instances will be created by default.
  Required if `vpcId` is set to a VPC other than the default VPC, and forbidden if it isn't.
* `hostnameTheme` - (Optional) Keyword representing the naming scheme that will be used for instance hostnames within this stack.
* `manageBerkshelf` - (Optional) Boolean value controlling whether Opsworks will run Berkshelf for this stack.
* `tags` - (Optional) A map of tags to assign to the resource.
  If configured with a provider [`defaultTags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.
* `useCustomCookbooks` - (Optional) Boolean value controlling whether the custom cookbook settings are enabled.
* `useOpsworksSecurityGroups` - (Optional) Boolean value controlling whether the standard OpsWorks security groups apply to created instances.
* `vpcId` - (Optional) ID of the VPC that this stack belongs to.
  Defaults to the region's default VPC.
* `customJson` - (Optional) Custom JSON attributes to apply to the entire stack.

The `customCookbooksSource` block supports the following arguments:

* `type` - (Required) The type of source to use. For example, "archive".
* `url` - (Required) The URL where the cookbooks resource can be found.
* `username` - (Optional) Username to use when authenticating to the source.
* `password` - (Optional) Password to use when authenticating to the source. Terraform cannot perform drift detection of this configuration.
* `sshKey` - (Optional) SSH key to use when authenticating to the source. Terraform cannot perform drift detection of this configuration.
* `revision` - (Optional) For sources that are version-aware, the revision to use.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `id` - The id of the stack.
* `tagsAll` - A map of tags assigned to the resource, including those inherited from the provider [`defaultTags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import OpsWorks stacks using the `id`. For example:

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { OpsworksStack } from "./.gen/providers/aws/opsworks-stack";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    OpsworksStack.generateConfigForImport(
      this,
      "bar",
      "00000000-0000-0000-0000-000000000000"
    );
  }
}

```

Using `terraform import`, import OpsWorks stacks using the `id`. For example:

```console
% terraform import aws_opsworks_stack.bar 00000000-0000-0000-0000-000000000000
```

<!-- cache-key: cdktf-0.20.8 input-6af1072741f186fcddc321c71a350298ac67484fcba6100a3bc7dffe500e612b -->