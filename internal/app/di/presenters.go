package di

import "salehelp/internal/presenter/product"

type Presenters struct {
	Product *product.Presenter
}

func (prov *Provider) InitPresenters() {
	prov.Presenters = Presenters{
		Product: product.NewPresenter(),
	}
}
