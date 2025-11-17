package di

import (
	"salehelp/internal/storage/pggorm/repository/cityprice"
	"salehelp/internal/storage/pggorm/repository/image"
	"salehelp/internal/storage/pggorm/repository/merchant"
	"salehelp/internal/storage/pggorm/repository/offer"
	"salehelp/internal/storage/pggorm/repository/product"
	"salehelp/internal/storage/pggorm/repository/stock"
)

type Repositories struct {
	Product   *product.Repository
	Stock     *stock.Repository
	Image     *image.Repository
	Merchant  *merchant.Repository
	CityPrice *cityprice.Repository
	Offer     *offer.Repository
}

func (prov *Provider) InitRepositories() *Repositories {
	return &Repositories{
		Product:   product.NewRepository(prov.log, prov.Database.GormDB),
		Stock:     stock.NewRepository(prov.log, prov.Database.GormDB),
		Image:     image.NewRepository(prov.log, prov.Database.GormDB),
		Merchant:  merchant.NewRepository(prov.log, prov.Database.GormDB),
		CityPrice: cityprice.NewRepository(prov.log, prov.Database.GormDB),
		Offer:     offer.NewRepository(prov.log, prov.Database.GormDB),
	}
}
