// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resiliencehub_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/resiliencehub"
	"github.com/aws/aws-sdk-go-v2/service/resiliencehub/types"
	awstypes "github.com/aws/aws-sdk-go-v2/service/resiliencehub/types"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	tfresiliencehub "github.com/hashicorp/terraform-provider-aws/internal/service/resiliencehub"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccResilienceHubResiliencyPolicy_basic(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var policy resiliencehub.DescribeResiliencyPolicyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_resiliencehub_resiliency_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionNot(t, names.USGovCloudPartitionID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ResilienceHubServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResiliencyPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResiliencyPolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResiliencyPolicyExists(ctx, resourceName, &policy),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, names.ResilienceHubServiceID, regexache.MustCompile(`resiliency-policy/.+$`)),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckNoResourceAttr(resourceName, names.AttrDescription),
					resource.TestCheckResourceAttr(resourceName, "tier", "NotApplicable"),
					resource.TestCheckResourceAttr(resourceName, "data_location_constraint", "AnyLocation"),
					resource.TestCheckResourceAttr(resourceName, "policy.az.rpo_in_secs", "3600"),
					resource.TestCheckResourceAttr(resourceName, "policy.az.rto_in_secs", "3600"),
					resource.TestCheckResourceAttr(resourceName, "policy.hardware.rpo_in_secs", "3600"),
					resource.TestCheckResourceAttr(resourceName, "policy.hardware.rto_in_secs", "3600"),
					resource.TestCheckNoResourceAttr(resourceName, "policy.region.rpo_in_secs"),
					resource.TestCheckNoResourceAttr(resourceName, "policy.region.rto_in_secs"),
					resource.TestCheckResourceAttr(resourceName, "policy.software.rpo_in_secs", "3600"),
					resource.TestCheckResourceAttr(resourceName, "policy.software.rto_in_secs", "3600"),
				),
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateIdFunc:                    testAccAttrImportStateIdFunc(resourceName, names.AttrARN),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: names.AttrARN,
			},
		},
	})
}

func TestAccResilienceHubResiliencyPolicy_update(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var policy1, policy2 resiliencehub.DescribeResiliencyPolicyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_resiliencehub_resiliency_policy.test"

	updatedPolicyObjValue := "86400"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionNot(t, names.USGovCloudPartitionID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ResilienceHubServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResiliencyPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResiliencyPolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResiliencyPolicyExists(ctx, resourceName, &policy1),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckNoResourceAttr(resourceName, names.AttrDescription),
					resource.TestCheckResourceAttr(resourceName, "tier", "NotApplicable"),
					resource.TestCheckResourceAttr(resourceName, "data_location_constraint", "AnyLocation"),
					resource.TestCheckResourceAttr(resourceName, "policy.az.rpo_in_secs", "3600"),
					resource.TestCheckResourceAttr(resourceName, "policy.az.rto_in_secs", "3600"),
					resource.TestCheckResourceAttr(resourceName, "policy.hardware.rpo_in_secs", "3600"),
					resource.TestCheckResourceAttr(resourceName, "policy.hardware.rto_in_secs", "3600"),
					resource.TestCheckNoResourceAttr(resourceName, "policy.region.rpo_in_secs"),
					resource.TestCheckNoResourceAttr(resourceName, "policy.region.rto_in_secs"),
					resource.TestCheckResourceAttr(resourceName, "policy.software.rpo_in_secs", "3600"),
					resource.TestCheckResourceAttr(resourceName, "policy.software.rto_in_secs", "3600"),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.Ct0),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, names.ResilienceHubServiceID, regexache.MustCompile(`resiliency-policy/.+`)),
				),
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateIdFunc:                    testAccAttrImportStateIdFunc(resourceName, names.AttrARN),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: names.AttrARN,
			},
			{
				Config: testAccResiliencyPolicyConfig_updatePolicy(rName, updatedPolicyObjValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResiliencyPolicyExists(ctx, resourceName, &policy2),
					resource.TestCheckResourceAttr(resourceName, "policy.az.rpo_in_secs", updatedPolicyObjValue),
					resource.TestCheckResourceAttr(resourceName, "policy.az.rto_in_secs", updatedPolicyObjValue),
					resource.TestCheckResourceAttr(resourceName, "policy.hardware.rpo_in_secs", updatedPolicyObjValue),
					resource.TestCheckResourceAttr(resourceName, "policy.hardware.rto_in_secs", updatedPolicyObjValue),
					resource.TestCheckResourceAttr(resourceName, "policy.region.rpo_in_secs", updatedPolicyObjValue),
					resource.TestCheckResourceAttr(resourceName, "policy.region.rto_in_secs", updatedPolicyObjValue),
					resource.TestCheckResourceAttr(resourceName, "policy.software.rpo_in_secs", updatedPolicyObjValue),
					resource.TestCheckResourceAttr(resourceName, "policy.software.rto_in_secs", updatedPolicyObjValue),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, names.ResilienceHubServiceID, regexache.MustCompile(`resiliency-policy/.+`)),
				),
			},
		},
	})
}

