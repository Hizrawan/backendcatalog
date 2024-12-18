package public

import (
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/app"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/controllers"
	httperr "github.com/xinchuantw/hoki-tabloid-backend/internal/errors"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/responses"
	stringsutil "github.com/xinchuantw/hoki-tabloid-backend/utils/strings"
)

type ContactController struct {
	controllers.Controller
}

func NewContactController(app *app.Registry) *ContactController {
	return &ContactController{
		Controller: controllers.Controller{
			App: app,
		},
	}
}

func (c *ContactController) CheckDuplicate(w http.ResponseWriter, r *http.Request) {
	actor := r.URL.Query().Get("actor")
	if actor == "" {
		panic(validation.Errors{"actor": validation.NewError("invalid_actor", "actor is required")})
	}

	value := r.URL.Query().Get("value")
	if value == "" {
		panic(validation.Errors{"value": validation.NewError("invalid_value", "value is required")})
	}

	typeParam := r.URL.Query().Get("type")
	if typeParam == "" {
		panic(validation.Errors{"type": validation.NewError("invalid_type", "type is required")})
	}

	medium := r.URL.Query().Get("medium")
	if medium == "" {
		panic(validation.Errors{"medium": validation.NewError("invalid_medium", "medium is required")})
	}

	// Validation based on type
	var isValid bool
	switch strings.ToLower(medium) {
	case "phone":
		isValid = stringsutil.ValidateTaiwanPhone(value)
	case "landline":
		isValid = stringsutil.ValidateTaiwanLandline(value)
	case "email":
		isValid = stringsutil.ValidateEmail(value)
	case "line_id":
		isValid = stringsutil.ValidateLineID(value)
	default:
		panic(httperr.NewErrUnprocessableEntity(
			"invalid_medium",
			"unsupported medium", nil,
		))
	}

	if !isValid {
		panic(httperr.NewErrUnprocessableEntity(
			"invalid_value",
			"value does not match the required format for the specified type", nil,
		))
	}

	isExist := true

	if err := responses.JSON(w, 200, struct {
		Data bool `json:"data"`
	}{
		Data: isExist,
	}); err != nil {
		panic(err)
	}
}
