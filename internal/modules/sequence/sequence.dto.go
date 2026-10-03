package sequence

import (
	"time"
	"uuid"
)

type InsertSequenceParams struct {
	Name           string
	Description    string
	OrganizationId string
	Owner          string
	SequenceSteps  []InsertSequenceStepParams
	TimezoneMode   SequenceTimezoneMode
	Timezone       *string
}

type InsertSequenceStepParams struct {
	Name          string
	StepOrder     int
	Type          StepType
	ScheduleType  *StepScheduleType
	Subject       *string
	Body          *string
	Variables     []InsertSequenceStepVariableParams
	WaitSeconds   *int64
	ScheduledDay  *int16
	ScheduledTime *string
	ScheduledAt   *time.Time
}

type InsertSequenceStepVariableParams struct {
	VariableId uuid.UUID
}

type CreateSequenceRequest struct {
	Name          string
	Description   string
	SequenceSteps *[]CreateSequenceStepRequest
	TimezoneMode  *SequenceTimezoneMode
	Timezone      *string
}

type CreateSequenceStepRequest struct {
	Name          string
	StepOrder     int
	Subject       *string
	Body          *string
	Type          StepType
	ScheduleType  *StepScheduleType
	WaitSeconds   *int64
	ScheduledDay  *int16
	ScheduledTime *string
	ScheduledAt   *time.Time
}
