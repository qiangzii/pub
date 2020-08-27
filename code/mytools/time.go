// Copyright (c) 2019 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package descriptor

import (
	"bufio"
	"math"
	"strconv"

	// "log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"go.ligato.io/cn-infra/v2/logging"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	"go.ligato.io/vpp-agent/v3/plugins/time/descriptor/adapter"
	"go.ligato.io/vpp-agent/v3/proto/ligato/time"
)

const (
	// TimeDescriptorName is the name of the descriptor time.
	TimeDescriptorName = "time"
)

// TimeDescriptor is only a user of a descriptor, which can be used
// as a starting point to build a new descriptor from.
type TimeDescriptor struct {
	log logging.Logger
}

// NewTimeDescriptor creates a new instance of the descriptor.
func NewTimeDescriptor(log logging.PluginLogger) *kvs.KVDescriptor {
	// descriptors are supposed to be stateless, so use the structure only
	// as a context for things that do not change once the descriptor is
	// constructed - e.g. a reference to the logger to use within the descriptor
	descrCtx := &TimeDescriptor{
		log: log.NewLogger("time-descriptor"),
	}

	// use adapter to convert typed descriptor into generic descriptor API
	typedDescr := &adapter.TimeDescriptor{
		Name:          TimeDescriptorName,
		NBKeyPrefix:   time.ModelTime.KeyPrefix(),
		ValueTypeName: time.ModelTime.ProtoName(),
		KeySelector:   time.ModelTime.IsKeyValid,
		KeyLabel:      time.ModelTime.StripKeyPrefix,
		// ValueComparator: descrCtx.EquivalentValues,
		// Validate: descrCtx.Validate,
		Create: descrCtx.Create,
		Delete: descrCtx.Delete,
		Update: descrCtx.Update,
		// UpdateWithRecreate: descrCtx.UpdateWithRecreate,
		// Retrieve:           descrCtx.Retrieve,
		// IsRetriableFailure: descrCtx.IsRetriableFailure,
		DerivedValues: descrCtx.DerivedValues,
		Dependencies:  descrCtx.Dependencies,
	}
	return adapter.NewTimeDescriptor(typedDescr)
}

// // Validate validates value before it is applied.
// func (d *TimeDescriptor) Validate(key string, value *time.Time) error {
// 	return nil
// }

// Create creates new value.
func (d *TimeDescriptor) Create(key string, value *time.Time) (metadata interface{}, err error) {
	d.log.Debugf("Descriptor Add Time %v", value)

	if value.Enable {
		//开启ntp模式：修改ntp的配置文件中的server和interval
		newIntervar := int64(math.Floor(math.Log2(float64(value.Ntp.Interval * 60)))) //初期传的是换算后的结果，后期直接传指数项
		mainserver := strings.Join([]string{"server", value.Ntp.Mainserver, "prefer", "minpoll", strconv.FormatInt(value.Ntp.Interval, 10), "maxpoll", strconv.FormatInt(newIntervar, 10), "\n"}, " ")
		secondserver := strings.Join([]string{"server", value.Ntp.Secondserver, "minpoll", strconv.FormatInt(value.Ntp.Interval, 10), "maxpoll", strconv.FormatInt(newIntervar, 10), "\n"}, " ")

		err := UpdateNtpConfig(mainserver, secondserver)
		if err != nil {
			d.log.Error(err)
			return nil, err
		}
		cmd1 := exec.Command("systemctl", "stop", "ntpd")
		err = cmd1.Run()
		if err != nil {
			return nil, err
		}

		cmd2 := exec.Command("systemctl", "start", "ntpd")
		err = cmd2.Run()
		if err != nil {
			return nil, err
		}
	} else {
		//开启手动配置，配置系统时间和日期,需要先关闭ntp服务
		// timedate := strings.Join([]string{value.Date, value.Time}, " ")
		// timezone := value.Timezone
		// err := SetSysTime(timedate, timezone)
		// if err != nil {
		// 	d.log.Error(err)
		// 	return nil, err
		// }
	}
	return nil, nil
}

