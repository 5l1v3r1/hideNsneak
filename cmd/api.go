// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rmikehodges/hideNsneak/deployer"

	"github.com/spf13/cobra"
)

var apiProvider string
var targetURI string
var apiIndices string

var api = &cobra.Command{
	Use:   "api",
	Short: "API Gateway",
	Long:  `API Gateway parent command. Use -h to see options`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'api --help' for usage.")
	},
}

var apiDeploy = &cobra.Command{
	//TODO: need to trim spaces
	Use:   "deploy",
	Short: "deploy an API Gateway",
	Long:  `deploy an API Gateway`,
	Args: func(cmd *cobra.Command, args []string) error {
		deployer.InitializeTerraformFiles()
		if !deployer.ProviderCheck(instanceProviders) {
			return fmt.Errorf("invalid providers specified: %v", instanceProviders)
		}
		r, _ := regexp.Compile(`http[s]{0,1}\:\/\/[a-zA-Z0-9]+\.[a-z]+`)
		if !r.MatchString(targetURI) {
			return fmt.Errorf("the target uri is formatted incorrectly i.e. http[s]://example.com")
		}
		if string(targetURI[len(targetURI)-1]) != "/" {
			targetURI = targetURI + "/"
		}

		return nil

	},
	Run: func(cmd *cobra.Command, args []string) {

		marshalledState := deployer.TerraformStateMarshaller()
		wrappers := deployer.CreateWrappersFromState(marshalledState)

		apiProvider = strings.ToUpper(apiProvider)

		wrappers = deployer.APIDeploy(apiProvider, targetURI, wrappers)

		mainFile := deployer.CreateMasterFile(wrappers)

		deployer.CreateTerraformMain(mainFile)

		deployer.TerraformApply()
	},
}

var apiDestroy = &cobra.Command{
	Use:   "destroy",
	Short: "destroys an API Gateway",
	Long:  `destroys an API Gateway via an index`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := deployer.IsValidNumberInput(apiIndices)

		if err != nil {
			return err
		}

		expandedNumIndex := deployer.ExpandNumberInput(apiIndices)

		err = deployer.ValidateNumberOfInstances(expandedNumIndex, "api")

		if err != nil {
			return err
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		marshalledState := deployer.TerraformStateMarshaller()
		list := deployer.ListAPIs(marshalledState)

		var namesToDelete []string

		expandedNumIndex := deployer.ExpandNumberInput(instanceDestroyIndices)

		for _, numIndex := range expandedNumIndex {
			namesToDelete = append(namesToDelete, list[numIndex].Name)
		}

		deployer.TerraformDestroy(namesToDelete)
		if len(apiIndices) > 2 {
			fmt.Println("Destroying multiple API gateways a few minutes...")
		}
		return
	},
}

var apiList = &cobra.Command{
	Use:   "list",
	Short: "list of API Gateways",
	Long:  `list API Gateways and show their target URIs, invoke URIs, providers and names`,
	Run: func(cmd *cobra.Command, args []string) {
		marshalledState := deployer.TerraformStateMarshaller()

		apiList := deployer.ListAPIs(marshalledState)

		for index, item := range apiList {
			fmt.Print(index)
			fmt.Println(item.String())
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(api)
	api.AddCommand(apiDeploy, apiList, apiDestroy)

	apiDeploy.PersistentFlags().StringVarP(&apiProvider, "provider", "p", "", "[Required] the cloud provider to use. AWS is the only currently supported provider")
	apiDeploy.MarkPersistentFlagRequired("providers")

	apiDeploy.PersistentFlags().StringVarP(&targetURI, "target", "t", "", "[Required] the target URL of the endpoint i.e. https://google.com")
	apiDeploy.MarkPersistentFlagRequired("target")

	apiDestroy.PersistentFlags().StringVarP(&apiIndices, "input", "i", "", "[Required] the indices of API Gateways to destroy i.e 0,2-4")
	apiDestroy.MarkPersistentFlagRequired("input")
}
