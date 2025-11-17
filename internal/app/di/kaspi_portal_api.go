package di

import (
	portalAPI "salehelp/internal/adapter/kaspi/portal/connector"
	"salehelp/internal/adapter/kaspi/portal/connector/client"
)

func (prov *Provider) InitKaspiPortalAPI(MerchantRepo client.MerchantRepo) error {
	prov.KaspiPortalAPI = portalAPI.New(MerchantRepo)

	return nil
}
