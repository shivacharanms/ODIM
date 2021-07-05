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

//Package telemetry ...
package telemetry

// ---------------------------------------------------------------------------------------
// IMPORT Section
// ---------------------------------------------------------------------------------------
import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime"
	"strings"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-telemetry/tcommon"
	tlresp "github.com/ODIM-Project/ODIM/svc-telemetry/tlresponse"
	"github.com/ODIM-Project/ODIM/svc-telemetry/tmodel"
)

// GetTelemetryService defines the functionality for knowing whether
// the telemetry service is enabled or not
//
// As return parameters RPC response, which contains status code, message, headers and data,
// error will be passed back.
func (e *ExternalInterface) GetTelemetryService() response.RPC {
	commonResponse := response.Response{
		OdataType:    "#TelemetryService.v1_2.TelemetryService",
		OdataID:      "/redfish/v1/TelemetryService",
		OdataContext: "/redfish/v1/$metadata#TelemetryService.TelemetryService",
		ID:           "TelemetryService",
		Name:         "Telemetry Service",
	}
	var resp response.RPC

	isServiceEnabled := false
	serviceState := "Disabled"
	//Checks if TelemetryService is enabled and sets the variable isServiceEnabled to true add servicState to enabled
	for _, service := range config.Data.EnabledServices {
		if service == "TelemetryService" {
			isServiceEnabled = true
			serviceState = "Enabled"
		}
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	resp.Header = map[string]string{
		"Allow":         "GET",
		"Cache-Control": "no-cache",
		"Connection":    "Keep-alive",
		"Content-type":  "application/json; charset=utf-8",
		"Link": "	</redfish/v1/SchemaStore/en/TelemetryService.json>; rel=describedby",
		"Transfer-Encoding": "chunked",
		"X-Frame-Options":   "sameorigin",
	}

	commonResponse.CreateGenericResponse(resp.StatusMessage)
	commonResponse.Message = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""
	resp.Body = tlresp.TelemetryService{
		Response: commonResponse,
		Status: tlresp.Status{
			State:        serviceState,
			Health:       "OK",
			HealthRollup: "OK",
		},
		ServiceEnabled: isServiceEnabled,
		MetricDefinitions: &dmtf.Link{
			Oid: "/redfish/v1/TelemetryService/MetricDefinitions",
		},
		MetricReportDefinitions: &dmtf.Link{
			Oid: "/redfish/v1/TelemetryService/MetricReportDefinitions",
		},
		MetricReports: &dmtf.Link{
			Oid: "/redfish/v1/TelemetryService/MetricReports",
		},
		Triggers: &dmtf.Link{
			Oid: "/redfish/v1/TelemetryService/Triggers",
		},
	}

	return resp

}

// GetMetricDefinitionCollection is a functioanlity to retrive all the available inventory
// resources from the added BMC's
func (e *ExternalInterface) GetMetricDefinitionCollection(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	metricDefinitionCollection := tlresp.Collection{
		OdataContext: "/redfish/v1/$metadata#MetricDefinitionCollection.MetricDefinitionCollection",
		OdataID:      "/redfish/v1/TelemetryService/MetricDefinitionCollection",
		OdataType:    "#MetricDefinitionCollection.MetricDefinitionCollection",
		Description:  "MetricDefinition Collection view",
		Name:         "MetricDefinitionCollection",
	}

	members := []dmtf.Link{}
	metricDefinitionCollectionKeysArray, err := e.DB.GetAllKeysFromTable("MetricDefinition", common.InMemory)
	if err != nil || len(metricDefinitionCollectionKeysArray) == 0 {
		log.Warn("odimra doesnt have servers")
	}

	for _, key := range metricDefinitionCollectionKeysArray {
		members = append(members, dmtf.Link{Oid: key})
	}
	metricDefinitionCollection.Members = members
	metricDefinitionCollection.MembersCount = len(members)
	resp.Body = metricDefinitionCollection
	resp.StatusCode = http.StatusOK
	return resp
}

// GetMetricReportDefinitionCollection is a functioanlity to retrive all the available inventory
// resources from the added BMC's
func (e *ExternalInterface) GetMetricReportDefinitionCollection(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	metricReportDefinitionCollection := tlresp.Collection{
		OdataContext: "/redfish/v1/$metadata#MetricReportDefinition.MetricReportDefinition",
		OdataID:      "/redfish/v1/TelemetryService/MetricReportDefinition",
		OdataType:    "#MetricReportDefinitionCollection.MetricReportDefinitionCollection",
		Description:  "MetricReportDefinition Collection view",
		Name:         "MetricReportDefinition",
	}

	members := []dmtf.Link{}
	metricReportDefinitionCollectionKeysArray, err := e.DB.GetAllKeysFromTable("MetricReportDefinition", common.InMemory)
	if err != nil || len(metricReportDefinitionCollectionKeysArray) == 0 {
		log.Warn("odimra doesnt have servers")
	}

	for _, key := range metricReportDefinitionCollectionKeysArray {
		members = append(members, dmtf.Link{Oid: key})
	}
	metricReportDefinitionCollection.Members = members
	metricReportDefinitionCollection.MembersCount = len(members)
	resp.Body = metricReportDefinitionCollection
	resp.StatusCode = http.StatusOK
	return resp
}

// GetMetricReportCollection is a functioanlity to retrive all the available inventory
// resources from the added BMC's
func (e *ExternalInterface) GetMetricReportCollection(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	metricReportCollection := tlresp.Collection{
		OdataContext: "/redfish/v1/$metadata#MetricReportCollection.MetricReportCollection",
		OdataID:      "/redfish/v1/TelemetryService/MetricReport",
		OdataType:    "#MetricReportCollection.MetricReportCollection",
		Description:  "MetricReport Collection view",
		Name:         "MetricReportCollection",
	}

	members := []dmtf.Link{}
	metricReportCollectionKeysArray, err := e.DB.GetAllKeysFromTable("MetricReport", common.InMemory)
	if err != nil || len(metricReportCollectionKeysArray) == 0 {
		log.Warn("odimra doesnt have servers")
	}

	for _, key := range metricReportCollectionKeysArray {
		members = append(members, dmtf.Link{Oid: key})
	}
	metricReportCollection.Members = members
	metricReportCollection.MembersCount = len(members)
	resp.Body = metricReportCollection
	resp.StatusCode = http.StatusOK
	return resp
}

// GetTriggerCollection is a functioanlity to retrive all the available inventory
// resources from the added BMC's
func (e *ExternalInterface) GetTriggerCollection(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	triggersCollection := tlresp.Collection{
		OdataContext: "/redfish/v1/$metadata#TriggerCollection.TriggerCollection",
		OdataID:      "/redfish/v1/TelemetryService/Triggers",
		OdataType:    "#TriggerCollection.TriggerCollection",
		Description:  "Triggers Collection view",
		Name:         "Triggers",
	}

	members := []dmtf.Link{}
	triggersCollectionKeysArray, err := e.DB.GetAllKeysFromTable("Triggers", common.InMemory)
	if err != nil || len(triggersCollectionKeysArray) == 0 {
		log.Warn("odimra doesnt have servers")
	}

	for _, key := range triggersCollectionKeysArray {
		members = append(members, dmtf.Link{Oid: key})
	}
	triggersCollection.Members = members
	triggersCollection.MembersCount = len(members)
	resp.Body = triggersCollection
	resp.StatusCode = http.StatusOK
	return resp
}

// GetMetricReportDefinition ...
func (e *ExternalInterface) GetMetricReportDefinition(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	data, gerr := e.DB.GetResource("MetricReportDefinition", req.URL, common.InMemory)
	if gerr != nil {
		log.Warn("Unable to get MetricReportDefinition details : " + gerr.Error())
		errorMessage := gerr.Error()
		if errors.DBKeyNotFound == gerr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"MetricReportDefinition", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp

}

// GetMetricReport is for to get metric report from southbound resource
func (e *ExternalInterface) GetMetricReport(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	var getDeviceInfoRequest = tcommon.ResourceInfoRequest{
		URL:                 req.URL,
		ContactClient:       e.External.ContactClient,
		DevicePassword:      e.External.DevicePassword,
		GetPluginStatus:     e.External.GetPluginStatus,
		GetAllKeysFromTable: e.DB.GetAllKeysFromTable,
		GetPluginData:       e.External.GetPluginData,
	}
	data, err := tcommon.GetResourceInfoFromDevice(getDeviceInfoRequest)
	if err != nil {
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"MetricReport", req.URL}, nil)
	}
	var resource map[string]interface{}
	json.Unmarshal(data, &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp

}

// GetMetricDefinition ...
func (e *ExternalInterface) GetMetricDefinition(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	data, gerr := e.DB.GetResource("MetricDefinition", req.URL, common.InMemory)
	if gerr != nil {
		log.Warn("Unable to get MetricDefinition details : " + gerr.Error())
		errorMessage := gerr.Error()
		if errors.DBKeyNotFound == gerr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"MetricDefinition", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp

}

// GetTrigger ...
func (e *ExternalInterface) GetTrigger(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	data, gerr := e.DB.GetResource("Triggers", req.URL, common.InMemory)
	if gerr != nil {
		log.Warn("Unable to get Triggers details : " + gerr.Error())
		errorMessage := gerr.Error()
		if errors.DBKeyNotFound == gerr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"Triggers", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp

}

// UpdateTrigger ...
func (e *ExternalInterface) UpdateTrigger(taskID string, sessionUserName string, req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	var percentComplete int32
	serverURI := req.URL
	log.Info("Request in telemetry service")
	log.Info(req.RequestBody)
	log.Info(serverURI)
	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: serverURI, UpdateTask: e.External.UpdateTask, TaskRequest: string(req.RequestBody)}

	//empty request check
	if isEmptyRequest(req.RequestBody) {
		errMsg := "empty request can not be processed"
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"request body"}, nil)
	}
	// parsing the reuqest body
	var trigger dmtf.Triggers
	err := json.Unmarshal(req.RequestBody, &trigger)
	if err != nil {
		errMsg := "unable to parse the request" + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.InternalError, errMsg, nil, nil)
	}
	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, trigger)
	if err != nil {
		errMsg := "Request parameters validaton failed: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	} else if invalidProperties != "" {
		errorMessage := "One or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Error(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		return resp
	}

	pluginList, err := tmodel.GetAllKeysFromTable("Plugin", common.OnDisk)
	if err != nil {
		errMsg := "Request parameters validaton failed: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	log.Info("Count of plugin list")
	log.Info(len(pluginList))
	targetList := formTargetList(pluginList)
	partialResultFlag := false
	subTaskChannel := make(chan int32, len(targetList))
	for _, target := range targetList {
		marshalBody, err := json.Marshal(trigger)
		if err != nil {
			errMsg := "Unable to parse the trigger update request" + err.Error()
			log.Warn(errMsg)
			return common.GeneralError(http.StatusBadRequest, response.InternalError, errMsg, nil, taskInfo)
		}
		triggerRequestBody := string(marshalBody)
		updateURL := req.URL
		go e.sendRequest(updateURL, taskID, target, triggerRequestBody, subTaskChannel, sessionUserName)
	}
	resp.StatusCode = http.StatusOK
	for i := 0; i < len(targetList); i++ {
		select {
		case statusCode := <-subTaskChannel:
			if statusCode != http.StatusOK {
				partialResultFlag = true
				if resp.StatusCode < statusCode {
					resp.StatusCode = statusCode
				}
			}
			if i < len(targetList)-1 {
				percentComplete := int32(((i + 1) / len(targetList)) * 100)
				var task = fillTaskData(taskID, serverURI, string(req.RequestBody), resp, common.Running, common.OK, percentComplete, http.MethodPatch)
				err := e.External.UpdateTask(task)
				if err != nil && err.Error() == common.Cancelling {
					task = fillTaskData(taskID, serverURI, string(req.RequestBody), resp, common.Cancelled, common.OK, percentComplete, http.MethodPatch)
					e.External.UpdateTask(task)
					runtime.Goexit()
				}
			}
		}
	}
	taskStatus := common.OK
	if partialResultFlag {
		taskStatus = common.Warning
	}
	percentComplete = 100
	if resp.StatusCode != http.StatusOK {
		errMsg := "One or more of the trigger update requests failed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID
		log.Warn(errMsg)
		switch resp.StatusCode {
		case http.StatusAccepted:
			return common.GeneralError(http.StatusAccepted, response.TaskStarted, errMsg, []interface{}{fmt.Sprintf("%v", targetList)}, taskInfo)
		case http.StatusUnauthorized:
			return common.GeneralError(http.StatusUnauthorized, response.ResourceAtURIUnauthorized, errMsg, []interface{}{fmt.Sprintf("%v", targetList)}, taskInfo)
		case http.StatusNotFound:
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"option", "Triggers"}, taskInfo)
		case http.StatusBadRequest:
			return common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errMsg, []interface{}{"Triggers"}, taskInfo)
		default:
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		}
	}

	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	log.Info("All Trigger updates requests successfully completed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID)
	resp.StatusMessage = response.Success
	resp.StatusCode = http.StatusOK
	args := response.Args{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully",
	}
	resp.Body = args.CreateGenericErrorResponse()

	var task = fillTaskData(taskID, serverURI, string(req.RequestBody), resp, common.Completed, taskStatus, percentComplete, http.MethodPatch)
	err = e.External.UpdateTask(task)
	if err != nil && err.Error() == common.Cancelling {
		task = fillTaskData(taskID, serverURI, string(req.RequestBody), resp, common.Cancelled, common.Critical, percentComplete, http.MethodPatch)
		e.External.UpdateTask(task)
		runtime.Goexit()
	}
	return resp
}

