package di

import (
	kaspiAPI "salehelp/internal/adapter/kaspi/api/connector"
)

func (prov *Provider) InitKaspiAPI() error {
	prov.KaspiAPI = kaspiAPI.New()

	return nil
}
