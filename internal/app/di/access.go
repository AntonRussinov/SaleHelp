package di

// var (
// 	ErrInitAccessManagerAdapter = errors.New("failed to initialize access manager database adapter")
// 	ErrInitAccessManager        = errors.New("failed to initialize access manager")
// )

// func (prov *Provider) InitAccessManager(ctx context.Context) error {
// 	adapter, err := access.GormAdapter(prov.Database.GormDB)
// 	if err != nil {
// 		return fmt.Errorf("%w: %w", ErrInitAccessManagerAdapter, err)
// 	}

// 	manager, err := access.NewManager(
// 		prov.log,
// 		access.WithAdapter(adapter),
// 		access.WithModel("configs/rbac_model.conf"),
// 		access.WithUserFunc(auth.GetUserFromContext))
// 	if err != nil {
// 		return fmt.Errorf("%w: %w", ErrInitAccessManager, err)
// 	}

// 	prov.AccessManager = manager

// 	return nil
// }
