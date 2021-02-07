## Terraform Provider

Terraform provider to list resources in an Azure resource group.

## Example

```terraform
resource "azureparse_resource_group" "example" {
  resource_group_name = "some-resource-group-name"
}
```

In this example, `azureparse_resource_group_parse.example` has two attributes `network_security_groups` and `route_tables`. Each attribute has `name` and `id` properties.
