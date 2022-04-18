package services

import (
	"context"
	"errors"
	"github.com/Rippling/gocode-template/mocks/repositories"
	"github.com/Rippling/gocode-template/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockCompanySvcDeps struct {
	*repositories.MockCompany
}

func getTestDeps(t *testing.T) (CompanySvcDeps, MockCompanySvcDeps, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockCompany := repositories.NewMockCompany(ctrl)
	return CompanySvcDeps{Company: mockCompany}, MockCompanySvcDeps{MockCompany: mockCompany}, ctrl
}

func Test_Company_Search(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name             string
		ctx              context.Context
		searchName       string
		companies        []*models.Company
		returnCompanyIdx int
		repoErr          error
		hasErr           bool
	}{
		{
			name:             "Empty company repo result",
			searchName:       "c1",
			companies:        nil,
			repoErr:          nil,
			returnCompanyIdx: -1,
			hasErr:           false},
		{
			name:       "Return list of companies",
			searchName: "c1",
			companies: []*models.Company{
				{
					ID:   "1",
					Name: "c1",
				},
				{
					ID:   "2",
					Name: "c2",
				},
			},
			repoErr: nil,
			hasErr:  false,
		},
		{
			name:             "Return error",
			searchName:       "c1",
			companies:        nil,
			returnCompanyIdx: -1,
			repoErr:          errors.New("failure"),
			hasErr:           true,
		},
	}
	for _, tc := range tcs {
		tc := tc // rebind tc into this lexical scope
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			deps, mockDeps, ctrl := getTestDeps(t)
			defer ctrl.Finish()
			c := New(deps)

			var returnCompany *models.Company
			if tc.returnCompanyIdx >= 0 {
				returnCompany = tc.companies[tc.returnCompanyIdx]
			}

			mockDeps.EXPECT().Search(tc.ctx, tc.searchName).Return(returnCompany, tc.repoErr)
			company, err := c.Search(tc.ctx, tc.searchName)
			if err != tc.repoErr {
				t.Errorf("expected error value same. got %d, exp %d", err, tc.repoErr)
			}
			if !assert.EqualValues(t, returnCompany, company) {
				t.Errorf("expected company value same. got %d, exp %d", err, tc.repoErr)
			}
		})
	}
}
