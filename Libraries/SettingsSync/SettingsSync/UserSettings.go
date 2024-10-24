/*******************************************************************************
 * Copyright 2023-2024 The V.I.S.O.R. authors
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 ******************************************************************************/

package SettingsSync

import (
	"Utils"
	"bytes"
	"errors"
	"time"
)

const _GET_SETTINGS_EACH_S int64 = 30

var last_crc16_GL []byte = nil
var stop_GL bool = false

/*
SyncUserSettings keeps synchronizing the remote user settings file with the local one in background.

-----------------------------------------------------------

– Params:
  - loop – if true, the function will keep running until it's stopped with StopUserSettingsSyncer()
*/
func SyncUserSettings() {
	go func() {
		var last_get_settings_when_s int64 = 0
		for {
			var update_settings bool = false
			if time.Now().Unix() >= last_get_settings_when_s + _GET_SETTINGS_EACH_S && Utils.IsCommunicatorConnectedSERVER() {
				update_settings = true

				last_get_settings_when_s = time.Now().Unix()
			}

			if update_settings {
				Utils.QueueMessageSERVER(false, Utils.NUM_LIB_SettingsSync, []byte("JSON|false|US"))
				var comms_map map[string]any = <- Utils.LibsCommsChannels_GL[Utils.NUM_LIB_SettingsSync]
				if comms_map == nil {
					return
				}
				map_value, ok := comms_map[Utils.COMMS_MAP_SRV_KEY]
				if !ok {
					continue
				}

				var new_crc16 []byte = map_value.([]byte)
				if !bytes.Equal(new_crc16, last_crc16_GL) {
					last_crc16_GL = new_crc16

					Utils.QueueMessageSERVER(false, Utils.NUM_LIB_SettingsSync, []byte("JSON|true|US"))
					comms_map = <- Utils.LibsCommsChannels_GL[Utils.NUM_LIB_SettingsSync]
					if comms_map == nil {
						return
					}

					var json []byte = []byte(Utils.DecompressString(comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)))

					_ = Utils.FromJsonGENERAL(json, &Utils.User_settings_GL)
				}
			}

			if Utils.WaitWithStopTIMEDATE(&stop_GL, 1) {
				return
			}
		}
	}()
}

/*
StopUserSettingsSyncer stops the user settings synchronizer.
 */
func StopUserSettingsSyncer() {
	stop_GL = true
}

/*
GetUserSettings returns the user settings in JSON format.

-----------------------------------------------------------

– Returns:
  - the user settings in JSON format
 */
func GetJsonUserSettings() string {
	return *Utils.ToJsonGENERAL(Utils.User_settings_GL)
}

/*
LoadUserSettings loads the user settings from the given JSON string.

-----------------------------------------------------------

– Params:
  - json – the JSON string to load the user settings from

– Returns:
  - true if the user settings were successfully loaded, false otherwise
 */
func LoadUserSettings(json string) error {
	if json == "" {
		return errors.New("empty JSON string")
	}

	if err := Utils.FromJsonGENERAL([]byte(json), &Utils.User_settings_GL); err != nil {
		return err
	}

	return nil
}

/*
IsWebsiteInfoEmpty returns true if the website domain and password are empty, false otherwise.

-----------------------------------------------------------

– Returns:
  - true if the website domain and password are empty, false otherwise
 */
func IsWebsiteInfoEmpty() bool {
	return Utils.User_settings_GL.General.Website_domain == "" && Utils.User_settings_GL.General.Website_pw == ""
}

/*
SetWebsiteInfo sets the website domain and password.

-----------------------------------------------------------

– Params:
  - website_domain – the domain of the VISOR website
  - website_password – the password for the VISOR website
 */
func SetWebsiteInfo(website_domain string, website_password string) {
	Utils.User_settings_GL.General.Website_domain = website_domain
	Utils.User_settings_GL.General.Website_pw = website_password
}
