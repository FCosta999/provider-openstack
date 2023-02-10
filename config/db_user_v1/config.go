package db_user_v1

import "github.com/upbound/upjet/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
    p.AddResourceConfigurator("openstack_db_user_v1", func(r *config.Resource) {
        // We need to override the default group that upjet generated for
        // this resource, which would be "github"
        r.ShortGroup = "db_user_v1"
    })
}
