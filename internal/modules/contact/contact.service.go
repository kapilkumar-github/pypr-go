package contact

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kapilkumar9395/pypr/internal/common"
)

type ContactService struct {
	repo *ContactRepository
}

func NewContactService(repo *ContactRepository) *ContactService {
	return &ContactService{
		repo,
	}
}

func (s *ContactService) CreateContact(ctx context.Context, payload CreateContactRequest, organizationId, userId string) *common.GenericResponse {
	if payload.EmailId == "" {
		return &common.GenericResponse{
			Error:   fmt.Errorf("Invalid emailId"),
			Message: "EmailId is mandatory",
			Data:    nil,
		}
	}
	params := &CreateContactParams{
		FirstName:      payload.FirstName,
		LastName:       payload.LastName,
		EmailId:        payload.EmailId,
		OrganizationId: organizationId,
		Owner:          userId,
	}

	err := s.repo.Create(ctx, *params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &common.GenericResponse{
				Error:   err,
				Message: "A contact with same EmailId already exists.",
				Data:    nil,
			}
		}
		return &common.GenericResponse{
			Error: err,
			Data:  nil,
		}
	}

	return &common.GenericResponse{
		Error:   nil,
		Message: "Contact created successfully",
		Data:    nil,
	}
}

func (s *ContactService) GetOrgAllContacts(ctx context.Context, orgId string) *common.GenericResponse {
	if orgId == "" {
		return &common.GenericResponse{
			Error:   fmt.Errorf("Invalid OrgId"),
			Message: "OrganizationId is mandatory",
			Data:    nil,
		}
	}

	data, err := s.repo.FetchOrgContactList(ctx, orgId)
	if err != nil {
		return &common.GenericResponse{
			Error: err,
			Data:  nil,
		}
	}

	return &common.GenericResponse{
		Error:   nil,
		Message: "Request completed successfully",
		Data:    data,
	}
}
