package auth

import (
	"context"
	"fmt"

	CommonService "github.com/kapilkumar9395/pypr/internal/common/service"
)

type AuthService struct {
	repo       *AuthRepository
	jwtService *JwtService
}

func NewAuthService(repo *AuthRepository, jwtService *JwtService) *AuthService {
	return &AuthService{repo, jwtService}
}

// UserLogin
func (s *AuthService) UserLogin(emailId, password string) (*LoginResponse, error) {
	ctx := context.Background()
	// Check if user exists
	userExists, err := s.repo.UserExists(ctx, emailId)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, fmt.Errorf("user with email %s does not exist", emailId)
	}

	// Validate password
	userProjection, err := s.repo.GetUserNCredByEmail(ctx, emailId)
	if err != nil {
		return nil, err
	}
	if CommonService.CompareHash(userProjection.PasswordHash, password) != nil {
		return nil, fmt.Errorf("invalid password for user with email %s", emailId)
	}

	// Return access token
	token, err := s.jwtService.GenerateToken(userProjection.OrganizationId, userProjection.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token, %s", emailId)
	}

	loginResponse := &LoginResponse{
		AccessToken: token,
	}

	return loginResponse, nil
}

func (s *AuthService) RegisterUser(emailId, firstName string, lastName string, password string, organizationId *string, invitationToken *string, timezone *string) error {
	if invitationToken != nil && *invitationToken != "" {
		// Register user with invitation token flow
		return s.RegisterUserWithInvitation(emailId, firstName, lastName, password, *invitationToken, timezone)
	}

	// Register user without invitation token flow
	return s.RegisterUserWithoutInvitation(emailId, firstName, password, lastName, organizationId, timezone)
}

func (s *AuthService) RegisterUserWithoutInvitation(emailId, firstName, password string, lastName string, organizationId *string, timezone *string) error {
	// Check if user already exists
	userExists, err := s.repo.UserExists(context.Background(), emailId)
	if err != nil {
		return err
	}
	if userExists {
		return fmt.Errorf("user with email %s already exists", emailId)
	}

	passwordHash, err := CommonService.HashString(password)
	if err != nil {
		return err
	}

	timezoneValue := "Asia/Kolkata" // Default timezone
	if timezone != nil && *timezone != "" {
		timezoneValue = *timezone
	}

	params := CreateUserParams{
		Email:        emailId,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: passwordHash,
		Timezone:     timezoneValue,
	}

	if organizationId != nil && *organizationId != "" {
		// New user with existing organization flow
		err, _ := s.repo.CreateUserFlowWithExistingOrg(context.Background(), params, *organizationId, nil)
		if err != nil {
			return err
		}
	} else {
		// New user with new organization flow
		return s.repo.CreateUserFlowWithNewOrg(params)
	}
	return nil
}

func (s *AuthService) RegisterUserWithInvitation(emailId, firstName string, lastName string, password, invitationToken string, timezone *string) error {
	// Check if user already exists
	userExists, err := s.repo.UserExists(context.Background(), emailId)
	if err != nil {
		return err
	}
	if userExists {
		return fmt.Errorf("user with email %s already exists", emailId)
	}

	passwordHash, err := CommonService.HashString(password)
	if err != nil {
		return err
	}

	timezoneValue := "Asia/Kolkata" // Default timezone
	if timezone != nil && *timezone != "" {
		timezoneValue = *timezone
	}

	params := CreateUserParams{
		Email:        emailId,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: passwordHash,
		Timezone:     timezoneValue,
	}

	// Get organization details using the invitation token
	organizationDetails, err := s.repo.GetOrganizationIdByInvitationToken(context.Background(), invitationToken, emailId)
	if err != nil {
		return fmt.Errorf("invalid or expired invitation token")
	}

	// Proceed with the existing organization flow
	err, _ = s.repo.CreateUserFlowWithExistingOrg(context.Background(), params, organizationDetails.OrganizationId, nil)
	if err != nil {
		return err
	}
	return nil
}