// Update updates existing value.
func (d *TimeDescriptor) Update(key string, oldTime, newTime *time.Time, oldMetadata interface{}) (newMetadata interface{}, err error) {
	d.log.Debugf("Descriptor update Time %v", newTime)

	if newTime.Enable {
		//开启ntp模式：修改ntp的配置文件中的server和interval
		newIntervar := int64(math.Floor(math.Log2(float64(newTime.Ntp.Interval * 60)))) //初期传的是换算后的结果，后期直接传指数项
		mainserver := strings.Join([]string{"server", newTime.Ntp.Mainserver, "prefer", "minpoll", strconv.FormatInt(newTime.Ntp.Interval, 10), "maxpoll", strconv.FormatInt(newIntervar, 10), "\n"}, " ")
		secondserver := strings.Join([]string{"server", newTime.Ntp.Secondserver, "minpoll", strconv.FormatInt(newTime.Ntp.Interval, 10), "maxpoll", strconv.FormatInt(newIntervar, 10), "\n"}, " ")

		err := UpdateNtpConfig(mainserver, secondserver)
		if err != nil {
			d.log.Error(err)
			return nil, err
		}

		cmd1 := exec.Command("systemctl", "stop", "ntpd")
		err = cmd1.Run()
		if err != nil {
			return nil, err
		}

		cmd2 := exec.Command("systemctl", "start", "ntpd")
		err = cmd2.Run()
		if err != nil {
			return nil, err
		}
	} else {
		//开启手动配置，配置系统时间和日期,需要先关闭ntp服务
		// timedate := strings.Join([]string{newTime.Date, newTime.Time}, " ")
		// timezone := newTime.Timezone
		// err := SetSysTime(timedate, timezone)
		// if err != nil {
		// 	d.log.Error(err)
		// 	return nil, err
		// }
	}
	return nil, nil
}

// // EquivalentValues compares two revisions of the same value for equality.
// func (d *TimeDescriptor) EquivalentValues(key string, old, new *time.Time) bool {
// 	// compare **non-primary** attributes here (none in the ValueUser)
// 	return true
// }

// Delete removes an existing value.
func (d *TimeDescriptor) Delete(key string, value *time.Time, metadata interface{}) error {
	return nil
}

// Retrieve retrieves values from SB.
// func (d *TimeDescriptor) Retrieve(correlate []adapter.TimeKVWithMetadata) (retrieved []adapter.TimeKVWithMetadata, err error) {
// 	return retrieved, nil
// }

// IsRetriableFailure returns true if the given error, returned by one of the CRUD
// operations, can be theoretically fixed by merely repeating the operation.
// func (d *TimeDescriptor) IsRetriableFailure(err error) bool {
// 	return true
// }

// DerivedValues breaks the value into multiple part handled/referenced
// separately.
func (d *TimeDescriptor) DerivedValues(key string, value *time.Time) (derived []kvs.KeyValuePair) {
	return derived
}

// Dependencies lists dependencies of the given value.
func (d *TimeDescriptor) Dependencies(key string, value *time.Time) (deps []kvs.Dependency) {
	return deps
}

// SetSysTime set time manual
func SetSysTime(timedate, timezone string) error {
	//TODO: 更新ntp服务管理方式
	//关闭ntp服务
	// cmd1 := exec.Command("timedatectl", "set-ntp", "false")
	cmd1 := exec.Command("systemctl", "stop", "ntpd")
	err := cmd1.Run()
	if err != nil {
		return err
	}

	//配置系统时间
	// timedate = fmt.Sprintf("%s", timedate)
	// cmd2 := exec.Command("timedatectl", "set-time", timedate)
	cmd2 := exec.Command("date", "-s", timedate) //虽然linux命令中需要双引号，但是这里作为参数时不要加引号
	err = cmd2.Run()
	if err != nil {
		return err
	}

	if timezone != "" {
		cmd3 := exec.Command("timedatectl", "set-timezone", timezone)
		err = cmd3.Run()
		if err != nil {
			return err
		}
	}

	return nil

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
