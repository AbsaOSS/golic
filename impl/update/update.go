package update

/*
Copyright 2022 Absa Group Limited

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/AbsaOSS/golic/utils/log"

	"github.com/denormal/go-gitignore"
	"github.com/enescakir/emoji"
	"github.com/logrusorgru/aurora"
)

type Update struct {
	opts     Options
	ctx      context.Context
	ignore   gitignore.GitIgnore
	cfg      *Config
	modified int
}

var logger = log.Log

func New(ctx context.Context, options Options) *Update {
	return &Update{
		ctx:  ctx,
		opts: options,
	}
}

func (u *Update) Run() (err error) {
	logger.Info().Msgf("%s reading %s", emoji.OpenBook, u.opts.LicIgnore)
	u.ignore, err = gitignore.NewFromFile(u.opts.LicIgnore)
	if err != nil {
		return err
	}
	logger.Info().Msgf("%s reading %s; use --verbose to see details", emoji.OpenBook, aurora.BrightCyan("master config"))
	if u.cfg, err = u.readCommonConfig(); err != nil {
		return
	}
	if _, err = os.Stat(u.opts.ConfigPath); !os.IsNotExist(err) {
		logger.Info().Msgf("%s reading %s", emoji.OpenBook, aurora.BrightCyan(u.opts.ConfigPath))
		logger.Info().Msgf("%s overriding %s with %s",
			emoji.ConstructionWorker, aurora.BrightCyan("master config"), aurora.BrightCyan(u.opts.ConfigPath))
		if u.cfg, err = u.readLocalConfig(); err != nil {
			return
		}
	} else {
		logger.Info().Msgf("%s skipping local %s", emoji.FileFolder, aurora.BrightCyan(u.opts.ConfigPath))
		err = nil
	}
	u.traverse()
	return
}

func (u *Update) String() string {
	switch u.opts.Type {
	case LicenseInject:
		return aurora.BrightCyan("inject").String()
	case LicenseRemove:
		return aurora.BrightCyan("remove").String()
	}
	return aurora.BrightRed("ERROR, unrecognised command").String()
}

func (u *Update) ExitCode() int {
	if u.opts.ModifiedExitStatus && u.modified != 0 {
		return 1
	}
	return 0
}

func read(f string) (s string, err error) {
	content, err := ioutil.ReadFile(f)
	if err != nil {
		return
	}
	// Convert []byte to string and print to screen
	return string(content), nil
}

func (u *Update) traverse() {
	skipped := 0
	visited := 0
	p := func(path string, i gitignore.GitIgnore, o Options, config *Config) (err error) {
		if !i.Ignore(path) {
			var skip bool
			symbol := ""
			cp := aurora.BrightYellow(path)
			visited++
			if err, skip = update(path, o, config); skip {
				symbol = "-> skip"
				cp = aurora.Magenta(path)
				skipped++
			}
			_, _ = emoji.Printf(" %s  %s %s  \n", emoji.Minus, cp, aurora.BrightMagenta(symbol))
		}
		return
	}

	err := filepath.Walk("./",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return p(path, u.ignore, u.opts, u.cfg)
			}
			return nil
		})
	if err != nil {
		logger.Err(err).Msg("")
	}
	u.modified = visited - skipped
	summary(skipped, visited)
}

func summary(skipped, visited int) {
	if skipped == visited {
		fmt.Printf("\n %s %v/%v %s\n\n", emoji.Ice, aurora.BrightCyan(visited-skipped), aurora.BrightWhite(visited), aurora.BrightCyan("changed"))
		return
	}
	fmt.Printf("\n %s %v/%v %s\n\n", emoji.Fire, aurora.BrightYellow(visited-skipped), aurora.BrightWhite(visited), aurora.BrightYellow("changed"))
}

func update(path string, o Options, config *Config) (err error, skip bool) {
	switch o.Type {
	case LicenseInject:
		return inject(path, o, config)
	case LicenseRemove:
		return remove(path, o, config)
	}
	return fmt.Errorf("invalid license type"), true
}

func inject(path string, o Options, config *Config) (err error, skip bool) {
	source, err := read(path)
	if err != nil {
		return err, false
	}
	rule := getRule(config, path)
	license, err := getCommentedLicense(config, o, rule)
	if err != nil {
		return err, false
	}
	// license is injected, continue
	if strings.Contains(source, license) {
		return nil, true
	}
	// split file to header and footer and extend with license
	header, footer := splitSource(source,config.Golic.Rules[rule].Under)
	if header != "" {
		header = header + "\n"
	}
	source = fmt.Sprintf("%s%s%s", header, license, footer)

	if !o.Dry {
		data := []byte(source)
		err = ioutil.WriteFile(path, data, os.ModeExclusive)
	}
	return
}

func remove(path string, o Options, config *Config) (err error, skip bool) {
	source, err := read(path)
	if err != nil {
		return err, false
	}
	rule := getRule(config, path)
	license, err := getCommentedLicense(config, o, rule)
	if err != nil {
		return err, false
	}
	if strings.Contains(source, license) {
		return RemoveFromFile(path, o, source, license, err), false
	}
	return nil, true
}

func RemoveFromFile(path string, o Options, source string, license string, err error) error {
	if !o.Dry {
		source = strings.Replace(source, license, "", 1)
		err = ioutil.WriteFile(path, []byte(source), os.ModeExclusive)
	}
	return err
}

func matchRule(config *Config, fileName string) (rule string, ok bool) {
	if _, ok = config.Golic.Rules[fileName]; ok {
		return fileName, ok
	}
	// if rule is pattern like Dockerfile*
	for k := range config.Golic.Rules {
		matched, _ := filepath.Match(k, fileName)
		if matched {
			return k, true
		}
	}
	return
}

func getCommentedLicense(config *Config, o Options, file string) (string, error) {
	var ok bool
	var template string
	var rule string
	if template, ok = config.Golic.Licenses[o.Template]; !ok {
		return "", fmt.Errorf("no license found for %s, check configuration (.golic.yaml)", o.Template)
	}
	//if _, ok =  config.Golic.Rules[rule]; !ok {
	if rule, ok = matchRule(config, file); !ok {
		return "", fmt.Errorf("no rule found for %s, check configuration (.golic.yaml)", rule)
	}
	template = strings.ReplaceAll(template, "{{copyright}}", o.Copyright)
	if config.IsWrapped(rule) {
		return fmt.Sprintf("%s\n%s%s\n",
				config.Golic.Rules[rule].Prefix,
				template,
				config.Golic.Rules[rule].Suffix),
			nil
	}
	// `\r\n` -> `\r\n #`, `\n` -> `\n #`
	content := strings.ReplaceAll(template, "\n", fmt.Sprintf("\n%s", config.Golic.Rules[rule].Prefix))
	content = strings.TrimSuffix(content, config.Golic.Rules[rule].Prefix)
	content = config.Golic.Rules[rule].Prefix + content
	// "# \n" -> "#\n" // "# \r\n" -> "#\r\n"; some environments automatically remove spaces in empty lines. This makes problems in license PR's
	cleanedPrefix := strings.TrimSuffix(config.Golic.Rules[rule].Prefix, " ")
	content = strings.ReplaceAll(content, fmt.Sprintf("%s \n", cleanedPrefix), fmt.Sprintf("%s\n", cleanedPrefix))
	content = strings.ReplaceAll(content, fmt.Sprintf("%s \r\n", cleanedPrefix), fmt.Sprintf("%s\r\n", cleanedPrefix))
	return content, nil
}

func splitSource(source string, rules []string) (header, footer string) {
	lines := strings.Split(source, "\n")
	if len(rules) == 0 {
		return "", source
	}
	for _, r := range rules {
		header, footer = findHeaderAndFooter(lines, r)
		if header != "" {
			return
		}
	}
	return
}

func findHeaderAndFooter(lines []string, match string) (header, footer string){
	for i, l := range lines {
		if isMatch(l, match) {
			header = strings.Join(lines[0:i+1], "\n")
			footer = strings.Join(lines[i+1:], "\n")
			return
		}
	}
	return "", strings.Join(lines, "\n")
}

func getRule(config *Config, path string) (rule string) {
	fileName := filepath.Base(path)
	for k := range config.Golic.Rules {
		matched, _ := filepath.Match(k, fileName)
		if matched {
			return k
		}
	}
	rule = filepath.Ext(path)
	if rule == "" {
		rule = filepath.Base(path)
	}
	return
}

func (u *Update) readLocalConfig() (*Config, error) {
	var c = &Config{}
	var rc = *u.cfg
	yamlFile, err := ioutil.ReadFile(u.opts.ConfigPath)
	if err != nil {
		return nil, nil
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, nil
	}
	for k, v := range c.Golic.Licenses {
		rc.Golic.Licenses[k] = v
	}
	for k, v := range c.Golic.Rules {
		rc.Golic.Rules[k] = v
	}
	return &rc, nil
}

func (u *Update) readCommonConfig() (c *Config, err error) {
	c = &Config{}
	err = yaml.Unmarshal([]byte(u.opts.MasterConfig), c)
	return
}
