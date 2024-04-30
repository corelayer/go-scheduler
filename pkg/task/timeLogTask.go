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

import (
	"reflect"
	"time"
)

type TimeLogTask struct {
	Timestamp time.Time
	status    Status
}

func (TimeLogTask) Name() string {
	return "timelog"
}

func (t TimeLogTask) Status() Status {
	return t.status
}

func (t TimeLogTask) Type() string {
	return reflect.TypeOf(t).String()
}

func (t TimeLogTask) SetStatus(s Status) Task {
	t.status = s
	return t
}

func (t TimeLogTask) WriteToPipeline() bool {
	return true
}
