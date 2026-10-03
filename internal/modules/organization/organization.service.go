package organization

import (
	"context"
	"log"
	"time"

	CommonService "github.com/kapilkumar9395/pypr/internal/common/service"
	emailInfra "github.com/kapilkumar9395/pypr/internal/infra/email"
)

type OrganizationService struct {
	OrganizationRepo OrganizationRepository
	EmailFactory     emailInfra.Factory
}

func NewOrganizationService(repo OrganizationRepository, emailFactory emailInfra.Factory) *OrganizationService {
	return &OrganizationService{
		OrganizationRepo: repo,
		EmailFactory:     emailFactory,
	}
}

func (s *OrganizationService) InviteUser(ctx context.Context, organizationId, userId, emailId string) (string, error) {
	token, err := CommonService.GenerateRandomToken(32) // Generate a random token of 32 characters
	if err != nil {
		return "", err
	}
	// Not hashing as of now
	// tokenHash, err := CommonService.HashString(token)
	// if err != nil {
	// 	return "", err
	// }
	params := InviteUserParams{
		OrganizationId: organizationId,
		Email:          emailId,
		EncryptedToken: token,
		TokenHash:      token,
		CreatedBy:      userId,
		ExpiresAt:      time.Now().Add(365 * 24 * time.Hour).Format(time.RFC3339), // Set expiration to 240 hours from now
	}

	// Proceed to create a new invitation
	err = s.OrganizationRepo.InviteUser(ctx, params)
	if err != nil {
		log.Printf("[OrganizationService][InviteUser] error while creating invite, %s", err)
		return "", err
	}

	// Send email to the invited user with the invitation link containing the token
	emailVariables := map[string]string{
		"OrganizationName": organizationId,
		"FirstName":        userId,
		"InviterName":      userId,
		"InvitationLink":   token,
	}
	emailBody, err := emailInfra.TemplateUserInvitation.ReplaceTemplateVariables(emailVariables)
	if err != nil {
		return "", err
	}
	emailSender, err := s.EmailFactory.GetSender(emailInfra.ProviderResend)
	if err != nil {
		return "", err
	}

	err = emailSender.Send(emailInfra.EmailMessage{
		From:    emailInfra.PyprHelloEmailId,
		To:      emailId,
		Subject: "You're invited to join " + organizationId,
		HTML:    emailBody,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
