package jsonutil

import (
	"encoding/json"
	"encoding/json/jsontext"
	"strings"
	"testing"

	"github.com/lvan100/golib/errutil"
)

func TestJSON(t *testing.T) {
	{
		s := `{
			"Int": 3,
			"IntPtr": 3,
			"Unknown": "abc"
		}`
		base := &Base{}
		r := strings.NewReader(s)
		d := jsontext.NewDecoder(r)
		if err := base.DecodeJSON(d); err != nil {
			t.Fatal(err)
		}
		buf, err := json.Marshal(base)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s", buf)
	}
	{
		s := `{
			"IntList": [3],
			"StringList": ["","null"]
		}`
		l := &List{}
		r := strings.NewReader(s)
		d := jsontext.NewDecoder(r)
		if err := l.DecodeJSON(d); err != nil {
			t.Fatal(err)
		}
		buf, err := json.Marshal(l)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s", buf)
	}
}

// HashKey returns a hash value for the given string.
// 业界推荐算法，对于短字符串，碰撞的概率很低很低.
func HashKey(s string) uint64 {
	const (
		offset = 14695981039346656037
		prime  = 1099511628211
	)
	h := uint64(offset)
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime
	}
	return h
}

type Base struct {
	Int       int
	IntPtr    *int
	String    string
	StringPtr *string
	Bytes     []byte
}

func (b *Base) DecodeJSON(d *jsontext.Decoder) error {
	const (
		hashInt       = 0x41a91f19c98dd49e // HashKey("Int")
		hashIntPtr    = 0x3305f2829a12fcb8 // HashKey("IntPtr")
		hashString    = 0x58b4b3ecd4eb6238 // HashKey("String")
		hashStringPtr = 0xe8751a6330efffa2 // HashKey("StringPtr")
		hashBytes     = 0xeeeea7adc131a244 // HashKey("Bytes")
	)

	if err := DecodeObjectBegin(d); err != nil {
		return err
	}

	// 设置默认值
	b.Int = 9

	// 记录必传字段
	var (
		foundInt bool
	)

	for {
		if d.PeekKind() == '}' {
			break
		}
		key, err := DecodeKey(d)
		if err != nil {
			return err
		}
		switch HashKey(key) {
		case hashInt:
			//if key != "Int" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
			if b.Int, err = DecodeInt[int](d); err != nil {
				return err
			}
			foundInt = true
		case hashIntPtr:
			//if key != "IntPtr" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
			if b.IntPtr, err = DecodeIntPtr[int](d); err != nil {
				return err
			}
		case hashString:
			//if key != "String" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
			if b.String, err = DecodeString(d); err != nil {
				return err
			}
		case hashStringPtr:
			//if key != "StringPtr" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
			if b.StringPtr, err = DecodeStringPtr(d); err != nil {
				return err
			}
		case hashBytes:
			//if key != "Bytes" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
			if b.Bytes, err = DecodeBytes(d); err != nil {
				return err
			}
		default:
			if err = d.SkipValue(); err != nil {
				return err
			}
		}
	}

	if err := DecodeObjectEnd(d); err != nil {
		return err
	}

	// 检查必传字段
	if !foundInt {
		return errutil.Explain(nil, "missing required field Int")
	}
	return nil
}

type List struct {
	IntList    []int
	StringList []string
	BytesList  [][]byte
}

func (l *List) DecodeJSON(d *jsontext.Decoder) error {
	const (
		hashIntList    = 0x9273f90a7d88b56a // HashKey("IntList")
		hashStringList = 0xc37ebdf18413dc00 // HashKey("StringList")
		hashBytesList  = 0x48c8b441d1a49dac // HashKey("BytesList")
	)

	if err := DecodeObjectBegin(d); err != nil {
		return err
	}

	for {
		if d.PeekKind() == '}' {
			break
		}
		key, err := DecodeKey(d)
		if err != nil {
			return err
		}
		switch HashKey(key) {
		case hashIntList:
			//if key != "IntList" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
			if l.IntList, err = DecodeInts[int](d); err != nil {
				return err
			}
		case hashStringList:
			//if key != "StringList" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
			if l.StringList, err = DecodeStrings(d); err != nil {
				return err
			}
		case hashBytesList:
			//if key != "BytesList" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
		default:
			if err = d.SkipValue(); err != nil {
				return err
			}
		}
	}

	if err := DecodeObjectEnd(d); err != nil {
		return err
	}

	return nil
}

type Map struct {
	IntMap     map[string]int
	StringMap  map[string]string
	BytesMap   map[string][]byte
	ListMap    map[string]List
	ListPtrMap map[string]*List
	MapPtrMap  map[string]*Map
}
