package main

import (
	"flag"
	"fmt"
	"ipfs-monitor/command"
	"ipfs-monitor/pinner"
	"ipfs-monitor/reporter"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron"
	"github.com/takama/daemon"
)

const (
	name        = "ipfs-monitor"
	description = "Monitor IPFS node and report status to IPHash server."
)

var stdlog, errlog *log.Logger

var ipfs_base_url = flag.String("ipfs_base_url", "http://127.0.0.1:5001", "Base URL of IPFS API, default value is http://127.0.0.1:5001")
var server_url = flag.String("server_url", "http://127.0.0.1", "Server URL for reporting status, default value is http://127.0.0.1")
var cron_expr = flag.String("cron_expr", "0 0/30 * * * *", "Cron expression for report IPFS node status")

// Service is the daemon service struct
type Service struct {
	daemon.Daemon
}

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {
	usage := "Usage: ipfs_monitor install | remove | start | stop | status"
	if len(os.Args) == 2 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}
	stdlog.Println("IPFS monitor starting...")
	stdlog.Printf("Use IPFS base URL: %s\n", *ipfs_base_url)
	stdlog.Printf("Use server URL: %s\n", *server_url)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	c := cron.New()
	c.AddFunc(*cron_expr, func() { reporter.Report() })
	c.Start()
	pinner.PinService()
	killSignal := <-interrupt
	stdlog.Println("Got signal:", killSignal)
	return "Service exited", nil
}

func main() {
	flag.Parse()
	command.Base_URL = *ipfs_base_url
	reporter.Report_URL = *server_url

	srv, err := daemon.New(name, description)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}

	fmt.Println(status)
}

func main1() {
	fmt.Println(len(os.Args))
	//r, _ := reporter.Report()
	// fmt.Printf("%d\n", len("abc"))
	// err := ioutil.WriteFile("/tmp/abc.txt", []byte("456456346"), 0644)
	// if os.IsNotExist(err) {
	// 	fmt.Println("123")
	// }
	// stat, fserr := disk.Usage("/home/aurawing")
	// if fserr != nil {
	// 	return
	// }
	// fmt.Println("Used:", stat.Used)

	// fs := syscall.Statfs_t{}
	// err := syscall.Statfs("/home", &fs)
	// if err != nil {
	// 	return
	// }
	// fmt.Println(fs.Bfree * uint64(fs.Bsize))
	// fmt.Println(fs.Blocks * uint64(fs.Bsize))
}
