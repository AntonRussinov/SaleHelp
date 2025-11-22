package access

// type (
// 	Manager struct {
// 		model    string
// 		adapter  any
// 		enforcer casbin.IEnforcer
// 		userFunc func(ctx context.Context) *domain.User
// 	}

// 	IManager interface {
// 		GetUserRoles(user *domain.User) ([]domain.UserRole, error)
// 		UpdateUserRole(user *domain.User) error
// 		Allowed(user *domain.User, res Resource, action Action) (bool, error)
// 		AddPolicy(role domain.UserRole, res Resource, act Action) error
// 		Middleware(res Resource, action Action) echo.MiddlewareFunc
// 	}

// 	Config struct {
// 		Model       string
// 		Adapter     any
// 		GetUserFunc func(ctx context.Context) *domain.User
// 	}

// 	Option func(m *Manager)
// )

// func GormAdapter(db *gorm.DB) (*gormadapter.Adapter, error) {
// 	const op = "access.GormAdapter"

// 	adapter, err := gormadapter.NewAdapterByDB(db)
// 	if err != nil {
// 		return nil, wrapError(op, err)
// 	}

// 	return adapter, nil
// }

// func NewManager(log *slog.Logger, opts ...Option) (*Manager, error) {
// 	const fn = "Access.NewManager"

// 	manager := &Manager{}

// 	for _, opt := range opts {
// 		opt(manager)
// 	}

// 	log.Debug("initializing casbin enforcer with model",
// 		slog.String("fn", fn),
// 		slog.String("model", manager.model))

// 	if manager.enforcer == nil {
// 		enforcer, err := casbin.NewEnforcer(manager.model, manager.adapter)
// 		if err != nil {
// 			return nil, wrapError(fn, err)
// 		}

// 		manager.enforcer = enforcer
// 	}

// 	err := manager.initPolicies()
// 	if err != nil {
// 		return nil, wrapError(fn, err)
// 	}

// 	return manager, nil
// }

// func (m *Manager) initPolicies() error {
// 	m.enforcer.ClearPolicy()

// 	err := m.enforcer.SavePolicy()
// 	if err != nil {
// 		return fmt.Errorf("failed to save policies: %w", err)
// 	}

// 	if err := m.enforcer.LoadPolicy(); err != nil {
// 		return err
// 	}

// 	policies := DefaultPolicies()

// 	for _, policy := range policies {
// 		err := m.AddPolicy(policy.Role, policy.Res, policy.Act)
// 		if err != nil {
// 			return fmt.Errorf("failed to add policy: %w", err)
// 		}
// 	}

// 	return nil
// }

// func (m *Manager) GetUserRoles(user *domain.User) ([]domain.UserRole, error) {
// 	id := user.ID.String()

// 	roles, err := m.enforcer.GetRolesForUser(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	domainRoles := make([]domain.UserRole, len(roles))
// 	for i, r := range roles {
// 		domainRoles[i] = domain.UserRole(r)
// 	}

// 	return domainRoles, nil
// }

// func (m *Manager) UpdateUserRole(user *domain.User) error {
// 	id := user.ID.String()

// 	newRole := string(user.Role)

// 	currentRoles, err := m.enforcer.GetRolesForUser(id)
// 	if err != nil {
// 		return err
// 	}

// 	hasRole := false

// 	for _, r := range currentRoles {
// 		if r == newRole {
// 			hasRole = true
// 			continue
// 		}

// 		_, err := m.enforcer.DeleteRoleForUser(id, r)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	if !hasRole {
// 		_, err = m.enforcer.AddRoleForUser(id, newRole)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (m *Manager) Allowed(user *domain.User, res Resource, action Action) (bool, error) {
// 	return m.enforcer.Enforce(user.ID.String(), string(res), string(action))
// }

// func (m *Manager) AddPolicy(role domain.UserRole, res Resource, act Action) error {
// 	hasPolicy, err := m.enforcer.HasPolicy(string(role), string(res), string(act))
// 	if err != nil {
// 		return err
// 	}

// 	if !hasPolicy {
// 		_, err := m.enforcer.AddPolicy(string(role), string(res), string(act))
// 		if err != nil {
// 			return fmt.Errorf("failed to add policy: %w", err)
// 		}
// 	}

// 	return nil
// }

// func (m *Manager) Middleware(res Resource, action Action) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			user := m.userFunc(c.Request().Context())
// 			if user == nil {
// 				return echo.ErrUnauthorized
// 			}

// 			ok, err := m.Allowed(user, res, action)
// 			if err != nil {
// 				return echo.ErrInternalServerError
// 			}

// 			if !ok {
// 				return echo.ErrForbidden
// 			}

// 			return next(c)
// 		}
// 	}
// }

// func wrapError(fn string, err error) error {
// 	return fmt.Errorf("%s: %w", fn, err)
// }
