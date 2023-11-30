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

package cron

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var validSecondOrMinute = `(?:^(?:[1-5]?\d){1}$)|(?:^(?:[1-5]?\d)(?:,(?:[1-5]?\d))+$)|(?:^(?:[1-5]?\d)-(?:[1-5]?\d)$)|(?:^(?:\*/[2-6]|\*/10|\*/12|\*/15|\*/20|\*/30)$)`
var validHour = `(?:^\*$^(?:(?:\d)|(?:1[0-9]{1})|(?:2[0-3]{1}))$)|(?:^(?:(?:\d)|(?:1[0-9]{1})|(?:2[0-3]{1}))(?:,(?:(?:\d)|(?:1[0-9]{1})|(?:2[0-3]{1})))*$)|(?:^(?:(?:\d)|(?:1[0-9]{1})|(?:2[0-3]{1}))-(?:(?:\d)|(?:1[0-9]{1})|(?:2[0-3]{1}))$)|(?:^(?:\*/[2-4]|\*/6|\*/8|\*/12)$)`
var validDayOfWeek = `(?:^(?:(?:[0-6]{1}))$)|(?:^(?:(?:[0-6]{1}))(?:,(?:(?:[0-6]{1})))*$)|(?:^(?:(?:[0-6]{1}))-(?:(?:[0-6]{1}))$)`
var validDayOfMonth = `(?:^(?:(?:[1-2]?[1-9]{1})|(?:3[0-1]{1}))$)|(?:^(?:(?:[1-2]?[1-9]{1})|(?:3[0-1]{1}))(?:,(?:(?:(?:[1-2]?[1-9]{1})|(?:3[0-1]{1}))))*$)|(?:^(?:(?:[1-2]?[1-9]{1})|(?:3[0-1]{1}))-(?:(?:[1-2]?[1-9]{1})|(?:3[0-1]{1}))$)`
var validMonth = `(?:^(?:(?:[1-9])|(?:1[0-2]{1}))$)|(?:^(?:(?:[1-9])|(?:1[0-2]{1}))(?:,(?:(?:[1-9])|(?:1[0-2]{1})))*$)|(?:^(?:(?:[1-9])|(?:1[0-2]{1}))-(?:(?:[1-9])|(?:1[0-2]{1}))$)`
var validYear = `(?:^\d+$)|(?:^(?:\d+)(?:,(?:\d+)+)$)|(?:^\d+-\d+$)|(?:^\*\/\d+$)`

var reSecondOrMinute = regexp.MustCompile(validSecondOrMinute)
var reHour = regexp.MustCompile(validHour)
var reMonth = regexp.MustCompile(validMonth)
var reDayOfWeek = regexp.MustCompile(validDayOfWeek)
var reDayOfMonth = regexp.MustCompile(validDayOfMonth)
var reYear = regexp.MustCompile(validYear)

func newElement(expression string, position position) (element, error) {
	e := element{
		expression: expression,
		p:          position,
		q:          qualificationNone,
	}

	return e, e.validate()
}

type element struct {
	expression string
	p          position
	q          qualification
}

func (e *element) validate() error {
	// First qualify the expression
	e.qualify()

	// No additional validation needed for expression
	if e.expression == "*" {
		return nil
	}

	var s []string
	switch e.p {
	case positionSecond:
		s = reSecondOrMinute.FindStringSubmatch(e.expression)
	case positionMinute:
		s = reSecondOrMinute.FindStringSubmatch(e.expression)
	case positionHour:
		s = reHour.FindStringSubmatch(e.expression)
	case positionDay:
		s = reDayOfMonth.FindStringSubmatch(e.expression)
	case positionMonth:
		s = reMonth.FindStringSubmatch(e.expression)
	case positionWeekday:
		s = reDayOfWeek.FindStringSubmatch(e.expression)
	case positionYear:
		s = reYear.FindStringSubmatch(e.expression)
	}

	if len(s) == 0 {
		return fmt.Errorf("invalid expression in %s", e.p.String())
	}

	// No additional validation is necessary for simple expressions
	if e.q == qualificationSimple {
		return nil
	}
	// No need to check for errors, as the regex will only allow integers to match
	p, _ := e.parseExpression()

	// If there are multiple values, check if they are ascending
	if len(p) > 1 {
		for i := 0; i < len(p)-1; i++ {
			if p[i] > p[i+1] {
				return fmt.Errorf("invalid order of values in %s", e.p.String())
			}
		}
	}
	return nil
}

func (e *element) qualify() {
	chars := []string{",", "-", "/"}
	if !strings.ContainsAny(e.expression, strings.Join(chars, "")) {
		e.q = qualificationSimple
		return
	}
	for i, char := range chars {
		if strings.Contains(e.expression, char) {
			// qualification starts at 2 = qualificationSimple
			e.q = qualification(i + 2)
		}
	}
}

func (e *element) parseExpression() ([]int, error) {
	var (
		s   []string
		o   []int
		err error
	)

	if e.expression == "*" {
		return nil, nil
	}

	switch e.q {
	case qualificationSimple:
		o = make([]int, 1)
		o[0], err = strconv.Atoi(e.expression)
	case qualificationMulti:
		s = strings.Split(e.expression, ",")
		o = make([]int, len(s))
		for k, v := range s {
			o[k], err = strconv.Atoi(v)
			if err != nil {
				break
			}
		}
	case qualificationRange:
		s = strings.Split(e.expression, "-")
		o = make([]int, 2)
		for k, v := range s {
			o[k], err = strconv.Atoi(v)
			if err != nil {
				break
			}
		}
	case qualificationStep:
		s = strings.Split(e.expression, "*/")
		o = make([]int, 1)
		o[0], err = strconv.Atoi(s[1])
	default:
		err = fmt.Errorf("invalid qualification to parse")
	}

	return o, err
}

func (e *element) Trigger(t time.Time) bool {
	if e.expression == "*" {
		return true
	}

	intervals, err := e.parseExpression()
	if err != nil {
		return false
	}

	var input int
	switch e.p {
	case positionSecond:
		input = t.Second()
	case positionMinute:
		input = t.Minute()
	case positionHour:
		input = t.Hour()
	case positionDay:
		input = t.Day()
	case positionMonth:
		input = int(t.Month())
	case positionWeekday:
		input = int(t.Weekday())
	case positionYear:
		input = t.Year()
	}
	return e.isDue(intervals, input)
}

func (e *element) isDue(inputs []int, t int) bool {
	if inputs == nil || len(inputs) == 0 {
		return false
	}

	switch e.q {
	case qualificationSimple:
		if inputs[0] == t {
			return true
		}
		return false
	case qualificationMulti:
		for _, i := range inputs {
			if i == t {
				return true
			}
		}
		return false
	case qualificationRange:
		if t >= inputs[0] && t <= inputs[1] {
			return true
		}
		return false
	case qualificationStep:
		if t%inputs[0] == 0 {
			return true
		}
		return false
	default:
		return false
	}
}
