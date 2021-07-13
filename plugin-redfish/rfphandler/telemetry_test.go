//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.
// Packahe rfphandler ...
package rfphandler

import (
	"encoding/json"
	"fmt"
	"github.com/ODIM-Project/ODIM/plugin-redfish/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpresponse"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpmodel"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"net/http"
	"testing"
)

func mockDeviceUpdate(device rfpmodel.Device, uri string)*http.Response{
	var response http.Response
	response.StatusCode = http.StatusOK
	return &response
}

func TestUpdateTrigger(t *testing.T) {
	config.SetUpMockConfig(t)

	deviceHost := "localhost"
	devicePort := "1234"
	e := ExternalInterface{
		TokenValidation: tokenValidationMock,
		SendRequestToDevice:   mockDeviceUpdate,
	}
	rfpmodel.DeviceInventory.Device["0e343dc6-f5f3-425a-9503-4a3c799579c8"] = rfpmodel.DeviceData{
		Address:  "172.16.1.205",
		UserName: "admin",
		Password: []byte("Admin123"),
	}
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/ODIM/v1")

	redfishRoutes.Patch("/TelemetryService/Triggers/{id}", e.UpdateTrigger)
	rfpresponse.PluginToken = "token"
	test := httptest.New(t, mockApp)
	attributes := map[string]interface{}{"EventTriggers": []string{"Alert"}}
	attributeByte, _ := json.Marshal(attributes)
	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
		"PostBody":       attributeByte,
	}
	test.PATCH("/ODIM/v1/TelemetryService/Triggers/sample").WithJSON(requestBody).Expect().Status(http.StatusOK)
}
