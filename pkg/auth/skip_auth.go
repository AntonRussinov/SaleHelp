package auth

// type SkipAuthenticator struct{}

// func NewSkipAuthenticator() *SkipAuthenticator {
// 	return &SkipAuthenticator{}
// }

// func DevUser() *domain.User {
// 	userdata := DevUserData()

// 	return &domain.User{
// 		ID:    userdata.ID,
// 		Name:  userdata.Name,
// 		Email: userdata.Email,
// 		Roles: userdata.Roles,
// 	}
// }

// func DevUserContext(ctx context.Context) context.Context {
// 	return SetUserToContext(ctx, DevUser())
// }

// func DevUserData() *UserData {
// 	return &UserData{
// 		ID:    uuid.MustParse("4f7cd908-c2ab-4e19-8a57-959b4e59d367"),
// 		Name:  "Dev User",
// 		Email: "dev-user@test.test",
// 		Roles: []domain.UserRole{domain.UserRoleSuperAdmin},
// 	}
// }

// func (*SkipAuthenticator) Authenticate(_ context.Context, _ string) (*UserData, string, error) {
// 	user := DevUserData()

// 	return user, "", nil
// }

// func (s *SkipAuthenticator) AuthenticateWithHeader(_ *http.Request) (*UserData, string, error) {
// 	return s.Authenticate(context.TODO(), "")
// }

// func (s *SkipAuthenticator) AuthenticateWithQuery(_ *http.Request) (*UserData, string, error) {
// 	return s.Authenticate(context.TODO(), "")
// }

// func (s *SkipAuthenticator) AuthenticateWithCookie(_ *http.Request) (*UserData, string, error) {
// 	return s.Authenticate(context.TODO(), "")
// }
