package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

// ErrInvalidParse не правильный аргумент для преобразования
var ErrInvalidParse = errors.New("wrong argument for parsing")

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, time.Duration(0), errors.New("wrong count of arguments")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, time.Duration(0), fmt.Errorf("%w: %w", ErrInvalidParse, err)
	}
	if steps <= 0 {
		return 0, time.Duration(0), errors.New("steps must be greater then 0")
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, time.Duration(0), fmt.Errorf("%w: %w", ErrInvalidParse, err)
	}
	if duration <= time.Duration(0) {
		return 0, time.Duration(0), errors.New("duration must be greater then 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	distanceKm := (float64(steps) * stepLength) / mInKm

	cal, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return fmt.Sprintf("Количество шагов: %d.\n"+
		"Дистанция составила %.2f км.\n"+
		"Вы сожгли %.2f ккал.\n", steps, distanceKm, cal)
}
