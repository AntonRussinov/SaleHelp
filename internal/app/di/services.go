package di

import "salehelp/internal/services/product"

type Services struct {
	//User *user.Service
	//Limit         *limit.Service
	Product *product.Service
}

func (prov *Provider) InitServices() {
	prov.initIndependentServices()
	prov.initDependentServices()
}

func (prov *Provider) initIndependentServices() {

}

func (prov *Provider) initDependentServices() {
	// prov.Services.Product = product.NewService(
	// 	prov.log, prov.Repositories.Stock, prov.Repositories.Product,
	// 	prov.Repositories.CityPrice, prov.Repositories.Offer,
	// 	prov.Repositories.Merchant, prov.Repositories.Image,
	// 	prov.KaspiPortalAPI)
}
