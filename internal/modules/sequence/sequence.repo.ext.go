package sequence

import (
	"fmt"
	"strings"
	"uuid"

	"github.com/kapilkumar9395/pypr/internal/common"
)

func (r *SequenceRepo) buildInsertStepsQuery(stepCount int) (string, int) {
	columns := []string{
		"sequence_id",
		"name",
		"step_order",
		"type",
		"schedule_type",
		"wait_seconds",
		"scheduled_day",
		"scheduled_time",
		"scheduled_at",
		"subject",
		"body",
	}
	colLen := len(columns)
	return fmt.Sprintf(`
		INSERT INTO sequence_steps (
			%s
		)
		VALUES %s
		RETURNING id, step_order;
	`, strings.Join(columns, ", "), common.GeneratePlaceholders(stepCount, colLen)), colLen
}
func (r *SequenceRepo) buildInsertSequenceStepVariablesQuery(
	steps []InsertSequenceStepParams,
	stepIDs map[int]uuid.UUID,
) (string, []any) {

	count := 0
	for _, step := range steps {
		count += len(step.Variables)
	}

	if count == 0 {
		return "", nil
	}

	columns := []string{
		"sequence_step_id",
		"variable_id",
	}
	colLen := len(columns)
	query := fmt.Sprintf(`
		INSERT INTO sequence_step_variables (
			%s
		)
		VALUES %s
		ON CONFLICT DO NOTHING;
	`, strings.Join(columns, ","), common.GeneratePlaceholders(count, colLen))

	args := make([]any, 0, count*colLen)

	for _, step := range steps {
		stepID := stepIDs[step.StepOrder]

		for _, variable := range step.Variables {
			variableID := variable.VariableId
			args = append(args, stepID, variableID)
		}
	}

	return query, args
}
