package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
)

// UserCredits godoc
// @Summary Get Redis
// @Description Get Redis
// @Tags Admin
// @Accept  json
// @Produce  json
// @Success 200 {object} response.JsonResponse{data=[]response.RedisResponse}
// @Router /admin/redis [GET]
func (h *httpDelivery) getRedisKeys(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getRedisKeys", r)
	defer h.Tracer.FinishSpan(span, log )

	res, err := h.Usecase.GetAllRedis(span)

	if err != nil {
		log.Error("h.Usecase.GetRedis", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, res, "")
}


// UserCredits godoc
// @Summary Get Redis
// @Description Get Redis
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param key path string true "Redis key"
// @Success 200 {object} response.JsonResponse{data=response.RedisResponse}
// @Router /admin/redis/{key} [GET]
func (h *httpDelivery) getRedis(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getRedis", r)
	defer h.Tracer.FinishSpan(span, log )

	var err error
	vars := mux.Vars(r)
	redisKey := vars["key"]
	
	res, err := h.Usecase.GetRedis(span, redisKey)

	if err != nil {
		log.Error("h.Usecase.GetRedis", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, &response.RedisResponse{Value: res}, "")
}

// UserCredits godoc
// @Summary Upsert Redis
// @Description Upsert Redis
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param request body request.UpsertRedisRequest true "Upsert redis key"
// @Success 200 {object} response.JsonResponse{data=response.RedisResponse}
// @Router /admin/redis [POST]
func (h *httpDelivery) upsertRedis(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("upsertRedis", r)
	defer h.Tracer.FinishSpan(span, log )

	var reqBody request.UpsertRedisRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	res, err := h.Usecase.UpsertRedis(span, reqBody.Key, reqBody.Value)
	if err != nil {
		log.Error("h.Usecase.UpsertRedis", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, &response.RedisResponse{Value: res}, "")
}

// UserCredits godoc
// @Summary Delete Redis
// @Description Delete Redis
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param key path string true "Redis key"
// @Success 200 {object} response.JsonResponse{data=string}
// @Router /admin/redis/{key} [DELETE]
func (h *httpDelivery) deleteRedis(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("deleteRedis", r)
	defer h.Tracer.FinishSpan(span, log )

	var err error
	vars := mux.Vars(r)
	redisKey := vars["key"]
	
	err = h.Usecase.DeleteRedis(span, redisKey)

	if err != nil {
		log.Error("h.Usecase.DeleteRedis", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "", "")
}
