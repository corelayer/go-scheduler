/*
 * Copyright 2023 CoreLayer BV
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

	"github.com/corelayer/go-scheduler/pkg/status"
)

type EmptyTask struct {
	status status.Status
}

func (t EmptyTask) Name() string {
	return "empty"
}

func (t EmptyTask) Status() status.Status {
	return t.status
}

func (t EmptyTask) Type() string {
	return reflect.TypeOf(EmptyTask{}).String()
}

func (t EmptyTask) SetStatus(s status.Status) Task {
	t.status = s
	return t
}

func (t EmptyTask) WriteToPipeline() bool {
	return true
}
