package azureparse

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-06-01/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
)

type client struct {
	StopContext context.Context

	PrivateEndpointsClient *network.PrivateEndpointsClient
	PrivateDNSZonesClient  *network.PrivateDNSZoneGroupsClient
	ResourceGroupsClient   *resources.GroupsClient
	RouteTablesClient      *network.RouteTablesClient
	SecurityGroupsClient   *network.SecurityGroupsClient
}

func buildClient(config *authentication.Config) (*client, error) {

	if config == nil {
		return nil, fmt.Errorf("error build config is nil: %v", config)
	}

	sender := sender.BuildSender("tfazureparse")

	env, err := authentication.DetermineEnvironment(config.Environment)
	if err != nil {
		return nil, fmt.Errorf("error determining environment: %v", err)
	}

	oauthConfig, err := config.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error building oauth config: %v", err)
	}

	auth, err := config.GetAuthorizationToken(sender, oauthConfig, env.TokenAudience)
	if err != nil {
		return nil, fmt.Errorf("error retrieving auth token: %v", err)
	}

	subscriptionID := config.SubscriptionID

	privateDNSZonesClient := network.NewPrivateDNSZoneGroupsClient(subscriptionID)
	privateEndpointsClient := network.NewPrivateEndpointsClient(subscriptionID)
	resourceGroupsClient := resources.NewGroupsClient(subscriptionID)
	routeTablesClient := network.NewRouteTablesClient(subscriptionID)
	securityGroupsClient := network.NewSecurityGroupsClient(subscriptionID)

	privateDNSZonesClient.Authorizer = auth
	privateEndpointsClient.Authorizer = auth
	resourceGroupsClient.Authorizer = auth
	routeTablesClient.Authorizer = auth
	securityGroupsClient.Authorizer = auth

	return &client{
		PrivateDNSZonesClient:  &privateDNSZonesClient,
		PrivateEndpointsClient: &privateEndpointsClient,
		ResourceGroupsClient:   &resourceGroupsClient,
		RouteTablesClient:      &routeTablesClient,
		SecurityGroupsClient:   &securityGroupsClient,
	}, nil
}
