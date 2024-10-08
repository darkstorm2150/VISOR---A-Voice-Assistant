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

package GPTComm

import (
	"Utils"
	"strings"
	"time"
)

var time_begin_ms_GL int64 = -1
var curr_entry_time_ms_GL int64 = -1
var last_speech_GL string = ""
var last_idx_begin_GL int = 0

const END_ENTRY string = "[3234_END]"
const ALL_DEVICES_ID string = "3234_ALL"

/*
SetTimeBegin sets the time to begin searching for the next speech.

-----------------------------------------------------------

– Params:
  - time_begin_ms – the time to begin searching for the next speech in milliseconds
 */
func SetTimeBegin(time_begin_ms int64) {
	time_begin_ms_GL = time_begin_ms
}

/*
GetNextSpeechSentence gets the next sentence to be spoken of the most recent speech.

THIS FUNCTION MUST BE IN LOOP BEFORE CALLING SendText()!!!

Each time the function is called, a new sentence is returned, until the end of the text file is reached, in which case
the function will return END_ENTRY.

In case a new speech is added to the text file, the function will continue the speech it was on until its end.

The function will wait until the time of the next speech is reached.

-----------------------------------------------------------

– Returns:
  - the next sentence to be spoken (sometimes may return an empty string - ignore) or END_ENTRY if the end of the text
    file is reached
 */
func GetNextSpeechSentence() string {
	if curr_entry_time_ms_GL == -1 {
		var comms_map map[string]any = <- Utils.LibsCommsChannels_GL[Utils.NUM_LIB_GPTComm]
		if comms_map == nil {
			return END_ENTRY
		}

		var response string = string(comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte))
		if response == "start" {
			var entry *_Entry = getEntry(-1, -1)
			var device_id string = entry.getDeviceID()
			if entry.getTime() >= time_begin_ms_GL && (device_id == Utils.Device_settings_GL.Device_ID || device_id == ALL_DEVICES_ID) {
				curr_entry_time_ms_GL = entry.getTime()
				time_begin_ms_GL = curr_entry_time_ms_GL + 1
				last_speech_GL = ""
			}
		} else if response == "true" || response == "false" {
			gpt_ready_GL = response

			return END_ENTRY
		}
	}

	//log.Println("JJJJJJJJJJJJJJJJJJJJJJJJJJJJJ")
	//log.Println("curr_entry_time_ms_GL:", curr_entry_time_ms_GL)
	//log.Println("time_begin_ms_GL:", time_begin_ms_GL)

	var sentence string = ""
	for {
		var entry *_Entry = getEntry(curr_entry_time_ms_GL, -1)
		var text = entry.getText()

		//log.Println("--------------------------")
		//log.Println("text: \"" + text + "\"")

		text = strings.Replace(text, "\n", ". ", -1)
		text = strings.Replace(text, END_ENTRY, ". " + END_ENTRY, 1)
		text = strings.Replace(text, "...", ".", -1)
		//log.Println("text: \"" + text + "\"")
		if last_idx_begin_GL != 0 && last_idx_begin_GL >= len(text) {
			sentence = ""

			break
		}

		//E se ainda não houver mais texto e isto já tiver tentado ir buscar...? Não vai haver texto, isto vai sair do ciclo e vai retornar END_ENTRY.
		//	Isto tem de esperar até encontrar o 3234_END!

		var dot_idx = strings.Index(text[last_idx_begin_GL:], ". ")
		var dot_idx2 = strings.IndexAny(text[last_idx_begin_GL:], "!?")
		if dot_idx2 != -1 && (dot_idx == -1 || dot_idx2 < dot_idx) {
			dot_idx = dot_idx2
		}

		//log.Println("dot_idx:", dot_idx)
		//log.Println("last_idx_begin_GL:", last_idx_begin_GL)
		//log.Println("text[last_idx_begin_GL:]:", text[last_idx_begin_GL:])

		// If the last dot index is not found, it means that the sentence is not finished yet. So, we must wait for the
		// next entry to be added to the text file.
		if dot_idx != -1 {
			sentence = text[last_idx_begin_GL : last_idx_begin_GL + dot_idx + 2]
			sentence = strings.Trim(sentence, " ")

			last_idx_begin_GL += dot_idx + 2

			break
		}

		if strings.Contains(text[last_idx_begin_GL:], END_ENTRY) {
			sentence = END_ENTRY
			curr_entry_time_ms_GL = -1
			last_idx_begin_GL = 0

			break
		} else {
			time.Sleep(1 * time.Second)
		}
	}

	if !strings.ContainsAny(sentence, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		sentence = ""
	}

	//log.Println("sentence: \"" + sentence + "\"")

	if sentence != "" {
		if last_speech_GL != "" {
			last_speech_GL += " "
		}
		last_speech_GL += sentence
	}

	//log.Println("last_speech_GL: \"" + last_speech_GL + "\"")

	return sentence
}

func GetLastSpeech() string {
	return last_speech_GL
}
