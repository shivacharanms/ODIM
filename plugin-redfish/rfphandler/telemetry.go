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

//Package rfphandler ...
package rfphandler

import (
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpmodel"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfputilities"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// UpdateTrigger updates the trigger parameters with read-write enabled
func UpdateTrigger(ctx iris.Context) {
	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	//Validating the token
	if token != "" {
		flag := TokenValidation(token)
		if !flag {
			log.Error("Invalid/Expired X-Auth-Token")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return
		}
	}
	var deviceDetails rfpmodel.Device
	uri := ctx.Request().RequestURI
	//Get device details from request
	err := ctx.ReadJSON(&deviceDetails)
	if err != nil {
		errMsg := "Unable to collect data from request: " + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString(errMsg)
		return
	}

	// prepare the device data
	var devices []rfpmodel.Device
	rfpmodel.GetAllDevicesInInventory(&devices)
	var successCount int
	var notFoundCount int
	var internalErrorCount int
	for _, device := range devices {
		device.PostBody = deviceDetails.PostBody
		response := sendRequestToDevice(device, uri)
		switch response.StatusCode {
		case http.StatusOK:
			successCount++
		case http.StatusNotFound:
			notFoundCount++
		default:
			internalErrorCount++
		}

	}
	if successCount > 0 {
		ctx.StatusCode(http.StatusOK)
		return
	}
}

func sendRequestToDevice(device rfpmodel.Device, uri string) *http.Response {
	redfishClient, err := rfputilities.GetRedfishClient()
	if err != nil {
		errMsg := "While trying to get the redfish client, got: " + err.Error()
		log.Error(errMsg)
		var resp *http.Response
		resp.StatusCode = http.StatusInternalServerError
		return resp
	}
	deviceDetails := &rfputilities.RedfishDevice{
		Host:     device.Host,
		Username: device.Username,
		Password: string(device.Password),
		PostBody: device.PostBody,
	}
	resp, err := redfishClient.DeviceCall(deviceDetails, uri, http.MethodPatch)
	if err != nil {
		errorMessage := "While trying to patch triggers, got:" + err.Error()
		log.Error(errorMessage)
	}
	return resp
}
