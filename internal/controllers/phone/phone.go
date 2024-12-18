package controller

import (
	"net/http"
	"strconv"

	"github.com/xinchuantw/hoki-tabloid-backend/internal/app"
	controllers "github.com/xinchuantw/hoki-tabloid-backend/internal/controllers"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type PhoneController struct {
	controllers.Controller
}

func NewPhoneController(app *app.Registry) *PhoneController {
	return &PhoneController{controllers.Controller{App: app}}
}

// GetPhones retrieves all phone records with pagination, sorting, and filtering
func (c *PhoneController) GetPhones(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	sortBy := r.URL.Query().Get("sortBy")
	order := r.URL.Query().Get("order")
	filterBy := r.URL.Query().Get("filterBy")
	filterValue := r.URL.Query().Get("filterValue")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10 // Default limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0 // Default offset
	}

	phones, err := models.GetPhones(c.App.DB, limit, offset, sortBy, order, filterBy, filterValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, phones)
}

// GetPhone retrieves a single phone record by ID
func (c *PhoneController) GetPhone(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "PhoneID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	phone, err := models.GetPhone(c.App.DB, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, phone)
}

// CreatePhone creates a new phone record and inserts installment values
func (c *PhoneController) CreatePhone(w http.ResponseWriter, r *http.Request) {
	var phone models.Phone
	if err := render.Bind(r, &phone); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx := c.App.DB.MustBegin()
	err := phone.Insert(tx)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create phone record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, phone)
}

// UpdatePhone updates an existing phone record by ID and records price changes
func (c *PhoneController) UpdatePhone(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "PhoneID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var phone models.Phone
	if err := render.Bind(r, &phone); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	phone.ID = id

	tx := c.App.DB.MustBegin()
	err = phone.Update(tx)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update phone record", http.StatusInternalServerError)
		return
	}

	installment := models.CalculateInstallments(phone.Price)
	installment.PhoneID = phone.ID
	// Delete existing installment records
	_, err = tx.Exec("DELETE FROM installments WHERE phone_id = ?", phone.ID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to delete existing installment records", http.StatusInternalServerError)
		return
	}

	err = installment.Insert(tx)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create installment record", http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, phone)
}

// DeletePhone deletes a phone record by ID
func (c *PhoneController) DeletePhone(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "PhoneID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	phone := models.Phone{ID: id}

	tx := c.App.DB.MustBegin()
	err = phone.Delete(tx)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to delete phone record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
