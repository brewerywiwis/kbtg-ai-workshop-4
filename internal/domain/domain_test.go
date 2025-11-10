package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validation(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		isValid bool
	}{
		{
			name: "valid user",
			user: User{
				Name:  "Test User",
				Email: "test@example.com",
				Phone: "081-234-5678",
			},
			isValid: true,
		},
		{
			name: "empty name",
			user: User{
				Name:  "",
				Email: "test@example.com",
				Phone: "081-234-5678",
			},
			isValid: false,
		},
		{
			name: "empty email",
			user: User{
				Name:  "Test User",
				Email: "",
				Phone: "081-234-5678",
			},
			isValid: false,
		},
		{
			name: "empty phone",
			user: User{
				Name:  "Test User",
				Email: "test@example.com",
				Phone: "",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestTransferStatus_Constants(t *testing.T) {
	assert.Equal(t, TransferStatus("pending"), TransferStatusPending)
	assert.Equal(t, TransferStatus("processing"), TransferStatusProcessing)
	assert.Equal(t, TransferStatus("completed"), TransferStatusCompleted)
	assert.Equal(t, TransferStatus("failed"), TransferStatusFailed)
	assert.Equal(t, TransferStatus("cancelled"), TransferStatusCancelled)
	assert.Equal(t, TransferStatus("reversed"), TransferStatusReversed)
}

func TestTransfer_Validation(t *testing.T) {
	tests := []struct {
		name     string
		transfer Transfer
		isValid  bool
	}{
		{
			name: "valid transfer",
			transfer: Transfer{
				FromUserID: 1,
				ToUserID:   2,
				Amount:     500,
			},
			isValid: true,
		},
		{
			name: "same user transfer",
			transfer: Transfer{
				FromUserID: 1,
				ToUserID:   1,
				Amount:     500,
			},
			isValid: false,
		},
		{
			name: "zero amount",
			transfer: Transfer{
				FromUserID: 1,
				ToUserID:   2,
				Amount:     0,
			},
			isValid: false,
		},
		{
			name: "negative amount",
			transfer: Transfer{
				FromUserID: 1,
				ToUserID:   2,
				Amount:     -100,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.transfer.Validate()
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
