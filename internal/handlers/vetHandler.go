package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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

func (handler *VetHandler) list(w http.ResponseWriter, r *http.Request) {

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
	Visits    string
}

func (handler *VetHandler) detail(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	vetId := atoiWithDefault(params.ByName("id"), 0)

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

	visitsJson, err := json.Marshal(visits)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	data := handler.NewTemplateData(r)
	data.Form = vetDetailForm{
		Id:        vet.Id,
		FirstName: vet.FirstName,
		LastName:  vet.LastName,
		Visits:    string(visitsJson),
	}

	handler.Render(w, r, http.StatusOK, "vet-detail.html", data)
}

type createVetForm struct {
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	validator.Validator `form:"-"`
}

func (handler *VetHandler) createPage(w http.ResponseWriter, r *http.Request) {
	data := handler.NewTemplateData(r)
	data.Form = createVetForm{}

	handler.Render(w, r, http.StatusOK, "vet-create.html", data)
}

func (handler *VetHandler) createPost(w http.ResponseWriter, r *http.Request) {
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

func (handler *VetHandler) edit(w http.ResponseWriter, r *http.Request) {
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

func (handler *VetHandler) editPost(w http.ResponseWriter, r *http.Request) {
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

func (handler *VetHandler) remove(w http.ResponseWriter, r *http.Request) {
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

type searchVetForm struct {
	LastName string `json:"lastName"`
}

func (handler *VetHandler) getByLastName(w http.ResponseWriter, r *http.Request) {
	var form searchVetForm

	data := handler.NewTemplateData(r)

	err := json.NewDecoder(r.Body).Decode(&form)
	err = errors.New("Test error")
	if err != nil {
		data.Alert = app.Alert{MsgType: alertConstants.DANGER, Msg: "Error while parsing search data."}
		handler.RenderPartial(w, r, http.StatusBadRequest, "alert.html", data)
		return
	}

	fmt.Println(form)
}
