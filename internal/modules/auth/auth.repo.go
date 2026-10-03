package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kapilkumar9395/pypr/internal/modules/user"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUserFlowWithExistingOrg(
	ctx context.Context,
	params CreateUserParams,
	organizationID string,
	tx pgx.Tx,
) (error, *user.UserModel) {
	var err error
	ownTx := false

	if tx == nil {
		tx, err = r.db.Begin(ctx)
		if err != nil {
			return fmt.Errorf("[CreateUserFlowWithExistingOrg] begin transaction: %w", err), nil
		}

		ownTx = true
		defer tx.Rollback(ctx)
	}

	var user user.UserModel

	err = tx.QueryRow(
		ctx,
		`
        WITH new_user AS (
            INSERT INTO users (
                email,
                first_name,
                last_name,
                timezone
            )
            VALUES ($1, $2, $3, $4)
            RETURNING id
        ),
        new_credentials AS (
            INSERT INTO user_credentials (
                user_id,
                password_hash
            )
            SELECT id, $5
            FROM new_user
        )
        INSERT INTO organization_members (
            organization_id,
            user_id,
            role
        )
        SELECT $6, id, 'MEMBER'
        FROM new_user
        RETURNING user_id
        `,
		params.Email,
		params.FirstName,
		params.LastName,
		params.Timezone,
		params.PasswordHash,
		organizationID,
	).Scan(&user.Id)

	if err != nil {
		return fmt.Errorf(
			"[CreateUserFlowWithExistingOrg] create user flow: %w",
			err,
		), nil
	}

	if ownTx {
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf(
				"[CreateUserFlowWithExistingOrg] commit transaction: %w",
				err,
			), nil
		}
	}

	return nil, &user
}
func (r *AuthRepository) CreateUserFlowWithNewOrg(params CreateUserParams) error {
	// Create organization then call the existing flow to create user and add to organization
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("[CreateUserFlowWithNewOrg] begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var organizationId string

	start := time.Now()
	err = tx.QueryRow(
		ctx,
		`
		INSERT INTO organizations (
			name,
			type,
			timezone
		)
		VALUES ($1, 'PERSONAL', $2)
		RETURNING id
		`,
		params.FirstName+"'s Organization",
		params.Timezone,
	).Scan(&organizationId)

	if err != nil {
		return fmt.Errorf("[CreateUserFlowWithNewOrg] create organization: %w", err)
	}

	log.Printf("[CreateUserFlowWithNewOrg] Organization created in %v", time.Since(start))

	err, _ = r.CreateUserFlowWithExistingOrg(ctx, params, organizationId, tx)
	if err != nil {
		return fmt.Errorf("[CreateUserFlowWithNewOrg] create user with new organization: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("[CreateUserFlowWithNewOrg] commit transaction: %w", err)
	}

	return nil
}

func (r *AuthRepository) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool

	err := r.db.QueryRow(
		ctx,
		"SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)",
		email,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *AuthRepository) GetOrganizationIdByInvitationToken(ctx context.Context, invitationToken, email string) (OrganizationByTokenProjection, error) {
	var organizationDetails OrganizationByTokenProjection

	log.Printf("InvitationToken: %s, Email: %s", invitationToken, email)
	err := r.db.QueryRow(
		ctx,
		"SELECT organizations.id, organizations.name FROM organizations JOIN organization_invites ON organizations.id = organization_invites.organization_id WHERE organization_invites.token_hash = $1 AND organization_invites.email = $2",
		invitationToken,
		email,
	).Scan(&organizationDetails.OrganizationId, &organizationDetails.Name)

	if err != nil {
		return OrganizationByTokenProjection{}, err
	}

	return organizationDetails, nil
}

func (r *AuthRepository) GetUserNCredByEmail(ctx context.Context, email string) (*UserLoginProjection, error) {
	var userModel UserLoginProjection

	err := r.db.QueryRow(
		ctx,
		`SELECT u.id, u.email, u.first_name, u.last_name, u.timezone, uc.password_hash, o.organization_id FROM users as u
		JOIN user_credentials as uc ON u.id = uc.user_id
		JOIN organization_members as o ON u.id = o.user_id
		WHERE email = $1`,
		email,
	).Scan(&userModel.UserId, &userModel.Email, &userModel.FirstName, &userModel.LastName, &userModel.Timezone, &userModel.PasswordHash, &userModel.OrganizationId)

	if err != nil {
		return nil, err
	}

	return &userModel, nil
}
