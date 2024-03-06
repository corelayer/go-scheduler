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

package job

import "fmt"

type Error struct {
	kind errKind
	err  error
}

func (e Error) Error() string {
	return e.kind.String()
}

type errKind int

func (k errKind) String() string {
	return fmt.Sprintf("%d", k)
}

const (
	ErrKindNotFound errKind = iota
	ErrKindExist
)

var (
	ErrNotFound = Error{kind: ErrKindNotFound}
	ErrExist    = Error{kind: ErrKindExist}
)
