package main

import (
	"bufio"
	"os"
	"regexp"
)

func main() {
	sourceURL = ""
	dstURL = ""
	os.
}

// UpdateNtpConfig update ntp.conf
func UpdateNtpConfig(mainserver, secondserver string) error {
	fileName := "/etc/ntp.conf"
	regStr := `^server.*`

	confStr, err := ReadConfFile(fileName, regStr)
	confStr = append(confStr, mainserver, secondserver)

	if err != nil {
		return err
	}

	//fmt.Println(err)
	err = WriteConfFile(fileName, confStr)

	return err
}

// ReadConfFile read config file into a slice
func ReadConfFile(fileName, regStr string) ([]string, error) {
	var confStr []string
	// open config file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	reg, err := regexp.Compile(regStr)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(file)
	scanner := bufio.NewScanner(buf)

	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Println(scanner.Text())
		if reg.MatchString(line) {
			continue
		}
		confStr = append(confStr, line, "\n")
	}
	// fmt.Println(confStr)
	defer file.Close()
	return confStr, err
}

//WriteConfFile overwrite config file
func WriteConfFile(fileName string, confStr []string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, v := range confStr {
		_, err := file.WriteString(v)
		if err != nil {
			return err
		}
	}
	return err
}
