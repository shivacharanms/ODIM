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
package handle
import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	telemetryproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// TelemetryRPCs used to define the service RPC function
type TelemetryRPCs struct {
	GetTelemetryServiceRPC               func(telemetryproto.TelemetryRequest) (*telemetryproto.TelemetryResponse, error)
}

// GetTelemetryService is the handler for getting TelemetryService details
func (a *TelemetryRPCs) GetTelemetryService(ctx iris.Context) {
	req := telemetryproto.TelemetryRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := a.GetTelemetryServiceRPC(req)
	if err != nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}