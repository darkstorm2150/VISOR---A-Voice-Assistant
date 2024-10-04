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

package ModsFileInfo

// Mod12GenInfo is the format of the custom generated information about this specific module.
type Mod12GenInfo struct {
	// User_location is the location of the user
	User_location UserLocation
}

type UserLocation struct {
	// Last_known_location is the last known location of the user
	Last_known_location string
	// Curr_location is the current location of the user
	Curr_location string
	// Last_time_checked_s is the last time the current location was checked in Unix time
	Last_time_checked_s int64
	// Prev_location is the previous location of the user
	Prev_location string
	// Prev_last_time_checked_s is the last time the previous location was checked in Unix time
	Prev_last_time_checked_s int64
}

///////////////////////////////////////////////////////////////////////////////

// Mod12UserInfo is the format of the custom information file about this specific module.
type Mod12UserInfo struct {
	// AlwaysWith_device is true if the device is always with the user
	AlwaysWith_device bool
	// Locs_info is the information about the locations
	Locs_info []_LocInfo
}

type _LocInfo struct {
	// Type is the type of the location "detector" (e.g. wifi)
	Type string
	// Name is the name of the detection (e.g. the wifi SSID)
	Name string
	// Address is the address of the detection (e.g. the wifi BSSID) in the format XX:XX:XX:XX:XX:XX
	Address string
	// Last_detection_s is the maximum amount of time in seconds without checking in which the device may still be in the
	// specified location
	Last_detection_s int64
	// Max_distance is the maximum distance in meters in which the device is in the specified location
	Max_distance int
	// Location is where the device is (e.g. "home")
	Location string
}
