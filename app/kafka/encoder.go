package kafka

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type Encoder interface {
	Encode(in []byte, out interface{}) error
	Decode(interface{}) ([]byte, error)
}

type jsonEncoder struct{}

func (*jsonEncoder) Encode(in []byte, out interface{}) error {
	err := json.Unmarshal(in, out)
	if err != nil {
		return err
	}
	return nil
}

func (*jsonEncoder) Decode(t interface{}) ([]byte, error) {
	msg, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

type stringEncoder struct{}

func (*stringEncoder) Encode(in []byte, out interface{}) error {
	switch out.(type) {
	case *string:
		*out.(*string) = string(in)
	default:
		return errors.New("incorrect type for string encoder with value")
	}
	return nil
}

func (*stringEncoder) Decode(t interface{}) ([]byte, error) {
	return []byte(t.(string)), nil
}

type byteEncoder struct{}

func (*byteEncoder) Encode(in []byte, out interface{}) error {
	switch out.(type) {
	case *[]byte:
		*out.(*[]byte) = in
	default:
		return errors.Errorf("incorrect type for byte encoder")
	}
	return nil
}

func (*byteEncoder) Decode(t interface{}) ([]byte, error) {
	switch t.(type) {
	case []byte:
		return t.([]byte), nil
	default:
		return nil, errors.Errorf("incorrect type for byte decoder")
	}
}

func ByteEncoder() Encoder {
	return &byteEncoder{}
}

func JsonEncoder() Encoder {
	return &jsonEncoder{}
}

func StringEncoder() Encoder {
	return &stringEncoder{}
}
