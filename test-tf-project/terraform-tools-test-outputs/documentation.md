## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 6.1.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_bastion"></a> [bastion](#module\_bastion) | umotif-public/bastion/aws | ~> 2.1.0 |
| <a name="module_moose_fargate"></a> [moose\_fargate](#module\_moose\_fargate) | ./fargate | n/a |
| <a name="module_opensearch"></a> [opensearch](#module\_opensearch) | ./opensearch | n/a |

## Resources

| Name | Type |
|------|------|
| [aws_cloudwatch_event_bus.syncer](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_event_bus) | resource |
| [aws_cloudwatch_event_bus.syncer_local](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_event_bus) | resource |
| [aws_dynamodb_table.moose](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table) | resource |
| [aws_iam_policy.lambda_opensearch_assume_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy) | resource |
| [aws_iam_policy.lambda_ssm_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy) | resource |
| [aws_iam_policy.moose_instance_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy) | resource |
| [aws_iam_policy.moose_instance_policy_s3](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy) | resource |
| [aws_iam_policy.moose_instance_policy_sqs](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy) | resource |
| [aws_iam_policy.moose_instance_policy_textract](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy) | resource |
| [aws_iam_role.lambda](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy_attachment.attach_lambda_assume_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.lambda-attach](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.lambda-attach-basicexecution](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.lambda-attach-s3](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.lambda-attach-ssm](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.lambda-attach-textract](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.lambda-attach-vpc](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_route53_record.moose_endpoint](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) | resource |
| [aws_route53_record.moose_endpoint_db](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) | resource |
| [aws_s3_bucket.presigned](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket) | resource |
| [aws_s3_bucket_cors_configuration.presigned](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_cors_configuration) | resource |
| [aws_s3_bucket_notification.presigned_eventbridge](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_notification) | resource |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |
| [aws_s3_bucket.rs_chainlang](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/s3_bucket) | data source |
| [aws_s3_bucket.textract_bucket](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/s3_bucket) | data source |
| [aws_sqs_queue.textract_sqs](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/sqs_queue) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_account_number"></a> [account\_number](#input\_account\_number) | description | `string` | `""` | no |
| <a name="input_bastion_ssh_key_name"></a> [bastion\_ssh\_key\_name](#input\_bastion\_ssh\_key\_name) | n/a | `string` | n/a | yes |
| <a name="input_build"></a> [build](#input\_build) | deploy build | `string` | `"0.1.0"` | no |
| <a name="input_cognito_identity_pool_id"></a> [cognito\_identity\_pool\_id](#input\_cognito\_identity\_pool\_id) | n/a | `string` | n/a | yes |
| <a name="input_cognito_identity_pool_role_arn"></a> [cognito\_identity\_pool\_role\_arn](#input\_cognito\_identity\_pool\_role\_arn) | n/a | `string` | n/a | yes |
| <a name="input_cognito_role_arn"></a> [cognito\_role\_arn](#input\_cognito\_role\_arn) | n/a | `string` | n/a | yes |
| <a name="input_cognito_user_pool_id"></a> [cognito\_user\_pool\_id](#input\_cognito\_user\_pool\_id) | n/a | `string` | n/a | yes |
| <a name="input_db_engine_version"></a> [db\_engine\_version](#input\_db\_engine\_version) | n/a | `string` | `"15.3"` | no |
| <a name="input_db_instance_class"></a> [db\_instance\_class](#input\_db\_instance\_class) | n/a | `string` | `"db.t3.small"` | no |
| <a name="input_db_password"></a> [db\_password](#input\_db\_password) | n/a | `string` | n/a | yes |
| <a name="input_db_username"></a> [db\_username](#input\_db\_username) | n/a | `string` | `"admin"` | no |
| <a name="input_domain_name"></a> [domain\_name](#input\_domain\_name) | n/a | `string` | n/a | yes |
| <a name="input_ecr_repo"></a> [ecr\_repo](#input\_ecr\_repo) | n/a | `string` | n/a | yes |
| <a name="input_env"></a> [env](#input\_env) | environment to deploy to | `string` | n/a | yes |
| <a name="input_fip_moose_api_key"></a> [fip\_moose\_api\_key](#input\_fip\_moose\_api\_key) | n/a | `string` | n/a | yes |
| <a name="input_grafana_url"></a> [grafana\_url](#input\_grafana\_url) | n/a | `string` | n/a | yes |
| <a name="input_hosted_zone_id"></a> [hosted\_zone\_id](#input\_hosted\_zone\_id) | n/a | `string` | n/a | yes |
| <a name="input_is_env_moosedataprod"></a> [is\_env\_moosedataprod](#input\_is\_env\_moosedataprod) | env dynamodb for development | `bool` | `false` | no |
| <a name="input_is_local_resource"></a> [is\_local\_resource](#input\_is\_local\_resource) | n/a | `bool` | `true` | no |
| <a name="input_moose_db_endpoint"></a> [moose\_db\_endpoint](#input\_moose\_db\_endpoint) | n/a | `string` | n/a | yes |
| <a name="input_moose_endpoint"></a> [moose\_endpoint](#input\_moose\_endpoint) | n/a | `string` | n/a | yes |
| <a name="input_opensearch_domain_cluster"></a> [opensearch\_domain\_cluster](#input\_opensearch\_domain\_cluster) | n/a | `string` | n/a | yes |
| <a name="input_opensearch_endpoint"></a> [opensearch\_endpoint](#input\_opensearch\_endpoint) | n/a | `string` | n/a | yes |
| <a name="input_rs_chain_repo"></a> [rs\_chain\_repo](#input\_rs\_chain\_repo) | n/a | `string` | n/a | yes |
| <a name="input_user_arn"></a> [user\_arn](#input\_user\_arn) | description | `string` | `"arn:aws:iam::752380136218:user/epuerta"` | no |

## Outputs

No outputs.
