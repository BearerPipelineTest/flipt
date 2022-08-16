package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	flipt "go.flipt.io/flipt/rpc/flipt"
	"go.flipt.io/flipt/storage"
)

func TestGetFlag(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
		req = &flipt.GetFlagRequest{Key: "foo"}
	)

	store.On("GetFlag", mock.Anything, "foo").Return(&flipt.Flag{
		Key:     req.Key,
		Enabled: true,
	}, nil)

	got, err := s.GetFlag(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)
	assert.Equal(t, "foo", got.Key)
	assert.Equal(t, true, got.Enabled)
}

func TestListFlags(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
	)

	store.On("ListFlags", mock.Anything, mock.Anything).Return(
		[]*flipt.Flag{
			{
				Key: "foo",
			},
		}, nil)

	got, err := s.ListFlags(context.TODO(), &flipt.ListFlagRequest{})
	require.NoError(t, err)

	assert.NotEmpty(t, got.Flags)
}

func TestListFlags_Pagination(t *testing.T) {
	tests := []struct {
		name     string
		offset   int32
		limit    int32
		expected storage.QueryParams
	}{
		{
			name: "default/no pagination",
			expected: storage.QueryParams{
				Offset: 0,
				Limit:  20,
			},
		},
		{
			name:   "negative offset",
			offset: -1,
			expected: storage.QueryParams{
				Offset: 0,
				Limit:  20,
			},
		},
		{
			name:  "negative limit",
			limit: -1,
			expected: storage.QueryParams{
				Offset: 0,
				Limit:  20,
			},
		},
		{
			name:  "zero limit",
			limit: 0,
			expected: storage.QueryParams{
				Offset: 0,
				Limit:  20,
			},
		},
		{
			name:  "max limit",
			limit: 100,
			expected: storage.QueryParams{
				Offset: 0,
				Limit:  50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				store = &storeMock{}
				s     = &Server{
					logger: logger,
					store:  store,
				}
			)

			store.On("CountFlags", mock.Anything).Return(uint(3), nil)
			store.On("ListFlags", mock.Anything, mock.MatchedBy(func(opts []storage.QueryOption) bool {
				params := storage.QueryParams{}
				for _, o := range opts {
					o(&params)
				}

				return assert.Equal(t, tt.expected, params)
			})).Return(
				[]*flipt.Flag{
					{
						Key: "foo",
					},
					{
						Key: "bar",
					},
					{
						Key: "baz",
					},
				}, nil)

			got, err := s.ListFlags(context.TODO(), &flipt.ListFlagRequest{
				Offset: tt.offset,
				Limit:  tt.limit,
			})

			require.NoError(t, err)
			assert.NotEmpty(t, got.Flags)
		})
	}
}

func TestCreateFlag(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
		req = &flipt.CreateFlagRequest{
			Key:         "key",
			Name:        "name",
			Description: "desc",
			Enabled:     true,
		}
	)

	store.On("CreateFlag", mock.Anything, req).Return(&flipt.Flag{
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
	}, nil)

	got, err := s.CreateFlag(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)
}

func TestUpdateFlag(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
		req = &flipt.UpdateFlagRequest{
			Key:         "key",
			Name:        "name",
			Description: "desc",
			Enabled:     true,
		}
	)

	store.On("UpdateFlag", mock.Anything, req).Return(&flipt.Flag{
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
	}, nil)

	got, err := s.UpdateFlag(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)
}

func TestDeleteFlag(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
		req = &flipt.DeleteFlagRequest{
			Key: "key",
		}
	)

	store.On("DeleteFlag", mock.Anything, req).Return(nil)

	got, err := s.DeleteFlag(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)
}

func TestCreateVariant(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
		req = &flipt.CreateVariantRequest{
			FlagKey:     "flagKey",
			Key:         "key",
			Name:        "name",
			Description: "desc",
		}
	)

	store.On("CreateVariant", mock.Anything, req).Return(&flipt.Variant{
		Id:          "1",
		FlagKey:     req.FlagKey,
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Attachment:  req.Attachment,
	}, nil)

	got, err := s.CreateVariant(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)
}

func TestUpdateVariant(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
		req = &flipt.UpdateVariantRequest{
			Id:          "1",
			FlagKey:     "flagKey",
			Key:         "key",
			Name:        "name",
			Description: "desc",
		}
	)

	store.On("UpdateVariant", mock.Anything, req).Return(&flipt.Variant{
		Id:          req.Id,
		FlagKey:     req.FlagKey,
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Attachment:  req.Attachment,
	}, nil)

	got, err := s.UpdateVariant(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)
}

func TestDeleteVariant(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
		req = &flipt.DeleteVariantRequest{
			Id: "1",
		}
	)

	store.On("DeleteVariant", mock.Anything, req).Return(nil)

	got, err := s.DeleteVariant(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)
}