func (e *ExternalInterface) sendRequest(serverURI, taskID string, plugin tmodel.Plugin, updateRequestBody string, subTaskChannel chan<- int32, sessionUserName string) {
	log.Info("INSIDE send request")
	var resp response.RPC
	subTaskURI, err := e.External.CreateChildTask(sessionUserName, taskID)
	if err != nil {
		subTaskChannel <- http.StatusInternalServerError
		log.Warn("Unable to create sub task")
		return
	}
	var subTaskID string
	strArray := strings.Split(subTaskURI, "/")
	if strings.HasSuffix(subTaskURI, "/") {
		subTaskID = strArray[len(strArray)-2]
	} else {
		subTaskID = strArray[len(strArray)-1]
	}
	taskInfo := &common.TaskUpdateInfo{TaskID: subTaskID, TargetURI: serverURI, UpdateTask: e.External.UpdateTask, TaskRequest: updateRequestBody}

	var percentComplete int32
	var contactRequest tcommon.PluginContactRequest
	contactRequest.ContactClient = e.External.ContactClient
	contactRequest.Plugin = plugin

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		contactRequest.HTTPMethodType = http.MethodPatch
		contactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.OID = "/ODIM/v1/Sessions"
		_, token, getResponse, err := e.External.ContactPlugin(contactRequest, "error while creating session with the plugin: ")

		if err != nil {
			subTaskChannel <- getResponse.StatusCode
			errMsg := err.Error()
			log.Warn(errMsg)
			common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo)
			return
		}
		contactRequest.Token = token
	} else {
		contactRequest.BasicAuth = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	var target tmodel.Target
	target.PostBody = []byte(updateRequestBody)
	contactRequest.DeviceInfo = target
	contactRequest.OID = serverURI
	contactRequest.HTTPMethodType = http.MethodPatch
	_, _, getResponse, err := e.External.ContactPlugin(contactRequest, "error while performing trigger update action: ")
	if err != nil {
		subTaskChannel <- getResponse.StatusCode
		errMsg := err.Error()
		log.Warn(errMsg)
		common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo)
		return
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	resp.StatusCode = http.StatusOK
	percentComplete = 100
	subTaskChannel <- int32(getResponse.StatusCode)
	var task = fillTaskData(subTaskID, serverURI, updateRequestBody, resp, common.Completed, common.OK, percentComplete, http.MethodPatch)
	err = e.External.UpdateTask(task)
	if err != nil && err.Error() == common.Cancelling {
		var task = fillTaskData(subTaskID, serverURI, updateRequestBody, resp, common.Cancelled, common.Critical, percentComplete, http.MethodPatch)
		e.External.UpdateTask(task)
	}
	return
}

func isEmptyRequest(requestBody []byte) bool {
	var updateRequest map[string]interface{}
	json.Unmarshal(requestBody, &updateRequest)
	if len(updateRequest) <= 0 {
		return true
	}
	return false
}

func formTargetList(keys []string) []tmodel.Plugin {
	var plugins []tmodel.Plugin
	for _, key := range keys {
		plugin, err := tmodel.GetPluginData(key)
		if err != nil {
			log.Error("failed to get details of " + key + " plugin: " + err.Error())
			continue
		}
		plugins = append(plugins, plugin)
	}
	return plugins
}
