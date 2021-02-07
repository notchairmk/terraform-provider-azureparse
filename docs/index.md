---
page_title: "azureparse Provider"
subcategory: ""
description: |-
  The azureparse provider can be used to list resources in a resource group.
---

# azureparse Provider

The azureparse provider can be used to list resources in an Azure resource group.

Some Azure services create additional resources which Terraform shouldn't manage, but can and should update. The azureparse provider creates a no-op shim resource to handle dependency tree operations, which outputs information about resources within that resource group.

A specific use case could be creating a network security rule on a network security group created by another Azure service (e.g. AKS).

## Example Usage

* post-hoc AKS UDR to route all cluster traffic through an appliance (https://docs.microsoft.com/en-us/azure/aks/limit-egress-traffic)

```terraform
locals {
    my_appliance_addr = "..."
}

resource "azurerm_kubernetes_cluster" "cluster" {
  ...
}

resource "azureparse_resource_group" "example" {
    resource_group_name = azurerm_kubernetes_cluster.cluster.node_resource_group
}

resource "azurerm_route" "udr" {
    name                = "udr"
    resource_group_name = azurerm_kubernetes_cluster.cluster.node_resource_group
    route_table_name    = azureparse_resource_group.example.route_tables[0].name
    address_prefix      = "0.0.0.0/0"
    next_hop_type       = "VirtualAppliance"
    next_hop_in_address = local.my_appliance_addr
}
```

## Schema

* `client_id` - (Optional)

* `client_secret` - (Optional)

* `subscription_id` - (Optional)

* `tenant_id` - (Optional)
