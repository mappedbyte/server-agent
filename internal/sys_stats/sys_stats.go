package sys_stats

import (
	"context"
	"fmt"
	"github.com/mappedbyte/server-agent/proto"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

func ReportSysStats(ip string, port, factor int) {
	ipAddress := fmt.Sprintf("%s:%d", ip, port)
	conn, err := grpc.Dial(ipAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Connection failed, please check IP and port")
		return
	}
	statsServiceClient := proto.NewSysStatsServiceClient(conn)
	ticker := time.NewTicker(time.Duration(factor) * time.Second)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func(statsServiceClient proto.SysStatsServiceClient, ticker *time.Ticker, waitGroup *sync.WaitGroup) {

		defer waitGroup.Done()
		for {
			select {
			case <-ticker.C:
				statsServiceClient.ReportSysStats(context.Background(), GetSysStats())
				statsServiceClient.ReportHostInfo(context.Background(), GetHostInfo())
			}
		}

	}(statsServiceClient, ticker, &waitGroup)

	waitGroup.Wait()
}
