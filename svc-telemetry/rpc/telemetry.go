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

package rpc

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// SESSAUTHFAILED string constant to raise errors
const SESSAUTHFAILED string = "Unable to authenticate session"

// GetTelemetryService is an rpc handler, it gets invoked during GET on TelemetryService API (/redfis/v1/TelemetryService/)
func (a *Telemetry) GetTelemetryService(ctx context.Context, req *teleproto.TelemetryRequest, resp *teleproto.TelemetryResponse) error {
	fillProtoResponse(resp, a.connector.GetTelemetryService())
	return nil
}

// GetMetricDefinitionCollection an rpc handler which is invoked during GET on MetricDefinition Collection
func (a *Telemetry) GetMetricDefinitionCollection(ctx context.Context, req *teleproto.TelemetryRequest, resp *teleproto.TelemetryResponse) error {
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn(SESSAUTHFAILED)
		fillProtoResponse(resp, authResp)
		return nil
	}
	fillProtoResponse(resp, a.connector.GetMetricDefinitionCollection(req))
	return nil
}

// GetMetricReportDefinitionCollection is an rpc handler which is invoked during GET on MetricReportDefinition Collection
func (a *Telemetry) GetMetricReportDefinitionCollection(ctx context.Context, req *teleproto.TelemetryRequest, resp *teleproto.TelemetryResponse) error {
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn(SESSAUTHFAILED)
		fillProtoResponse(resp, authResp)
		return nil
	}
	fillProtoResponse(resp, a.connector.GetMetricReportDefinitionCollection(req))
	return nil
}

// GetMetricReportCollection is an rpc handler which is invoked during GET on MetricReport Collection
func (a *Telemetry) GetMetricReportCollection(ctx context.Context, req *teleproto.TelemetryRequest, resp *teleproto.TelemetryResponse) error {
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn(SESSAUTHFAILED)
		fillProtoResponse(resp, authResp)
		return nil
	}
	fillProtoResponse(resp, a.connector.GetMetricReportCollection(req))
	return nil
}

// GetTriggerCollection is an rpc handler which is invoked during GET on TriggerCollection
func (a *Telemetry) GetTriggerCollection(ctx context.Context, req *teleproto.TelemetryRequest, resp *teleproto.TelemetryResponse) error {
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn(SESSAUTHFAILED)
		fillProtoResponse(resp, authResp)
		return nil
	}
	fillProtoResponse(resp, a.connector.GetTriggerCollection(req))
	return nil
}

// GetMetricDefinition is an rpc handler which is invoked during GET on MetricDefinition
func (a *Telemetry) GetMetricDefinition(ctx context.Context, req *teleproto.TelemetryRequest, resp *teleproto.TelemetryResponse) error {
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn(SESSAUTHFAILED)
		fillProtoResponse(resp, authResp)
		return nil
	}
	fillProtoResponse(resp, a.connector.GetMetricDefinition(req))
	return nil
}

// GetMetricReportDefinition is an rpc handler which is invoked during GET on MetricReportDefinition
func (a *Telemetry) GetMetricReportDefinition(ctx context.Context, req *teleproto.TelemetryRequest, resp *teleproto.TelemetryResponse) error {
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn(SESSAUTHFAILED)
		fillProtoResponse(resp, authResp)
		return nil
	}
	fillProtoResponse(resp, a.connector.GetMetricReportDefinition(req))
	return nil
}

// GetMetricReport is an rpc handler which is invoked during GET on MetricReport
func (a *Telemetry) GetMetricReport(ctx context.Context, req *teleproto.TelemetryRequest, resp *teleproto.TelemetryResponse) error {
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn(SESSAUTHFAILED)
		fillProtoResponse(resp, authResp)
		return nil
	}
	fillProtoResponse(resp, a.connector.GetMetricReport(req))
	return nil
}

// GetTrigger is an rpc handler which is invoked during GET on Triggers
func (a *Telemetry) GetTrigger(ctx context.Context, req *teleproto.TelemetryRequest, resp *teleproto.TelemetryResponse) error {
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn(SESSAUTHFAILED)
		fillProtoResponse(resp, authResp)
		return nil
	}
	fillProtoResponse(resp, a.connector.GetTrigger(req))
	return nil
}

// UpdateTrigger is an rpc handler which is invoked during update on Trigger
func (a *Telemetry) UpdateTrigger(ctx context.Context, req *teleproto.TelemetryRequest, resp *teleproto.TelemetryResponse) error {
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn(SESSAUTHFAILED)
		fillProtoResponse(resp, authResp)
		return nil
	}
	sessionUserName, err := a.connector.External.GetSessionUserName(req.SessionToken)
	if err != nil {
		errMsg := "error while trying to get the session username: " + err.Error()
		generateRPCResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil), resp)
		log.Warn(errMsg)
		return nil
	}
	taskURI, err := a.connector.External.CreateTask(sessionUserName)
	if err != nil {
		errMsg := "error while trying to create task: " + err.Error()
		generateRPCResponse(common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), resp)
		log.Warn(errMsg)
		return nil
	}
	strArray := strings.Split(taskURI, "/")
	var taskID string
	if strings.HasSuffix(taskURI, "/") {
		taskID = strArray[len(strArray)-2]
	} else {
		taskID = strArray[len(strArray)-1]
	}
	err = a.connector.External.UpdateTask(common.TaskData{
		TaskID:          taskID,
		TargetURI:       taskURI,
		TaskState:       common.Running,
		TaskStatus:      common.OK,
		PercentComplete: 0,
		HTTPMethod:      http.MethodPost,
	})
	if err != nil {
		log.Warn("error while contacting task-service with UpdateTask RPC : " + err.Error())
	}
	go a.connector.UpdateTrigger(taskID, sessionUserName, req)
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Content-type": "application/json; charset=utf-8",
			"Location":     "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	generateRPCResponse(rpcResp, resp)
	//fillProtoResponse(resp, a.connector.UpdateTrigger(req))
	return nil
}
