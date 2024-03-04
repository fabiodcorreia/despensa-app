package services

import (
	"errors"
	"testing"

	"github.com/fabiodcorreia/despensa-app/internal/models"
	"github.com/fabiodcorreia/despensa-app/internal/storage"
	"github.com/fabiodcorreia/despensa-app/test/mocks"
	"go.uber.org/mock/gomock"
)

func TestAddLocation(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	store := mocks.NewMockLocationStore(c)
	service := NewLocationService(store)

	tests := []struct {
		name            string
		serviceInput    models.Location
		serviceErr      error
		serviceErrCheck bool
		mockErr         error
		mockDepInput    string
		mockDepOutput   models.Location
		mockDepErr      error
	}{
		{
			name:            "Adds location successfully",
			serviceInput:    models.NewLocationWithID("1", "Congelador - Gaveta 1"),
			serviceErr:      nil,
			serviceErrCheck: false,
			mockErr:         nil,
			mockDepOutput:   models.Location{},
			mockDepErr:      storage.ErrNotFound,
		},
		{
			name:            "Returns error if location already exists",
			serviceInput:    models.NewLocationWithID("2", "Congelador - Gaveta 2"),
			serviceErr:      ErrLocationExists,
			serviceErrCheck: true,
			mockErr:         nil,
			mockDepOutput:   models.NewLocationWithID("2", "Congelador - Gaveta 2"),
			mockDepErr:      nil,
		},
		{
			name:            "Returns error if location id is too big",
			serviceInput:    models.NewLocationWithID("333333333333", "Congelador - Gaveta 3"),
			serviceErr:      models.ErrLocationIDTooLong,
			serviceErrCheck: true,
			mockErr:         nil,
			mockDepOutput:   models.NewLocationWithID("333333333333", "Congelador - Gaveta 3"),
			mockDepErr:      nil,
		},
		{
			name:            "Returns error if location name is too big",
			serviceInput:    models.NewLocationWithID("4", "Congelador - Gaveta 44444"),
			serviceErr:      models.ErrLocationNameTooLong,
			serviceErrCheck: true,
			mockErr:         nil,
			mockDepOutput:   models.NewLocationWithID("4", "Congelador - Gaveta 44444"),
			mockDepErr:      nil,
		},
		{
			name:            "Returns error if location id is empty",
			serviceInput:    models.NewLocationWithID("", "Congelador - Gaveta 5"),
			serviceErr:      models.ErrLocationIDEmpty,
			serviceErrCheck: true,
			mockErr:         nil,
			mockDepOutput:   models.NewLocationWithID("", "Congelador - Gaveta 5"),
			mockDepErr:      nil,
		},
		{
			name:            "Returns error if location name is empty",
			serviceInput:    models.NewLocationWithID("6", ""),
			serviceErr:      models.ErrLocationNameEmpty,
			serviceErrCheck: true,
			mockErr:         nil,
			mockDepOutput:   models.NewLocationWithID("6", ""),
			mockDepErr:      nil,
		},
	}

	for _, tt := range tests {
		store.EXPECT().GetLocationByID(tt.serviceInput.ID).Return(tt.mockDepOutput, tt.mockDepErr).AnyTimes()
		store.EXPECT().AddLocation(tt.serviceInput).Return(tt.mockErr).AnyTimes()
		err := service.AddLocation(tt.serviceInput)

		if tt.serviceErr != nil {
			if err == nil {
				t.Errorf("%s: expected error %v, got nil", tt.name, tt.serviceErr)
			} else if tt.serviceErrCheck && err.Error() != tt.serviceErr.Error() {
				t.Errorf("%s: expected error %v, got %v", tt.name, tt.serviceErr, err)
			}
		} else {
			if err != nil {
				t.Errorf("%s: unexpected error: %v", tt.name, err)
			}
		}
	}
}

func TestGetLocation(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	store := mocks.NewMockLocationStore(c)
	service := NewLocationService(store)

	tests := []struct {
		name        string
		input       string
		output      interface{}
		expectedErr error
		specificErr bool
		storeErr    error
	}{
		{
			name:        "Gets location successfully",
			input:       "H0qD56o8DC",
			output:      models.NewLocationWithID("H0qD56o8DC", "Congelador - Gaveta 1"),
			expectedErr: nil,
			specificErr: false,
		},
		{
			name:        "Returns ErrLocationNotFound when id not found in store",
			input:       "not-found-id",
			output:      models.Location{},
			expectedErr: ErrLocationNotFound,
			specificErr: true,
			storeErr:    storage.ErrNotFound,
		},
		{
			name:        "Returns ErrLocationNotFound when id is empty",
			input:       "",
			output:      models.Location{},
			expectedErr: ErrLocationNotFound,
			specificErr: true,
			storeErr:    storage.ErrNotFound,
		},
		{
			name:        "Returns other errors from store",
			input:       "error-id",
			output:      models.Location{},
			expectedErr: errors.New("other store error"),
			specificErr: false,
			storeErr:    errors.New("other store error"),
		},
	}

	for _, tt := range tests {
		store.EXPECT().GetLocationByID(tt.input).Return(tt.output, tt.storeErr).AnyTimes()
		output, err := service.GetLocation(tt.input)

		if tt.expectedErr != nil {
			if err == nil {
				t.Errorf("%s: expected error %v, got nil", tt.name, tt.expectedErr)
			} else if tt.specificErr && err.Error() != tt.expectedErr.Error() {
				t.Errorf("%s: expected error %v, got %v", tt.name, tt.expectedErr, err)
			}
		} else {
			if err != nil {
				t.Errorf("%s: unexpected error: %v", tt.name, err)
			}
			if output != tt.output {
				t.Errorf("%s: expected location %v, got %v", tt.name, tt.output, output)
			}
		}
	}
}

/*
Generate a list of test cases for the following function in table format, show include the most possible scenarios, including security cases and fuzzy data.

func (s LocationService) AddLocation(loc models.Location) error {
}

*/
