package main

import (
	"fmt"
	"strings"
)

type constraint struct {
	Name   string
	IsSet  func(*schema) bool
	String func(*schema) string
}

var (
	constraintsTable = []constraint{
		{
			Name: "Enum",
			IsSet: func(sch *schema) bool {
				return len(sch.Enum) > 0
			},
			String: func(sch *schema) string {
				var vals []string
				for _, e := range sch.Enum {
					vals = append(vals, e.String())
				}

				return fmt.Sprintf("[`%s`]", strings.Join(vals, "`, `"))
			},
		},
		{
			Name: "MinLength",
			IsSet: func(sch *schema) bool {
				return sch.MinLength != nil
			},
			String: func(sch *schema) string {
				return fmt.Sprintf("%d", *sch.MinLength)
			},
		},
		{
			Name: "MaxLength",
			IsSet: func(sch *schema) bool {
				return sch.MaxLength != nil
			},
			String: func(sch *schema) string {
				return fmt.Sprintf("%d", *sch.MaxLength)
			},
		},
		{
			Name: "Minimum",
			IsSet: func(sch *schema) bool {
				return sch.Minimum != nil
			},
			String: func(sch *schema) string {
				return fmt.Sprintf("%d", *sch.Minimum)
			},
		},
		{
			Name: "Maximum",
			IsSet: func(sch *schema) bool {
				return sch.Maximum != nil
			},
			String: func(sch *schema) string {
				return fmt.Sprintf("%d", *sch.Maximum)
			},
		},
		{
			Name: "MinItems",
			IsSet: func(sch *schema) bool {
				return sch.MinItems != nil
			},
			String: func(sch *schema) string {
				return fmt.Sprintf("%d", *sch.MinItems)
			},
		},
		{
			Name: "MaxItems",
			IsSet: func(sch *schema) bool {
				return sch.MaxItems != nil
			},
			String: func(sch *schema) string {
				return fmt.Sprintf("%d", *sch.MaxItems)
			},
		},
		{
			Name: "MultipleOf",
			IsSet: func(sch *schema) bool {
				return sch.MultipleOf != nil
			},
			String: func(sch *schema) string {
				return fmt.Sprintf("%d", *sch.MultipleOf)
			},
		},
		{
			Name: "Pattern",
			IsSet: func(sch *schema) bool {
				return sch.Pattern != ""
			},
			String: func(sch *schema) string {
				return fmt.Sprintf("`%s`", sch.Pattern)
			},
		},
	}
)

func printConstraints(sch *schema) string {

	var constraints []string

	for _, constraint := range constraintsTable {
		if constraint.IsSet(sch) {
			constraints = append(constraints, fmt.Sprintf("%s: %s", constraint.Name, constraint.String(sch)))
		}
	}

	return strings.Join(constraints, ", ")
}
