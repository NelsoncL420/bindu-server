package chart

import (
	"net/http"
	"strconv"

	"github.com/factly/bindu-server/config"
	"github.com/factly/bindu-server/model"
	"github.com/factly/bindu-server/util"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// details - Get chart by id
// @Summary Show a chart by id
// @Description Get chart by ID
// @Tags Chart
// @ID get-chart-by-id
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param chart_id path string true "Chart ID"
// @Success 200 {object} model.Chart
// @Router /charts/{chart_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	chartID := chi.URLParam(r, "chart_id")
	id, err := strconv.Atoi(chartID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	oID, err := util.GetOrganisation(r.Context())

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	result := &model.Chart{}

	result.ID = uint(id)

	err = config.DB.Model(&model.Chart{}).Where(&model.Chart{
		OrganisationID: uint(oID),
	}).Preload("Medium").Preload("Theme").Preload("Tags").Preload("Categories").First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	renderx.JSON(w, http.StatusOK, result)
}
