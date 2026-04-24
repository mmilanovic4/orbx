package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type UnitGroup struct {
	Base  string
	Units map[string]float64
	Order []string
}

var groups = map[string]UnitGroup{
	"length": {
		Base:  "m",
		Order: []string{"in", "cm", "ft", "m", "km", "mi"},
		Units: map[string]float64{
			"in": 0.0254,
			"cm": 0.01,
			"ft": 0.3048,
			"m":  1,
			"km": 1000,
			"mi": 1609.344,
		},
	},
	"weight": {
		Base:  "kg",
		Order: []string{"g", "oz", "lb", "kg"},
		Units: map[string]float64{
			"g":  0.001,
			"oz": 0.0283495,
			"lb": 0.453592,
			"kg": 1,
		},
	},
	"storage": {
		Base:  "b",
		Order: []string{"b", "kb", "mb", "gb", "tb"},
		Units: map[string]float64{
			"b":  1,
			"kb": 1024,
			"mb": 1024 * 1024,
			"gb": 1024 * 1024 * 1024,
			"tb": 1024 * 1024 * 1024 * 1024,
		},
	},
	"time": {
		Base:  "s",
		Order: []string{"ms", "s", "min", "h", "day", "week"},
		Units: map[string]float64{
			"ms":   0.001,
			"s":    1,
			"min":  60,
			"h":    3600,
			"day":  86400,
			"week": 604800,
		},
	},
}

var tempUnits = map[string]bool{"c": true, "f": true, "k": true}
var tempOrder = []string{"c", "f", "k"}

// Ordered by decreasing length to avoid prefix collisions during unit matching (e.g. "mb" matching "m")
var allUnits = []string{"week", "day", "min", "tb", "gb", "mb", "kb", "in", "cm", "ft", "km", "mi", "ms", "oz", "lb", "kg", "b", "m", "g", "h", "s", "c", "f", "k"}

func convertTemp(value float64, from, to string) (float64, bool) {
	if !tempUnits[from] || !tempUnits[to] {
		return 0, false
	}
	if from == to {
		return value, true
	}

	var celsius float64
	switch from {
	case "c":
		celsius = value
	case "f":
		celsius = (value - 32) * 5 / 9
	case "k":
		celsius = value - 273.15
	}

	switch to {
	case "c":
		return celsius, true
	case "f":
		return celsius*9/5 + 32, true
	case "k":
		return celsius + 273.15, true
	}

	return 0, false
}

func findGroup(unit string) (UnitGroup, bool) {
	for _, group := range groups {
		if _, ok := group.Units[unit]; ok {
			return group, true
		}
	}
	return UnitGroup{}, false
}

func convert(value float64, from, to string) (float64, bool) {
	fromGroup, ok1 := findGroup(from)
	toGroup, ok2 := findGroup(to)

	if !ok1 || !ok2 {
		return 0, false
	}
	if fromGroup.Base != toGroup.Base {
		return 0, false
	}

	base := value * fromGroup.Units[from]
	return base / toGroup.Units[to], true
}

func parseInput(input string) (float64, string, error) {
	input = strings.ToLower(strings.TrimSpace(input))

	for _, unit := range allUnits {
		before, ok := strings.CutSuffix(input, unit)
		if !ok || before == "" {
			continue
		}
		val, err := strconv.ParseFloat(before, 64)
		if err != nil {
			continue
		}
		return val, unit, nil
	}

	return 0, "", fmt.Errorf("invalid format (e.g. 180cm, 6ft, 70kg, 37c, 1gb, 90min)")
}

var convertCmd = &cobra.Command{
	Use:   "convert [value][unit]",
	Short: "Convert units: length, weight, temperature, storage, time",
	Long: fmt.Sprintf("Supported units: %s.",
		strings.Join(allUnits, ", "),
	),
	GroupID: "util",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		val, from, err := parseInput(args[0])
		if err != nil {
			return err
		}

		if tempUnits[from] {
			fmt.Printf("%.2f%s", val, from)
			for _, unit := range tempOrder {
				if unit == from {
					continue
				}
				result, _ := convertTemp(val, from, unit)
				fmt.Printf(" = %.2f%s", result, unit)
			}
			fmt.Println()
			return nil
		}

		group, ok := findGroup(from)
		if !ok {
			return fmt.Errorf("unknown unit: %s", from)
		}

		fmt.Printf("%.4g%s", val, from)
		for _, unit := range group.Order {
			if unit == from {
				continue
			}
			result, _ := convert(val, from, unit)
			fmt.Printf(" = %.4g%s", result, unit)
		}
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
