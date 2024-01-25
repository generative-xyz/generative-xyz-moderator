package http

import (
	"context"
	"net/http"
	"os"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
)

func (h *httpDelivery) ModularInscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	//vars := mux.Vars(r)
	//genNFTAddr := vars["genNFTAddr"]

	f := structure.FilterTokens{}
	err := f.CreateFilter(r)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	genNFTAddr := os.Getenv("MODULAR_PROJECT_ID")
	f.GenNFTAddr = &genNFTAddr
	bf, err := h.BaseFilters(r)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f.BaseFilters = *bf
	resp, err := h.Usecase.ListModulars(ctx, f)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	//
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}
