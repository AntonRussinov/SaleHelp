package di

type Repositories struct {
	//Application *application.Repositories
}

func (prov *Provider) InitRepositories() *Repositories {
	return &Repositories{
		//Application:            application.NewRepositories(log, db),
	}
}
