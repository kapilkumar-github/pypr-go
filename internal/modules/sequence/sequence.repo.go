package sequence

import (
	"context"
	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SequenceRepo struct {
	db *pgxpool.Pool
}

func NewSequenceRepo(db *pgxpool.Pool) *SequenceRepo {
	return &SequenceRepo{
		db,
	}
}

func (r *SequenceRepo) InsertSequence(ctx context.Context, params InsertSequenceParams) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Create sequence
	var sequenceId uuid.UUID

	err = tx.QueryRow(
		ctx,
		INSERT_SEQUENCE,
		params.OrganizationId,
		params.Owner,
		params.Name,
		params.Description,
		params.TimezoneMode,
		params.Timezone,
	).Scan(&sequenceId)

	if err != nil {
		return err
	}

	// 2. Bulk insert steps
	// Build args from steps
	// RETURNING id, step_order
	steps := len(params.SequenceSteps)
	if steps != 0 {

		sequenceStepsQuery, colLen := r.buildInsertStepsQuery(steps)
		stepsArgs := make([]any, 0, steps*colLen)

		for _, step := range params.SequenceSteps {
			stepsArgs = append(stepsArgs,
				sequenceId,
				step.Name,
				step.StepOrder,
				step.Type,
				step.ScheduleType,
				step.WaitSeconds,
				step.ScheduledDay,
				step.ScheduledTime,
				step.ScheduledAt,
				step.Subject,
				step.Body,
			)
		}

		rows, err := tx.Query(ctx, sequenceStepsQuery, stepsArgs...)
		if err != nil {
			return err
		}
		stepsMap := make(map[int]uuid.UUID)
		for rows.Next() {
			var id uuid.UUID
			var stepOrder int

			if err := rows.Scan(&id, &stepOrder); err != nil {
				rows.Close()
				return err
			}

			stepsMap[stepOrder] = id
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return err
		}

		rows.Close()
		// 3. Bulk insert sequence_step_variables
		// using the returned step IDs
		variablesQuery, args := r.buildInsertSequenceStepVariablesQuery(params.SequenceSteps, stepsMap)
		_, err = tx.Exec(ctx, variablesQuery, args...)
		if err != nil {
			return nil
		}
	}
	// 4. Commit
	return tx.Commit(ctx)
}
