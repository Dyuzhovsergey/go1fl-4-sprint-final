package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("invalid data format")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("steps less than 0")
	}
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("duration less than 0")
	}
	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
	lenghtStep := float64(height) * stepLengthCoefficient
	return (float64(steps) * lenghtStep) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)
	durationHour := duration.Hours()
	return dist / durationHour
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)

	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calories float64

	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration.Hours(), dist, speed, calories), nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	switch {
	case steps <= 0:
		return 0, errors.New("invalid input: steps must be > 0")
	case weight <= 0:
		return 0, errors.New("invalid input: weight must be > 0")
	case height <= 0:
		return 0, errors.New("invalid input: height must be > 0")
	case duration <= 0:
		return 0, errors.New("invalid input: duration must be > 0")
	}

	speed := meanSpeed(steps, height, duration)
	durationMin := duration.Minutes()
	calories := (weight * speed * durationMin) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	switch {
	case steps <= 0:
		return 0, errors.New("invalid input: steps must be > 0")
	case weight <= 0:
		return 0, errors.New("invalid input: weight must be > 0")
	case height <= 0:
		return 0, errors.New("invalid input: height must be > 0")
	case duration <= 0:
		return 0, errors.New("invalid input: duration must be > 0")
	}

	speed := meanSpeed(steps, height, duration)
	durationMin := duration.Minutes()
	calories := (weight * speed * durationMin) / minInH
	return calories * walkingCaloriesCoefficient, nil
}
