/*
 * Copyright 2024 CoreLayer BV
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package task

import "log/slog"

type Message struct {
	Message string
	Task    string
	Type    MessageType
	Data    interface{}
}

type LogData struct {
	Level slog.Level
	Attrs []slog.Attr
}

func NewLogMessage(message string, t Task, data interface{}) Message {
	return Message{
		Message: message,
		Task:    t.Type(),
		Type:    LogMessage,
		Data:    data,
	}
}

func NewErrorMessage(message string, t Task, err error) Message {
	return Message{
		Message: message,
		Task:    t.Type(),
		Type:    ErrorMessage,
		Data:    err,
	}
}
