package di

import "salehelp/internal/storage/pggorm/repository/product"

type Repositories struct {
	Product *product.Repository
}

func (prov *Provider) InitRepositories() *Repositories {
	return &Repositories{
		Product: product.NewRepository(prov.log, prov.Database.GormDB),
	}
}
