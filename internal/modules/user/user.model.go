package user

type UserModel struct {
	OrganizationId string `json:"organizationId" db:"organization_id"`
	Id             string `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
}
