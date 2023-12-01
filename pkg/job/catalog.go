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

import "github.com/google/uuid"

type Catalog interface {
	CatalogReader
	CatalogWriter
}

type CatalogReader interface {
	All() []Job
	Schedulable(limit int) []Job
}

type CatalogWriter interface {
	Activate(uuid uuid.UUID)
	Add(job Job)
	Update(job Job)
	Delete(uuid uuid.UUID)
}