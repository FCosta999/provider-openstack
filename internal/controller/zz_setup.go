/*
Copyright 2021 Upbound Inc.
*/

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/upbound/upjet/pkg/controller"

	userv1 "github.com/fcosta999/provider-openstack/internal/controller/db_user_v1/userv1"
	applicationcredentialv3 "github.com/fcosta999/provider-openstack/internal/controller/identity_application_credential_v3/applicationcredentialv3"
	providerconfig "github.com/fcosta999/provider-openstack/internal/controller/providerconfig"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		userv1.Setup,
		applicationcredentialv3.Setup,
		providerconfig.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
