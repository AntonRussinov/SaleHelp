package access_test

// type (
// 	initPoliciesTestCase struct {
// 		name          string
// 		savePolicyErr error
// 		loadPolicyErr error
// 		addPolicyErr  error
// 		wantErr       bool
// 	}
// 	updateUserRoleTestCase struct {
// 		name           string
// 		paramUser      *domain.User
// 		getRolesReturn []string
// 		getRolesErr    error
// 		deleteRoleErr  error
// 		addRoleErr     error
// 		wantErr        bool
// 	}

// 	middlewareTestCase struct {
// 		name                string
// 		userFunc            func(ctx context.Context) *domain.User
// 		enforceCalled       bool
// 		enforceParam1       string
// 		enforceParam2       string
// 		enforceParam3       string
// 		enforceReturn       bool
// 		enforceErr          error
// 		wantErr             bool
// 		wantErrUnauthorized bool
// 		wantErrForbidden    bool
// 		wantErrInternal     bool
// 	}
// )

// const (
// 	validModel   = "../../test/data/access/rbac_valid.conf"
// 	invalidModel = "../../test/data/access/rbac_invalid.conf"
// )

// func absFilePath(t *testing.T, file string) string {
// 	t.Helper()

// 	abspath, err := filepath.Abs(file)

// 	require.NoError(t, err)

// 	return abspath
// }

// func tempPolicyFile(t *testing.T) string {
// 	t.Helper()

// 	dir := t.TempDir()

// 	file, err := os.Create(dir + "/data.csv")

// 	require.NoError(t, err)

// 	return file.Name()
// }

// func TestAccessGormAdapter(t *testing.T) {
// 	t.Parallel()

// 	t.Run("success", func(t *testing.T) {
// 		t.Parallel()

// 		conn, mock, err := sqlmock.New()
// 		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
// 		mock.ExpectExec("CREATE").WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))
// 		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
// 		mock.ExpectExec("CREATE").WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))
// 		require.NoError(t, err)

// 		dialector := postgres.New(postgres.Config{Conn: conn})

// 		db, err := gorm.Open(dialector)
// 		require.NoError(t, err)

// 		adapter, err := access.GormAdapter(db)
// 		require.NoError(t, err)
// 		require.NotNil(t, adapter)
// 	})

// 	t.Run("error", func(t *testing.T) {
// 		t.Parallel()

// 		conn, _, err := sqlmock.New()
// 		require.NoError(t, err)

// 		dialector := postgres.New(postgres.Config{Conn: conn})

// 		db, err := gorm.Open(dialector)
// 		require.NoError(t, err)

// 		adapter, err := access.GormAdapter(db)
// 		require.Error(t, err)
// 		require.Nil(t, adapter)
// 	})
// }

// func TestAccessNewManager(t *testing.T) {
// 	t.Parallel()

// 	log := slogdiscard.NewDiscardLogger()

// 	t.Run("wrong params", func(t *testing.T) {
// 		t.Parallel()

// 		manager, err := access.NewManager(log, access.WithAdapter(""))
// 		require.Error(t, err)
// 		require.Errorf(t, err, "access.New: invalid parameters for enforcer")
// 		require.Nil(t, manager)
// 	})

// 	t.Run("invalid model", func(t *testing.T) {
// 		t.Parallel()

// 		manager, err := access.NewManager(log,
// 			access.WithModel(absFilePath(t, invalidModel)),
// 			access.WithAdapter(tempPolicyFile(t)))
// 		require.Error(t, err)
// 		require.Errorf(t, err, "access.New: invalid parameters for enforcer")
// 		require.Nil(t, manager)
// 	})

// 	t.Run("success", func(t *testing.T) {
// 		t.Parallel()

// 		manager, err := access.NewManager(log,
// 			access.WithModel(absFilePath(t, validModel)),
// 			access.WithAdapter(tempPolicyFile(t)))
// 		require.NoError(t, err)
// 		require.NotNil(t, manager)
// 	})
// }

// func TestAccessManagerUpdateUserRole(t *testing.T) {
// 	t.Parallel()

// 	manager, err := access.NewManager(
// 		slogdiscard.NewDiscardLogger(),
// 		access.WithModel(absFilePath(t, validModel)),
// 		access.WithAdapter(tempPolicyFile(t)))
// 	require.NoError(t, err)
// 	require.NotNil(t, manager)

