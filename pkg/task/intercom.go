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
	"sync"
)

type IntercomMessage struct {
	Name    string
	Content Message
	Log     *LogData
}

func NewIntercom(name string, chOut chan IntercomMessage) *Intercom {
	return &Intercom{
		name:     name,
		chOut:    chOut,
		messages: make([]Message, 0),
		mux:      &sync.Mutex{},
	}
}

type Intercom struct {
	name     string
	chOut    chan IntercomMessage
	messages []Message
	mux      *sync.Mutex
}

func (c *Intercom) Add(m Message, l *LogData) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.messages = append(c.messages, m)
	c.chOut <- IntercomMessage{
		Name:    c.name,
		Content: m,
		Log:     l,
	}
}

func (c *Intercom) CountErrorMessages() int {
	return len(c.Get(ErrorMessage))
}

func (c *Intercom) Get(t MessageType) []Message {
	c.mux.Lock()
	defer c.mux.Unlock()

	var messages = make([]Message, 0)
	for _, m := range c.messages {
		if m.Type == t {
			messages = append(messages, m)
		}
	}
	return messages
}

func (c *Intercom) GetAll() []Message {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.messages
}

func (c *Intercom) GetErrors() []error {
	// DO NOT LOCK HERE, locking is done in c.Get()
	var errors = make([]error, 0)
	for _, m := range c.Get(ErrorMessage) {
		errors = append(errors, m.Data.(error))
	}
	return errors
}

func (c *Intercom) HasErrors() bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	for _, m := range c.messages {
		if m.Type == ErrorMessage {
			return true
		}
	}
	return false
}

func (c *Intercom) Reset() {
	c.mux.Lock()
	defer c.mux.Unlock()

	clear(c.messages)
}
