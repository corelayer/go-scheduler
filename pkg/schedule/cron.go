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

var reSpace = regexp.MustCompile(`\s+`)
var reYear = regexp.MustCompile(`\d{4}`)

var reSecond = regexp.MustCompile(`^([1-5]?\d)(,([1-5]?\d))*$|^\*/[2-6]$|^\*/10$|^\*/12$|^\*/15$|^\*/20$|^\*/30$`)
var reMinute = regexp.MustCompile(`^([1-5]?\d)(,([1-5]?\d))*$|^\*/[2-6]$|^\*/10$|^\*/12$|^\*/15$|^\*/20$|^\*/30$`)

var cronWeekdayLiterals = strings.NewReplacer(
	"SUN", "0",
	"MON", "1",
	"TUE", "2",
	"WED", "3",
	"THU", "4",
	"FRI", "5",
	"SAT", "6")

var cronMonthLiterals = strings.NewReplacer(
	"JAN", "1",
	"FEB", "2",
	"MAR", "3",
	"APR", "4",
	"MAY", "5",
	"JUN", "6",
	"JUL", "7",
	"AUG", "8",
	"SEP", "9",
	"OCT", "10",
	"NOV", "11",
	"DEC", "12")

var cronTemplates = map[string]string{
	"@yearly":      "0 0 1 1 *",
	"@annually":    "0 0 1 1 *",
	"@monthly":     "0 0 1 * *",
	"@weekly":      "0 0 * * 0",
	"@daily":       "0 0 * * *",
	"@hourly":      "0 * * * *",
	"@always":      "* * * * *",
	"@5minutes":    "*/5 * * * *",
	"@10minutes":   "*/10 * * * *",
	"@15minutes":   "*/15 * * * *",
	"@30minutes":   "0,30 * * * *",
	"@everysecond": "* * * * * *",
}

type CronPosition int

const (
	CronSecond CronPosition = iota
	CronMinute
	CronHour
	CronDay
	CronMonth
	CronWeekday
	CronYear
)

type CronExpression int

const (
	CronEveryExpression CronExpression = iota
	CronSimpleExpression
	CronMultiExpression
	CronRangeExpression
	CronStepExpression
)

func NewCronSegment(expression string, position CronPosition) (CronSegment, error) {
	var (
		err error
		s   CronSegment
	)
	s = CronSegment{
		expression: expression,
		position:   position,
	}

	err = s.init()
	return s, err
}

type CronSegment struct {
	expression     string
	position       CronPosition
	expressionType CronExpression
}

func (c CronSegment) init() error {
	var (
		err error
	)

	err = c.setExpressionType()

	return err
}

func (c CronSegment) setExpressionType() error {
	if strings.Contains(c.expression, "*") {
		c.expressionType = CronEveryExpression
		return nil
	}

	if strings.Contains(c.expression, ",") {
		c.expressionType = CronMultiExpression
		return nil
	}

	if strings.Contains(c.expression, "-") {
		c.expressionType = CronRangeExpression
		return nil
	}

	if strings.Contains(c.expression, "/") {
		c.expressionType = CronStepExpression
		return nil
	}

	_, err := strconv.Atoi(c.expression)
	if err != nil {
		return err
	}

	c.expressionType = CronSimpleExpression
	return nil
}

func (c CronSegment) IsDue(t time.Time) bool {
	switch c.position {
	case CronSecond:
		return true
	case CronMinute:
		return c.minuteDue(t)
	case CronHour:
		return true
	case CronDay:
		return true
	case CronMonth:
		return true
	case CronWeekday:
		return true
	case CronYear:
		return true
	default:
		return false
	}
}

func (c CronSegment) minuteDue(t time.Time) bool {
	var (
		err   error
		value int
	)
	minute := t.Minute()
	value, err = strconv.Atoi(c.expression)
	if err != nil {
		value = 0
	}
	if minute == value {
		return true
	}
	return true
}

func NewCron(s string) (Cron, error) {
	c := Cron{
		expression: s,
		segments:   make([]CronSegment, 6),
	}
	c.replaceCronTemplates()
	c.normalize()
	err := c.parse()
	if err != nil {
		return Cron{}, err
	}
	return c, nil
}

type Cron struct {
	expression string
	segments   []CronSegment
}

func (c Cron) replaceCronTemplates() {
	// Replace template expression with a valid cron expression
	if t, ok := cronTemplates[c.expression]; ok {
		c.expression = t
	}
}

func (c Cron) normalize() {
	// Replace all spaces with a single space
	c.expression = reSpace.ReplaceAllString(c.expression, " ")

	// Transform string to uppercase before replacing characters to numbers
	c.expression = strings.ToUpper(c.expression)

	// Replace weekday literals to numbers
	c.expression = cronWeekdayLiterals.Replace(c.expression)
	// Replace month literals to numbers
	c.expression = cronMonthLiterals.Replace(c.expression)
}

func (c Cron) standardize() ([]string, error) {
	segments := strings.Split(c.expression, " ")
	count := len(segments)

	// Expect at least 5 segments: minute, hour, day, month, weekday
	// Maximum 7 segments: second, minute, hour ,day, month, weekday, year
	if count < 5 || count > 7 {
		return nil, fmt.Errorf("invalid segment count, got %d, expected 5-7 segments separated by space", count)
	}

	if (count == 5) || (count == 6 && reYear.MatchString(segments[5])) {
		segments = append([]string{"0"}, segments...)
	}

	return segments, nil
}

func (c Cron) parse() error {
	var (
		segments []string
		err      error
	)

	segments, err = c.standardize()
	if err != nil {
		return err
	}

	if len(segments) == 7 {
		c.segments = make([]CronSegment, 7)
	}

	for position, expression := range segments {
		var s CronSegment
		s, err = NewCronSegment(expression, CronPosition(position))
		if err != nil {
			break
		}
		c.segments[position] = s
	}
	return err
}
