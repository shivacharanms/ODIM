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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ODIM-Project/ODIM/plugin-redfish/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpresponse"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"io/ioutil"
	"net/http"
	"testing"
)

func mockUpdateTriggerHandler(username, password, url string, w http.ResponseWriter) {
	resp, err := mockChangeBiosSettings(username, url)
	if err != nil && resp == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(resp.StatusCode)
}

func mockUpdateTrigger(username, url string) (*http.Response, error) {
	if url == "/ODIM/v1/TelemetryService/Triggers/sample" && username == "admin" {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Success")),
		}, nil
	}
	return nil, fmt.Errorf("Error")
}

func TestUpdateTrigger(t *testing.T) {
	config.SetUpMockConfig(t)

	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockUpdateTriggerHandler)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/ODIM/v1")

	redfishRoutes.Patch("/TelemetryService/Triggers/{id}", UpdateTrigger)
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
