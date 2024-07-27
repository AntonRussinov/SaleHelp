package di

type Presenters struct {
	//Application     *application.Presenters
}

func (prov *Provider) InitPresenters() *Presenters {
	return &Presenters{
		//Application:     application.NewPresenters(),
	}
}
