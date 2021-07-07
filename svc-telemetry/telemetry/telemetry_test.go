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
package telemetry

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-telemetry/tcommon"
	"github.com/ODIM-Project/ODIM/svc-telemetry/tlresponse"
	"github.com/ODIM-Project/ODIM/svc-telemetry/tmodel"
	"github.com/stretchr/testify/assert"
)

func mockIsAuthorized(sessionToken string, privileges, oemPrivileges []string) response.RPC {
	if sessionToken != "validToken" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "error while trying to authenticate session", nil, nil)
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil)
}

func mockContactClient(url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {
	return nil, fmt.Errorf("InvalidRequest")
}

func mockGetResource(table, key string, dbType common.DbType) (string, *errors.Error) {
	if (key == "/redfish/v1/TelemetryService/Triggers/invalidID") || (key == "/redfish/v1/TelemetryService/MetricDefinition/invalidID") {
		return "", errors.PackError(errors.DBKeyNotFound, "not found")
	}
	return "body", nil
}

func mockGetAllKeysFromTable(table string, dbType common.DbType) ([]string, error) {
	return []string{"/redfish/v1/TelemetryService/FirmwareInentory/uuid:1"}, nil
}

func mockGetTarget(id string) (*tmodel.Target, *errors.Error) {
	var target tmodel.Target
	target.PluginID = id
	target.DeviceUUID = "uuid"
	target.UserName = "admin"
	target.Password = []byte("password")
	target.ManagerAddress = "ip"
	return &target, nil
}

func mockGetPluginData(id string) (tmodel.Plugin, *errors.Error) {
	var plugin tmodel.Plugin
	plugin.IP = "ip"
	plugin.Port = "port"
	plugin.Username = "plugin"
	plugin.Password = []byte("password")
	plugin.PluginType = "Redfish"
	plugin.PreferredAuthType = "basic"
	return plugin, nil
}

func mockContactPlugin(req tcommon.PluginContactRequest, errorMessage string) ([]byte, string, tcommon.ResponseStatus, error) {
	var responseStatus tcommon.ResponseStatus

	return []byte(`{"Attributes":"sample"}`), "token", responseStatus, nil
}

func stubDevicePassword(password []byte) ([]byte, error) {
	return password, nil
}

func stubGenericSave(reqBody []byte, table string, uuid string) error {
	return nil
}

func mockGetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		External: External{
			Auth:            mockIsAuthorized,
			ContactClient:   mockContactClient,
			GetTarget:       mockGetTarget,
			GetPluginData:   mockGetPluginData,
			ContactPlugin:   mockContactPlugin,
			DevicePassword:  stubDevicePassword,
			CreateChildTask: mockCreateChildTask,
			UpdateTask:      mockUpdateTask,
			GenericSave:     stubGenericSave,
		},
		DB: DB{
			GetAllKeysFromTable: mockGetAllKeysFromTable,
			GetResource:         mockGetResource,
		},
	}
}

func mockCreateChildTask(sessionID, taskID string) (string, error) {
	switch taskID {
	case "taskWithoutChild":
		return "", fmt.Errorf("subtask cannot created")
	case "subTaskWithSlash":
		return "someSubTaskID/", nil
	default:
		return "someSubTaskID", nil
	}
}

func mockUpdateTask(task common.TaskData) error {
	if task.TaskID == "invalid" {
		return fmt.Errorf("task with this ID not found")
	}
	return nil
}

