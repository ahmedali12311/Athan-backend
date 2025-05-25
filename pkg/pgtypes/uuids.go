package pgtypes

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type UUIDS []uuid.UUID

func (u *UUIDS) Scan(src any) error {
	uuidList := UUIDS{}
	switch src := src.(type) {
	case string:
		s := strings.ReplaceAll(src, "{", "")
		s = strings.ReplaceAll(s, "}", "")
		sarr := strings.Split(s, ",")
		for _, v := range sarr {
			val, err := uuid.Parse(v)
			if err != nil {
				return err
			}
			uuidList = append(uuidList, val)
		}
	case []uint8:
		src2 := string(src)
		s := strings.ReplaceAll(src2, "[", "")
		s = strings.ReplaceAll(s, "]", "")
		s = strings.ReplaceAll(s, " ", "")
		s = strings.ReplaceAll(s, "\"", "")
		sarr := strings.Split(s, ",")
		for _, v := range sarr {
			val, err := uuid.Parse(v)
			if err != nil {
				return err
			}
			uuidList = append(uuidList, val)
		}
	case nil:
		*u = nil
		return nil
	default:
		*u = nil
		return errors.New("unknown uuid array")
	}
	*u = uuidList
	return nil
}
