package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"semay.com/models"
	"semay.com/models/controlers"
)

// ##########################################################################
var testsRolesPatchID = []struct {
	name         string           //name of string
	description  string           // description of the test case
	route        string           // route path to test
	role_id      string           //path param
	patch_data   models.RolePatch // patch_data
	expectedCode int              // expected HTTP status code
}{
	// First test case
	{
		name:        "patch Roles By ID check - 1",
		description: "patch Single Role by ID",
		route:       "/admin/role/:role_id",
		role_id:     "1",
		patch_data: models.RolePatch{
			Name:        "Name one Patched",
			Description: "Description of Name one",
		},
		expectedCode: 200,
	},

	// Second test case
	{
		name:        "get Role By ID check - 2",
		description: "get HTTP status 404, when Role Does not exist",
		route:       "/admin/role/:role_id",
		role_id:     "100",
		patch_data: models.RolePatch{
			Name:        "Name one",
			Description: "Description of Name one",
		},
		expectedCode: 404,
	},
}

func TestPatchRolesByID(t *testing.T) {

	// loading env file
	godotenv.Load(".test.env")

	// Setup Test APP
	TestApp := echo.New()

	// Iterate through test single test cases
	for _, test := range testsRolesPatchID {
		t.Run(test.name, func(t *testing.T) {
			//  changing post data to json
			patch_data, _ := json.Marshal(test.patch_data)

			req := httptest.NewRequest(http.MethodPatch, test.route, bytes.NewReader(patch_data))

			// Add specfic headers if needed as below
			req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx := TestApp.NewContext(req, resp)
			echo_contx.SetPath(test.route)

			// seting path paramenters
			echo_contx.SetParamNames("role_id")
			echo_contx.SetParamValues(test.role_id)

			// Now testing the GetRoles funciton
			controlers.GetRoleByID(echo_contx)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test.description)

		})
	}

}

// ##########################################################################
// Define a structure for specifying input and output data
// of a single test case
var testsRolesGet = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get Roles working - 1",
		description:  "get HTTP status 200",
		route:        "/admin/role?page=1&size=10",
		expectedCode: 200,
	},
	// Second test case
	{
		name:         "get Roles Working - 2",
		description:  "get HTTP status 404, when Role Does not exist",
		route:        "/admin/role?page=1&size=10",
		expectedCode: 200,
	},
}

func TestGetRoles(t *testing.T) {
	// loading env file
	godotenv.Load(".test.env")

	// Setup Test APP
	TestApp := echo.New()

	// Iterate through test single test cases
	for _, test := range testsRolesGet {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			// Add specfic headers if needed as below
			req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx := TestApp.NewContext(req, resp)
			echo_contx.SetPath("/admin/role")
			// Now testing the GetRoles funciton
			controlers.GetRoles(echo_contx)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test.description)

		})
	}

}

// ##############################################################

var testsRolesGetByID = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	role_id      string // path parm
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get Roles By ID check - 1",
		description:  "get Single Role by ID",
		route:        "/admin/role/:role_id",
		role_id:      "1",
		expectedCode: 200,
	},

	// Second test case
	{
		name:         "get Role By ID check - 2",
		description:  "get HTTP status 404, when Role Does not exist",
		route:        "/admin/role/:role_id",
		role_id:      "100",
		expectedCode: 404,
	},
}

func TestGetRolesByID(t *testing.T) {

	// loading env file
	godotenv.Load(".test.env")

	// Setup Test APP
	TestApp := echo.New()

	// Iterate through test single test cases
	for _, test := range testsRolesGetByID {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			// Add specfic headers if needed as below
			req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx := TestApp.NewContext(req, resp)
			echo_contx.SetPath(test.route)

			// seting path paramenters
			echo_contx.SetParamNames("role_id")
			echo_contx.SetParamValues(test.role_id)

			// Now testing the GetRoles funciton
			controlers.GetRoleByID(echo_contx)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test.description)

		})
	}

}
