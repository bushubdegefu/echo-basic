package controlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"semay.com/common"
	"semay.com/database"
	"semay.com/models"
)

// GetRoles is a function to get a Roles by ID
// @Summary Get Roles
// @Description Get Roles
// @Tags Role
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security Refresh
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]RoleGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /roles [get]
func GetRoles(contx echo.Context) error {

	//  parsing Query Prameters
	Page, _ := strconv.Atoi(contx.QueryParam("page"))
	Limit, _ := strconv.Atoi(contx.QueryParam("size"))
	//  checking if query parameters  are correct
	if Page == 0 || Limit == 0 {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: "Not Allowed, Bad request",
			Data:    nil,
		})
	}

	//  Getting Database connection
	db := database.ReturnSession()

	//  querying result with pagination using gorm function
	result, err := common.PaginationPureModel(db, models.Role{}, []models.Role{}, uint(Page), uint(Limit))
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: true,
			Message: "Success get all roles.",
			Data:    "something",
		})
	}

	// returning result if all the above completed successfully
	return contx.JSON(http.StatusOK, result)
}

// GetRoleByID is a function to get a Roles by ID
// @Summary Get Role by ID
// @Description Get role by ID
// @Tags Role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} common.ResponseHTTP{data=RoleGet}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /roles/{role_id} [get]
func GetRoleByID(contx echo.Context) error {

	//  parsing Query Prameters
	id, err := strconv.Atoi(contx.Param("role_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  Getting Database connection
	db := database.ReturnSession()

	// Preparing and querying database using Gorm
	var roles_get models.RoleGet
	var roles models.Role
	if res := db.Model(&models.Role{}).Preload(clause.Associations).Where("id = ?", id).First(&roles); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
				Success: false,
				Message: "Role not found",
				Data:    nil,
			})
		}
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: "Error retrieving role",
			Data:    nil,
		})
	}

	// filtering response data according to filtered defined struct
	mapstructure.Decode(roles, &roles_get)

	//  Finally returing response if All the above compeleted successfully
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success got one role.",
		Data:    &roles_get,
	})
}

// Add Role to data
// @Summary Add a new Role
// @Description Add Role
// @Tags Role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role body RolePost true "Add Role"
// @Success 200 {object} common.ResponseHTTP{data=RolePost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /roles [post]
func PostRole(contx echo.Context) error {
	//  parsing Query Prameters
	db := database.ReturnSession()

	// validator initialization
	validate := validator.New()

	//validating post data
	posted_role := new(models.RolePost)

	//first parse request data
	if err := contx.Bind(&posted_role); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(posted_role); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  initiate -> role
	role := new(models.Role)
	role.Name = posted_role.Name
	role.Description = posted_role.Description

	//  start transaction to database
	tx := db.Begin()

	// add  data using transaction if values are valid
	if err := tx.Create(&role).Error; err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: "Role Creation Failed",
			Data:    err,
		})
	}

	// close transaction
	tx.Commit()

	// return data if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Role created successfully.",
		Data:    role,
	})
}

// Patch Role to data
// @Summary Patch Role
// @Description Patch Role
// @Tags Role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role body RolePost true "Patch Role"
// @Param id path int true "Role ID"
// @Success 200 {object} common.ResponseHTTP{data=RolePost}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /role/{role_id} [patch]
func PatchRole(contx echo.Context) error {

	// Get database connection
	db := database.ReturnSession()

	//  initialize data validator
	validate := validator.New()

	// validate path params
	id, err := strconv.Atoi(contx.Param("role_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate data struct
	patch_role := new(models.RolePatch)
	if err := contx.Bind(&patch_role); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validating
	if err := validate.Struct(patch_role); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// startng update transaction
	var role models.Role
	role.ID = uint(id)
	tx := db.Begin()

	// Check if the record exists
	if err := db.First(&role, role.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If the record doesn't exist, return an error response
			tx.Rollback()
			return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
				Success: false,
				Message: "Role not found",
				Data:    nil,
			})
		}
		// If there's an unexpected error, return an internal server error response
		tx.Rollback()
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Update the record
	if err := db.Model(&role).UpdateColumns(*patch_role).Error; err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Return  success response
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Role updated successfully.",
		Data:    role,
	})
}

// DeleteRoles function removes a role by ID
// @Summary Remove Role by ID
// @Description Remove role by ID
// @Tags Role
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /role/{role_id} [delete]
func DeleteRole(contx echo.Context) error {

	// get deleted role attributes to return
	var role models.Role

	// validate path params
	id, err := strconv.Atoi(contx.Param("role_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Getting Database connection
	db := database.ReturnSession()

	// perform delete operation if the object exists
	tx := db.Begin()

	// first getting role and checking if it exists
	if err := db.Where("id = ?", id).First(&role).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
				Success: false,
				Message: "Role not found",
				Data:    nil,
			})
		}
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: "Error retrieving role",
			Data:    nil,
		})
	}

	// Delete the role
	if err := db.Delete(&role).Error; err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: "Error deleting role",
			Data:    nil,
		})
	}

	// Commit the transaction
	tx.Commit()

	// Return success respons
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Role deleted successfully.",
		Data:    role,
	})
}
