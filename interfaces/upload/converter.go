package fileupload

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseTime(t string) (time.Duration, error) {
	var sec, mins int
	var err error

	parts := strings.SplitN(t, ".", 2)

	switch len(parts) {
	case 1:
		sec, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
	case 2:
		mins, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}

		sec, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
	default:
		return 0, fmt.Errorf("invalid time: %s", t)
	}

	if sec > 59 || sec < 0 || mins > 59 || mins < 0 {
		return 0, fmt.Errorf("invalid time: %s", t)
	}

	return time.Duration(mins)*time.Minute + time.Duration(sec)*time.Second, nil
}
