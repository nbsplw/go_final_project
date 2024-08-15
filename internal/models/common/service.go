package common

import (
	"errors"
	"fmt"
	"main/pkg"
	"time"
)

func (t *AddTask) CheckTask() error {
	if t.Title == "" {
		return fmt.Errorf("не указан заголовок задачи")
	}
	if t.Date == "" {
		t.Date = time.Now().Format("20060102")
	}
	date, err := time.Parse("20060102", t.Date)
	if err != nil {
		return fmt.Errorf("дата представлена в формате, отличном от 20060102")
	}

	if t.Repeat != "" && t.Repeat[0] != 'd' && t.Repeat[0] != 'w' && t.Repeat[0] != 'm' && t.Repeat[0] != 'y' {
		return errors.New("неверное правило повторения")
	}

	if len(t.Repeat) > 0 {
		if t.Repeat[0] != 'd' && t.Repeat[0] != 'w' && t.Repeat[0] != 'm' && t.Repeat[0] != 'y' {
			return errors.New("неверное правило повторения")
		}
		if t.Repeat[0] == 'd' || t.Repeat[0] == 'w' || t.Repeat[0] == 'm' {
			if len(t.Repeat) < 3 {
				return errors.New("неверное правило повторения")
			}
		}
	}

	if date.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
		if t.Repeat == "" {
			t.Date = time.Now().Format("20060102")
		}
	}

	if date.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
		if t.Repeat != "" {
			t.Date, err = pkg.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				return fmt.Errorf("ошибка при вычислении следующей даты: %v", err)
			}
		}
	}

	return nil
}
