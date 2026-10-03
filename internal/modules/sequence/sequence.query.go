package sequence

var (
	INSERT_SEQUENCE = `
		INSERT INTO sequences (
			organization_id,
			owner,
			name,
			description,
			timezone_mode,
			timezone
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
)
