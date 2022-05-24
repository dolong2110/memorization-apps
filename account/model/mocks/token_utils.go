package mocks

import (
	"crypto/rsa"
	"github.com/dolong2110/Memoirization-Apps/account/model"
	"github.com/stretchr/testify/mock"
)

// MockTokenUtils is a mock type for utils.token
type MockTokenUtils struct {
	mock.Mock
}

// ValidateIDToken test function in utils folder
func (m *MockTokenUtils) ValidateIDToken(tokenString string, key *rsa.PublicKey) (*model.IDTokenCustomClaims, error) {
	ret := m.Called(tokenString, key)

	var r0 *model.IDTokenCustomClaims
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.IDTokenCustomClaims)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// ValidateRefreshToken test function in utils folder
func (m *MockTokenUtils) ValidateRefreshToken(tokenString string, key string) (*model.User, error) {
	ret := m.Called(tokenString, key)

	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
