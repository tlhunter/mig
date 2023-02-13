package migrations

import (
	"bufio"
	"errors"
	"os"
)

type MigrationPair struct {
	Up     string
	Down   string
	UpTx   bool
	DownTx bool
}

const (
	STATE_START  = 0 // ignore content
	STATE_UP     = 1 // capture Up queries
	STATE_MIDDLE = 2 // ignore content
	STATE_DOWN   = 3 // capture Down queries
	STATE_FINISH = 4 // ignore content

	DELIM_BEGIN_UP         = "--BEGIN MIGRATION UP--"
	DELIM_BEGIN_UP_NO_TX   = "--BEGIN MIGRATION UP NO TRANSACTION--"
	DELIM_END_UP           = "--END MIGRATION UP--"
	DELIM_BEGIN_DOWN       = "--BEGIN MIGRATION DOWN--"
	DELIM_BEGIN_DOWN_NO_TX = "--BEGIN MIGRATION DOWN NO TRANSACTION--"
	DELIM_END_DOWN         = "--END MIGRATION DOWN--"
)

// Opens a migration file then steps through it looking for an up and down block.
// Queries within the two blocks are then returned.
// Lines that fall outside of the blocks are ignored.
// If it doesn't find a well formed up them down block an error is returned.
// This is because any poorly-formed comments should not be mis-interpreted.
func GetQueriesFromFile(filename string) (MigrationPair, error) {
	pair := MigrationPair{
		Up:     "",
		Down:   "",
		UpTx:   true,
		DownTx: true,
	}

	state := STATE_START

	file, err := os.Open(filename)
	if err != nil {
		return pair, err
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		switch line {
		case DELIM_BEGIN_UP:
			if state != STATE_START {
				return pair, errors.New("invalid begin up delimiter")
			}
			pair.UpTx = true
			state = STATE_UP

		case DELIM_BEGIN_UP_NO_TX:
			if state != STATE_START {
				return pair, errors.New("invalid begin up no transaction delimiter")
			}
			pair.UpTx = false
			state = STATE_UP

		case DELIM_END_UP:
			if state != STATE_UP {
				return pair, errors.New("invalid end up delimiter")
			}
			state = STATE_MIDDLE

		case DELIM_BEGIN_DOWN:
			if state != STATE_MIDDLE {
				return pair, errors.New("invalid begin down delimiter")
			}
			pair.DownTx = true
			state = STATE_DOWN

		case DELIM_BEGIN_DOWN_NO_TX:
			if state != STATE_MIDDLE {
				return pair, errors.New("invalid begin down no transaction delimiter")
			}
			pair.DownTx = false
			state = STATE_DOWN

		case DELIM_END_DOWN:
			if state != STATE_DOWN {
				return pair, errors.New("invalid end down delimiter")
			}
			state = STATE_FINISH

		default:
			if state == STATE_UP {
				pair.Up += line + "\n"
			} else if state == STATE_DOWN {
				pair.Down += line + "\n"
			}
			// ignore content outside of delimiters
		}
	}

	if state != STATE_FINISH {
		return pair, errors.New("failed to parse migration file")
	}

	return pair, nil
}
