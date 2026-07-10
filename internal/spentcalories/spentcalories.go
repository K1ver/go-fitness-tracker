package spentcalories

import (
	"errors"
	"fmt"
	"log"
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

var (
	// ErrInvalidParse ошибка при преобразовании
	ErrInvalidParse = errors.New("wrong argument for parsing")
	// ErrWrongArgument не правильный аргумент
	ErrWrongArgument = errors.New("wrong argument")
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", time.Duration(0), errors.New("wrong count of arguments")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", time.Duration(0), fmt.Errorf("%w: %w", ErrInvalidParse, err)
	}
	if steps <= 0 {
		return 0, "", time.Duration(0), errors.New("steps must be greater then 0")
	}

	typeOfActivity := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", time.Duration(0), fmt.Errorf("%w: %w", ErrInvalidParse, err)
	}
	if duration <= time.Duration(0) {
		return 0, "", time.Duration(0), errors.New("duration must be greater then 0")
	}

	return steps, typeOfActivity, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceKm := (float64(steps) * stepLength) / mInKm
	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= time.Duration(0) {
		return 0
	}
	distanceKm := distance(steps, height)
	return distanceKm / duration.Hours()
}

// TrainingInfo формирует информацию по тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, typeOfActivity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var spentCalories float64
	switch typeOfActivity {
	case "Ходьба":
		spentCalories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		spentCalories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	if err != nil {
		log.Println(err)
		return "", err
	}
	return fmt.Sprintf("Тип тренировки: %s\n"+
		"Длительность: %.2f ч.\n"+
		"Дистанция: %.2f км.\n"+
		"Скорость: %.2f км/ч\n"+
		"Сожгли калорий: %.2f\n", typeOfActivity, float64(duration.Minutes())/60.0, distance(steps, height), meanSpeed(steps, height, duration), spentCalories), nil
}

// RunningSpentCalories функция подсчета калорий при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("%w: steps must be greater than zero", ErrWrongArgument)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("%w: weight must be greater than zero", ErrWrongArgument)
	}
	if height <= 0 {
		return 0, fmt.Errorf("%w: height must be greater than zero", ErrWrongArgument)
	}
	if duration <= time.Duration(0) {
		return 0, fmt.Errorf("%w: duration must be greater than zero", ErrWrongArgument)
	}

	avgSpeed := meanSpeed(steps, height, duration)

	spentCalories := (weight * avgSpeed * duration.Minutes()) / float64(minInH)
	return spentCalories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("%w: steps must be greater than zero", ErrWrongArgument)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("%w: weight must be greater than zero", ErrWrongArgument)
	}
	if height <= 0 {
		return 0, fmt.Errorf("%w: height must be greater than zero", ErrWrongArgument)
	}
	if duration <= time.Duration(0) {
		return 0, fmt.Errorf("%w: duration must be greater than zero", ErrWrongArgument)
	}

	avgSpeed := meanSpeed(steps, height, duration)

	spentCalories := (weight * avgSpeed * duration.Minutes()) / minInH * walkingCaloriesCoefficient
	return spentCalories, nil
}
