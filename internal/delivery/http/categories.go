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
// @Summary Get Categorys
// @Description Get Categorys
// @Tags Categories
// @Accept  json
// @Produce  json
// @Success 200 {object} response.JsonResponse{data=response.CategoryResp}
// @Router /categories [GET]
func (h *httpDelivery) getCategories(w http.ResponseWriter, r *http.Request) {

	f := structure.FilterCategories{}
	baseF, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	f.BaseFilters =*baseF

	data, err := h.Usecase.GetCategories(f)
	if err != nil {
		h.Logger.Error("h.Usecase.GetCategorys", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := []response.CategoryResp{}
	iConfs := data.Result
	confs := iConfs.([]entity.Categories)

	for _, conf := range confs  {
		respItem := &response.CategoryResp{}
		err := response.CopyEntityToRes(respItem, &conf)
		if err != nil {
			h.Logger.Error("response.CopyEntityToRes", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	resp = append(resp, *respItem)
	}

	
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(data, resp), "")
}

// UserCredits godoc
// @Summary create Category
// @Description create Category
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param request body request.CreateCategoryRequest true "Create a Category"
// @Success 200 {object} response.JsonResponse{data=response.CategoryResp}
// @Router /categories [POST]
func (h *httpDelivery) createCategory(w http.ResponseWriter, r *http.Request) {

	var reqBody request.CreateCategoryRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	err = reqBody.Validate()
	if err != nil {
		h.Logger.Error("reqBody.Validate", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	category, err := h.Usecase.CreateCategory(structure.CategoryData{
		Name: *reqBody.Name,
	})

	if err != nil {
		h.Logger.Error("h.Usecase.CreateCategory", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	resp := &response.CategoryResp{}
	response.CopyEntityToRes(resp, category)

	
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary delete Category
// @Description delete Category
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {object} response.JsonResponse{data=response.CategoryResp}
// @Router /categories/{id} [DELETE]
func (h *httpDelivery) deleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := h.Usecase.DeleteCategory(id)
	if err != nil {
		h.Logger.Error("h.Usecase.DeleteCategory", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}

// UserCredits godoc
// @Summary update Category
// @Description update Category
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param request body request.CreateCategoryRequest true "Create a Category"
// @Param id path string true "Category ID"
// @Success 200 {object} response.JsonResponse{data=response.CategoryResp}
// @Router /categories/{id} [PUT]
func (h *httpDelivery) updateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var reqBody request.CreateCategoryRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	cat,  err := h.Usecase.UpdateCategory(structure.UpdateCategoryData{
		Name: reqBody.Name,
		ID: &id,
	})

	if err != nil {
		h.Logger.Error("h.Usecase.UpdateCategory", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := &response.CategoryResp{}
	response.CopyEntityToRes(resp, cat)

	
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary get one Category
// @Description get one Category
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {object} response.JsonResponse{data=response.CategoryResp}
// @Router /categories/{id} [GET]
func (h *httpDelivery) getCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	category, err := h.Usecase.GetCategory(id)
	if err != nil {
		h.Logger.Error("h.Usecase.GetCategory", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	resp := &response.CategoryResp{}
	response.CopyEntityToRes(resp, category)
	
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

