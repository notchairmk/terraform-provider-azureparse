package loafsley

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-06-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func resourceGroupParse() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerRead,
		Delete: resourceServerRead,

		Schema: map[string]*schema.Schema{
			"resource_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"network_security_groups": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
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

			"private_dns_zones": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
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

			"private_endpoints": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
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

			"route_tables": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
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
	nsgClient := m.(*Client).SecurityGroupsClient
	privateDNSZonesClient := m.(*Client).PrivateDNSZonesClient
	privateEndpointsClient := m.(*Client).PrivateEndpointsClient
	routeTableClient := m.(*Client).RouteTablesClient
	ctx := m.(*Client).StopContext

	resourceGroupName := d.Get("resource_group_name").(string)
	nsgList, err := nsgClient.List(ctx, resourceGroupName)
	if err != nil {
		return fmt.Errorf("error listing security groups: %v", err)
	}
	err = d.Set("network_security_groups", flattenNetworkSecurityGroups(nsgList))
	if err != nil {
		return fmt.Errorf("error setting state network security groups: %v", err)
	}

	privateEndpointList, err := privateEndpointsClient.List(ctx, resourceGroupName)
	if err != nil {
		return fmt.Errorf("error listing private endpoints: %v", err)
	}
	privateEndpoints := flattenPrivateEndpoints(privateEndpointList)
	err = d.Set("private_endpoints", privateEndpoints)
	if err != nil {
		return fmt.Errorf("error setting private endpoint state: %v", err)
	}

	if len(privateEndpoints) > 0 {
		privateEndpoint := privateEndpoints[0].(map[string]interface{})
		dnsZones, err := privateDNSZonesClient.List(ctx, privateEndpoint["name"].(string), resourceGroupName)
		if err != nil {
			return fmt.Errorf("error listing dns zones: %v", err)
		}
		err = d.Set("private_dns_zones", flattenPrivateDNSZones(dnsZones))
		if err != nil {
			return fmt.Errorf("error setting private dns zone state: %v", err)
		}
	}

	routeTablesList, err := routeTableClient.List(ctx, resourceGroupName)
	if err != nil {
		return fmt.Errorf("error listing route tables: %v", err)
	}
	err = d.Set("route_tables", flattenRouteTables(routeTablesList))
	if err != nil {
		return fmt.Errorf("error setting state route tables: %v", err)
	}

	return nil
}

func flattenNetworkSecurityGroups(groupList network.SecurityGroupListResultPage) []interface{} {
	groups := groupList.Values()
	nsgs := make([]interface{}, 0)
	for _, n := range groups {
		nsg := make(map[string]string)
		if n.ID != nil {
			nsg["id"] = *n.ID
		}

		if n.Name != nil {
			nsg["name"] = *n.Name
		}
		nsgs = append(nsgs, nsg)
	}

	return nsgs
}

func flattenPrivateDNSZones(zoneList network.PrivateDNSZoneGroupListResultPage) []interface{} {
	z := zoneList.Values()
	zones := make([]interface{}, 0)
	if z == nil {
		return zones
	}
	for _, item := range z {
		zone := make(map[string]string)
		if item.ID != nil {
			zone["id"] = *item.ID
		}

		if item.Name != nil {
			zone["name"] = *item.Name
		}
		zones = append(zones, zone)
	}

	return zones
}

func flattenPrivateEndpoints(endpointList network.PrivateEndpointListResultPage) []interface{} {
	e := endpointList.Values()
	endpoints := make([]interface{}, 0)
	if e == nil {
		return endpoints
	}
	for _, item := range e {
		newEndpoint := make(map[string]string)
		if item.ID != nil {
			newEndpoint["id"] = *item.ID
		}

		if item.Name != nil {
			newEndpoint["name"] = *item.Name
		}
		endpoints = append(endpoints, newEndpoint)
	}
	return endpoints
}

func flattenRouteTables(tableList network.RouteTableListResultPage) []interface{} {
	tables := tableList.Values()
	routeTables := make([]interface{}, 0)
	if tables == nil {
		return routeTables
	}
	for _, n := range tables {
		routeTable := make(map[string]string)
		if n.ID != nil {
			routeTable["id"] = *n.ID
		}

		if n.Name != nil {
			routeTable["name"] = *n.Name
		}
		routeTables = append(routeTables, routeTable)
	}

	return routeTables
}
