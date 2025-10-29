package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/application/service"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	domainService "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/service"
	starterdto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/starter"
)

type StarterHandler struct {
	starterSvc        *service.StarterApplicationService
	enrichmentService *domainService.StarterEnrichmentService
}

func NewStarterHandler(
	starterSvc *service.StarterApplicationService,
	enrichmentService *domainService.StarterEnrichmentService,
) *StarterHandler {
	return &StarterHandler{
		starterSvc:        starterSvc,
		enrichmentService: enrichmentService,
	}
}

// ListStarters GET /api/v1/starters
func (sh *StarterHandler) ListStarters(ctx *gin.Context) {
	httputil.Wrap(sh.listStarters)(ctx)
}

func (sh *StarterHandler) listStarters(ctx *gin.Context) (res interface{}, err error) {
	var req starterdto.ListStartersRequest
	if err := httputil.ValidateQuery(ctx, &req); err != nil {
		return nil, err
	}
	req.SetDefaults()

	rawResult, err := sh.starterSvc.ListStarters(ctx, req.ToQuery())
	if err != nil {
		return nil, err
	}

	enrichedDomain, err := sh.enrichmentService.EnrichStarters(ctx, rawResult.Data)
	if err != nil {
		return nil, err
	}

	enrichedDTO := starterdto.FromDomainEnrichment(enrichedDomain)
	responseData := starterdto.FromStartersEnriched(rawResult.Data, enrichedDTO)

	return &httputil.PaginatedResult[*starterdto.StarterResponse]{
		Data:       responseData,
		Pagination: httputil.CursorPagination(ctx, rawResult.Pagination),
	}, nil
}

// CreateStarter POST /api/v1/starters
func (sh *StarterHandler) CreateStarter(ctx *gin.Context) {
	httputil.Wrap(sh.createStarter)(ctx)
}

func (sh *StarterHandler) createStarter(ctx *gin.Context) (res interface{}, err error) {
	var req starterdto.CreateStarterRequest
	if err := httputil.ValidateBody(ctx, &req); err != nil {
		return nil, err
	}

	starter, err := sh.starterSvc.CreateStarter(ctx, req.ToCommand())
	if err != nil {
		return nil, err
	}

	enrichedDomain, err := sh.enrichmentService.EnrichStarters(ctx, []*model.Starter{starter})
	if err != nil {
		return nil, err
	}

	enrichedDTO := starterdto.FromDomainEnrichment(enrichedDomain)
	return starterdto.FromDomainEnriched(starter, enrichedDTO), nil
}

// Find GET /api/v1/starters/{domain}
func (sh *StarterHandler) Find(ctx *gin.Context) {
	httputil.Wrap(sh.find)(ctx)
}

func (sh *StarterHandler) find(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		Domain string `uri:"domain" binding:"required"`
	}
	if err := httputil.ValidateURI(ctx, &uriReq); err != nil {
		return nil, err
	}

	starter, err := sh.starterSvc.GetStarterByDomain(ctx, uriReq.Domain)
	if err != nil {
		return nil, err
	}
	enrichedDomain, err := sh.enrichmentService.EnrichStarters(ctx, []*model.Starter{starter})
	if err != nil {
		return nil, err
	}

	enrichedDTO := starterdto.FromDomainEnrichment(enrichedDomain)
	return starterdto.FromDomainEnriched(starter, enrichedDTO), nil
}

// UpdateStarter PATCH /api/v1/starters/{domain}
func (sh *StarterHandler) UpdateStarter(ctx *gin.Context) {
	httputil.Wrap(sh.updateStarter)(ctx)
}

func (sh *StarterHandler) updateStarter(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		Domain string `uri:"domain" binding:"required"`
	}
	if err := httputil.ValidateURI(ctx, &uriReq); err != nil {
		return nil, err
	}
	var req starterdto.UpdateStarterRequest
	if err := httputil.ValidateBody(ctx, &req); err != nil {
		return nil, err
	}

	starter, err := sh.starterSvc.UpdateStarter(ctx, req.ToCommand(uriReq.Domain))
	if err != nil {
		return nil, err
	}
	enrichedDomain, err := sh.enrichmentService.EnrichStarters(ctx, []*model.Starter{starter})
	if err != nil {
		return nil, err
	}

	enrichedDTO := starterdto.FromDomainEnrichment(enrichedDomain)
	return starterdto.FromDomainEnriched(starter, enrichedDTO), nil
}

// SoftDeleteStarter DELETE /api/v1/starters/{domain}
func (sh *StarterHandler) SoftDeleteStarter(ctx *gin.Context) {
	httputil.Wrap(sh.softDeleteStarter)(ctx)
}

func (sh *StarterHandler) softDeleteStarter(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		Domain string `uri:"domain" binding:"required"`
	}
	if err := httputil.ValidateURI(ctx, &uriReq); err != nil {
		return nil, err
	}

	if err := sh.starterSvc.SoftDeleteStarter(ctx, uriReq.Domain); err != nil {
		return nil, err
	}

	return gin.H{
		"message": "Starter soft deleted successfully",
		"domain":  uriReq.Domain,
	}, nil
}
