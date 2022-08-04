package config

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type ServerConfig struct {
	Name                string `cfg:"name"`
	IpVer               string `cfg:"ip_ver"`
	Ip                  string `cfg:"ip"`
	Port                int    `cfg:"port"`
	MaxConn             int    `cfg:"max_connection"`
	MaxPacketSize       int    `cfg:"max_packet_size"`
	MaxPoolCapacitySize int    `cfg:"max_pool_capacity_size"`
	LimitTask           int    `cfg:"limit_task"`
}

var ServerCon *ServerConfig

func init() {
	//无配置文件下的初始化 默认配置
	ServerCon = &ServerConfig{
		Name:                "tcp-server",
		IpVer:               "tcp4",
		Ip:                  "localhost",
		Port:                1234,
		MaxConn:             100,
		MaxPacketSize:       1024,
		MaxPoolCapacitySize: 300,
		LimitTask:           100,
	}
}

//LoadConfig 根据路径加载配置文件
func LoadConfig(configPath string) {
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	err = parse(file)
	if err != nil {
		log.Fatalln(err)
	}
}

//解析配置文件到 ServerConfig 结构体中
func parse(file *os.File) error {
	//
	var (
		err   error
		key   string
		value string
	)
	rawMap := make(map[string]string)
	scanner := bufio.NewScanner(file)

	//一行一行的读取文件
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] == '#' {
			continue
		}

		pivot := strings.IndexAny(line, ":")
		if pivot > 0 && pivot < len(line)-1 {
			key = line[:pivot]
			value = line[pivot+1:]
			rawMap[strings.ToLower(key)] = value
		}
	}
	if err = scanner.Err(); err != nil {
		return err
	}
	//通过反射来获取tag 字段
	config := &ServerConfig{}
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)

	for i := 0; i < v.Elem().NumField(); i++ {
		filed := t.Elem().Field(i)
		filedVal := v.Elem().Field(i)

		//指定tag
		var ok bool
		key, ok = filed.Tag.Lookup("cfg")
		if !ok {
			key = filed.Name
		}

		//取出value
		var exist bool
		value, exist = rawMap[strings.ToLower(key)]
		//存在值 给结构体字段赋值
		if filedVal.CanSet() && exist {
			switch filed.Type.Kind() {
			case reflect.String:
				filedVal.SetString(value)
			case reflect.Int:
				if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
					filedVal.SetInt(intValue)
				}
			}
		}
	}
	ServerCon = config
	return nil
}
