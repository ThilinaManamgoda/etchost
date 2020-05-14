/*
 * Copyright (c) 2020 Thilina Manamgoda.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package parser parses the hosts file.
package parser

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strings"
)

// eTCHosts represents the host mapping.
type eTCHosts map[string][]string

// Parser represents the Host file Parser.
type Parser struct {
	Path          string
	content       []byte
	parsedContent eTCHosts
}

type hostEntry struct {
	ip      string
	mapping string
	comment string
}

// parse parses the given hosts file.
func parse(data []byte) (eTCHosts, error) {
	hostsMap := map[string][]string{}
	for _, line := range strings.Split(strings.Trim(string(data), " \t\r\n"), "\n") {
		line = strings.Replace(strings.Trim(line, " \t"), "\t", " ", -1)
		if len(line) == 0 || line[0] == ';' || line[0] == '#' {
			continue
		}
		pieces := strings.SplitN(line, " ", 2)
		if len(pieces) > 1 && len(pieces[0]) > 0 {
			if names := strings.Fields(pieces[1]); len(names) > 0 {
				if _, ok := hostsMap[pieces[0]]; ok {
					hostsMap[pieces[0]] = append(hostsMap[pieces[0]], names...)
				} else {
					hostsMap[pieces[0]] = names
				}
			}
		}
	}
	return hostsMap, nil
}

// readFile reads the given file.
func readFile(path string) ([]byte, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (p *Parser) addNewHostMapping(entry hostEntry) error {
	var file, err = os.OpenFile(p.Path, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return errors.Wrap(err, "Unable to open file")
	}
	defer func() {
		file.Close()
	}()

	if entry.comment != "" {
		_, err = file.WriteString("# " + entry.comment + "\n")
		if err != nil {
			return errors.Wrap(err, "Unable to add the comment")
		}
	}
	_, err = file.WriteString(entry.mapping + "\n")
	if err != nil {
		return errors.Wrap(err, "Unable to add the host mapping")
	}
	err = file.Sync()
	if err != nil {
		return errors.Wrap(err, "Unable save")
	}
	return nil
}

func (p *Parser) isIPExists(ip string) bool {
	_, ok := p.parsedContent[ip]
	return ok
}

func (p *Parser) isDomainsExist(domains []string) (bool, string) {
	for _, domain := range domains {
		exists, _, _ := p.isDomainExists(domain)
		if exists {
			return true, domain
		}
	}
	return false, ""
}

func (p *Parser) isDomainExists(domain string) (bool, string, int) {
	for i, val := range p.parsedContent {
		for dIndex, domainT := range val {
			if domain == domainT {
				return true, i, dIndex
			}
		}
	}
	return false, "", 0
}

// AddNewMapping method adds a new host entry.
func (p *Parser) AddNewMapping(ip string, domains []string, comment string) error {
	exists, d := p.isDomainsExist(domains)
	if exists {
		return errors.New(fmt.Sprintf("%s already exists", d))
	}
	if p.isIPExists(ip) {
		hosts := append(domains, p.parsedContent[ip]...)
		err := p.updateHostMapping(false, hostEntry{
			ip:      ip,
			mapping: constructMapping(ip, hosts),
			comment: comment,
		})
		if err != nil {
			return errors.Wrap(err, "Unable update the mapping")
		}
	} else {
		err := p.addNewHostMapping(hostEntry{
			mapping: constructMapping(ip, domains),
			comment: comment,
		})
		if err != nil {
			return errors.Wrap(err, "Unable add the mapping")
		}
	}
	return nil
}

// RemoveDomainFromHostMapping removes a domain from an entry.
func (p *Parser) RemoveDomainFromHostMapping(domain string) error {
	exists, ip, dIndex := p.isDomainExists(domain)
	if !exists {
		return errors.New(fmt.Sprintf("%s doesn't exists", domain))
	}
	mapping := p.parsedContent[ip]
	if len(mapping) == 1 {
		err := p.RemoveHostMapping(ip)
		if err != nil {
			return err
		}
	} else {
		copy(mapping[dIndex:], mapping[dIndex+1:])
		mapping[len(mapping)-1] = ""
		mapping = mapping[:len(mapping)-1]
		err := p.updateHostMapping(false, hostEntry{
			ip:      ip,
			mapping: constructMapping(ip, mapping),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveHostMapping removes an entry.
func (p *Parser) RemoveHostMapping(ip string) error {
	err := p.updateHostMapping(true, hostEntry{ip: ip})
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) updateHostMapping(remove bool, entry hostEntry) error {
	lines := strings.Split(string(p.content), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, entry.ip) {
			if remove {
				index := i
				if ((i - 1) > 0) && strings.HasPrefix(lines[i-1], "#") {
					copy(lines[index-1:], lines[index+1:])
					lines[len(lines)-1] = ""
					lines[len(lines)-2] = ""
					lines = lines[:len(lines)-2]
				} else {
					copy(lines[index:], lines[index+1:])
					lines[len(lines)-1] = ""
					lines = lines[:len(lines)-1]
				}
			} else {
				lines[i] = entry.mapping
			}
		}
	}
	output := strings.Join(lines, "\n")
	err := ioutil.WriteFile(p.Path, []byte(output), 0644)
	if err != nil {
		return err
	}
	return nil
}

// Init initialize the parser.
func (p *Parser) Init() error {
	bs, err := readFile(p.Path)
	if err != nil {
		return errors.Wrap(err, "Unable to parse read the file")
	}
	p.content = bs

	p.parsedContent, err = parse(bs)
	if err != nil {
		return errors.Wrap(err, "Unable to parse host file")
	}
	return nil
}

func constructMapping(ip string, domains []string) string {
	return ip + " " + strings.Join(domains, " ")
}
