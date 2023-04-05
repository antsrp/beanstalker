package events

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/antsrp/beanstalker/config"
	ievents "github.com/antsrp/beanstalker/interfaces/events"
	"github.com/antsrp/beanstalker/interfaces/queue"
)

type Command string

const (
	SET  Command = "set"
	LIST Command = "list"
	PUT  Command = "put"
)

type Option string

const (
	TUBE     Option = "tube"
	TTR      Option = "ttr"
	DELAY    Option = "delay"
	PRIORITY Option = "priority"
)

const (
	UNKNOWN_COMMAND = "unknown command"
	EMPTY_COMMAND   = "empty command"

	UNKNOWN_OPTION = "unknown option"
	EMPTY_OPTION   = "empty option"

	EMPTY_OPTION_VALUE = "empty value"
	BAD_OPTION_VALUE   = "bad option value"

	NO_OPTION_SET = "no option"

	INT_VALUE_CONSTRAINT = "value must be of integer type and more than zero"
	JSON_CONSTRAINT      = "value must be of json type"
)

var (
	ErrUnknownCommand = errors.New(UNKNOWN_COMMAND)
	ErrEmptyCommand   = errors.New(EMPTY_COMMAND)

	ErrUnknownOption = errors.New(UNKNOWN_OPTION)
	ErrEmptyOption   = errors.New(EMPTY_OPTION)

	ErrEmptyOptionValue = errors.New(EMPTY_OPTION_VALUE)
	ErrBadOptionValue   = errors.New(BAD_OPTION_VALUE)

	ErrNoOptionSet = errors.New(NO_OPTION_SET)
)

const (
	PUT_FROM_CONSOLE = "-d"
	PUT_FROM_FILE    = "-f"
)

type Listener struct {
	producer queue.Producer
	cfg      *config.Config
}

var _ ievents.Controller = Listener{}

func NewListener(qp queue.Producer, cfg *config.Config) *Listener {
	return &Listener{
		producer: qp,
		cfg:      cfg,
	}
}

func (l Listener) updateQueueSettings(sets *config.Settings, changed map[Option]bool) {
	for k, v := range changed {
		if !v { // not changed
			continue
		}
		switch k {
		case TUBE:
			l.cfg.ChangeTube(sets.Tube)
			l.producer.SetTube(sets.Tube)
		case DELAY:
			l.cfg.ChangeDelay(sets.Delay)
			l.producer.SetDelay(sets.Delay)
		case PRIORITY:
			l.cfg.ChangePriority(sets.Priority)
			l.producer.SetPriority(sets.Priority)
		case TTR:
			l.cfg.ChangeTTR(sets.TTR)
			l.producer.SetTTR(sets.TTR)
		}
	}
}

func (l Listener) readDataFromFile(path string) ([]byte, error) {
	rel, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(rel)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func (l Listener) Parse(s string) (*string, error) {
	s = strings.Trim(s, "\r\n")
	arr := strings.Split(s, " ")
	ln := len(arr)

	if ln == 0 {
		return nil, ErrEmptyCommand
	}

	cmd := strings.ToLower(arr[0])

	switch Command(cmd) {
	case LIST:
		list, err := l.producer.ListTubes()
		return &list, err
	case SET:
		if ln == 1 {
			return nil, fmt.Errorf("%w for operation %s", ErrNoOptionSet, cmd)
		}
		sets := &config.Settings{}

		checker := map[Option]bool{
			TUBE:     false,
			DELAY:    false,
			PRIORITY: false,
			TTR:      false,
		}

		for i := 1; i < ln; i++ {
			parts := strings.Split(arr[i], "=")
			lp := len(parts)
			if lp == 0 {
				return nil, ErrEmptyOption
			}
			key := strings.ToLower(parts[0])
			optkey := Option(key)
			if lp == 1 {
				if _, contains := checker[optkey]; contains {
					return nil, fmt.Errorf("%w for option: %s", ErrEmptyOptionValue, key)
				} else {
					return nil, fmt.Errorf("%w: %s", ErrUnknownOption, key)
				}
			}
			if parts[1] == "" {
				return nil, fmt.Errorf("%w for option: %s", ErrEmptyOptionValue, key)
			}
			switch optkey {
			case TUBE:
				sets.Tube = parts[1]
			case DELAY:
				var delay int
				if _, err := fmt.Sscan(parts[1], &delay); err != nil || delay < 0 {
					return nil, fmt.Errorf("%w %s for key %s: %s", ErrBadOptionValue, parts[1], key, INT_VALUE_CONSTRAINT)
				} else {
					sets.Delay = delay
				}
			case PRIORITY:
				var priority uint32
				if _, err := fmt.Sscan(parts[1], &priority); err != nil {
					return nil, fmt.Errorf("%w %s for key %s: %s", ErrBadOptionValue, parts[1], key, INT_VALUE_CONSTRAINT)
				} else {
					sets.Priority = priority
				}
			case TTR:
				var ttr int
				if _, err := fmt.Sscan(parts[1], &ttr); err != nil || ttr < 0 {
					return nil, fmt.Errorf("%w %s for key %s: %s", ErrBadOptionValue, parts[1], key, INT_VALUE_CONSTRAINT)
				} else {
					sets.TTR = ttr
				}
			default:
				return nil, fmt.Errorf("%w: %s", ErrUnknownOption, key)
			}
			checker[optkey] = true
		}
		l.updateQueueSettings(sets, checker)
		return nil, nil
	case PUT:
		if ln <= 2 {
			return nil, ErrEmptyOption
		}
		key := strings.ToLower(arr[1])
		switch key {
		case PUT_FROM_CONSOLE:
			value := strings.Join(arr[2:], " ")
			data := []byte(value)
			fmt.Println("string: ", string(data))
			if res, err := l.producer.Put(data); err != nil {
				return nil, err
			} else {
				return &res, nil
			}
		case PUT_FROM_FILE:
			data, err := l.readDataFromFile(arr[2])
			if err != nil {
				return nil, fmt.Errorf("%w %s for key %s: %w", ErrBadOptionValue, arr[2], key, err)
			}
			if res, err := l.producer.Put(data); err != nil {
				return nil, err
			} else {
				return &res, nil
			}
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnknownOption, key)
		}
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownCommand, cmd)
	}
}
