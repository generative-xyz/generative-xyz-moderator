package response

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

type IResponse interface {
	SetID(string)
	GetID() string
}

type BaseResponse struct {
	ID string `json:"id"`
}

func (p *BaseResponse) SetID(ID string) {
	p.ID = ID
}

func (p BaseResponse) GetID() string {
	return p.ID
}

type IHttpResponse interface {
	RespondWithError(w http.ResponseWriter, httpCode int, appCode int, payload error)
	RespondSuccess(w http.ResponseWriter, httpCode int, appCode int, payload interface{}, customerMessage string)
	RespondWithoutContainer(w http.ResponseWriter, httpCode int, payload interface{})
}

type JsonResponse struct {
	Error  *RespondErr `json:"error"`
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

type RespondErr struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"code"`
}

type httpResponse struct {
}

func NewHttpResponse() *httpResponse {
	return new(httpResponse)
}

func (h *httpResponse) RespondWithError(w http.ResponseWriter, httpCode int, appCode int, payload error) {
	h.respondWithJSON(w, payload, httpCode, appCode, payload, "")
}

func (h *httpResponse) RespondSuccess(w http.ResponseWriter, httpCode int, appCode int, payload interface{}, customerMessage string) {
	h.respondWithJSON(w, nil, httpCode, appCode, payload, customerMessage)
}

func (h *httpResponse) respondWithJSON(w http.ResponseWriter, respErr error, httpCode int, appCode int, payload interface{}, customerMessage string) {

	code := ResponseMessage[appCode].Code
	//message := ResponseMessage[appCode].Message

	if customerMessage != "" {
		//message = customerMessage
	}

	jsr := JsonResponse{
		Data:   payload,
		Status: true,
	}

	if respErr != nil {
		errMessage := &RespondErr{}
		errMessage.Message = respErr.Error()
		errMessage.ErrorCode = code
		jsr.Error = errMessage
		jsr.Status = false
	}

	response, _ := json.Marshal(jsr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	_, err := w.Write(response)
	if err != nil {
		panic(err)
	}
}

func (h *httpResponse) RespondWithoutContainer(w http.ResponseWriter, httpCode int, payload interface{}) {

	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	_, err := w.Write(response)
	if err != nil {
		panic(err)
	}
}

func CopyEntityToRes(toValue IResponse, from entity.IEntity) error {
	err := copier.Copy(toValue, from)
	if err != nil {
		return err
	}
	toValue.SetID(from.GetID())
	return nil
}

func CopyEntityToResNoID(toValue IResponse, from entity.IEntityNoID) error {
	err := copier.Copy(toValue, from)
	if err != nil {
		return err
	}
	return nil
}

// HandlerFunc --
type HandlerFunc func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error)

type restHandlerTemplate struct {
	handlerFunc HandlerFunc
	httpResp    *httpResponse
}

// NewRESTHandlerTemplate --
func NewRESTHandlerTemplate(handlerFunc HandlerFunc) http.Handler {
	return &restHandlerTemplate{
		handlerFunc: handlerFunc,
		httpResp:    NewHttpResponse(),
	}
}

func (h *restHandlerTemplate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	userUuid, ok := ctx.Value(utils.SIGNED_USER_ID).(string)
	if ok && userUuid != "" {
		vars[utils.SIGNED_USER_ID] = userUuid
	}
	userWalletAddress, ok := ctx.Value(utils.SIGNED_WALLET_ADDRESS).(string)
	if ok && userWalletAddress != "" {
		vars[utils.SIGNED_WALLET_ADDRESS] = userWalletAddress
	}
	item, err := h.handlerFunc(ctx, r, vars)
	if err != nil {
		h.httpResp.RespondWithError(w, http.StatusBadRequest, Error, err)
		return
	}
	h.httpResp.RespondSuccess(w, http.StatusOK, Success, item, "")
}
