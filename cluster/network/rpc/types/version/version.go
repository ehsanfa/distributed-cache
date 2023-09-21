package version

import (
	"bytes"
	"dbcache/cluster/version"
	"encoding/gob"
)

type Version struct {
	Ver version.Version
}

type marshal struct {
	Number uint64
}

func (v *Version) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(marshal{
		Number: v.Ver.Number(),
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (v *Version) UnmarshalBinary(data []byte) error {
	m := &marshal{}
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&m); err != nil {
		return err
	}
	v.Ver = version.CreateGenClockVersion(m.Number)
	return nil
}
