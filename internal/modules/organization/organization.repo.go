package organization

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationRepository struct {
	db *pgxpool.Pool
}

func NewOrganizationRepository(db *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) InviteUser(ctx context.Context, params InviteUserParams) error {
	_, err := r.db.Exec(
		ctx,
		INVITE_USER,
		params.OrganizationId,
		params.Email,
		"MEMBER",
		params.EncryptedToken,
		params.TokenHash,
		params.ExpiresAt,
		params.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("[InviteUser] error while creating invite, %s", err)
	}
	return nil
}
