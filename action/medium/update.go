package medium

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/bindu-server/config"
	"github.com/factly/bindu-server/model"
	"github.com/factly/bindu-server/util"
	"github.com/factly/bindu-server/util/slug"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// update - Update medium by id
// @Summary Update a medium by id
// @Description Update medium by ID
// @Tags Medium
// @ID update-medium-by-id
// @Produce json
// @Consume json
// @Param X-User header string true "User ID"
// @Param medium_id path string true "Medium ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param Medium body medium false "Medium"
// @Success 200 {object} model.Medium
// @Router /media/{medium_id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	mediumID := chi.URLParam(r, "medium_id")
	id, err := strconv.Atoi(mediumID)

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
	medium := &medium{}
	err = json.NewDecoder(r.Body).Decode(&medium)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(medium)

	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &model.Medium{}
	result.ID = uint(id)

	// check record exists or not
	err = config.DB.Where(&model.Medium{
		OrganisationID: uint(oID),
	}).First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	var mediumSlug string

	if result.Slug == medium.Slug {
		mediumSlug = result.Slug
	} else if medium.Slug != "" && slug.Check(medium.Slug) {
		mediumSlug = slug.Approve(medium.Slug, oID, config.DB.NewScope(&model.Medium{}).TableName())
	} else {
		mediumSlug = slug.Approve(slug.Make(medium.Name), oID, config.DB.NewScope(&model.Medium{}).TableName())
	}

	config.DB.Model(&result).Updates(model.Medium{
		Name: medium.Name,
		Slug: mediumSlug,
		Type: medium.Type,
		URL:  medium.URL,
	}).First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
