/*
Copyright © 2022 Polygon <engineering@polygon.technology>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// loadtestCmd represents the loadtest command
var loadtestCmd = &cobra.Command{
	Use:   "loadtest",
	Short: "A simple script for quickly running a load test",
	Long:  `Loadtest gives us a simple way to run a generic load test against an eth/EVM style json RPC endpoint`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("loadtest called")
	},
}

type (
	loadTestParams struct {
		Requests    *int64
		Concurrency *int64
		TimeLimit   *int64
		Timeout     *int64
		PostFile    *string
		ContentType *string
		Verbosity   *int64
		Auth        *string
		Proxy       *string
		ProxyAuth   *string
		KeepAlive   *bool
	}
)

func init() {
	rootCmd.AddCommand(loadtestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadtestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loadtestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	ltp := new(loadTestParams)
	// Apache Bench Parameters
	ltp.Requests = loadtestCmd.PersistentFlags().Int64P("requests", "n", 1, "Number of requests to perform for the benchmarking session. The default is to just perform a single request which usually leads to non-representative benchmarking results.")
	ltp.Concurrency = loadtestCmd.PersistentFlags().Int64P("concurrency", "c", 1, "Number of multiple requests to perform at a time. Default is one request at a time.")
	ltp.TimeLimit = loadtestCmd.PersistentFlags().Int64P("time-limit", "t", -1, "Maximum number of seconds to spend for benchmarking. Use this to benchmark within a fixed total amount of time. Per default there is no timelimit.")
	ltp.Timeout = loadtestCmd.PersistentFlags().Int64P("timeout", "s", 30, "Maximum number of seconds to wait before the socket times out. Default is 30 seconds.")
	ltp.PostFile = loadtestCmd.PersistentFlags().StringP("post-file", "p", "", "File containing data to POST.")
	ltp.ContentType = loadtestCmd.PersistentFlags().StringP("content-type", "T", "application/json", "Content-type header to use for POST/PUT data, eg. application/x-www-form-urlencoded.")
	// https://logging.apache.org/log4j/2.x/manual/customloglevels.html
	ltp.Verbosity = loadtestCmd.PersistentFlags().Int64P("verbosity", "v", 200, "0 - Silent\n100 Fatals\n200 Errors\n300 Warnings\n400 INFO\n500 Debug\n600 Trace")
	ltp.Auth = loadtestCmd.PersistentFlags().StringP("auth", "A", "", "username:password used for www basic auth")
	ltp.Proxy = loadtestCmd.PersistentFlags().StringP("proxy", "X", "", "proxy:port combination to use a proxy server for the requests.")
	ltp.ProxyAuth = loadtestCmd.PersistentFlags().StringP("proxy-auth", "P", "", "Supply BASIC Authentication credentials to a proxy en-route. The username and password are separated by a single : and sent on the wire base64 encoded. The string is sent regardless of whether the proxy needs it (i.e., has sent an 407 proxy authentication needed).")
	ltp.KeepAlive = loadtestCmd.PersistentFlags().BoolP("keep-alive", "k", true, "Enable the HTTP KeepAlive feature, i.e., perform multiple requests within one HTTP session.")

}
