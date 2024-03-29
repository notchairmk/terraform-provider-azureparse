package azureparse

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-06-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerRead,
		Delete: resourceServerRead,

		Schema: map[string]*schema.Schema{
			"resource_group_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"network_security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"route_tables": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	id := uuid.NewV4()
	d.SetId(id.String())
	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	nsgClient := m.(*client).SecurityGroupsClient
	resourceGroupsClient := m.(*client).ResourceGroupsClient
	routeTableClient := m.(*client).RouteTablesClient
	ctx := m.(*client).StopContext

	resourceGroupName := d.Get("resource_group_name").(string)
	// validate resource group
	res, err := resourceGroupsClient.Get(ctx, resourceGroupName)
	if err != nil {
		return fmt.Errorf("error checking resource group: %v", err)
	}
	if res.ID != nil {
		err = d.Set("resource_group_id", *res.ID)
		if err != nil {
			return fmt.Errorf("error setting resource group ID: %v", err)
		}
	}

	nsgRes, err := nsgClient.List(ctx, resourceGroupName)
	if err != nil {
		return fmt.Errorf("error listing network security groups: %v", err)
	}

	nsgList, _ := flattenNetworkSecurityGroups(nsgRes)
	err = d.Set("network_security_groups", nsgList)
	if err != nil {
		return fmt.Errorf("error setting state network security groups: %v", err)
	}

	routeTableRes, err := routeTableClient.List(ctx, resourceGroupName)
	if err != nil {
		return fmt.Errorf("error listing route tables: %v", err)
	}

	routeTableList, _ := flattenRouteTables(routeTableRes)
	err = d.Set("route_tables", routeTableList)
	if err != nil {
		return fmt.Errorf("error setting route table list: %v", err)
	}

	return nil
}

func flattenNetworkSecurityGroups(groupList network.SecurityGroupListResultPage) ([]interface{}, map[string]interface{}) {
	groups := groupList.Values()
	nsgList := make([]interface{}, 0)
	nsgMap := make(map[string]interface{})
	for _, n := range groups {
		nsgName := ""
		nsg := make(map[string]string, 2)
		nsg["id"] = *n.ID
		nsg["name"] = *n.Name
		nsgName = *n.Name
		nsgList = append(nsgList, nsg)
		nsgMap[nsgName] = nsg
	}

	return nsgList, nsgMap
}

func flattenRouteTables(tableList network.RouteTableListResultPage) ([]interface{}, map[string]interface{}) {
	tables := tableList.Values()
	routeTableList := make([]interface{}, 0)
	routeTableMap := make(map[string]interface{})
	if tables == nil {
		return routeTableList, nil
	}
	for _, n := range tables {
		routeTable := make(map[string]string, 2)
		routeTable["id"] = *n.ID
		routeTable["name"] = *n.Name
		name := *n.Name
		routeTableList = append(routeTableList, routeTable)
		routeTableMap[name] = routeTable
	}

	return routeTableList, routeTableMap
}
