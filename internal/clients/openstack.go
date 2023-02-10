/*
Copyright 2021 Upbound Inc.
*/

package clients

import (
	"context"
	"encoding/json"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/upbound/upjet/pkg/terraform"

	"github.com/fcosta999/provider-openstack/apis/v1beta1"
)

const (
    KeyAuthUrl    = "auth_url"
    KeyPassword   = "password"
    KeyRegion     = "region"
    KeyTenantName = "tenant_name"
    KeyUserName   = "user_name"
)

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		configRef := mg.GetProviderConfigReference()
		if configRef == nil {
			return ps, errors.New(errNoProviderConfig)
		}
		pc := &v1beta1.ProviderConfig{}
		if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
			return ps, errors.Wrap(err, errGetProviderConfig)
		}

		t := resource.NewProviderConfigUsageTracker(client, &v1beta1.ProviderConfigUsage{})
		if err := t.Track(ctx, mg); err != nil {
			return ps, errors.Wrap(err, errTrackUsage)
		}

		data, err := resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, client, pc.Spec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}
		creds := map[string]string{}
		if err := json.Unmarshal(data, &creds); err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		// set provider configuration
        ps.Configuration = map[string]any{}
        if v, ok := creds[KeyAuthUrl]; ok {
            ps.Configuration[KeyAuthUrl] = v
        }
        if v, ok := creds[KeyPassword]; ok {
            ps.Configuration[KeyPassword] = v
        }
        if v, ok := creds[KeyRegion]; ok {
            ps.Configuration[KeyRegion] = v
        }
        if v, ok := creds[KeyTenantName]; ok {
            ps.Configuration[KeyTenantName] = v
        }
        if v, ok := creds[KeyUserName]; ok {
            ps.Configuration[KeyUserName] = v
        }
		return ps, nil
	}
}
