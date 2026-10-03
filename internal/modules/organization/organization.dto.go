package organization

type InviteUserParams struct {
	OrganizationId string `json:"organizationId"`
	Email          string `json:"email"`
	EncryptedToken string `json:"encryptedToken"`
	TokenHash      string `json:"tokenHash"`
	CreatedBy      string `json:"createdBy"`
	ExpiresAt      string `json:"expiresAt"`
}

type InviteUserRequest struct {
	EmailId string `json:"emailId" binding:"required,email"`
	Role    string `json:"role" binding:"required"`
}
