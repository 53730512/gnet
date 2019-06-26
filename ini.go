package gnet

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

//Config ...
type configST struct {
	filepath string
	conflist []map[string]map[string]string
}

func newConfig() *configST {
	return &configST{}
}

func (v *configST) Load(filepath string) bool {
	v.filepath = filepath
	v.conflist = v.conflist[0:0]
	v.conflist = v.ReadList()
	return true
}

//GetValue key values:string
func (c *configST) GetValue(section, name string) string {
	c.ReadList()
	conf := c.ReadList()
	for _, v := range conf {
		for key, value := range v {
			if key == section {
				return value[name]
			}
		}
	}
	return ""
}

//GetValueInt key values:int
func (c *configST) GetValueInt(section, name string) int {
	c.ReadList()
	conf := c.ReadList()
	for _, v := range conf {
		for key, value := range v {
			if key == section {
				val, _ := strconv.Atoi(value[name])
				return val
			}
		}
	}
	return 0
}

//GetValueInt32 key values:int
func (c *configST) GetValueInt32(section, name string) int32 {
	c.ReadList()
	conf := c.ReadList()
	for _, v := range conf {
		for key, value := range v {
			if key == section {
				val, _ := strconv.ParseInt(value[name], 10, 32)
				return int32(val)
			}
		}
	}
	return 0
}

//GetValueInt64 key values:int
func (c *configST) GetValueInt64(section, name string) int64 {
	c.ReadList()
	conf := c.ReadList()
	for _, v := range conf {
		for key, value := range v {
			if key == section {
				val, _ := strconv.ParseInt(value[name], 10, 64)
				return val
			}
		}
	}
	return 0
}

//GetValueArray key values:[]int,split by ","
func (c *configST) GetValueArray(section, name string) []string {
	c.ReadList()
	conf := c.ReadList()
	for _, v := range conf {
		for key, value := range v {
			if key == section {
				arr := strings.Split(value[name], ",")
				return arr
			}
		}
	}
	return nil
}

//GetValueIntArray key values:[]int,split by ","
func (c *configST) GetValueIntArray(section, name string) []int {
	c.ReadList()
	conf := c.ReadList()
	for _, v := range conf {
		for key, value := range v {
			if key == section {
				arr := strings.Split(value[name], ",")
				arrValue := []int{}
				for _, str := range arr {
					val, _ := strconv.Atoi(str)
					arrValue = append(arrValue, val)
				}
				return arrValue
			}
		}
	}
	return nil
}

//SetValue Set the corresponding value of the key value, if not add, if there is a key change
func (c *configST) SetValue(section, key, value string) bool {
	c.ReadList()
	data := c.conflist
	var ok bool
	var index = make(map[int]bool)
	var conf = make(map[string]map[string]string)
	for i, v := range data {
		_, ok = v[section]
		index[i] = ok
	}

	i, ok := func(m map[int]bool) (i int, v bool) {
		for i, v := range m {
			if v == true {
				return i, true
			}
		}
		return 0, false
	}(index)

	if ok {
		c.conflist[i][section][key] = value
		return true
	}

	conf[section] = make(map[string]string)
	conf[section][key] = value
	c.conflist = append(c.conflist, conf)
	return true

}

//DeleteValue Delete the corresponding key values
func (c *configST) DeleteValue(section, name string) bool {
	c.ReadList()
	data := c.conflist
	for i, v := range data {
		for key := range v {
			if key == section {
				delete(c.conflist[i][key], name)
				return true
			}
		}
	}
	return false
}

//ReadList List all the configuration file
func (c *configST) ReadList() []map[string]map[string]string {

	file, err := os.Open(c.filepath)
	if err != nil {
		c.CheckErr(err)
	}
	defer file.Close()
	var data map[string]map[string]string
	var section string
	buf := bufio.NewReader(file)
	for {
		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				c.CheckErr(err)
			}
			if len(line) == 0 {
				break
			}
		}
		switch {
		case len(line) == 0:
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
			data = make(map[string]map[string]string)
			data[section] = make(map[string]string)
		default:
			i := strings.IndexAny(line, "=")
			value := strings.TrimSpace(line[i+1 : len(line)])
			data[section][strings.TrimSpace(line[0:i])] = value
			if c.uniquappend(section) == true {
				c.conflist = append(c.conflist, data)
			}
		}

	}

	return c.conflist
}

//CheckErr ...
func (c *configST) CheckErr(err error) string {
	if err != nil {
		return fmt.Sprintf("Error is :'%s'", err.Error())
	}
	return "Notfound this error"
}

//Ban repeated appended to the slice method
func (c *configST) uniquappend(conf string) bool {
	for _, v := range c.conflist {
		for k := range v {
			if k == conf {
				return false
			}
		}
	}
	return true
}
