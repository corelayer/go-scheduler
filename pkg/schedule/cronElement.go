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

package schedule

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

func NewCronElement(expression string, position CronElementType) (CronElement, error) {
	var (
		err error
		s   CronElement
	)
	s = CronElement{
		expression:  expression,
		elementType: position,
	}

	err = s.init()
	return s, err
}

type CronElement struct {
	expression     string
	elementType    CronElementType
	expressionType CronExpressionType
}

func (c CronElement) init() error {
	var (
		err error
	)

	err = c.validate()

	if err != nil {
		c.expressionType = CronInvalidExpression
		return err
	}

	c.setExpressionType()
	return nil
}

func (c CronElement) validate() error {
	var s []string
	switch c.elementType {
	case CronSecond:
		s = reSecondOrMinute.FindStringSubmatch(c.expression)
	case CronMinute:
		s = reSecondOrMinute.FindStringSubmatch(c.expression)
	case CronHour:
		s = reHour.FindStringSubmatch(c.expression)
	case CronDay:
		s = reDayOfMonth.FindStringSubmatch(c.expression)
	case CronMonth:
		s = reMonth.FindStringSubmatch(c.expression)
	case CronWeekday:
		s = reDayOfWeek.FindStringSubmatch(c.expression)
	}

	switch len(s) {
	case 0:
		return fmt.Errorf("invalid expression in %s", c.elementType.String())
	case 1:
		return nil
	default:
		return fmt.Errorf("invalid match count in %s", c.elementType.String())
	}
}

func (c CronElement) setExpressionType() {
	chars := []string{"*", ",", "-", "/"}
	if !strings.ContainsAny(c.expression, strings.Join(chars, "")) {
		c.expressionType = CronSimpleExpression
		return
	}
	for i, char := range chars {
		if strings.Contains(c.expression, char) {
			c.expressionType = CronExpressionType(i + 1)
			return
		}
	}
}

func (c CronElement) splitExpression() ([]int, error) {
	var (
		s   []string
		o   []int
		err error
	)

	switch c.expressionType {
	case CronSimpleExpression:
		o = make([]int, 1)
		o[0], err = strconv.Atoi(c.expression)
	case CronMultiExpression:
		s = strings.Split(c.expression, ",")
		o = make([]int, len(s))
		for k, v := range s {
			o[k], err = strconv.Atoi(v)
			if err != nil {
				break
			}
		}
	case CronRangeExpression:
		s = strings.Split(c.expression, "-")
		o = make([]int, 2)
		for k, v := range s {
			o[k], err = strconv.Atoi(v)
			if err != nil {
				break
			}
		}
	case CronStepExpression:
		s = strings.Split(c.expression, "*/")
		o = make([]int, 1)
		o[0], err = strconv.Atoi(s[1])
	}

	return o, err
}

func (c CronElement) IsDue(t time.Time) bool {
	switch c.elementType {
	case CronSecond:
		return c.isExpressionDue(t.Second())
	case CronMinute:
		return c.isExpressionDue(t.Minute())
	case CronHour:
		return c.isExpressionDue(t.Hour())
	case CronDay:
		return c.isExpressionDue(t.Day())
	case CronMonth:
		return c.isExpressionDue(int(t.Month()))
	case CronWeekday:
		return c.isExpressionDue(int(t.Weekday()))
	case CronYear:
		return c.isExpressionDue(t.Year())
	default:
		return false
	}
}

func (c CronElement) isExpressionDue(t int) bool {
	var (
		err       error
		intervals []int
	)
	if c.expressionType == CronSimpleExpression && c.expression == "*" {
		return true
	}

	intervals, err = c.splitExpression()
	if err != nil {
		return false
	}

	switch c.expressionType {
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
