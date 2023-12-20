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

package job

import "testing"

func TestJob_IsPending1(t *testing.T) {
	j := Job{
		Enabled: false,
		Status:  StatusNone,
	}

	if j.IsPending() {
		t.Errorf("job is pending, expected %s", StatusNone)
	}
}

func TestJob_IsPending2(t *testing.T) {
	j := Job{
		Enabled: true,
		Status:  StatusPending,
	}

	if !j.IsPending() {
		t.Errorf("job is %s, expected %s", j.Status, StatusPending)
	}
}
