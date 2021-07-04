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
	"net/http"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-telemetry/tcommon"
	"github.com/ODIM-Project/ODIM/svc-telemetry/tmodel"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
)

// ExternalInterface struct holds the structs to which hold function pointers to outboud calls
type ExternalInterface struct {
	External External
	DB       DB
}

// External struct holds the function pointers all outboud services
type External struct {
	ContactClient      func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	Auth               func(string, []string, []string) response.RPC
	DevicePassword     func([]byte) ([]byte, error)
	GetPluginData      func(string) (tmodel.Plugin, *errors.Error)
	ContactPlugin      func(tcommon.PluginContactRequest, string) ([]byte, string, tcommon.ResponseStatus, error)
	GetTarget          func(string) (*tmodel.Target, *errors.Error)
	CreateTask         func(string) (string, error)
	CreateChildTask    func(string, string) (string, error)
	GetSessionUserName func(string) (string, error)
	GenericSave        func([]byte, string, string) error
	GetPluginStatus    func(tmodel.Plugin) bool
	UpdateTask         func(common.TaskData) error
}

type responseStatus struct {
	StatusCode    int32
	StatusMessage string
	MsgArgs       []interface{}
}

// DB struct holds the function pointers to database operations
type DB struct {
	GetAllKeysFromTable func(string, common.DbType) ([]string, error)
	GetResource         func(string, string, common.DbType) (string, *errors.Error)
}

// GetExternalInterface retrieves all the external connections update package functions uses
func GetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		External: External{
			ContactClient:      pmbhandle.ContactPlugin,
			Auth:               services.IsAuthorized,
			DevicePassword:     common.DecryptWithPrivateKey,
			GetPluginData:      tmodel.GetPluginData,
			ContactPlugin:      tcommon.ContactPlugin,
			GetTarget:          tmodel.GetTarget,
			GetSessionUserName: services.GetSessionUserName,
			GenericSave:        tmodel.GenericSave,
			GetPluginStatus:    tcommon.GetPluginStatus,
			CreateTask:         services.CreateTask,
			UpdateTask:         TaskData,
			CreateChildTask:    services.CreateChildTask,
		},
		DB: DB{
			GetAllKeysFromTable: tmodel.GetAllKeysFromTable,
			GetResource:         tmodel.GetResource,
		},
	}
}

func fillTaskData(taskID, targetURI, request string, resp response.RPC, taskState string, taskStatus string, percentComplete int32, httpMethod string) common.TaskData {
	return common.TaskData{
		TaskID:          taskID,
		TargetURI:       targetURI,
		TaskRequest:     request,
		Response:        resp,
		TaskState:       taskState,
		TaskStatus:      taskStatus,
		PercentComplete: percentComplete,
		HTTPMethod:      httpMethod,
	}
}

// TaskData update the task with the given data
func TaskData(taskData common.TaskData) error {
	respBody, _ := json.Marshal(taskData.Response.Body)
	payLoad := &taskproto.Payload{
		HTTPHeaders:   taskData.Response.Header,
		HTTPOperation: taskData.HTTPMethod,
		JSONBody:      taskData.TaskRequest,
		StatusCode:    taskData.Response.StatusCode,
		TargetURI:     taskData.TargetURI,
		ResponseBody:  respBody,
	}

	err := services.UpdateTask(taskData.TaskID, taskData.TaskState, taskData.TaskStatus, taskData.PercentComplete, payLoad, time.Now())
	if err != nil && (err.Error() == common.Cancelling) {
		// We cant do anything here as the task has done it work completely, we cant reverse it.
		//Unless if we can do opposite/reverse action for delete server which is add server.
		services.UpdateTask(taskData.TaskID, common.Cancelled, taskData.TaskStatus, taskData.PercentComplete, payLoad, time.Now())
		if taskData.PercentComplete == 0 {
			return fmt.Errorf("error while starting the task: %v", err)
		}
		runtime.Goexit()
	}
	return nil
}