package dnsconfig

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/abh/errorutil"
)

type objMap map[string]interface{}

func jsonLoader(fileName string, objmap objMap, fn func() error) error {
	fh, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fh.Close()

	decoder := json.NewDecoder(fh)
	if err = decoder.Decode(&objmap); err != nil {
		extra := ""
		if serr, ok := err.(*json.SyntaxError); ok {
			if _, seekErr := fh.Seek(0, io.SeekStart); seekErr != nil {
				return fmt.Errorf("seek error in %s: %w", fh.Name(), seekErr)
			}
			line, col, highlight := errorutil.HighlightBytePosition(fh, serr.Offset)
			extra = fmt.Sprintf(":\nError at line %d, column %d (file offset %d):\n%s",
				line, col, serr.Offset, highlight)
		}
		return fmt.Errorf("error parsing JSON object in config file %s%s\n%v",
			fh.Name(), extra, err)
	}

	return fn()
}

func toInt(i interface{}) (int, error) {
	switch i.(type) {
	case string:
		return strconv.Atoi(i.(string))
	case float64:
		return int(i.(float64)), nil
	case nil:
		return 0, nil
	}
	return 0, fmt.Errorf("Unknown type %T", i)
}

func toBool(i interface{}) (bool, error) {
	if b, ok := i.(bool); ok {
		return b, nil
	}
	n, err := toInt(i)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
