package di

type Services struct {
	//User *user.Service
	//Limit         *limit.Service

}

func (prov *Provider) InitServices() {
	prov.initIndependentServices()
	prov.initDependentServices()
}

func (prov *Provider) initIndependentServices() {

}

func (prov *Provider) initDependentServices() {

}