// 	user1 := domain.User{
// 		ID:   uuid.New(),
// 		Role: domain.UserRoleAdmin,
// 	}

// 	err = manager.UpdateUserRole(&user1)
// 	require.NoError(t, err)

// 	roles, err := manager.GetUserRoles(&user1)
// 	require.NoError(t, err)
// 	require.ElementsMatch(t, []domain.UserRole{domain.UserRoleAdmin}, roles)

// 	user1.Role = domain.UserRoleUser
// 	err = manager.UpdateUserRole(&user1)
// 	require.NoError(t, err)

// 	roles, err = manager.GetUserRoles(&user1)
// 	require.NoError(t, err)
// 	require.ElementsMatch(t, []domain.UserRole{domain.UserRoleUser}, roles)

// 	err = manager.UpdateUserRole(&user1)
// 	require.NoError(t, err)

// 	roles, err = manager.GetUserRoles(&user1)
// 	require.NoError(t, err)
// 	require.ElementsMatch(t, []domain.UserRole{domain.UserRoleUser}, roles)
// }

// func updateUserRolesTestCases() []updateUserRoleTestCase {
// 	faker := gofakeit.New(1)
// 	userID := uuid.MustParse(faker.UUID())

// 	return []updateUserRoleTestCase{
// 		{
// 			name:        "returns error if get roles failed",
// 			paramUser:   &domain.User{ID: userID},
// 			getRolesErr: assert.AnError,
// 			wantErr:     true,
// 		},
// 		{
// 			name:           "returns error if failed to delete existing user role",
// 			paramUser:      &domain.User{ID: userID, Role: domain.UserRoleAdmin},
// 			getRolesReturn: []string{"SandStorm.Admin", "SandStorm.Editor"},
// 			deleteRoleErr:  assert.AnError,
// 			wantErr:        true,
// 		},
// 		{
// 			name:           "returns error if failed to save new role",
// 			paramUser:      &domain.User{ID: userID, Role: domain.UserRoleAdmin},
// 			getRolesReturn: []string{"Sandstorm.Editor"},
// 			addRoleErr:     assert.AnError,
// 			wantErr:        true,
// 		},
// 		{
// 			name:           "returns no error on success",
// 			paramUser:      &domain.User{ID: userID, Role: domain.UserRoleAdmin},
// 			getRolesReturn: []string{"Sandstorm.Editor"},
// 		},
// 	}
// }

// func TestAccessManagerUpdateUserRoleErrors(t *testing.T) {
// 	t.Parallel()

// 	tests := updateUserRolesTestCases()

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			enforcer := casbin_mocks.NewMockIEnforcer(t)

// 			initPoliciesSuccessExpectations(t, enforcer)

// 			manager, err := access.NewManager(slogdiscard.NewDiscardLogger(), access.WithCustomEnforcer(enforcer))
// 			require.NoError(t, err)

// 			updateUserRoleExpectations(t, enforcer, tt)

// 			err = manager.UpdateUserRole(tt.paramUser)
// 			if tt.wantErr {
// 				require.Error(t, err)
// 			} else {
// 				require.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func updateUserRoleExpectations(t *testing.T, enforcer *casbin_mocks.MockIEnforcer, tt updateUserRoleTestCase) {
// 	t.Helper()

// 	enforcer.EXPECT().GetRolesForUser(tt.paramUser.ID.String()).Return(tt.getRolesReturn, tt.getRolesErr)

// 	if tt.getRolesErr != nil {
// 		return
// 	}

// 	newRole := string(tt.paramUser.Role)

// 	hasRole := false

// 	for _, r := range tt.getRolesReturn {
// 		if r == newRole {
// 			hasRole = true
// 			continue
// 		}

// 		enforcer.EXPECT().DeleteRoleForUser(tt.paramUser.ID.String(), r).Return(true, tt.deleteRoleErr)

// 		if tt.deleteRoleErr != nil {
// 			return
// 		}
// 	}

// 	if !hasRole {
// 		enforcer.EXPECT().AddRoleForUser(tt.paramUser.ID.String(), newRole).Return(true, tt.addRoleErr)
// 	}
// }

// func TestAccessManagerGetUserRoles(t *testing.T) {
// 	t.Parallel()

// 	faker := gofakeit.New(1)
// 	userID := uuid.MustParse(faker.UUID())

