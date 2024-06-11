package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
var testsRolesPostID = []struct {
	name         string          //name of string
	description  string          // description of the test case
	route        string          // route path to test
	role_id      string          //path param
	post_data    models.RolePost // patch_data
	expectedCode int             // expected HTTP status code
}{
	// First test case
	{
		name:        "patch Roles By ID check - 1",
		description: "patch Single Role by ID",
		route:       "/admin/role",
		post_data: models.RolePost{
			Name:        "New one Posted 3",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},

	// Second test case
	{
		name:        "get Role By ID check - 2",
		description: "get HTTP status 404, when Role Does not exist",
		route:       "/admin/role",
		post_data: models.RolePost{
			Name:        "Name one",
			Description: "Description of Name one",
		},
		expectedCode: 500,
	},
}

func TestPostRolesByID(t *testing.T) {

	// loading env file
	godotenv.Load(".test.env")

	// Setup Test APP
	TestApp := echo.New()

	// Iterate through test single test cases
	for _, test := range testsRolesPostID {
		t.Run(test.name, func(t *testing.T) {
			//  changing post data to json
			post_data, _ := json.Marshal(test.post_data)

			req := httptest.NewRequest(http.MethodPost, test.route, bytes.NewReader(post_data))

			// Add specfic headers if needed as below
			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx := TestApp.NewContext(req, resp)
			echo_contx.SetPath(test.route)

			// Now testing the GetRoles funciton
			controlers.PostRole(echo_contx)

			var responseMap map[string]interface{}
			body, _ := io.ReadAll(resp.Body)
			uerr := json.Unmarshal(body, &responseMap)
			if uerr != nil {
				// fmt.Printf("Error marshaling response : %v", uerr)
				fmt.Println()
			}

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test.description)
			//  running delete test if post is success
			if resp.Result().StatusCode == 200 {
				t.Run("Checking the Delete Request Path for Roles", func(t *testing.T) {

					test_route := fmt.Sprintf("%v/:%v", test.route, "role_id")

					req_delete := httptest.NewRequest(http.MethodDelete, test.route, bytes.NewReader(post_data))

					// Add specfic headers if needed as below
					req_delete.Header.Set("Content-Type", "application/json")

					//  this is the response recorder
					resp_delete := httptest.NewRecorder()

					//  create echo context to test the app function
					echo_contx_del := TestApp.NewContext(req_delete, resp_delete)
					echo_contx_del.SetPath(test_route)

					// seting path paramenters
					path_value := fmt.Sprintf("%v", responseMap["data"].(map[string]interface{})["id"])
					echo_contx_del.SetParamNames("role_id")
					echo_contx_del.SetParamValues(path_value)

					// Now testing the GetRoles funciton
					controlers.DeleteRole(echo_contx_del)
					assert.Equalf(t, 200, resp.Result().StatusCode, test.description+"deleteing")
				})
			}
		})
	}

}

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
			Name:        "Name one eight",
			Description: "Description of Name one for test one",
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
			Name:        "Name one eight",
			Description: "Description of Name one for test 2",
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
			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx := TestApp.NewContext(req, resp)
			echo_contx.SetPath(test.route)

			// seting path paramenters
			echo_contx.SetParamNames("role_id")
			echo_contx.SetParamValues(test.role_id)

			// Now testing the GetRoles funciton
			controlers.PatchRole(echo_contx)

			// fmt.Println("########")
			// fmt.Println(resp.Result().StatusCode)
			// body, _ := io.ReadAll(resp.Result().Body)
			// fmt.Println(string(body))
			// fmt.Println("########")

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
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx := TestApp.NewContext(req, resp)
			echo_contx.SetPath(test.route)
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
			// req.Header.Set("X-APP-TOKEN", "hi")

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
