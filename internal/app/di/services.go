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
	prov.Services.Product = product.NewService(prov.log, prov.Repositories.Product, prov.KaspiPortalAPI)
}