// 	t.Run("returns error if enforcer get roles failes", func(t *testing.T) {
// 		t.Parallel()

// 		enforcer := casbin_mocks.NewMockIEnforcer(t)
// 		initPoliciesSuccessExpectations(t, enforcer)
// 		manager, err := access.NewManager(slogdiscard.NewDiscardLogger(), access.WithCustomEnforcer(enforcer))
// 		require.NoError(t, err)

// 		enforcer.EXPECT().GetRolesForUser(userID.String()).Return(nil, assert.AnError)

// 		roles, err := manager.GetUserRoles(&domain.User{ID: userID})
// 		require.Error(t, err)
// 		require.Nil(t, roles)
// 	})

// 	t.Run("returns domain matched user roles on success", func(t *testing.T) {
// 		t.Parallel()

// 		enforcer := casbin_mocks.NewMockIEnforcer(t)
// 		initPoliciesSuccessExpectations(t, enforcer)
// 		manager, err := access.NewManager(slogdiscard.NewDiscardLogger(), access.WithCustomEnforcer(enforcer))
// 		require.NoError(t, err)

// 		enforcer.EXPECT().GetRolesForUser(userID.String()).Return([]string{"SandStorm.Admin", "SandStorm.Editor"}, nil)

// 		roles, err := manager.GetUserRoles(&domain.User{ID: userID})
// 		require.NoError(t, err)
// 		require.ElementsMatch(t, []domain.UserRole{domain.UserRoleAdmin, domain.UserRoleEditor}, roles)
// 	})
// }

// func TestAccessManagerAllowed(t *testing.T) {
// 	t.Parallel()

// 	manager, err := access.NewManager(
// 		slogdiscard.NewDiscardLogger(),
// 		access.WithModel(absFilePath(t, validModel)),
// 		access.WithAdapter(tempPolicyFile(t)))
// 	require.NoError(t, err)
// 	require.NotNil(t, manager)

// 	user1 := domain.User{
// 		ID:   uuid.New(),
// 		Role: domain.UserRoleAdmin,
// 	}

// 	ok, err := manager.Allowed(&user1, access.ResourceUsers, access.ActionDelete)
// 	require.NoError(t, err)

// 	require.False(t, ok)
// }

// func initPoliciesTestCases() []initPoliciesTestCase {
// 	return []initPoliciesTestCase{
// 		{
// 			name:          "returns error if failed to clear existing policies",
// 			savePolicyErr: assert.AnError,
// 			wantErr:       true,
// 		},
// 		{
// 			name:          "returns error if failed to load policy after clearing",
// 			savePolicyErr: nil,
// 			loadPolicyErr: assert.AnError,
// 			wantErr:       true,
// 		},
// 		{
// 			name:          "returns error if failed to add default policies",
// 			savePolicyErr: nil,
// 			loadPolicyErr: nil,
// 			addPolicyErr:  assert.AnError,
// 			wantErr:       true,
// 		},
// 		{
// 			name: "returns no error on success",
// 		},
// 	}
// }

// func TestAccessManagerInitPolicies(t *testing.T) {
// 	t.Parallel()

// 	tests := initPoliciesTestCases()

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			enforcer := casbin_mocks.NewMockIEnforcer(t)

// 			initPoliciesExpectations(t, enforcer, tt)

// 			manager, err := access.NewManager(slogdiscard.NewDiscardLogger(), access.WithCustomEnforcer(enforcer))

// 			if tt.wantErr {
// 				require.Error(t, err)
// 				require.Nil(t, manager)
// 			} else {
// 				require.NoError(t, err)
// 				require.NotNil(t, manager)
// 			}
// 		})
// 	}
// }

// func initPoliciesSuccessExpectations(t *testing.T, enforcer *casbin_mocks.MockIEnforcer) {
// 	t.Helper()

// 	initPoliciesExpectations(t, enforcer, initPoliciesTestCase{})
// }

// func initPoliciesExpectations(t *testing.T, enforcer *casbin_mocks.MockIEnforcer, tt initPoliciesTestCase) {
// 	t.Helper()

// 	enforcer.EXPECT().ClearPolicy()

// 	enforcer.EXPECT().SavePolicy().Return(tt.savePolicyErr)

// 	if tt.savePolicyErr != nil {
// 		return
// 	}

// 	enforcer.EXPECT().LoadPolicy().Return(tt.loadPolicyErr)

