// +build linux

package ps

import (
	"fmt"
	"strconv"
	"io/ioutil"
	"strings"
)

// Refresh reloads all the data associated with this process.
func (p *UnixProcess) Refresh() error {
	statPath := fmt.Sprintf("/proc/%d/stat", p.pid)
	dataBytes, err := ioutil.ReadFile(statPath)
	if err != nil {
		return err
	}

	// First, parse out the image name
	data := string(dataBytes)	
	fields := strings.Fields(data)		
	vszBytes, err := strconv.ParseInt(fields[22], 10, 64)
	if err != nil {
		p.vsize = -1
	} else {
		p.vsize = vszBytes/1024/1024
	}
	
	binStart := strings.IndexRune(data, '(') + 1	
	binEnd := strings.IndexRune(data[binStart:], ')')
	p.binary = data[binStart : binStart+binEnd]
    
	// Move past the image name and start parsing the rest
	data = data[binStart+binEnd+2:]
	_, err = fmt.Sscanf(data,
		"%c %d %d %d",
		&p.state,
		&p.ppid,
		&p.pgrp,
		&p.sid)

	
	return err
}
