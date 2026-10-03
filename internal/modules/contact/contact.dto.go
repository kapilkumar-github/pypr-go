package contact

type CreateContactParams struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	EmailId        string `json:"emailId"`
	OrganizationId string `json:"organizationId"`
	Owner          string `json:"owner"`
}

type CreateContactRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	EmailId   string `json:"emailId"`
}

type OrgContactListProjection struct {
	Id             string  `json:"id"`
	OrganizationId string  `json:"organizationId"`
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	EmailId        string  `json:"emailId"`
	Company        *string `json:"company"`
	JobTitle       *string `json:"jobTitle"`
	OwnerId        string  `json:"ownerId"`
	OwnerFirstName string  `json:"ownerFirstName"`
	OwnerLastName  string  `json:"ownerLastName"`
}
