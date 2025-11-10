package mockgen

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpsertUser(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	// Given
	tests := []struct {
		name                 string
		user                 User
		specifyFunctionCalls func(mock *MockIUserRepo)
		expectedError        error
	}{
		{
			user: User{ID: 1, Name: "User 1"},
			name: "Should insert",
			specifyFunctionCalls: func(mockRepo *MockIUserRepo) {
				mockRepo.EXPECT().GetUserByID(1).Return(nil, nil).Times(1)
				mockRepo.EXPECT().Insert(User{ID: 1, Name: "User 1"}).Return(nil).Times(1)
			},
		},
		{
			name: "User existed - Should update",
			user: User{ID: 1, Name: "New User Name"},
			specifyFunctionCalls: func(mockRepo *MockIUserRepo) {
				mockRepo.EXPECT().GetUserByID(1).Return(&User{ID: 1, Name: "User 1"}, nil).Times(1)
				mockRepo.EXPECT().Update(1, User{ID: 1, Name: "New User Name"}).Return(nil).Times(1)
			},
		},
		{
			expectedError: invalidUserIDError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := NewMockIUserRepo(mockCtl)
			if test.specifyFunctionCalls != nil {
				test.specifyFunctionCalls(mockRepo)
			}
			// When
			userService := UserService{repo: mockRepo}
			err := userService.Upsert(test.user)
			// Then
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	// Given
	mockRepo := NewMockIUserRepo(mockCtl)
	mockRepo.EXPECT().Delete(1).Do(func(id int) {
		t.Logf("Deleting user with id: %d", id)
	}).Return(errors.New("Fail to delete")).Times(1)
	// When
	userService := UserService{repo: mockRepo}
	err := userService.DeleteUserByID(1)
	// Then
	assert.EqualError(t, err, "Fail to delete")

}
