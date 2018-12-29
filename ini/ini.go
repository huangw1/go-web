package ini

import (
	"runtime"
	"io"
	"os"
	"bytes"
	"strconv"
	"regexp"
	"errors"
	"sync"
	"fmt"
	"bufio"
	"strings"
)

var (
	DefaultSection = "default"
	LineBreak = "\n"
	varPattern = regexp.MustCompile(`\<([^>]+)\>`)
	PrettyFormat = true
)

func init()  {
	if  runtime.GOOS == "windows" {
		LineBreak = "\r\n"
	}
}

func inSlice(str string, slice []string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

type dataSource interface {
	Reader() (io.Reader, error)
}

type sourceFile struct {
	name string
}

func (source sourceFile) Reader() (io.Reader, error) {
	return os.Open(source.name)
}

type sourceData struct {
	data []byte
}

func (source *sourceData) Reader() (io.Reader, error) {
	return bytes.NewReader(source.data), nil
}

/**
	Key
 */
type Key struct {
	s *Section
	name string
	value string
	Comment string
	isAutoIncrement bool
}

func (k *Key) Name() string {
	return k.name
}

func (k *Key) Value() string {
	return k.value
}

func (k *Key) String() string {
	// todo 处理变量
	return k.value
}

func (k *Key) Bool() (bool, error) {
	return strconv.ParseBool(k.String())
}

func (k *Key) Float64() (float64, error) {
	return strconv.ParseFloat(k.String(), 64)
}

func (k *Key) Int() (int, error) {
	return strconv.Atoi(k.String())
}

func (k *Key) Int64() (int64, error) {
	return strconv.ParseInt(k.String(), 10, 64)
}

func (k *Key) MustBool(defaultVal ...bool) bool {
	val, err := k.Bool()
	if len(defaultVal) > 0 && err != nil {
		return defaultVal[0]
	}
	return val
}

func (k *Key) MustFloat64(defaultVal ...float64) float64 {
	val, err := k.Float64()
	if len(defaultVal) > 0 && err != nil {
		return defaultVal[0]
	}
	return val
}

func (k *Key) MustInt(defaultVal ...int) int {
	val, err := k.Int()
	if len(defaultVal) > 0 && err != nil {
		return defaultVal[0]
	}
	return val
}

func (k *Key) MustInt64(defaultVal ...int64) int64 {
	val, err := k.Int64()
	if len(defaultVal) > 0 && err != nil {
		return defaultVal[0]
	}
	return val
}

/**
	Section
 */
type Section struct {
	f *File
	Comment string
	name string
	keys map[string]*Key
	keyList []string
	keysHash map[string]string
}

func newSection(f *File, name string) *Section {
	return &Section{
		f: f,
		Comment: "",
		name: name,
		keys: make(map[string]*Key),
		keyList: make([]string, 0, 10),
		keysHash: make(map[string]string),
	}
}

func (s *Section) Name() string {
	return s.name
}

func (s *Section) NewKey(name, val string) (*Key, error) {
	if len(name) == 0 {
		return nil, errors.New("key is required")
	}
	if s.f.BlockMode {
		s.f.mutex.Lock()
		defer s.f.mutex.Unlock()
	}
	if inSlice(name, s.keyList) {
		s.keys[name].value = val
		s.keysHash[name] = val
		return s.keys[name], nil
	}
	key := &Key{
		s: s,
		name: name,
		value: val,
		Comment: "",
		isAutoIncrement: false,
	}
	s.keyList = append(s.keyList, name)
	s.keys[name] = key
	s.keysHash[name] = val
	return key, nil
}

func (s *Section) GetKey(name string) (*Key, error) {
	if s.f.BlockMode {
		s.f.mutex.Lock()
		defer s.f.mutex.Unlock()
	}
	key := s.keys[name]
	if key == nil {
		return nil, errors.New(fmt.Sprintf("%s does not exists", name))
	}
	return key, nil
}

func (s *Section) Key(name string) *Key {
	key, err := s.GetKey(name)
	if err != nil {
		return &Key{}
	}
	return key
}

func (s *Section) DeleteKey(name string) {
	if s.f.BlockMode {
		s.f.mutex.Lock()
		defer s.f.mutex.Unlock()
	}
	for i, v := range s.keyList {
		if v == name {
			delete(s.keys, name)
			delete(s.keysHash, name)
			s.keyList = append(s.keyList[:i], s.keyList[i:]...)
		}
	}
}

/**
	File
 */
type File struct {
	mutex sync.RWMutex
	BlockMode bool
	dataSources []dataSource
	sections map[string]*Section
	sectionList []string
}

func NewFile(sources []dataSource) *File {
	return &File{
		BlockMode: true,
		dataSources: sources,
		sections: make(map[string]*Section),
		sectionList: make([]string, 0, 10),
	}
}

func parseSource(source interface{}) (dataSource, error) {
	switch s := source.(type) {
	case string:
		return sourceFile{s}, nil
	case []byte:
		return &sourceData{s}, nil
	default:
		return nil, errors.New("unknown ini source")
	}
}

func Load(source interface{}, others ...interface{}) (_ *File, err error) {
	sources := make([]dataSource, len(others) + 1)
	sources[0], err = parseSource(source)
	if err != nil {
		return nil, err
	}
	for i := range others {
		sources[i + 1], err = parseSource(others[i])
		if err != nil {
			return nil, err
		}
	}
	file := NewFile(sources)
	return  file, file.Reload()
}

func (f *File) Reload() error {
	for _, s := range f.dataSources {
		r, err := s.Reader()
		if err != nil {
			return err
		}
		if err = f.parse(r); err != nil {
			return err
		}
	}
	return nil
}

func (f *File) parse(reader io.Reader) error {
	buf := bufio.NewReader(reader)
	// BOM
	mask, err := buf.Peek(3)
	if err == nil && len(mask) >= 3 && mask[0] == 239 && mask[1] == 187 && mask[2] == 191 {
		buf.Read(mask)
	}
	section, err := f.NewSection(DefaultSection)
	if err != nil {
		return err
	}
	isEnd := false
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		length := len(line)
		if err != nil {
			if err != io.EOF {
				return errors.New(fmt.Sprintf("error in reading next line %v", err))
			}
			if length == 0 {
				break
			}
			isEnd = true
		}
		if length == 0 {
			continue
		}

		switch {
		case line[0] == '[' && line[length - 1] == ']':
			name := strings.TrimSpace(line[1:length - 1])
			section, err = f.NewSection(name)
			if err != nil {
				return err
			}
			continue
		}

		pos := strings.IndexAny(line, "=")
		if pos > 0 {
			name := strings.TrimSpace(line[:pos])
			val := strings.TrimSpace(line[pos + 1:])
			_, err := section.NewKey(name, val)
			if err != nil {
				return err
			}
		}

		if isEnd {
			break
		}
	}
	return nil
}

func (f *File) NewSection(name string) (*Section, error) {
	if len(name) == 0 {
		return nil, errors.New("section is required")
	}
	if f.BlockMode {
		f.mutex.Lock()
		defer f.mutex.Unlock()
	}
	if inSlice(name, f.sectionList) {
		return f.sections[name], nil
	}
	f.sectionList = append(f.sectionList, name)
	f.sections[name] = newSection(f, name)
	return f.sections[name], nil
}

func (f *File) GetSection(name string) (*Section, error) {
	if len(name) == 0 {
		name = DefaultSection
	}
	if f.BlockMode {
		f.mutex.RLock()
		defer f.mutex.RUnlock()
	}
	sec := f.sections[name]
	if sec == nil {
		return nil, errors.New(fmt.Sprintf("not found section %s", name))
	}
	return sec, nil
}

func (f *File) Section(name string) *Section {
	sec, err := f.GetSection(name)
	if err != nil {
		return newSection(f, name)
	}
	return sec
}

