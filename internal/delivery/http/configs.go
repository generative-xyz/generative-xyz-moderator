package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
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

	data, err := h.Usecase.GetConfigs(span, structure.FilterConfigs{})
	if err != nil {
		log.Error("h.Usecase.GetConfigs", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := []response.ConfigResp{}
	iConfs := data.Result
	confs := iConfs.([]entity.Configs)

	for _, conf := range confs  {
		respItem := &response.ConfigResp{}
		err := response.CopyEntityToRes(respItem, &conf)
		if err != nil {
			log.Error("response.CopyEntityToRes", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	
		resp = append(resp, *respItem)
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(data, resp), "")
}

// UserCredits godoc
// @Summary create config
// @Description create config
// @Tags Configs
// @Accept  json
// @Produce  json
// @Param request body request.CreateConfigRequest true "Create a config"
// @Success 200 {object} response.JsonResponse{data=response.ConfigResp}
// @Router /configs [POST]
func (h *httpDelivery) createConfig(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("createConfig", r)
	defer h.Tracer.FinishSpan(span, log )

	var reqBody request.CreateConfigRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	err = reqBody.Validate()
	if err != nil {
		log.Error("reqBody.Validate", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	config, err := h.Usecase.CreateConfig(span, structure.ConfigData{
		Key: *reqBody.Key,
		Value: *reqBody.Value,
	})

	if err != nil {
		log.Error("h.Usecase.CreateConfig", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	resp := &response.ConfigResp{}
	response.CopyEntityToRes(resp, config)

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
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
	vars := mux.Vars(r)
	key := vars["key"]
	err := h.Usecase.DeleteConfig(span, key)
	if err != nil {
		log.Error("h.Usecase.DeleteConfig", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
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
	vars := mux.Vars(r)
	key := vars["key"]
	config, err := h.Usecase.GetConfig(span, key)
	if err != nil {
		log.Error("h.Usecase.GetConfig", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	resp := &response.ConfigResp{}
	response.CopyEntityToRes(resp, config)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

