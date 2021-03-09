/*
Copyright 2021 ABSA Group Limited

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
package cmd

import (
	"net/url"
	"os"

	"github.com/AbsaOSS/golic/impl/inject"
	"github.com/spf13/cobra"
)

var injectOptions inject.Options

var injectCmd = &cobra.Command{
	Use:   "inject",
	Short: "",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(injectOptions.LicIgnore); os.IsNotExist(err) {
			logger.Error().Msgf("invalid license path '%s'",injectOptions.LicIgnore)
			_ = cmd.Help()
			os.Exit(0)
		}
		if _,err := url.Parse(injectOptions.ConfigURL); err != nil {
			logger.Error().Msgf("invalid config.yaml url '%s'",injectOptions.ConfigURL)
			_ = cmd.Help()
			os.Exit(0)
		}
		i := inject.New(ctx, injectOptions)
		Command(i).MustRun()
	},
}

func init() {
	injectCmd.Flags().StringVarP(&injectOptions.LicIgnore, "licignore", "l", ".licignore", ".licignore path")
	injectCmd.Flags().StringVarP(&injectOptions.Template, "template", "t", "apache2", "license key")
	injectCmd.Flags().StringVarP(&injectOptions.Copyright, "copyright", "c", "2021 MyCompany",
		"company initials entered into license")
	injectCmd.Flags().BoolVarP(&injectOptions.Dry, "dry", "d", false, "dry run")
	injectCmd.Flags().StringVarP(&injectOptions.ConfigURL, "config-url", "u", "https://raw.githubusercontent.com/AbsaOSS/golic/main/config.yaml", "config URL")
	rootCmd.AddCommand(injectCmd)
}