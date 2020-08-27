package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/apparentlymart/go-cidr/cidr"
)

func main() {
	shellpath := "./getsource.sh" //脚本文件位置
	cmd := exec.Command("/bin/bash", shellpath)
	// output, _ := cmd.Output()
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(output) //这里不会实时输出
	err = cmd.Wait() //等到cmd执行结束
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(output))

	cidrDir := "/var/isp/cidr/"       //cidr格式原数据文件夹
	rangeDir := "/var/isp/range/"     //range格式原数据文件夹
	resultsDir := "/var/isp/results/" //结果保存文件夹

	err = AllFilesToWrite(cidrDir, resultsDir, "cidr")
	err = AllFilesToWrite(rangeDir, resultsDir, "range")
	if err != nil {
		fmt.Println(err)
	}
}

//将所有该目录下的原数据文件转换格式并写到结果文件夹
func AllFilesToWrite(sourceDir, resultsDir, types string) error {
	err := filepath.Walk(sourceDir, func(path string, f os.FileInfo, err error) error {
		// fileList = append(fileList,file)
		_, file := filepath.Split(path)
		sourceFile := path
		resultFile := resultsDir + file + "Addr.json"
		if file != "" {
			ChangeFormatAndWrite(sourceFile, resultFile, types)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil

}

//单个文件的处理
func ChangeFormatAndWrite(sourceFile, resultFile, types string) error {
	// var results []string

	file1, err := os.Open(sourceFile)
	file2, err := os.Create(resultFile)

	_, file := filepath.Split(sourceFile)

	bufread := bufio.NewReader(file1)
	scanner := bufio.NewScanner(bufread)

	var buf bytes.Buffer

	type rg struct {
		Start uint32 `json:"start"`
		End   uint32 `json:"end"`
	}

	buf.WriteString("/vnf-agent/vpp1/config/vpp/v1/ippool/" + file + "\n")

	buf.WriteString(`{ "name": "` + file + `"`)
	buf.WriteByte(',')

	buf.WriteString(`"desc": "ip ranges of ` + file + `"`)
	buf.WriteByte(',')

	buf.WriteString(`"refcnt": `)
	buf.WriteByte('0')
	buf.WriteByte(',')

	buf.WriteString(`"type": "ISP"`)
	buf.WriteByte(',')

	buf.WriteString(`"ranges": [`)

	var firstIP, lastIP net.IP

	for scanner.Scan() {
		line := scanner.Text()
		//cidr类型处理
		if types == "cidr" {
			_, ipnet, _ := net.ParseCIDR(line)
			firstIP, lastIP = cidr.AddressRange(ipnet)
			// firstStr :=first.String()
			// lastStr :=last.String()

		} else if types == "range" { //range类型处理
			r := strings.Split(line, "-")
			firstStr := strings.TrimSpace(r[0])
			lastStr := strings.TrimSpace(r[1])
			firstIP = net.ParseIP(firstStr)
			lastIP = net.ParseIP(lastStr)
		}

		firstU32 := IPtoU32(firstIP)
		lastU32 := IPtoU32(lastIP)
		r := rg{
			Start: firstU32,
			End:   lastU32,
		}

		jsonr, err := json.Marshal(r)
		buf.Write(jsonr)
		buf.WriteByte(',')
		if err != nil {
			return err
		}

	}

	buf.Truncate(buf.Len() - 1)
	buf.WriteString("]}")
	file2.Write(buf.Bytes())

	return err

}

func IPtoU32(ipnr net.IP) uint32 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum uint32

	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)

	return sum

}
