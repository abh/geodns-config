package dnsconfig

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/abh/errorutil"
)

func loadJSONFile(fileName string, v interface{}) error {
	fh, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fh.Close()

	decoder := json.NewDecoder(fh)
	if err := decoder.Decode(v); err != nil {
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
	return nil
}

// flexBool decodes JSON booleans, numbers (0/non-zero), and numeric strings
// ("0", "1") into a bool — matching the historical config format.
type flexBool bool

func (b *flexBool) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	v, err := toBool(raw)
	if err != nil {
		return err
	}
	*b = flexBool(v)
	return nil
}

// flexInt decodes JSON numbers and numeric strings into an int.
type flexInt int

func (n *flexInt) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	v, err := toInt(raw)
	if err != nil {
		return err
	}
	*n = flexInt(v)
	return nil
}

func toInt(i interface{}) (int, error) {
	switch v := i.(type) {
	case string:
		return strconv.Atoi(v)
	case float64:
		return int(v), nil
	case nil:
		return 0, nil
	}
	return 0, fmt.Errorf("unknown type %T", i)
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
