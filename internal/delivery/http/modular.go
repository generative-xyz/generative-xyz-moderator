package http

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
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
	resp, err := h.Usecase.GroupListModulars(ctx, f)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	//
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) PreviewModularInscriptions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inscriptionID := vars["tokenID"]

	html, err := h.Usecase.PreviewTokenByTokenID(inscriptionID)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
	return
}

func (h *httpDelivery) LoadContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inscriptionID := vars["inscriptionID"]

	html, err := h.Usecase.LoadContent(inscriptionID)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
	return
}
