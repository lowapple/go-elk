package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"
)

// ServerConfig Type
type ServerConfig struct {
	RunMode string `yaml:"RunMode"`
}

// LogConfig Type
type LogConfig struct {
	Level         string `yaml:"Level"`
	WriteFile     bool   `yaml:"WriteFile"`
	Path          string `yaml:"Path"`
	FileName      string `yaml:"FileName"`
	RotatePattern string `yaml:"RotatePattern"`
	RotateMaxAge  int    `yaml:"RotateMaxAge"`
}

// HttpServerConfig Type
type HttpServerConfig struct {
	Port         int           `yaml:"Port"`
	ReadTimeout  time.Duration `yaml:"ReadTimeout"`  // sec
	WriteTimeout time.Duration `yaml:"WriteTimeout"` // sec
}

// DatabaseConfig Type
type DatabaseConfig struct {
	Type            string        `yaml:"Type"`
	Host            string        `yaml:"Host"`
	Port            int           `yaml:"Port"`
	Name            string        `yaml:"Name"`
	User            string        `yaml:"User"`
	Password        string        `yaml:"Password"`
	MaxIdleConns    int           `yaml:"MaxIdleConns"`
	MaxOpenConns    int           `yaml:"MaxOpenConns"`
	ConnMaxLifetime time.Duration `yaml:"ConnMaxLifetime"` // sec
}

// RedisConfig Type
type RedisConfig struct {
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
}

// PasswordConfig Type
type PasswordConfig struct {
	HashCost int `yaml:"HashCost"`
}

// AwsKmsConfig Type
type AwsKmsConfig struct {
	KeyId string `yaml:"KeyId"`
}

// KmsConfig Type
type KmsConfig struct {
	Iv  string `yaml:"Iv"`
	Key string `yaml:"Key"`
}

// Config Type
type Config struct {
	//Server     ServerConfig     `yaml:"Server"`
	Log LogConfig `yaml:"Log"`
	//HttpServer HttpServerConfig `yaml:"HttpServer"`
	//Database   DatabaseConfig   `yaml:"Database"`
	//Redis      RedisConfig      `yaml:"Redis"`
	//Password   PasswordConfig   `yaml:"Password"`
	//AwsKms     AwsKmsConfig     `yaml:"AwsKms"`
	//Kms        KmsConfig        `yaml:"Kms"`
}

var (
	configFilePath = "config/.config.yml"
	Conf           Config
)

// Setup config 셋업 (설정파일 로딩), 파일명: config/.config.yml
func Setup() bool {

	ok := true
	filename, _ := filepath.Abs(configFilePath)

	if !CheckFileExist(filename) {
		msg := fmt.Sprintf("config file not exist. [%s]", filename)
		log.Print(msg)
		return false
	}

	yamlFile, _ := ioutil.ReadFile(filename)

	if err := yaml.Unmarshal(yamlFile, &Conf); err != nil {
		msg := fmt.Sprintf("config file load error. file path=[%s], error=%s", configFilePath, err.Error())
		log.Print(msg)
		ok = false
	}
	return ok
}

// CheckFileExist 파일이 있는지 체크
func CheckFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else {
		return false
	}
}

// CheckNotExist check if the file exists
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

// GetSize get the file size
func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)

	return len(content), err
}

// GetExt get the file ext
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// CheckPermission check if the file has permission
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

// IsNotExistMkDir create a directory if it does not exist
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// MkDir create a directory
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Open a file according to a specific mode
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// MustOpen maximize trying to open the file
func MustOpen(filePath, fileName string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	lDir := dir + "/" + filePath
	perm := CheckPermission(lDir)
	if perm {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", lDir)
	}

	err = IsNotExistMkDir(lDir)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", lDir, err)
	}

	f, err := Open(lDir+"/"+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to OpenFile :%v", err)
	}

	return f, nil
}

// GetWorkDirectory get working directory
func GetWorkDirectory() string {
	curDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("[info] working directory %s", curDir)
	return curDir
}
