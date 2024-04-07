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

type VetHandler struct {
	*app.App
}

func NewVetHandler(app *app.App) *VetHandler {
	return &VetHandler{
		app,
	}
}

type vetListForm struct {
	PageLen  []int
	PageSize int
	Vets     []models.Vet
	LastName string
}

func (handler *VetHandler) list(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	pageSize := atoiWithDefault(query.Get("pageSize"), 5)
	page := atoiWithDefault(query.Get("page"), 0)
	isSideList := query.Get("sideList") == "true"

	pageLen, err := handler.Vets.GetAllVetsPageLen(pageSize)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	vets, err := handler.Vets.GetVets(page, pageSize)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	data := handler.NewTemplateData(r)
	data.Form = vetListForm{
		PageLen:  make([]int, pageLen),
		PageSize: pageSize,
		Vets:     vets,
	}

	if isSideList {
		handler.RenderPartial(w, r, http.StatusOK, "vet-list-side.html", data)
		return
	}

	handler.Render(w, r, http.StatusOK, "vet-list-link.html", data)
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

func (handler *VetHandler) getByLastName(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	page := atoiWithDefault(query.Get("page"), 0)
	pageSize := atoiWithDefault(query.Get("pageSize"), 3)
	lastName := query.Get("lastName")

	data := handler.NewTemplateData(r)

	vets, err := handler.Vets.GetVetsByLastName(lastName, page, pageSize)
	if err != nil {
		handler.Logger.Error(err.Error())
		data.Alert = app.Alert{MsgType: alertConstants.DANGER, Msg: "Error while searching vets."}
		handler.RenderPartial(w, r, http.StatusBadRequest, "alert.html", data)
		return
	}

	pageLen, err := handler.Vets.GetVetsPageLenLastName(pageSize, lastName)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	data.Form = vetListForm{
		PageLen:  make([]int, pageLen),
		PageSize: pageSize,
		Vets:     vets,
		LastName: lastName,
	}

	handler.RenderPartial(w, r, http.StatusOK, "vet-list-side.html", data)
}
