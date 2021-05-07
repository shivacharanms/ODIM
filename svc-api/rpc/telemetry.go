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

//Package rpc ...
package rpc

import (
	"context"
	"fmt"

	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

// DoGetTelemetryService defines the RPC call function for
// the GetTelemetryService from telemetry micro service
func DoGetTelemetryService(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	telemetry := teleproto.NewTelemetryService(services.Telemetry, services.Service.Client())

	resp, err := telemetry.GetTelemetryService(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetMetricDefinitionCollection defines the RPC call function for
// the GetMetricDefinitionCollectionRPC from telemetry micro service
func DoGetMetricDefinitionCollection(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	telemetry := teleproto.NewTelemetryService(services.Telemetry, services.Service.Client())

	resp, err := telemetry.GetMetricDefinitionCollectionRPC(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetMetricReportDefinitionCollection defines the RPC call function for
// the GetMetricReportDefinitionCollectionRPC from telemetry micro service
func DoGetMetricReportDefinitionCollection(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	telemetry := teleproto.NewTelemetryService(services.Telemetry, services.Service.Client())

	resp, err := telemetry.GetMetricReportDefinitionCollectionRPC(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetMetricReportCollection defines the RPC call function for
// the GetMetricReportCollectionRPC from telemetry micro service
func DoGetMetricReportCollection(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	telemetry := teleproto.NewTelemetryService(services.Telemetry, services.Service.Client())

	resp, err := telemetry.GetMetricReportCollectionRPC(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetTriggerCollection defines the RPC call function for
// the GetTriggerCollectionRPC from telemetry micro service
func DoGetTriggerCollection(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	telemetry := teleproto.NewTelemetryService(services.Telemetry, services.Service.Client())

	resp, err := telemetry.GetTriggerCollectionRPC(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetMetricDefinition defines the RPC call function for
// the GetMetricDefinitionRPC from telemetry micro service
func DoGetMetricDefinition(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	telemetry := teleproto.NewTelemetryService(services.Telemetry, services.Service.Client())

	resp, err := telemetry.GetMetricDefinitionRPC(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetMetricReportDefinition defines the RPC call function for
// the GetMetricReportDefinitionRPC from telemetry micro service
func DoGetMetricReportDefinition(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	telemetry := teleproto.NewTelemetryService(services.Telemetry, services.Service.Client())

	resp, err := telemetry.GetMetricReportDefinitionRPC(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetMetricReport defines the RPC call function for
// the GetMetricReportRPC from telemetry micro service
func DoGetMetricReport(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	telemetry := teleproto.NewTelemetryService(services.Telemetry, services.Service.Client())

	resp, err := telemetry.GetMetricReportRPC(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetTrigger defines the RPC call function for
// the GetTriggerRPC from telemetry micro service
func DoGetTrigger(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	telemetry := teleproto.NewTelemetryService(services.Telemetry, services.Service.Client())

	resp, err := telemetry.GetTriggerRPC(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoUpdateTrigger defines the RPC call function for
// the UpdateTriggerRPC from telemetry micro service
func DoUpdateTrigger(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	telemetry := teleproto.NewTelemetryService(services.Telemetry, services.Service.Client())

	resp, err := telemetry.UpdateTriggerRPC(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}