func TestAccResilienceHubResiliencyPolicy_dataLocationConstraint(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var policy1, policy2 resiliencehub.DescribeResiliencyPolicyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_resiliencehub_resiliency_policy.test"

	expectNoARNChange := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionNot(t, names.USGovCloudPartitionID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ResilienceHubServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResiliencyPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResiliencyPolicyConfig_dataLocationConstraint(rName, awstypes.DataLocationConstraintSameContinent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResiliencyPolicyExists(ctx, resourceName, &policy1),
					resource.TestCheckResourceAttr(resourceName, "data_location_constraint", string(awstypes.DataLocationConstraintSameContinent)),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					expectNoARNChange.AddStateValue(resourceName, tfjsonpath.New(names.AttrARN)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					},
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateIdFunc:                    testAccAttrImportStateIdFunc(resourceName, names.AttrARN),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: names.AttrARN,
			},
			{
				Config: testAccResiliencyPolicyConfig_dataLocationConstraint(rName, awstypes.DataLocationConstraintSameCountry),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResiliencyPolicyExists(ctx, resourceName, &policy2),
					resource.TestCheckResourceAttr(resourceName, "data_location_constraint", string(awstypes.DataLocationConstraintSameCountry)),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					expectNoARNChange.AddStateValue(resourceName, tfjsonpath.New(names.AttrARN)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateIdFunc:                    testAccAttrImportStateIdFunc(resourceName, names.AttrARN),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: names.AttrARN,
			},
		},
	})
}

func TestAccResilienceHubResiliencyPolicy_description(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var policy1, policy2 resiliencehub.DescribeResiliencyPolicyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_resiliencehub_resiliency_policy.test"

	const (
		initial = "initial"
		updated = "updated"
	)

	expectNoARNChange := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionNot(t, names.USGovCloudPartitionID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ResilienceHubServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResiliencyPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResiliencyPolicyConfig_description(rName, initial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResiliencyPolicyExists(ctx, resourceName, &policy1),
					resource.TestCheckResourceAttr(resourceName, names.AttrDescription, initial),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					expectNoARNChange.AddStateValue(resourceName, tfjsonpath.New(names.AttrARN)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					},
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateIdFunc:                    testAccAttrImportStateIdFunc(resourceName, names.AttrARN),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: names.AttrARN,
			},
			{
				Config: testAccResiliencyPolicyConfig_description(rName, updated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResiliencyPolicyExists(ctx, resourceName, &policy2),
					resource.TestCheckResourceAttr(resourceName, names.AttrDescription, updated),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					expectNoARNChange.AddStateValue(resourceName, tfjsonpath.New(names.AttrARN)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateIdFunc:                    testAccAttrImportStateIdFunc(resourceName, names.AttrARN),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: names.AttrARN,
			},
		},
	})
}

func TestAccResilienceHubResiliencyPolicy_name(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var policy1, policy2 resiliencehub.DescribeResiliencyPolicyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rNameUpdated := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_resiliencehub_resiliency_policy.test"

	expectNoARNChange := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionNot(t, names.USGovCloudPartitionID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ResilienceHubServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResiliencyPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResiliencyPolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResiliencyPolicyExists(ctx, resourceName, &policy1),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					expectNoARNChange.AddStateValue(resourceName, tfjsonpath.New(names.AttrARN)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					},
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateIdFunc:                    testAccAttrImportStateIdFunc(resourceName, names.AttrARN),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: names.AttrARN,
			},
			{
				Config: testAccResiliencyPolicyConfig_basic(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResiliencyPolicyExists(ctx, resourceName, &policy2),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rNameUpdated),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					expectNoARNChange.AddStateValue(resourceName, tfjsonpath.New(names.AttrARN)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateIdFunc:                    testAccAttrImportStateIdFunc(resourceName, names.AttrARN),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: names.AttrARN,
			},
		},
	})
}

