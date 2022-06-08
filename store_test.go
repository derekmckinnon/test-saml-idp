package idp

import (
	"errors"
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlidp"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_GetUser(t *testing.T) {
	original := &samlidp.User{Name: "Test"}
	store := &Store{}

	err := store.AddUser(original)
	require.Nil(t, err)

	result, err := store.GetUser(original.Name)
	require.Nil(t, err)
	require.Equal(t, original.Name, result.Name)
}

func TestStore_GetUserMissing(t *testing.T) {
	store := &Store{}

	user, err := store.GetUser("")
	require.NotNil(t, err)
	require.Nil(t, user)
}

func TestStore_GetUsers(t *testing.T) {
	user1 := &samlidp.User{Name: "Test1"}
	user2 := &samlidp.User{Name: "Test2"}
	store := &Store{}

	err := store.AddUser(user1)
	require.Nil(t, err)
	err = store.AddUser(user2)
	require.Nil(t, err)

	users, err := store.GetUsers()
	require.Nil(t, err)
	require.Len(t, users, 2)
}

func TestStore_GetUsersEmpty(t *testing.T) {
	store := &Store{}

	users, err := store.GetUsers()
	require.Nil(t, err)
	require.Len(t, users, 0)
}

func TestStore_GetServiceProvider(t *testing.T) {
	original := &samlidp.Service{
		Metadata: saml.EntityDescriptor{
			EntityID: "Test",
		},
	}
	store := &Store{}

	err := store.AddServiceProvider(original)
	require.Nil(t, err)

	result, err := store.GetServiceProvider(original.Metadata.EntityID)
	require.Nil(t, err)
	require.Equal(t, original.Metadata.EntityID, result.Metadata.EntityID)
}

func TestStore_GetServiceProviderMissing(t *testing.T) {
	store := &Store{}

	sp, err := store.GetServiceProvider("")
	require.NotNil(t, err)
	require.Nil(t, sp)
}

func TestStore_GetServiceProviders(t *testing.T) {
	sp1 := &samlidp.Service{
		Metadata: saml.EntityDescriptor{
			EntityID: "Test1",
		},
	}
	sp2 := &samlidp.Service{
		Metadata: saml.EntityDescriptor{
			EntityID: "Test2",
		},
	}
	store := &Store{}

	err := store.AddServiceProvider(sp1)
	require.Nil(t, err)
	err = store.AddServiceProvider(sp2)
	require.Nil(t, err)

	services, err := store.GetServiceProviders()
	require.Nil(t, err)
	require.Len(t, services, 2)
}

func TestStore_GetServiceProvidersEmpty(t *testing.T) {
	store := &Store{}

	services, err := store.GetServiceProviders()
	require.Nil(t, err)
	require.Len(t, services, 0)
}

func TestStore_GetSession(t *testing.T) {
	original := &saml.Session{ID: "Test"}
	store := &Store{}

	err := store.AddSession(original)
	require.Nil(t, err)

	result, err := store.GetSession(original.ID)
	require.Nil(t, err)
	require.Equal(t, original.ID, result.ID)
}

func TestStore_GetSessionMissing(t *testing.T) {
	store := &Store{}

	sp, err := store.GetSession("")
	require.NotNil(t, err)
	require.Nil(t, sp)
}

func TestStore_GetSessions(t *testing.T) {
	sp1 := &saml.Session{ID: "Test1"}
	sp2 := &saml.Session{ID: "Test2"}
	store := &Store{}

	err := store.AddSession(sp1)
	require.Nil(t, err)
	err = store.AddSession(sp2)
	require.Nil(t, err)

	services, err := store.GetSessions()
	require.Nil(t, err)
	require.Len(t, services, 2)
}

func TestStore_GetSessionsEmpty(t *testing.T) {
	store := &Store{}

	services, err := store.GetSessions()
	require.Nil(t, err)
	require.Len(t, services, 0)
}

func Test_GetResourcesEmptyList(t *testing.T) {
	store := &Store{}

	users, err := getResources[samlidp.User](store, usersPrefix, func(string) (*samlidp.User, error) {
		return nil, nil
	})
	require.Nil(t, err)
	require.Empty(t, users)
}

func Test_GetResourcesGetterError(t *testing.T) {
	store := &Store{}
	_ = store.AddUser(&samlidp.User{Name: "Test"})

	users, err := getResources[samlidp.User](store, usersPrefix, func(string) (*samlidp.User, error) {
		return nil, errors.New("foobar")
	})
	require.Nil(t, users)
	require.Error(t, err)
}
