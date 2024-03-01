package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"pet-clinic.bonglee.com/internal/app"
	"pet-clinic.bonglee.com/internal/constants/alertConstants"
	"pet-clinic.bonglee.com/internal/models"
	"pet-clinic.bonglee.com/internal/validator"
)

type VetHandler struct {
	*app.App
}

func NewVetHandler(app *app.App) *VetHandler {
	return &VetHandler{
		app,
	}
}

type vetListForm struct {
	PageLen int
	Vets    []models.Vet
}

func (handler *VetHandler) vetList(w http.ResponseWriter, r *http.Request) {

	pageSizeInt := atoiWithDefault(r.URL.Query().Get("pageSize"), 10)
	pageInt := atoiWithDefault(r.URL.Query().Get("page"), 1)

	pageLen, err := handler.Vets.GetVetsPageLen(pageSizeInt)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	vets, err := handler.Vets.GetVets(pageInt, pageSizeInt)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	data := handler.NewTemplateData(r)
	data.Form = vetListForm{
		PageLen: pageLen,
		Vets:    vets,
	}

	handler.Render(w, r, http.StatusOK, "vet-list.html", data)
}

type vetDetailForm struct {
	Id        int
	FirstName string
	LastName  string
	Visits    []models.VisitDetailDto
}

func (handler *VetHandler) vetDetail(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	vetId, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		handler.ClientError(w, r, http.StatusBadRequest)
		return
	}

	vet, err := handler.Vets.GetById(vetId)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	visits, err := handler.Visits.GetByVetId(vetId)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	data := handler.NewTemplateData(r)
	data.Form = vetDetailForm{
		Id:        vet.Id,
		FirstName: vet.FirstName,
		LastName:  vet.LastName,
		Visits:    visits,
	}

	handler.Render(w, r, http.StatusOK, "vet-detail.html", data)
}

type createVetForm struct {
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	validator.Validator `form:"-"`
}

func (handler *VetHandler) vetCreate(w http.ResponseWriter, r *http.Request) {
	data := handler.NewTemplateData(r)
	data.Form = createVetForm{}

	handler.Render(w, r, http.StatusOK, "vet-create.html", data)
}

func (handler *VetHandler) vetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form createVetForm

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	form.Validator.CheckField(validator.NotBlank(form.FirstName), "firstName", "First name cannot be empty")
	form.Validator.CheckField(validator.NotBlank(form.LastName), "lastName", "Last name cannot be empty")

	if !form.Validator.Valid() {
		data := handler.NewTemplateData(r)
		data.Form = form
		handler.Render(w, r, http.StatusUnprocessableEntity, "vet-create.html", data)
		return
	}

	vetId, err := handler.Vets.Insert(form.FirstName, form.LastName)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	w.Header().Add("HX-Redirect", fmt.Sprintf("/vet/detail/%v", vetId))
}

type editVetForm struct {
	Vet                 models.Vet `json:"vet"`
	validator.Validator `form:"-"`
}

func (handler *VetHandler) vetEdit(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	vetId := atoiWithDefault(params.ByName("id"), 0)

	vet, err := handler.Vets.GetById(vetId)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	data := handler.NewTemplateData(r)
	data.Form = editVetForm{
		Vet: vet,
	}

	handler.Render(w, r, http.StatusOK, "vet-edit.html", data)
}

func (handler *VetHandler) vetEditPut(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id := atoiWithDefault(params.ByName("id"), 0)
	var form editVetForm

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	updatedVet := form.Vet

	form.Validator.CheckField(validator.NotBlank(updatedVet.FirstName), "firstName", "First name cannot be blank")
	form.Validator.CheckField(validator.NotBlank(updatedVet.LastName), "lastName", "Last name cannot be blank")

	if !form.Validator.Valid() {
		data := handler.NewTemplateData(r)
		data.Form = form
		handler.Render(w, r, http.StatusUnprocessableEntity, "vet-edit.html", data)
		return
	}

	err = handler.Vets.Update(id, updatedVet.FirstName, updatedVet.LastName)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	w.Header().Add("HX-Redirect", fmt.Sprintf("/vet/detail/%v", id))
}

func (handler *VetHandler) vetRemove(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := atoiWithDefault(params.ByName("id"), 0)

	err := handler.Vets.Remove(id)
	if err != nil || id == 0 {
		handler.ServerError(w, r, err)
		return
	}

	handler.Session.Put(r.Context(), alertConstants.FLASH_MSG, "Vet removed")
	w.Header().Add("HX-Redirect", "/vet")
}

type createVisitForm struct {
	models.CreateVisitDto `json:"visit"`
	Visits                []models.VisitDetailDto
	validator.Validator   `form:"-"`
}

func (handler *VetHandler) vetCreateVisitPost(w http.ResponseWriter, r *http.Request) {
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

	data.Form = &createVisitForm{
		Visits: visits,
	}

	data.Alert = app.Alert{MsgType: alertConstants.SUCCESS, Msg: "Appointment created."}
	handler.RenderPartial(w, r, http.StatusOK, "appt-form.html", data)
}