func TestAccResilienceHubResiliencyPolicy_tier(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var policy1, policy2 resiliencehub.DescribeResiliencyPolicyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_resiliencehub_resiliency_policy.test"

	expectNoARNChange := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionNot(t, names.USGovCloudPartitionID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ResilienceHubServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResiliencyPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResiliencyPolicyConfig_tier(rName, awstypes.ResiliencyPolicyTierMissionCritical),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResiliencyPolicyExists(ctx, resourceName, &policy1),
					resource.TestCheckResourceAttr(resourceName, "tier", string(awstypes.ResiliencyPolicyTierMissionCritical)),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					expectNoARNChange.AddStateValue(resourceName, tfjsonpath.New(names.AttrARN)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					},
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateIdFunc:                    testAccAttrImportStateIdFunc(resourceName, names.AttrARN),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: names.AttrARN,
			},
			{
				Config: testAccResiliencyPolicyConfig_tier(rName, awstypes.ResiliencyPolicyTierCoreServices),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResiliencyPolicyExists(ctx, resourceName, &policy2),
					resource.TestCheckResourceAttr(resourceName, "tier", string(awstypes.ResiliencyPolicyTierCoreServices)),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					expectNoARNChange.AddStateValue(resourceName, tfjsonpath.New(names.AttrARN)),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateIdFunc:                    testAccAttrImportStateIdFunc(resourceName, names.AttrARN),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: names.AttrARN,
			},
		},
	})
}

func TestAccResilienceHubResiliencyPolicy_tags(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var policy1, policy2 resiliencehub.DescribeResiliencyPolicyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_resiliencehub_resiliency_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionNot(t, names.USGovCloudPartitionID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ResilienceHubServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResiliencyPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResiliencyPolicyConfig_tags1(rName, acctest.CtKey1, acctest.CtValue1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResiliencyPolicyExists(ctx, resourceName, &policy1),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.Ct1),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey1, acctest.CtValue1),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, names.ResilienceHubServiceID, regexache.MustCompile(`resiliency-policy/.+`)),
				),
			},
			{
				Config: testAccResiliencyPolicyConfig_tag2(rName, acctest.CtKey1, acctest.CtValue1Updated, acctest.CtKey2, acctest.CtValue2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResiliencyPolicyExists(ctx, resourceName, &policy2),
					testAccCheckResiliencyPolicyNotRecreated(&policy1, &policy2),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.Ct2),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey1, acctest.CtValue1Updated),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey2, acctest.CtValue2),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, names.ResilienceHubServiceID, regexache.MustCompile(`resiliency-policy/.+`)),
				),
			},
		},
	})
}

func TestAccResilienceHubResiliencyPolicy_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var policy resiliencehub.DescribeResiliencyPolicyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_resiliencehub_resiliency_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionNot(t, names.USGovCloudPartitionID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ResilienceHubServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResiliencyPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResiliencyPolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResiliencyPolicyExists(ctx, resourceName, &policy),
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfresiliencehub.ResourceResiliencyPolicy, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckResiliencyPolicyDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ResilienceHubClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_resiliencehub_resiliency_policy" {
				continue
			}

			input := &resiliencehub.DescribeResiliencyPolicyInput{
				PolicyArn: aws.String(rs.Primary.Attributes[names.AttrARN]),
			}
			_, err := conn.DescribeResiliencyPolicy(ctx, input)
			if errs.IsA[*types.ResourceNotFoundException](err) {
				return nil
			}
			if err != nil {
				return create.Error(names.ResilienceHub, create.ErrActionCheckingDestroyed, tfresiliencehub.ResNameResiliencyPolicy, rs.Primary.ID, err)
			}

			return create.Error(names.ResilienceHub, create.ErrActionCheckingDestroyed, tfresiliencehub.ResNameResiliencyPolicy, rs.Primary.ID, errors.New("not destroyed"))
		}

		return nil
	}
}