// 	if tt.loadPolicyErr != nil {
// 		return
// 	}

// 	defaultPolicies := access.DefaultPolicies()

// 	// switching true | false for has policy every time
// 	hasPolicy := false
// 	for _, pol := range defaultPolicies {
// 		enforcer.EXPECT().HasPolicy(string(pol.Role), string(pol.Res), string(pol.Act)).Return(hasPolicy, nil)

// 		if !hasPolicy {
// 			enforcer.EXPECT().AddPolicy(string(pol.Role), string(pol.Res), string(pol.Act)).Return(true, tt.addPolicyErr)
// 		}

// 		if tt.addPolicyErr != nil {
// 			return
// 		}

// 		hasPolicy = !hasPolicy
// 	}
// }

// func middlewareTestCases() []middlewareTestCase {
// 	return []middlewareTestCase{
// 		{
// 			name:                "returns unauthorized error if user func returns nil user",
// 			userFunc:            func(ctx context.Context) *domain.User { return nil },
// 			enforceCalled:       false,
// 			wantErr:             true,
// 			wantErrUnauthorized: true,
// 		},
// 		{
// 			name:            "returns internal server error if enforcer fails to enforce rule",
// 			userFunc:        func(ctx context.Context) *domain.User { return auth.DevUser() },
// 			enforceCalled:   true,
// 			enforceParam1:   auth.DevUser().ID.String(),
// 			enforceParam2:   string(access.ResourceAll),
// 			enforceParam3:   string(access.ActionManage),
// 			enforceReturn:   false,
// 			enforceErr:      assert.AnError,
// 			wantErr:         true,
// 			wantErrInternal: true,
// 		},
// 		{
// 			name:             "returns forbidden error if enforcer returns false",
// 			userFunc:         func(ctx context.Context) *domain.User { return auth.DevUser() },
// 			enforceCalled:    true,
// 			enforceParam1:    auth.DevUser().ID.String(),
// 			enforceParam2:    string(access.ResourceAll),
// 			enforceParam3:    string(access.ActionManage),
// 			enforceReturn:    false,
// 			enforceErr:       nil,
// 			wantErr:          true,
// 			wantErrForbidden: true,
// 		},
// 		{
// 			name:          "returns no error and calls handler on success",
// 			userFunc:      func(ctx context.Context) *domain.User { return auth.DevUser() },
// 			enforceCalled: true,
// 			enforceParam1: auth.DevUser().ID.String(),
// 			enforceParam2: string(access.ResourceAll),
// 			enforceParam3: string(access.ActionManage),
// 			enforceReturn: true,
// 		},
// 	}
// }

// func TestAccessManagerMiddleware(t *testing.T) {
// 	t.Parallel()

// 	tests := middlewareTestCases()

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			enforcer := casbin_mocks.NewMockIEnforcer(t)

// 			initPoliciesSuccessExpectations(t, enforcer)

// 			manager, err := access.NewManager(slogdiscard.NewDiscardLogger(), access.WithCustomEnforcer(enforcer),
// 				access.WithUserFunc(tt.userFunc))
// 			require.NoError(t, err)

// 			middleware := manager.Middleware(access.ResourceAll, access.ActionManage)

// 			req := httptest.NewRequest(http.MethodGet, "/", nil)
// 			rec := httptest.NewRecorder()
// 			c := echo.New().NewContext(req, rec)

// 			called := false
// 			handler := middleware(func(c echo.Context) error {
// 				if tt.wantErr {
// 					t.Fatal("handler should not be called if middleware returns err")
// 				}

// 				called = true

// 				return nil
// 			})

// 			if tt.enforceCalled {
// 				enforcer.EXPECT().Enforce(tt.enforceParam1, tt.enforceParam2, tt.enforceParam3).
// 					Return(tt.enforceReturn, tt.enforceErr).Once()
// 			}

// 			err = handler(c)

// 			switch {
// 			case tt.wantErrUnauthorized:
// 				require.ErrorIs(t, err, echo.ErrUnauthorized)
// 			case tt.wantErrForbidden:
// 				require.ErrorIs(t, err, echo.ErrForbidden)
// 			case tt.wantErrInternal:
// 				require.ErrorIs(t, err, echo.ErrInternalServerError)
// 			}

// 			if !tt.wantErr {
// 				require.True(t, called, "handler should be called if no error in middleware")
// 			}
// 		})
// 	}
// }
