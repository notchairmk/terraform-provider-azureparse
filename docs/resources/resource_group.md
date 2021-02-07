---
page_title: "azureparse_resource_group Resource - terraform-provider-azureparse"
subcategory: ""
description: |-
  Parse Azure service managed resource group to modify internal resources
---

# Resource `azureparse_resource`

Parse Azure service managed resource group to modify internal resources

## Example Usage

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

## Argument Reference

- `resource_group_name` - (Required)

## Attributes Reference

- `resource_group_id` - The ID of the resource group.

- `network_security_groups` - Zero or more network_security blocks as defined below.

- `route_tables` - Zero or more route_table blocks as defined below.

---

A `route_table` block exports the following:

* `id` - The route table resource ID.

* `name` - The route table name.

---

A `network_security_group` block exports the following:

* `id` - The network security group resource ID.

* `name` - The network security group name.
