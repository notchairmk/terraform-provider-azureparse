## azureparse Terraform Provider

The azureparse provider can be used to list resources in an Azure resource group.

Some Azure services create additional resources which Terraform shouldn't manage, but can and should update. The azureparse provider creates a no-op shim resource to handle dependency tree operations, which outputs information about resources within that resource group.

A specific use case could be creating a network security rule on a network security group created by another Azure service (e.g. AKS).

## Example

```hcl
resource "azureparse_resource_group" "example" {
  resource_group_name = "some-resource-group-name"
}
```

In this example, `azureparse_resource_group_parse.example` has two attributes `network_security_groups` and `route_tables`. Each attribute has `name` and `id` properties.
