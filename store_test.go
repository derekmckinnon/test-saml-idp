package main

import (
	"github.com/crewjam/saml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryStore_GetUser_ReturnsNilForMissingKey(t *testing.T) {
	sut := MemoryStore{}

	result, ok := sut.GetUser("foobar")

	assert.Nil(t, result)
	assert.False(t, ok)
}

func TestMemoryStore_PutUser_SetsCorrectData(t *testing.T) {
	sut := MemoryStore{}

	user := User{
		Id:       "foobar",
		Email:    "foo@bar.com",
		Name:     "Foo",
		Password: "foobar2",
	}

	sut.PutUser(&user)

	result, ok := sut.GetUser(user.Id)

	if assert.NotNil(t, result) {
		assert.Equal(t, user.Id, result.Id)
		assert.Equal(t, user.Email, result.Email)
		assert.Equal(t, user.Name, result.Name)
	}

	assert.True(t, ok)
}

func TestMemoryStore_GetServiceProvider_ReturnsNilForMissingKey(t *testing.T) {
	sut := MemoryStore{}

	result, ok := sut.GetServiceProvider("foobar")

	assert.Nil(t, result)
	assert.False(t, ok)
}

func TestMemoryStore_PutServiceProvider(t *testing.T) {
	sut := MemoryStore{}

	sp := ServiceProvider{
		Id: "foobar",
		EntityDescriptor: saml.EntityDescriptor{
			EntityID: "foobar",
		},
	}

	sut.PutServiceProvider(&sp)

	result, ok := sut.GetServiceProvider(sp.Id)

	if assert.NotNil(t, result) {
		assert.Equal(t, sp.Id, result.Id)
		assert.Equal(t, sp.EntityDescriptor.EntityID, result.EntityDescriptor.EntityID)
	}

	assert.True(t, ok)
}
