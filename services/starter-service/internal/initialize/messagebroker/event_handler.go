package messagebroker

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/pkg/httputil"
	starterquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/query"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	domainRepo "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	domainService "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/service"
)

type EventHandler struct {
	starterRepo       domainRepo.StarterRepository
	starterSearchRepo domainRepo.StarterSearchRepository
	enrichmentService *domainService.StarterEnrichmentService
}

func NewEventHandler(
	starterRepo domainRepo.StarterRepository,
	starterSearchRepo domainRepo.StarterSearchRepository,
	enrichmentService *domainService.StarterEnrichmentService,
) *EventHandler {
	return &EventHandler{
		starterRepo:       starterRepo,
		starterSearchRepo: starterSearchRepo,
		enrichmentService: enrichmentService,
	}
}

func (h *EventHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *EventHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *EventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		ctx := session.Context()

		// Parse event
		event, err := events.BytesToEvent(msg.Value)
		if err != nil {
			log.Printf("Failed to parse event: %v", err)
			session.MarkMessage(msg, "")
			continue
		}

		log.Printf("Processing event: Type=%s, ID=%s", event.Type, event.ID)

		var payload events.IndexStarterPayload
		if err := event.UnmarshalPayload(&payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		// TODO: switch to handle base on event type
		batchSize := 100
		page := 1

		emptyQuery := starterquery.ListStartersQuery{
			Pagination: httputil.ReqPagination{
				Page: &page, Limit: &batchSize,
			},
			Keyword:  payload.Domain,
			SearchBy: "domain",
		}

		starters, _, err := h.starterRepo.SearchByKeyword(ctx, &emptyQuery)

		if err != nil {
			return err
		}

		enriched, err := h.enrichmentService.EnrichStarters(ctx, starters)

		esDocs := make([]*model.StarterESDoc, len(starters))
		for i := range starters {
			esDoc := model.NewStarterESDocFromStarter(starters[i], enriched)
			esDocs[i] = esDoc
		}
		if err != nil {
			log.Printf("Failed to create starter: %v from message", err)
			return err
		}

		switch event.Type {
		case events.EventTypeStarterInsert, events.EventTypeStarterUpdate:
			err := h.starterSearchRepo.IndexStarter(ctx, esDocs[0])
			if err != nil {
				return err
			}
		case events.EventTypeStarterDelete:
			err := h.starterSearchRepo.DeleteFromIndex(ctx, starters[0].Domain())
			if err != nil {
				return err
			}
		default:
			log.Printf("Unknown event type: %s", event.Type)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
