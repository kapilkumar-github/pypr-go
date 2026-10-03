package sequence

type StepType string

const (
	StepTypeEmail StepType = "EMAIL"
	StepTypeWait  StepType = "WAIT"
)

type StepScheduleType string

const (
	StepScheduleTypeImmediate   StepScheduleType = "IMMEDIATE"
	StepScheduleTypeWeekdayTime StepScheduleType = "WEEKDAY_TIME"
	StepScheduleTypeExactDate   StepScheduleType = "EXACT_DATE"
)

type SequenceStatus string

const (
	SequenceStatusDraft    SequenceStatus = "DRAFT"
	SequenceStatusActive   SequenceStatus = "ACTIVE"
	SequenceStatusPaused   SequenceStatus = "PAUSED"
	SequenceStatusArchived SequenceStatus = "ARCHIVED"
)

type SequenceTimezoneMode string

const (
	SequenceTimezoneModeOwner   SequenceTimezoneMode = "OWNER"
	SequenceTimezoneModeContact SequenceTimezoneMode = "CONTACT"
	SequenceTimezoneModeCustom  SequenceTimezoneMode = "CUSTOM"
)
