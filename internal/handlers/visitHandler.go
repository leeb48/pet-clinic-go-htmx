package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"pet-clinic.bonglee.com/internal/app"
	"pet-clinic.bonglee.com/internal/constants/alertConstants"
	"pet-clinic.bonglee.com/internal/models"
	"pet-clinic.bonglee.com/internal/validator"
)

type VisitHandler struct {
	*app.App
}

func NewVisitHandler(app *app.App) *VisitHandler {
	return &VisitHandler{
		app,
	}
}

type VisitDetailForm struct {
	VisitDetail models.VisitDetailDto
	VisitJson   string
	VetPageSize int
}

func (handler *VisitHandler) visitDetail(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	visitId := atoiWithDefault(params.ByName("id"), 0)

	visit, err := handler.Visits.GetById(visitId)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	visitJson, err := json.Marshal(visit)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	data := handler.NewTemplateData(r)
	data.Form = &VisitDetailForm{
		VisitDetail: visit,
		VisitJson:   string(visitJson),
		VetPageSize: 3,
	}

	handler.Render(w, r, http.StatusOK, "visit-detail.html", data)
}

type createVisitForm struct {
	models.CreateVisitDto `json:"visit"`
	Visits                string
	VetPageSize           int
	validator.Validator   `form:"-"`
}

func (handler *VisitHandler) editVisitPage(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	visitId := atoiWithDefault(params.ByName("id"), 0)

	visit, err := handler.Visits.GetById(visitId)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	data := handler.NewTemplateData(r)
	data.Form = &VisitDetailForm{
		VisitDetail: visit,
		VetPageSize: 3,
	}

	handler.Render(w, r, http.StatusOK, "visit-edit.html", data)
}

func (handler *VisitHandler) createVisitPost(w http.ResponseWriter, r *http.Request) {
	var form createVisitForm

	data := handler.NewTemplateData(r)

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		data.Alert = app.Alert{MsgType: alertConstants.DANGER, Msg: "Please check to make sure all inputs are correct."}
		handler.RenderPartial(w, r, http.StatusBadRequest, "alert.html", data)
		return
	}

	form.CheckField(validator.NotNilId(form.PetId), "error", "Pet ID could not be found")
	form.CheckField(validator.NotNilId(form.VetId), "error", "Vet ID could not be found")

	if !form.Validator.Valid() {
		data.Alert = app.Alert{MsgType: alertConstants.DANGER,
			Msg: fmt.Sprintf("%v", form.Validator.FieldErrors["error"])}
		handler.RenderPartial(w, r, http.StatusBadRequest, "alert.html", data)
		return
	}

	err = handler.Visits.Create(form.PetId, form.VetId, form.Appointment, form.VisitReason, form.Duration)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	visits, err := handler.Visits.GetByVetId(form.VetId)
	if err != nil {
		data.Alert = app.Alert{MsgType: alertConstants.WARNING, Msg: "Failed to get latests visits."}
		handler.RenderPartial(w, r, http.StatusBadRequest, "alert.html", data)
		return
	}

	visitsJson, err := json.Marshal(visits)
	if err != nil {
		data.Alert = app.Alert{MsgType: alertConstants.WARNING, Msg: "Failed to get latests visits."}
		handler.RenderPartial(w, r, http.StatusBadRequest, "alert.html", data)
		return
	}

	data.Form = &createVisitForm{
		Visits: string(visitsJson),
	}

	data.Alert = app.Alert{MsgType: alertConstants.SUCCESS, Msg: "Appointment created."}
	handler.RenderPartial(w, r, http.StatusOK, "appt-form.html", data)
}

func (handler *VisitHandler) getVisitCalendarByVetId(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	vetId := atoiWithDefault(params.ByName("id"), 0)

	visits, err := handler.Visits.GetByVetId(vetId)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	visitJson, err := json.Marshal(visits)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	data := handler.NewTemplateData(r)
	data.Form = &VisitDetailForm{
		VisitJson: string(visitJson),
	}

	handler.RenderPartial(w, r, http.StatusOK, "visit-calendar.html", data)
}

func (handler *VisitHandler) removeVisit(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	visitId := atoiWithDefault(params.ByName("visitId"), 0)
	vetId := atoiWithDefault(params.ByName("vetId"), 0)

	err := handler.Visits.Remove(visitId)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	handler.Session.Put(r.Context(), alertConstants.FLASH_MSG, "Visit removed")
	w.Header().Add("HX-Redirect", fmt.Sprintf("/vet/detail/%v", vetId))
}
