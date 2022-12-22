package http

import (
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
)

// UserCredits godoc
// @Summary Get configs
// @Description Get configs
// @Tags Configs
// @Accept  json
// @Produce  json
// @Success 200 {object} response.JsonResponse{data=response.ConfigResp}
// @Router /configs [GET]
func (h *httpDelivery) getConfigs(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getConfigs", r)
	defer h.Tracer.FinishSpan(span, log )


	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}

// UserCredits godoc
// @Summary create config
// @Description create config
// @Tags Configs
// @Accept  json
// @Produce  json
// @Success 200 {object} response.JsonResponse{data=response.ConfigResp}
// @Router /configs [POST]
func (h *httpDelivery) createConfig(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("createConfig", r)
	defer h.Tracer.FinishSpan(span, log )


	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}

// UserCredits godoc
// @Summary delete config
// @Description delete config
// @Tags Configs
// @Accept  json
// @Produce  json
// @Param key path string true "config key"
// @Success 200 {object} response.JsonResponse{data=response.ConfigResp}
// @Router /configs/{key} [DELETE]
func (h *httpDelivery) deleteConfig(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("deleteConfigs", r)
	defer h.Tracer.FinishSpan(span, log )


	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}

// UserCredits godoc
// @Summary get one config
// @Description get one config
// @Tags Configs
// @Accept  json
// @Produce  json
// @Param key path string true "config key"
// @Success 200 {object} response.JsonResponse{data=response.ConfigResp}
// @Router /configs/{key} [GET]
func (h *httpDelivery) getConfig(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getConfig", r)
	defer h.Tracer.FinishSpan(span, log )


	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}