func testAccCheckResiliencyPolicyExists(ctx context.Context, name string, policy *resiliencehub.DescribeResiliencyPolicyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.ResilienceHub, create.ErrActionCheckingExistence, tfresiliencehub.ResNameResiliencyPolicy, name, errors.New("not found"))
		}

		if rs.Primary.Attributes[names.AttrARN] == "" {
			return create.Error(names.ResilienceHub, create.ErrActionCheckingExistence, tfresiliencehub.ResNameResiliencyPolicy, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ResilienceHubClient(ctx)
		resp, err := conn.DescribeResiliencyPolicy(ctx, &resiliencehub.DescribeResiliencyPolicyInput{
			PolicyArn: aws.String(rs.Primary.Attributes[names.AttrARN]),
		})

		if err != nil {
			return create.Error(names.ResilienceHub, create.ErrActionCheckingExistence, tfresiliencehub.ResNameResiliencyPolicy, rs.Primary.ID, err)
		}

		*policy = *resp

		return nil
	}
}

func testAccAttrImportStateIdFunc(resourceName, attrName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.Attributes[attrName], nil
	}
}

func testAccPreCheck(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ResilienceHubClient(ctx)

	input := &resiliencehub.ListResiliencyPoliciesInput{}
	_, err := conn.ListResiliencyPolicies(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}
	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccCheckResiliencyPolicyNotRecreated(before, after *resiliencehub.DescribeResiliencyPolicyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.ToString(before.Policy.PolicyArn), aws.ToString(after.Policy.PolicyArn); before != after {
			return create.Error(names.ResilienceHub, create.ErrActionCheckingNotRecreated, tfresiliencehub.ResNameResiliencyPolicy, before, errors.New("recreated"))
		}

		return nil
	}
}

func testAccResiliencyPolicyConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_resiliencehub_resiliency_policy" "test" {
  name = %[1]q

  tier = "NotApplicable"

  policy {
    az {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
    hardware {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
    software {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
  }
}
`, rName)
}

func testAccResiliencyPolicyConfig_dataLocationConstraint(rName string, dataLocationConstraint awstypes.DataLocationConstraint) string {
	return fmt.Sprintf(`
resource "aws_resiliencehub_resiliency_policy" "test" {
  name = %[1]q

  tier = "NotApplicable"

  data_location_constraint = %[2]q

  policy {
    az {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
    hardware {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
    software {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
  }
}
`, rName, dataLocationConstraint)
}

func testAccResiliencyPolicyConfig_description(rName, resPolicyDescValue string) string {
	return fmt.Sprintf(`
resource "aws_resiliencehub_resiliency_policy" "test" {
  name = %[1]q

  description = %[2]q

  tier = "NotApplicable"

  policy {
    az {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
    hardware {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
    software {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
  }
}
`, rName, resPolicyDescValue)
}

func testAccResiliencyPolicyConfig_tier(rName string, tier awstypes.ResiliencyPolicyTier) string {
	return fmt.Sprintf(`
resource "aws_resiliencehub_resiliency_policy" "test" {
  name = %[1]q

  tier = %[2]q

  data_location_constraint = "AnyLocation"

  policy {
    az {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
    hardware {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
    software {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
  }
}
`, rName, tier)
}

func testAccResiliencyPolicyConfig_updatePolicy(rName, resPolicyObjValue string) string {
	return fmt.Sprintf(`
resource "aws_resiliencehub_resiliency_policy" "test" {
  name = %[1]q

  description = %[1]q

  tier = "NotApplicable"

  data_location_constraint = "AnyLocation"

  policy {
    region {
      rpo_in_secs = %[2]q
      rto_in_secs = %[2]q
    }
    az {
      rpo_in_secs = %[2]q
      rto_in_secs = %[2]q
    }
    hardware {
      rpo_in_secs = %[2]q
      rto_in_secs = %[2]q
    }
    software {
      rpo_in_secs = %[2]q
      rto_in_secs = %[2]q
    }
  }

  tags = {
    Name  = %[1]q
    Value = "Other"
  }
}
`, rName, resPolicyObjValue)
}

func testAccResiliencyPolicyConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_resiliencehub_resiliency_policy" "test" {
  name = %[1]q

  description = %[1]q

  tier = "NotApplicable"

  data_location_constraint = "AnyLocation"

  policy {
    region {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
    az {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
    hardware {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
    software {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
  }

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccResiliencyPolicyConfig_tag2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_resiliencehub_resiliency_policy" "test" {
  name = %[1]q

  description = %[1]q

  tier = "NotApplicable"

  data_location_constraint = "AnyLocation"

  policy {
    region {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
    az {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
    hardware {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
    software {
      rpo_in_secs = 3600
      rto_in_secs = 3600
    }
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
