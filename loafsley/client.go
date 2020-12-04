package loafsley

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-06-01/network"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
)

type Client struct {
	StopContext context.Context

	PrivateEndpointsClient *network.PrivateEndpointsClient
	PrivateDNSZonesClient  *network.PrivateDNSZoneGroupsClient
	RouteTablesClient      *network.RouteTablesClient
	SecurityGroupsClient   *network.SecurityGroupsClient
}

func Build(config *authentication.Config) (*Client, error) {

	if config == nil {
		return nil, fmt.Errorf("error build config is nil: %v", config)
	}

	sender := sender.BuildSender("AzureRM")

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

	subscriptionId := config.SubscriptionID
	privateDNSZonesClient := network.NewPrivateDNSZoneGroupsClient(subscriptionId)
	privateEndpointsClient := network.NewPrivateEndpointsClient(subscriptionId)
	routeTablesClient := network.NewRouteTablesClient(subscriptionId)
	securityGroupsClient := network.NewSecurityGroupsClient(subscriptionId)

	privateDNSZonesClient.Authorizer = auth
	privateEndpointsClient.Authorizer = auth
	routeTablesClient.Authorizer = auth
	securityGroupsClient.Authorizer = auth

	return &Client{
		PrivateDNSZonesClient:  &privateDNSZonesClient,
		PrivateEndpointsClient: &privateEndpointsClient,
		RouteTablesClient:      &routeTablesClient,
		SecurityGroupsClient:   &securityGroupsClient,
	}, nil
}
