package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
)

func (h *httpDelivery) GetListModularWorkshop(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	f, err := h.BaseFilters(r)
	offset := f.Limit * (f.Page - 1)
	data, err := h.Usecase.Repo.GetListModularWorkShopByAddress(context.Background(), address, offset, f.Limit)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	total, err := h.Usecase.Repo.GetTotalModularWorkShopByAddress(context.Background(), address)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	totalPage := total / f.Limit
	if total%f.Limit != 0 {
		totalPage = totalPage + 1
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, response.PaginationResponse{
		Result:    data,
		PageSize:  f.Limit,
		Page:      f.Page,
		Total:     total,
		TotalPage: totalPage,
	}, "")
}

func (h *httpDelivery) GetModularWorkshopDetail(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("id")
	data, err := h.Usecase.Repo.GetModularWorkshopById(context.Background(), uuid)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, data, "")
}

func (h *httpDelivery) SaveModularWorkshop(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	var data *entity.ModularWorkshopEntity
	err = json.Unmarshal(body, &data)
	if data.ID.IsZero() {
		data.BaseEntity.SetID()
		data.BaseEntity.SetCreatedAt()
		err = h.Usecase.Repo.SaveModularWorkshop(context.Background(), data)
		if err != nil {
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	} else {
		data.BaseEntity.SetUpdatedAt()
		err = h.Usecase.Repo.UpdateModularWorkshop(context.Background(), data)
		if err != nil {
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "", "")
}

func (h *httpDelivery) RemoveModularWorkshop(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	address := r.URL.Query().Get("address")
	err := h.Usecase.Repo.RemoveModularWorkshop(context.Background(), uuid, address)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "", "")
}
