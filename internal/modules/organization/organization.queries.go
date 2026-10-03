package organization

var (
	INVITE_USER = `
		INSERT INTO organization_invites (
			organization_id,
			email,
			role,
			token_encrypted,
			token_hash,
			expires_at,
			created_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (organization_id, email)
		WHERE accepted_at IS NULL
		DO UPDATE SET
			role = EXCLUDED.role,
			token_encrypted = EXCLUDED.token_encrypted,
			token_hash = EXCLUDED.token_hash,
			expires_at = EXCLUDED.expires_at,
			created_by = EXCLUDED.created_by,
			created_at = NOW()
	`
)
