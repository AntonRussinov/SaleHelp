package product_test

// import (
// 	"salehelp/domain"
// 	"salehelp/internal/storage"
// 	"context"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// )

// func TestGetLimit(t *testing.T) {
// 	t.Parallel()

// 	lim := domain.Limit{User: &domain.User{ID: uuid.New()}, Active: true, Limit: 12}

// 	tests := []struct {
// 		name         string
// 		serviceParam uuid.UUID
// 		repoReturn   domain.Limit
// 		repoError    error
// 		wantErr      bool
// 		wantErrAs    any
// 		wantReturn   *domain.Limit
// 	}{
// 		{name: "success", serviceParam: lim.User.ID, repoReturn: lim, wantReturn: &lim},
// 		{name: "err not found", serviceParam: uuid.New(), repoError: storage.ErrNotFound, wantErr: true, wantErrAs: &domain.NotFoundError{}},
// 		{name: "other error", serviceParam: lim.User.ID, repoError: assert.AnError, wantErr: true, wantErrAs: &domain.APPError{}},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			test := setupServiceTest(t)

// 			repoReturn := tt.repoReturn
// 			test.repo.EXPECT().FindByID(mock.Anything, tt.serviceParam).Return(&repoReturn, tt.repoError).Once()
// 			res, err := test.srv.GetLimit(context.Background(), tt.serviceParam)

// 			if tt.wantErr {
// 				require.Error(t, err)
// 			} else {
// 				require.NoError(t, err)
// 			}

// 			wantErrAs := tt.wantErrAs
// 			if tt.wantErrAs != nil {
// 				require.ErrorAs(t, err, &wantErrAs)
// 			}

// 			assert.Equal(t, tt.wantReturn, res)
// 		})
// 	}
// }
