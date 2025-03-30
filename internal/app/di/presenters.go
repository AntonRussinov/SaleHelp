package di

import "salehelp/internal/presenter/product"

type Presenters struct {
	Product *product.Presenter
}

func (prov *Provider) InitPresenters() *Presenters {
	return &Presenters{
		Product: product.NewPresenter(),
	}
}
