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
	"github.com/spf13/cobra"
	"strings"
)

// netCmd represents the net command
var (
	netCmd = &cobra.Command{
		Use:   "net",
		Short: "Manage networks (add delete)",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(`You must specify "add" or "delete"`)
			return errMissingArgument
		},
	}
	netAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add New Network",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("net add called")
		},
	}
	netShowCmd = &cobra.Command{
		Use:   "show [HANDLE]",
		Short: "Show Network",
		Long:  ``,
		RunE:  showNetwork,
	}
	netDelCmd = &cobra.Command{
		Use:   "delete [HANDLE]",
		Short: "Delete Network",
		RunE:  deleteNetwork,
	}
)

func init() {
	rootCmd.AddCommand(netCmd)
	netCmd.AddCommand(netShowCmd)
	netCmd.AddCommand(netAddCmd)
	netCmd.AddCommand(netDelCmd)
	netDelCmd.Flags().Bool("custdel", false, "delete customer if network was Reassigned")
}

type NetBlock struct {
	XMLName      xml.Name `xml:"netBlock"`
	Xmlns        string   `xml:"xmlns,attr"`
	Type         string   `xml:"type"`
	Description  string   `xml:"description"`
	StartAddress string   `xml:"startAddress"`
	EndAddress   string   `xml:"endAddress"`
	CidrLength   string   `xml:"cidrLength"`
}

type Network struct {
	XMLName xml.Name `xml:"net"`
	Xmlns   string   `xml:"xmlns,attr"`
	Version string   `xml:"version"`
	Comment struct {
		Line []struct {
			Text   string `xml:",chardata"`
			Number string `xml:"number,attr"`
		} `xml:"line"`
	} `xml:"comment"`
	RegistrationDate string `xml:"registrationDate"`
	OrgHandle        string `xml:"orgHandle"`
	Handle           string `xml:"handle"`
	NetBlocks        struct {
		NetBlock []NetBlock `xml:"netBlock"`
	} `xml:"netBlocks"`
	CustomerHandle  string `xml:"customerHandle"`
	ParentNetHandle string `xml:"parentNetHandle"`
	NetName         string `xml:"netName"`
	OriginASes      struct {
		OriginAS []string `xml:"originAS"`
	} `xml:"originASes"`
	PocLinks struct {
		PocLinkRef []string `xml:"pocLinkRef"`
	} `xml:"pocLinks"`
}

var (
	NETTYPES = map[string]string{
		"A":  "Reallocation",
		"AF": "AFRINIC allocated",
		"AP": "APNIC allocated",
		"AR": "ARIN allocated",
		"AV": "ARIN early reservation",
		"DA": "Direct Allocation",
		"DS": "Direct Assignment",
		"FX": "AFRINIC transferred",
		"IR": "IANA reserved",
		"IU": "IANA special use",
		"LN": "LACNIC allocated",
		"LX": "LACNIC transferred",
		"PV": "APNIC early reservation",
		"PX": "APNIC early registration",
		"RD": "RIPE NCC allocated",
		"RN": "RIPE allocated",
		"RV": "RIPE early reservation",
		"RX": "RIPE NCC Transferred",
		"S":  "Reassigned",
	}
)

func (n *Network) String() string {
	a := make([]string, 0)
	a = append(a, fmt.Sprintf("%s: ", n.Handle))
	for _, b := range n.NetBlocks.NetBlock {
		a = append(a, fmt.Sprintf("\t%s: %s-%s [/%s] ", b.Description, b.StartAddress, b.EndAddress, b.CidrLength))
	}
	return strings.Join(a, "\n")
}

func showNetwork(cmd *cobra.Command, args []string) error {
	var (
		err  error
		info Network
	)
	if len(args) == 0 {
		return fmt.Errorf("missing handle")
	}
	err = restGet(context.Background(), makeUrl("rest/net", args[0]), &info)
	if err != nil {
		return err
	}
	fmt.Println(info.String())
	return nil
}

func deleteNetwork(cmd *cobra.Command, args []string) error {
	var (
		err  error
		info Network
	)
	if len(args) == 0 {
		return fmt.Errorf("missing handle")
	}
	// The DELETE request below will return a Ticket and we want the network info, so we get the network first
	err = restGet(context.Background(), makeUrl("rest/net", args[0]), &info)
	if err != nil {
		return err
	}
	err = restDelete(context.Background(), makeUrl("rest/net", args[0]), nil)
	if err != nil {
		return err
	}
	fmt.Printf("Network %s Deleted\n", info.Handle)
	if flagVerbose {
		fmt.Println(info.String())
	}
	if len(info.CustomerHandle) > 0 {
		if custdel, err := cmd.Flags().GetBool("custdel"); err == nil && custdel {
			if len(info.NetBlocks.NetBlock) > 0 && info.NetBlocks.NetBlock[0].Type == "S" {
				var cust Customer
				err = restDelete(context.Background(), makeUrl("rest/customer", info.CustomerHandle), &cust)
				if err != nil {
					return err
				}
				fmt.Printf("Customer %s Deleted\n", info.CustomerHandle)
				if flagVerbose {
					fmt.Println(cust.String())
				}
			}
		}
	}
	return nil
}
