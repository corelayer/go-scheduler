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
)

type IntercomMessageTask struct {
	Message string
	status  Status
}

func (t IntercomMessageTask) Name() string {
	return "intercom"
}

func (t IntercomMessageTask) Status() Status {
	return t.status
}

func (t IntercomMessageTask) Type() string {
	return reflect.TypeOf(t).String()
}

func (t IntercomMessageTask) SetStatus(s Status) Task {
	t.status = s
	return t
}

func (t IntercomMessageTask) WriteToPipeline() bool {
	return true
}
