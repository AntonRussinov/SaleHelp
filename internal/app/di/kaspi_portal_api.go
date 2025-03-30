package di

import (
	portalAPI "salehelp/internal/adapter/kaspi/portal/connector"
)

func (prov *Provider) InitKaspiPortalAPI() error {
	prov.KaspiPortalAPI = portalAPI.New()

	return nil
}
