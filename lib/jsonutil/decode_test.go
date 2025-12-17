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
		s := `[{
			"IntList": [3],
			"StringList": ["","null"],
			"IntPtrList": [[3,null]]
		}]`
		r := strings.NewReader(s)
		d := jsontext.NewDecoder(r)
		l, err := DecodeArray(DecodeObject(NewList))(d)
		if err != nil {
			t.Fatal(err)
		}
		buf, err := json.Marshal(l)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s", buf)
	}
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
			"StringList": ["","null"],
			"IntPtrList": [[3,null]]
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
	{
		s := `{
			"IntIntMap": {"3":3},
			"StrStrPtrMap": {"3":"3"},
			"StrListPtrMap": {"3": {
				"IntList": [3],
				"StringList": ["","null"],
				"IntPtrList": [[3,null]]
			}}
		}`
		m := &Map{}
		r := strings.NewReader(s)
		d := jsontext.NewDecoder(r)
		if err := m.DecodeJSON(d); err != nil {
			t.Fatal(err)
		}
		buf, err := json.Marshal(m)
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
			if b.Int, err = DecodeInt(d); err != nil {
				return err
			}
			foundInt = true
		case hashIntPtr:
			//if key != "IntPtr" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
			if b.IntPtr, err = DecodeIntPtr(d); err != nil {
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
	IntPtrList [][]*int
}

func NewList() *List {
	return &List{}
}

func (l *List) DecodeJSON(d *jsontext.Decoder) error {
	const (
		hashIntList    = 0x9273f90a7d88b56a // HashKey("IntList")
		hashStringList = 0xc37ebdf18413dc00 // HashKey("StringList")
		hashIntPtrList = 0x80c120a941785e80 // HashKey("IntPtrList")
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
			if l.IntList, err = DecodeArray(DecodeInt)(d); err != nil {
				return err
			}
		case hashStringList:
			//if key != "StringList" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
			if l.StringList, err = DecodeArray(DecodeString)(d); err != nil {
				return err
			}
		case hashIntPtrList:
			//if key != "IntPtrList" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
			if l.IntPtrList, err = DecodeArray(DecodeArray(DecodeIntPtr))(d); err != nil {
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

	return nil
}

type Map struct {
	IntIntMap       map[int64]int
	StrStrPtrMap    map[string]*string
	StrListPtrMap   map[string]*List
	IntMapStrIntMap map[int]map[string]int
}

func (m *Map) DecodeJSON(d *jsontext.Decoder) error {
	const (
		hashIntIntMap       = 0xc3172d7c329b5eef // HashKey("IntIntMap")
		hashStrStrPtrMap    = 0xc3447c678a33494b // HashKey("StrStrPtrMap")
		hashStrListPtrMap   = 0xcdf8c071ff4079ec // HashKey("StrListPtrMap")
		hashIntMapStrIntMap = 0x8fb20e202820cc8e // HashKey("IntMapStrIntMap")
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
		case hashIntIntMap:
			//if key != "IntIntMap" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
			if m.IntIntMap, err = DecodeMap(DecodeInt64, DecodeInt)(d); err != nil {
				return err
			}
		case hashStrStrPtrMap:
			//if key != "StrStrPtrMap" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
			if m.StrStrPtrMap, err = DecodeMap(DecodeString, DecodeStringPtr)(d); err != nil {
				return err
			}
		case hashStrListPtrMap:
			//if key != "StrListPtrMap" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
			if m.StrListPtrMap, err = DecodeMap(DecodeString, DecodeObject(NewList))(d); err != nil {
				return err
			}
		case hashIntMapStrIntMap:
			//if key != "IntMapStrIntMap" {
			//	return fmt.Errorf("unknown field name: %s", key)
			//}
			if m.IntMapStrIntMap, err = DecodeMap(DecodeInt, DecodeMap(DecodeString, DecodeInt))(d); err != nil {
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

	return nil
}
