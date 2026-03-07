package service

import (
	"enque-learning/constants"
	"fmt"
	"strconv"
	"strings"
)

// ProcessCalc processes a calculator command and returns the response message
func (s *Service) ProcessCalc(arguments []string) string {
	// Validate arguments
	if len(arguments) == 0 {
		return constants.CalcUsageMessage
	}

	expression := strings.Join(arguments, " ")

	// Calculate
	result, err := s.evaluate(expression)
	if err != nil {
		return fmt.Sprintf(constants.CalcErrorTemplate, err.Error())
	}

	return fmt.Sprintf(constants.CalcResultTemplate, expression, result)
}

// evaluate calculates a mathematical expression
func (s *Service) evaluate(expr string) (float64, error) {
	parts := strings.Fields(expr)

	if len(parts) != 3 {
		return 0, fmt.Errorf(constants.CalcInvalidFormat)
	}

	num1, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf(constants.CalcInvalidFirstNumber, parts[0])
	}

	num2, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, fmt.Errorf(constants.CalcInvalidSecondNumber, parts[2])
	}

	operator := parts[1]

	switch operator {
	case "+":
		return num1 + num2, nil
	case "-":
		return num1 - num2, nil
	case "*", "x", "×":
		return num1 * num2, nil
	case "/", "÷":
		if num2 == 0 {
			return 0, fmt.Errorf(constants.CalcDivisionByZero)
		}
		return num1 / num2, nil
	case "^", "**":
		result := 1.0
		for i := 0; i < int(num2); i++ {
			result *= num1
		}
		return result, nil
	default:
		return 0, fmt.Errorf(constants.CalcInvalidOperator, operator)
	}
}
