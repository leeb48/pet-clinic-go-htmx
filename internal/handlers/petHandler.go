package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"pet-clinic.bonglee.com/internal/app"
	"pet-clinic.bonglee.com/internal/models/customErrors"
	"pet-clinic.bonglee.com/internal/validator"
)

type PetHandler struct {
	*app.App
}

func NewPetHandler(app *app.App) *PetHandler {
	return &PetHandler{
		app,
	}
}

type newPetTypeForm struct {
	NewPetType          string `json:"newPetType"`
	validator.Validator `form:"-"`
}

func (handler *PetHandler) adminPage(w http.ResponseWriter, r *http.Request) {
	data := handler.NewTemplateData(r)
	data.Form = newPetTypeForm{}
	handler.Render(w, r, http.StatusOK, "admin.html", data)
}

func (handler *PetHandler) newPetTypePost(w http.ResponseWriter, r *http.Request) {
	var form newPetTypeForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		handler.ClientError(w, r, http.StatusBadRequest)
		return
	}

	form.NewPetType = strings.ToUpper(form.NewPetType)

	form.CheckField(validator.NotBlank(form.NewPetType), "newPetType", "Pet type cannot be blank")

	if !form.Valid() {
		data := handler.NewTemplateData(r)
		data.Form = form
		handler.Render(w, r, http.StatusUnprocessableEntity, "admin.html", data)
		return
	}

	err = handler.PetTypes.Insert(form.NewPetType)
	if err != nil {

		handler.Logger.Error(err.Error())

		if errors.Is(err, customErrors.ErrDuplicatePetType) {
			form.AddFieldError("newPetType", "Duplicate pet type")
			data := handler.NewTemplateData(r)
			data.Form = form
			handler.Render(w, r, http.StatusUnprocessableEntity, "admin.html", data)
		} else {
			handler.ServerError(w, r, err)
		}
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
