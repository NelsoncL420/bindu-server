package theme

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/bindu-server/config"
	"github.com/factly/bindu-server/model"
	"github.com/factly/bindu-server/util"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - Create theme
// @Summary Create theme
// @Description Create theme
// @Tags Theme
// @ID add-theme
// @Consume json
// @Produce json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param Theme body theme true "Theme Object"
// @Success 201 {object} model.Theme
// @Failure 400 {array} string
// @Router /themes [post]
func create(w http.ResponseWriter, r *http.Request) {

	oID, err := util.GetOrganisation(r.Context())

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	theme := &theme{}

	err = json.NewDecoder(r.Body).Decode(&theme)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(theme)

	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &model.Theme{
		Name:           theme.Name,
		Config:         theme.Config,
		OrganisationID: uint(oID),
	}

	err = config.DB.Model(&model.Theme{}).Create(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	renderx.JSON(w, http.StatusCreated, result)
}
