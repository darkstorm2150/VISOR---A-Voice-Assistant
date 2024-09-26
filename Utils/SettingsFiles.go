/*******************************************************************************
 * Copyright 2023-2024 Edw590
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

package Utils

import (
	"Utils/ModsFileInfo"
	"errors"
	"os"
	"strings"
)

const USER_SETTINGS_FILE string = "UserSettings_EOG.json"
const GEN_SETTINGS_FILE string = "GeneratedSettings_EOG.json"

var User_settings_GL UserSettings
var Gen_settings_GL GenSettings

type UserSettings struct {
	PersonalConsts _PersonalConsts
	MOD_2  ModsFileInfo.Mod2UserInfo
	MOD_4  ModsFileInfo.Mod4UserInfo
	MOD_6  ModsFileInfo.Mod6UserInfo
	MOD_7  ModsFileInfo.Mod7UserInfo
	MOD_10 ModsFileInfo.Mod10UserInfo
	MOD_12 ModsFileInfo.Mod12UserInfo
}

type GenSettings struct {
	MOD_2 ModsFileInfo.Mod2GenInfo
	MOD_5 ModsFileInfo.Mod5GenInfo
	MOD_9 ModsFileInfo.Mod9GenInfo
}

///////////////////////////////////////////////////////////////

type _PersonalConsts struct {
	// Device_ID is the device ID of the current device
	Device_ID string

	// VISOR_dir is the full path to the main directory of VISOR.
	VISOR_dir string
	// VISOR_server is an INTERNAL attribute to be filled INTERNALLY that indicates if the version running is the server
	// or the client version
	VISOR_server bool

	// VISOR_email_addr is VISOR's email address
	VISOR_email_addr string
	// VISOR_email_pw is VISOR's email password
	VISOR_email_pw string

	// User_email_addr is the email address of the user, used for all email communication
	User_email_addr string

	// Website_url is the URL of the VISOR website
	Website_url string
	// Website_pw is the password for the VISOR website
	Website_pw string
	// Website_dir is the full path to the directory of the VISOR website
	Website_dir string

	// WolframAlpha_AppID is the app ID for the Wolfram Alpha API
	WolframAlpha_AppID string

	// Picovoice_API_key is the API key for the Picovoice API
	Picovoice_API_key string
}

/*
LoadUserSettings is the function that initializes the global variables of the UserSettings struct.
*/
func LoadUserSettings(server bool) error {
	bytes, err := os.ReadFile(USER_SETTINGS_FILE)
	if err != nil {
		cwd, err := os.Getwd()
		if err != nil {
			cwd = "[ERROR]"
		}
		return errors.New("no " + USER_SETTINGS_FILE + " file found in the current working directory: \"" + cwd + "\" - aborting")
	}

	if err := FromJsonGENERAL(bytes, &User_settings_GL); err != nil {
		return err
	}

	User_settings_GL.PersonalConsts.VISOR_server = server

	if User_settings_GL.PersonalConsts.VISOR_server {
		if !strings.Contains(User_settings_GL.PersonalConsts.VISOR_email_addr, "@") || User_settings_GL.PersonalConsts.Device_ID == "" ||
			User_settings_GL.PersonalConsts.VISOR_email_pw == "" || !strings.Contains(User_settings_GL.PersonalConsts.User_email_addr, "@") ||
			!strings.Contains(User_settings_GL.PersonalConsts.Website_url, "http") || User_settings_GL.PersonalConsts.Website_pw == "" ||
			User_settings_GL.PersonalConsts.WolframAlpha_AppID == "" || User_settings_GL.PersonalConsts.Picovoice_API_key == "" {
			return errors.New("some fields in " + USER_SETTINGS_FILE + " are empty or incorrect - aborting")
		}
	} else {
		if !strings.Contains(User_settings_GL.PersonalConsts.User_email_addr, "@") ||
				!strings.Contains(User_settings_GL.PersonalConsts.Website_url, "http") ||
				User_settings_GL.PersonalConsts.Website_pw == "" {
			return errors.New("some fields in " + USER_SETTINGS_FILE + " are empty or incorrect - aborting")
		}
	}

	var visor_path GPath = PathFILESDIRS(true, "", User_settings_GL.PersonalConsts.VISOR_dir)
	if !visor_path.Exists() {
		return errors.New("the VISOR directory \"" + visor_path.GPathToStringConversion() + "\" does not exist - aborting")
	}
	if User_settings_GL.PersonalConsts.VISOR_server {
		var website_path GPath = PathFILESDIRS(true, "", User_settings_GL.PersonalConsts.Website_dir)
		if !website_path.Exists() {
			return errors.New("the website directory \"" + website_path.GPathToStringConversion() + "\" does not exist - aborting")
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////

/*
loadGenSettings is the function that initializes the global variables of the GenSettings struct.
*/
func loadGenSettings() error {
	bytes, err := os.ReadFile(GEN_SETTINGS_FILE)
	if err != nil {
		cwd, err := os.Getwd()
		if err != nil {
			cwd = "[ERROR]"
		}
		return errors.New("no " + GEN_SETTINGS_FILE + " file found in the current working directory: \"" + cwd + "\" - aborting")
	}

	if err := FromJsonGENERAL(bytes, &Gen_settings_GL); err != nil {
		return err
	}

	return nil
}

/*
saveGenSettings is the function that saves the global variables of the GenSettings struct to the GEN_SETTINGS_FILE file.
 */
func saveGenSettings() bool {
	var p_string *string = ToJsonGENERAL(Gen_settings_GL)
	if p_string == nil {
		return false
	}

	if err := os.WriteFile(GEN_SETTINGS_FILE, []byte(*p_string), 0777); err != nil {
		return false
	}

	return true
}