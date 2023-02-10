/*
Copyright 2021 Upbound Inc.
*/

package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/upbound/upjet/pkg/config"

	//"github.com/fcosta999/provider-openstack/config/null"
	"github.com/fcosta999/provider-openstack/config/db_user_v1"
	"github.com/fcosta999/provider-openstack/config/identity_application_credential_v3"
)

const (
	resourcePrefix = "openstack"
	modulePath     = "github.com/fcosta999/provider-openstack"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		db_user_v1.Configure,
		identity_application_credential_v3.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
