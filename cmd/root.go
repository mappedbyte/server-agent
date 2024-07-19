package cmd

import (
	"github.com/mappedbyte/server-agent/internal/sys_stats"
	"github.com/spf13/cobra"
	"log"
	"net"
)

var host string
var port int
var factor int

func init() {
	// 添加标志到根命令
	rootCmd.Flags().StringVarP(&host, "host", "H", "", "The host IP address")
	rootCmd.Flags().IntVarP(&port, "port", "p", -1, "The port number")
	rootCmd.Flags().IntVarP(&factor, "factor", "f", -1, "The factor value")
}

// 定义根命令
var rootCmd = &cobra.Command{
	Use:   "server-agent",
	Short: " push service information to the server.",
	Long:  `This is a tool that uses RPC (Remote Procedure Call) to push service information to the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		ip := net.ParseIP(host)
		if ip == nil {
			log.Fatalln("Invalid IP address")
			return
		}
		if port == -1 || port <= 0 {
			log.Fatalln("Please enter the correct port")
			return
		}
		if factor == -1 || factor <= 0 {
			log.Fatalln("Please enter the correct scheduling unit as seconds")
			return
		}
		log.Printf("Host: %s, Port: %d, Factor: %d\n", host, port, factor)
		sys_stats.ReportSysStats(host, port, factor)
	},
}

func Execute() error {
	return rootCmd.Execute()
}
