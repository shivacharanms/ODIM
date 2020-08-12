// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/system"
)

func TestAggregator_GetAggregationService(t *testing.T) {
	config.SetUpMockConfig(t)
	config.Data.EnabledServices = append(config.Data.EnabledServices, "AggregationService")
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive GetAggregationService",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.GetAggregationService(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.GetAggregationService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_AddCompute(t *testing.T) {
	config.SetUpMockConfig(t)
	addComputeRetrieval := config.AddComputeSkipResources{
		SystemCollection: []string{"Chassis", "LogServices"},
	}
	mockPluginData(t, "ILO")

	config.Data.AddComputeSkipResources = &addComputeRetrieval
	system.ActiveReqSet.ReqRecord = make(map[string]interface{})
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

	successReq, _ := json.Marshal(system.AddResourceRequest{
		ManagerAddress: "100.0.0.1:50000",
		UserName:       "admin",
		Password:       "password",
		Oem: &system.AddOEM{
			PluginID:          "GRF",
			PreferredAuthType: "BasicAuth",
			PluginType:        "RF-GENERIC",
		},
	})
	invalidReqBody, _ := json.Marshal(system.AddResourceRequest{
		ManagerAddress: ":50000",
		UserName:       "admin",
		Password:       "password",
		Oem: &system.AddOEM{
			PluginID:          "GRF",
			PreferredAuthType: "BasicAuth",
			PluginType:        "RF-GENERIC",
		},
	})
	missingparamReq, _ := json.Marshal(system.AddResourceRequest{})
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noDetailsToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noTaskToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "taskWithSlashToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "with invalid request",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: []byte("someData")},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Manager Address",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: invalidReqBody},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "with missing parameters",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: missingparamReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.AddCompute(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.AddCompute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_DeleteCompute(t *testing.T) {
	successReq, _ := json.Marshal(`map[string]interface{}{"parameters": []Parameters{{Name: "/redfish/v1/systems/ef83e569-7336-492a-aaee-31c02d9db831:1"}}}`)
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},

		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noDetailsToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noTaskToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "taskWithSlashToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.DeleteCompute(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.DeleteCompute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_deleteServer(t *testing.T) {
	type args struct {
		taskID    string
		targetURI string
		a         *Aggregator
		req       *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			args: args{
				taskID:    "someID",
				targetURI: "someURI",
				a:         &Aggregator{connector: connector},
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  []byte("someData"),
				},
			},
			wantErr: false,
		},
		{
			name: "task updation fails",
			args: args{
				taskID:    "invalid",
				targetURI: "someURI",
				a:         &Aggregator{connector: connector},
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  []byte("successReq"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := deleteServer(tt.args.taskID, tt.args.targetURI, tt.args.a, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("deleteServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_Reset(t *testing.T) {
	successReq, _ := json.Marshal(`map[string]interface{}{"parameters": []Parameters{{Name: "/redfish/v1/systems/ef83e569-7336-492a-aaee-31c02d9db831:1"}}}`)
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noDetailsToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noTaskToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "taskWithSlashToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.Reset(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.Reset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_reset(t *testing.T) {
	type args struct {
		ctx             context.Context
		taskID          string
		sessionUserName string
		req             *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				taskID: "someID",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  []byte("someData"),
				},
			},
			wantErr: false,
		},
		{
			name: "task updation fails",
			a:    &Aggregator{connector: connector},
			args: args{
				taskID: "invalid",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  []byte("successReq"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.reset(tt.args.ctx, tt.args.taskID, tt.args.sessionUserName, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.reset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_SetDefaultBootOrder(t *testing.T) {
	successReq, _ := json.Marshal(`map[string]interface{}{"parameters": []Parameters{{Name: "/redfish/v1/systems/ef83e569-7336-492a-aaee-31c02d9db831:1"}}}`)
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noDetailsToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noTaskToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "taskWithSlashToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.SetDefaultBootOrder(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.SetDefaultBootOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_RediscoverSystemInventory(t *testing.T) {
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.RediscoverSystemInventoryRequest
		resp *aggregatorproto.RediscoverSystemInventoryResponse
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.RediscoverSystemInventoryRequest{
					SystemID:  "someSystemID",
					SystemURL: "someURL",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.RediscoverSystemInventory(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.RediscoverSystemInventory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_ValidateManagerAddress(t *testing.T) {
	type args struct {
		name    string
		arg     string
		wanterr bool
	}
	tests := []struct {
		name    string
		arg     string
		wanterr bool
	}{
		{
			name:    "Valid manager address - IP",
			arg:     "127.0.0.1",
			wanterr: false,
		},
		{
			name:    "Valid manager address - IP and port",
			arg:     "127.0.0.1:1234",
			wanterr: false,
		},
		{
			name:    "Valid manager address - FQDN",
			arg:     "localhost",
			wanterr: false,
		},
		{
			name:    "Valid manager address - FQDN and Port",
			arg:     "localhost:1234",
			wanterr: false,
		},
		{
			name:    "Invalid manager address - IP",
			arg:     "a.b.c.d",
			wanterr: true,
		},
		{
			name:    "Invalid manager address - FQDN",
			arg:     "unknown",
			wanterr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateManagerAddress(tt.arg); (err != nil) != tt.wanterr {
				t.Errorf("validateManagerAddress error = %v, wantErr %v", err, tt.wanterr)
			}
		})
	}
}

func TestAggregator_AddAggreagationSource(t *testing.T) {
	config.SetUpMockConfig(t)
	addComputeRetrieval := config.AddComputeSkipResources{
		SystemCollection: []string{"Chassis", "LogServices"},
	}
	mockPluginData(t, "ILO")

	config.Data.AddComputeSkipResources = &addComputeRetrieval
	system.ActiveReqSet.UpdateMu.Lock()
	system.ActiveReqSet.ReqRecord = make(map[string]interface{})
	system.ActiveReqSet.UpdateMu.Unlock()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	successReq, _ := json.Marshal(system.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: "password",
		Links: &system.Links{
			Oem: &system.AddOEM{
				PluginID:          "ILO",
				PreferredAuthType: "BasicAuth",
				PluginType:        "Compute",
			},
		},
	})
	invalidReqBody, _ := json.Marshal(system.AggregationSource{
		HostName: ":50000",
		UserName: "admin",
		Password: "password",
		Links: &system.Links{
			Oem: &system.AddOEM{
				PluginID:          "ILO",
				PreferredAuthType: "BasicAuth",
				PluginType:        "Compute",
			},
		},
	})
	missingparamReq, _ := json.Marshal(system.AggregationSource{})
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noDetailsToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noTaskToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "taskWithSlashToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "with invalid request",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: []byte("someData")},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Manager Address",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: invalidReqBody},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "with missing parameters",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: missingparamReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.AddAggregationSource(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.AddAggreagationSource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func mockSystemResourceData(body []byte, table, key string) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create(table, key, string(body)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", table, err.Error())
	}
	return nil
}

func TestAggregator_CreateAggregate(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

	reqData, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"})
	err := mockSystemResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}

	reqData1, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1"})
	err = mockSystemResourceData(reqData1, "ComputerSystem", "/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}

	successReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	})
	successReq1, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{},
	})
	invalidReqBody, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/123456",
		},
	})
	missingparamReq, _ := json.Marshal(agmodel.Aggregate{})
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
	}
	tests := []struct {
		name           string
		a              *Aggregator
		args           args
		wantStatusCode int32
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "positive case with empty elements",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq1},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "with invalid request",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: []byte("someData")},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid System",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: invalidReqBody},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "with missing parameters",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: missingparamReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.a.CreateAggregate(tt.args.ctx, tt.args.req, tt.args.resp); tt.args.resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.CreateAggregate() status code = %v, wantStatusCode %v", tt.args.resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestAggregator_GetAllAggregates(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantStatusCode int32
	}{
		{
			name: "Positive cases",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantStatusCode: 0, // to be replaced http.StatusOK
		},
		{
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantStatusCode: 0, // to be replaced http.StatusUnauthorized
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.a.GetAllAggregates(tt.args.ctx, tt.args.req, tt.args.resp); tt.args.resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.GetAllAggregates() error = %v, wantStatusCode %v", tt.args.resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestAggregator_GetAggregate(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantStatusCode int32
	}{
		{
			name: "Positive cases",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", URL: "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantStatusCode: 0, // to be replaced http.StatusOK
		},
		{
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", URL: "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantStatusCode: 0, // to be replaced http.StatusUnauthorized
		},
		{
			name: "Invalid aggregate id",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", URL: "/redfish/v1/AggregationService/Aggregates/1"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantStatusCode: 0, // to be replaced http.StatusBadRequest
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.a.GetAggregate(tt.args.ctx, tt.args.req, tt.args.resp); tt.args.resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.GetAggregate() error = %v, wantStatusCode %v", tt.args.resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestAggregator_DeleteAggregate(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantStatusCode int32
	}{
		{
			name: "Positive cases",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", URL: "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantStatusCode: 0, // to be replaced http.StatusOK
		},
		{
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", URL: "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantStatusCode: 0, // to be replaced http.StatusUnauthorized
		},
		{
			name: "Invalid aggregate id",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", URL: "/redfish/v1/AggregationService/Aggregates/1"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantStatusCode: 0, // to be replaced http.StatusBadRequest
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.a.DeleteAggregate(tt.args.ctx, tt.args.req, tt.args.resp); tt.args.resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.DeleteAggregate() error = %v, wantStatusCode %v", tt.args.resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}