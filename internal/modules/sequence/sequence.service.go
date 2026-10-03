package sequence

import (
	"context"
	"fmt"
	"log"

	"github.com/kapilkumar9395/pypr/internal/common"
	"github.com/kapilkumar9395/pypr/internal/modules/variable"
)

type SequenceService struct {
	sequenceRepo *SequenceRepo
	variableRepo *variable.VariableRepo
}

func NewSequenceService(sequenceRepo *SequenceRepo, variableRepo *variable.VariableRepo) *SequenceService {
	return &SequenceService{
		sequenceRepo,
		variableRepo,
	}
}

func (s *SequenceService) CreateSequence(ctx context.Context, organizationId, userId string, payload CreateSequenceRequest) error {
	// Create Final Payload

	sequenceParams := InsertSequenceParams{
		Name:           payload.Name,
		Description:    payload.Description,
		OrganizationId: organizationId,
		Owner:          userId,
		TimezoneMode:   common.ThisOrDefault(payload.TimezoneMode, SequenceTimezoneModeOwner),
		Timezone:       payload.Timezone,
	}

	// Check if sequence steps exists
	sequenceStepParams := make([]InsertSequenceStepParams, 0, len(*payload.SequenceSteps))
	if payload.SequenceSteps != nil && len(*payload.SequenceSteps) != 0 {
		for _, step := range *payload.SequenceSteps {
			variables, err := s.extractVariables(ctx, step.Subject, step.Body, step.Type)
			if err != nil {
				fmt.Printf("%s", err)
				return err
			}
			sequenceStepParams = append(sequenceStepParams, InsertSequenceStepParams{
				Name:          step.Name,
				StepOrder:     step.StepOrder,
				Subject:       step.Subject,
				Body:          step.Body,
				Type:          step.Type,
				ScheduleType:  step.ScheduleType,
				Variables:     variables,
				WaitSeconds:   step.WaitSeconds,
				ScheduledDay:  step.ScheduledDay,
				ScheduledTime: step.ScheduledTime,
				ScheduledAt:   step.ScheduledAt,
			})
		}
	}

	sequenceParams.SequenceSteps = sequenceStepParams

	err := s.sequenceRepo.InsertSequence(ctx, sequenceParams)
	if err != nil {
		log.Printf("Error while creating sequence. Error: %s", err)
		return err
	}

	return nil

}

func (s *SequenceService) extractVariables(ctx context.Context, subject, body *string, stepType StepType) ([]InsertSequenceStepVariableParams, error) {
	if stepType == StepTypeEmail {
		if subject == nil || body == nil {
			return nil, fmt.Errorf("subject and body are required for email step")
		}
	} else {
		return nil, nil
	}
	codes := common.ExtractVariableCodes(*subject, *body)
	variableIds, err := s.variableRepo.GetIDsByCodes(ctx, codes)
	if err != nil {
		log.Printf("Error while fetching variables, e: %s", err)
		return nil, err
	}
	variables := []InsertSequenceStepVariableParams{}
	for _, code := range codes {
		variables = append(variables, InsertSequenceStepVariableParams{VariableId: variableIds[code]})
	}

	return variables, nil
}
