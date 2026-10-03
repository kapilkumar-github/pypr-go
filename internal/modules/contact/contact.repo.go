package contact

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactRepository struct {
	db *pgxpool.Pool
}

func NewContactRepository(db *pgxpool.Pool) *ContactRepository {
	return &ContactRepository{
		db,
	}
}

func (r *ContactRepository) Create(ctx context.Context, params CreateContactParams) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO contacts (
			organization_id,
			owner,
			first_name,
			last_name,
			email
		)
		VALUES ($1, $2, $3, $4, $5);
	`,
		params.OrganizationId,
		params.Owner,
		params.FirstName,
		params.LastName,
		params.EmailId,
	)

	return err
}

func (r *ContactRepository) FetchOrgContactList(ctx context.Context, orgId string) ([]OrgContactListProjection, error) {
	var list []OrgContactListProjection

	query := `
		SELECT 
		c.id::text,
		u.id::text, u.first_name, u.last_name,
		c.email, c.first_name, c.last_name,
		c.company, c.job_title
		from contacts as c
		left join users as u on u.id = c.owner
		where c.organization_id = $1
	`

	rows, err := r.db.Query(ctx, query, orgId)
	if err != nil {
		log.Printf("[ContactRepository][FetchOrgContactList] Error while fetching the contact list, Error: %v", err)
		return nil, err
	}
	for rows.Next() {
		var item OrgContactListProjection
		err := rows.Scan(
			&item.Id,
			&item.OwnerId,
			&item.OwnerFirstName,
			&item.OwnerLastName,
			&item.EmailId,
			&item.FirstName,
			&item.LastName,
			&item.Company,
			&item.JobTitle,
		)
		if err != nil {
			log.Printf("Error while scanning to OrgContactListProjection, %v", err)
			return nil, err
		}
		item.OrganizationId = orgId
		list = append(list, item)
	}

	return list, nil
}
