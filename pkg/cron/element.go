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
var validHour = `(?:^(?:(?:\d)|(?:1[0-9]{1})|(?:2[0-3]{1}))$)|(?:^(?:(?:\d)|(?:1[0-9]{1})|(?:2[0-3]{1}))(?:,(?:(?:\d)|(?:1[0-9]{1})|(?:2[0-3]{1})))*$)|(?:^(?:(?:\d)|(?:1[0-9]{1})|(?:2[0-3]{1}))-(?:(?:\d)|(?:1[0-9]{1})|(?:2[0-3]{1}))$)|(?:^(?:\*/[2-4]|\*/6|\*/8|\*/12)$)`
var validDayOfWeek = `(?:^(?:(?:[0-6]{1}))$)|(?:^(?:(?:[0-6]{1}))(?:,(?:(?:[0-6]{1})))*$)|(?:^(?:(?:[0-6]{1}))-(?:(?:[0-6]{1}))$)`
var validDayOfMonth = `(?:^(?:(?:[1-2]?[1-9]{1})|(?:3[0-1]{1}))$)|(?:^(?:(?:[1-2]?[1-9]{1})|(?:3[0-1]{1}))(?:,(?:(?:(?:[1-2]?[1-9]{1})|(?:3[0-1]{1}))))*$)|(?:^(?:(?:[1-2]?[1-9]{1})|(?:3[0-1]{1}))-(?:(?:[1-2]?[1-9]{1})|(?:3[0-1]{1}))$)`
var validMonth = `(?:^(?:(?:\d)|(?:1[0-2]{1}))$)|(?:^(?:(?:\d)|(?:1[0-2]{1}))(?:,(?:(?:\d)|(?:1[0-2]{1})))*$)|(?:^(?:(?:\d)|(?:1[0-2]{1}))-(?:(?:\d)|(?:1[0-2]{1}))$)`

var reSecondOrMinute = regexp.MustCompile(validSecondOrMinute)
var reHour = regexp.MustCompile(validHour)
var reMonth = regexp.MustCompile(validMonth)
var reDayOfWeek = regexp.MustCompile(validDayOfWeek)
var reDayOfMonth = regexp.MustCompile(validDayOfMonth)

func NewElement(expression string, position ElementType) (Element, error) {
	var (
		err error
		s   Element
	)
	s = Element{
		expression:  expression,
		elementType: position,
	}

	err = s.init()
	return s, err
}

type Element struct {
	expression     string
	elementType    ElementType
	expressionType ExpressionType
}

func (e Element) init() error {
	var (
		err error
	)

	err = e.validate()

	if err != nil {
		e.expressionType = CronInvalidExpression
		return err
	}

	e.setExpressionType()
	return nil
}

func (e Element) validate() error {
	var s []string
	switch e.elementType {
	case CronSecond:
		s = reSecondOrMinute.FindStringSubmatch(e.expression)
	case CronMinute:
		s = reSecondOrMinute.FindStringSubmatch(e.expression)
	case CronHour:
		s = reHour.FindStringSubmatch(e.expression)
	case CronDay:
		s = reDayOfMonth.FindStringSubmatch(e.expression)
	case CronMonth:
		s = reMonth.FindStringSubmatch(e.expression)
	case CronWeekday:
		s = reDayOfWeek.FindStringSubmatch(e.expression)
	}

	switch len(s) {
	case 0:
		return fmt.Errorf("invalid expression in %s", e.elementType.String())
	case 1:
		return nil
	default:
		return fmt.Errorf("invalid match count in %s", e.elementType.String())
	}
}

func (e Element) setExpressionType() {
	chars := []string{"*", ",", "-", "/"}
	if !strings.ContainsAny(e.expression, strings.Join(chars, "")) {
		e.expressionType = CronSimpleExpression
		return
	}
	for i, char := range chars {
		if strings.Contains(e.expression, char) {
			e.expressionType = ExpressionType(i + 1)
			return
		}
	}
}

func (e Element) splitExpression() ([]int, error) {
	var (
		s   []string
		o   []int
		err error
	)

	switch e.expressionType {
	case CronSimpleExpression:
		o = make([]int, 1)
		o[0], err = strconv.Atoi(e.expression)
	case CronMultiExpression:
		s = strings.Split(e.expression, ",")
		o = make([]int, len(s))
		for k, v := range s {
			o[k], err = strconv.Atoi(v)
			if err != nil {
				break
			}
		}
	case CronRangeExpression:
		s = strings.Split(e.expression, "-")
		o = make([]int, 2)
		for k, v := range s {
			o[k], err = strconv.Atoi(v)
			if err != nil {
				break
			}
		}
	case CronStepExpression:
		s = strings.Split(e.expression, "*/")
		o = make([]int, 1)
		o[0], err = strconv.Atoi(s[1])
	}

	return o, err
}

func (e Element) IsDue(t time.Time) bool {
	switch e.elementType {
	case CronSecond:
		return e.isExpressionDue(t.Second())
	case CronMinute:
		return e.isExpressionDue(t.Minute())
	case CronHour:
		return e.isExpressionDue(t.Hour())
	case CronDay:
		return e.isExpressionDue(t.Day())
	case CronMonth:
		return e.isExpressionDue(int(t.Month()))
	case CronWeekday:
		return e.isExpressionDue(int(t.Weekday()))
	case CronYear:
		return e.isExpressionDue(t.Year())
	default:
		return false
	}
}

func (e Element) isExpressionDue(t int) bool {
	var (
		err       error
		intervals []int
	)
	if e.expressionType == CronSimpleExpression && e.expression == "*" {
		return true
	}

	intervals, err = e.splitExpression()
	if err != nil {
		return false
	}

	switch e.expressionType {
	case CronSimpleExpression:
		if intervals[0] == t {
			return true
		}
	case CronMultiExpression:
		for _, i := range intervals {
			if i == t {
				return true
			}
		}
		return false
	case CronRangeExpression:
		if t >= intervals[0] && t <= intervals[1] {
			return true
		}
	}
	return false
}
