// Package cmd
/*
Copyright © 2021 Joseph Cracchiolo <joe@cracchiolo.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// netCmd represents the net command
var (
	custCmd = &cobra.Command{
		Use:   "customer",
		Short: "Manage Customers (add delete show)",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(`You must specify "add", "delete" or "show"`)
			return errMissingArgument
		},
	}
	custAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add New Customer",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("customer add called")
		},
	}
	custShowCmd = &cobra.Command{
		Use:   "show [HANDLE]",
		Short: "Show Customer",
		RunE:  showCustomer,
	}
	custDeleteCmd = &cobra.Command{
		Use:   "delete [HANDLE]",
		Short: "Delete Customer",
		RunE:  deleteCustomer,
	}
)

func init() {
	rootCmd.AddCommand(custCmd)
	custCmd.AddCommand(custAddCmd)
	custCmd.AddCommand(custShowCmd)
	custCmd.AddCommand(custDeleteCmd)
}

// Customer schema.  See https://www.arin.net/resources/manage/regrws/payloads/#introduction for more info
type Customer struct {
	XMLName      xml.Name `xml:"customer"`
	Text         string   `xml:",chardata"`
	Xmlns        string   `xml:"xmlns,attr"`
	CustomerName string   `xml:"customerName"`
	Iso31661     struct {
		Text  string `xml:",chardata"`
		Name  string `xml:"name"`
		Code2 string `xml:"code2"`
		Code3 string `xml:"code3"`
		E164  string `xml:"e164"`
	} `xml:"iso3166-1"`
	Handle        string `xml:"handle"`
	StreetAddress struct {
		Line []struct {
			Text   string `xml:",chardata"`
			Number string `xml:"number,attr"`
		} `xml:"line"`
	} `xml:"streetAddress"`
	City       string `xml:"city"`
	Iso31662   string `xml:"iso3166-2"`
	PostalCode string `xml:"postalCode"`
	Comment    struct {
		Line []struct {
			Text   string `xml:",chardata"`
			Number string `xml:"number,attr"`
		} `xml:"line"`
	} `xml:"comment"`
	ParentOrgHandle  string `xml:"parentOrgHandle"`
	RegistrationDate string `xml:"registrationDate"`
	PrivateCustomer  string `xml:"privateCustomer"`
}

func (c *Customer) String() string {
	a := make([]string, 0)
	a = append(a, fmt.Sprintf("Handle: %s", c.Handle))
	a = append(a, fmt.Sprintf("Name: %s", c.CustomerName))
	a = append(a, fmt.Sprintf("Address: %s %s", c.StreetAddress.Line[0].Number, c.StreetAddress.Line[0].Text))
	return strings.Join(a, "\n")
}

func showCustomer(cmd *cobra.Command, args []string) error {
	var (
		err  error
		info Customer
	)
	if len(args) == 0 {
		return fmt.Errorf("missing handle")
	}
	err = restGet(context.Background(), makeUrl("rest/customer", args[0]), &info)
	if err != nil {
		return err
	}
	fmt.Println(info.String())
	return nil
}

func deleteCustomer(cmd *cobra.Command, args []string) error {
	var (
		err  error
		info Customer
	)
	if len(args) == 0 {
		return fmt.Errorf("missing handle")
	}
	err = restDelete(context.Background(), makeUrl("rest/customer", args[0]), &info)
	if err != nil {
		return err
	}
	fmt.Printf("Customer %s Deleted\n", info.Handle)
	if flagVerbose {
		fmt.Println(info.String())
	}
	return nil
}
