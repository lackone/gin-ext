package formatter

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/lackone/gin-ext/framework/contract"
	"time"
)

func JsonFormatter(level contract.LogLevel, t time.Time, msg string, fields map[string]interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	fields["msg"] = msg
	fields["level"] = Prefix(level)
	fields["datetime"] = t.Format(time.RFC3339)
	c, err := json.Marshal(fields)
	if err != nil {
		return bf.Bytes(), errors.New("json format error")
	}

	bf.Write(c)
	return bf.Bytes(), nil
}