func TestGetTelemetryService(t *testing.T) {
	successResponse := response.Response{
		OdataType:    "#TelemetryService.v1_2.TelemetryService",
		OdataID:      "/redfish/v1/TelemetryService",
		OdataContext: "/redfish/v1/$metadata#TelemetryService.TelemetryService",
		ID:           "TelemetryService",
		Name:         "Telemetry Service",
	}
	successResponse.CreateGenericResponse(response.Success)
	successResponse.Message = ""
	successResponse.MessageID = ""
	successResponse.Severity = ""
	common.SetUpMockConfig()
	tests := []struct {
		name string
		want response.RPC
	}{
		{
			name: "telemetry service enabled",
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Allow":             "GET",
					"Cache-Control":     "no-cache",
					"Connection":        "Keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Link":              "</redfish/v1/SchemaStore/en/TelemetryService.json>; rel=describedby",
					"Transfer-Encoding": "chunked",
					"X-Frame-Options":   "sameorigin",
				},
				Body: tlresponse.TelemetryService{
					Response: successResponse,
					Status: tlresponse.Status{
						State:        "Enabled",
						Health:       "OK",
						HealthRollup: "OK",
					},
					ServiceEnabled: true,
					MetricDefinitions: dmtf.Link{
						Oid: "/redfish/v1/TelemetryService/MetricDefinitions",
					},
					MetricReportDefinitions: dmtf.Link{
						Oid: "/redfish/v1/TelemetryService/MetricReportDefinitions",
					},
					MetricReports: dmtf.Link{
						Oid: "/redfish/v1/TelemetryService/MetricReports",
					},
					Triggers: dmtf.Link{
						Oid: "/redfish/v1/TelemetryService/Triggers",
					},
				},
			},
		},
	}
	config.Data.EnabledServices = []string{"TelemetryService"}
	e := mockGetExternalInterface()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := e.GetTelemetryService()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTelemetryService() = %v, want %v", got, tt.want)
			}
		})
		config.Data.EnabledServices = []string{"XXXX"}
	}
}

func TestGetMetricDefinitionCollection(t *testing.T) {
	req := &teleproto.TelemetryRequest{}
	e := mockGetExternalInterface()
	response := e.GetMetricDefinitionCollection(req)

	telemetry := response.Body.(tlresponse.Collection)
	assert.Equal(t, int(response.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.Equal(t, telemetry.MembersCount, 1, "Member count does not match")
}

func TestGetMetricReportDefinitionCollection(t *testing.T) {
	req := &teleproto.TelemetryRequest{}
	e := mockGetExternalInterface()
	response := e.GetMetricReportDefinitionCollection(req)

	telemetry := response.Body.(tlresponse.Collection)
	assert.Equal(t, int(response.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.Equal(t, telemetry.MembersCount, 1, "Member count does not match")
}

func TestGetMetricReportCollection(t *testing.T) {
	req := &teleproto.TelemetryRequest{}
	e := mockGetExternalInterface()
	response := e.GetMetricReportCollection(req)

	telemetry := response.Body.(tlresponse.Collection)
	assert.Equal(t, int(response.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.Equal(t, telemetry.MembersCount, 1, "Member count does not match")
}

func TestGetTriggerCollection(t *testing.T) {
	req := &teleproto.TelemetryRequest{}
	e := mockGetExternalInterface()
	response := e.GetTriggerCollection(req)

	telemetry := response.Body.(tlresponse.Collection)
	assert.Equal(t, int(response.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.Equal(t, telemetry.MembersCount, 1, "Member count does not match")
}

func TestGetMetricReportDefinition(t *testing.T) {
	config.SetUpMockConfig(t)
	req := &teleproto.TelemetryRequest{
		ResourceID: "CPUDefinition",
	}
	e := mockGetExternalInterface()
	response := e.GetMetricReportDefinition(req)

	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

}

func TestGetMetricReportDefinitionInvalidID(t *testing.T) {
	req := &teleproto.TelemetryRequest{
		ResourceID: "invalidID",
	}
	e := mockGetExternalInterface()
	response := e.GetMetricReportDefinition(req)

	//Todo: update after individual get
	//assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound")
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusNotFound")
}

func TestGetMetricDefinition(t *testing.T) {
	config.SetUpMockConfig(t)
	req := &teleproto.TelemetryRequest{
		ResourceID: "CPUDefinition",
	}
	e := mockGetExternalInterface()
	response := e.GetMetricDefinition(req)

	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

}

func TestGetMetricDefinitionInvalidID(t *testing.T) {
	req := &teleproto.TelemetryRequest{
		ResourceID: "invalidID",
	}
	e := mockGetExternalInterface()
	response := e.GetMetricDefinition(req)

	//Todo: update after individual get
	//assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound")
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusNotFound")
}

func TestGetTrigger(t *testing.T) {
	config.SetUpMockConfig(t)
	req := &teleproto.TelemetryRequest{
		ResourceID: "CPUDefinition",
	}
	e := mockGetExternalInterface()
	response := e.GetTrigger(req)

	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

}

func TestGetTriggerInvalidID(t *testing.T) {
	req := &teleproto.TelemetryRequest{
		ResourceID: "invalidID",
	}
	e := mockGetExternalInterface()
	response := e.GetTrigger(req)

	//Todo: update after individual get
	//assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound")
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusNotFound")
}
