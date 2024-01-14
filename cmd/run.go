package cmd

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/masoncitemple4/certwatch/internal/defaults"
	"github.com/masoncitemple4/certwatch/internal/scheduler"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [host] [port]",
	Short: "The run command will start the scheduler service.",
	Long:  `Start the scheduler to begin checking for tasks to run. This is meant to be run as a daemon or background service. In the context of this certwatch CLI there is really only one task. That task is to check the SSL certificates specified, log and report back the results`,
	Run: func(cmd *cobra.Command, args []string) {
		hostsfile, _ := cmd.Flags().GetString("hostsfile")
		logfile, _ := cmd.Flags().GetString("logfile")

		if hostsfile == "" && len(args) == 0 {
			log.Fatal("no hosts file or host:port specified. exiting...")
		}

		hosts := make(map[string]scheduler.Host)
		if hostsfile != "" {
			readHostsFile(hostsfile, &hosts)
		}

		if len(args) > 0 && len(args) < 3 {
			if _, ok := hosts[args[0]]; !ok {
				hosts[args[0]] = scheduler.Host{
					Hostname: args[0],
					Port:     args[1],
				}
			}
		}

		if err := runDefault(logfile, &hosts); err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringP("hostsfile", "H", "", "The file or directory containing the hosts to check. Only supports json at the moment.")

	runCmd.PersistentFlags().StringP("logfile", "l", "", "Houses all logs for the scheduler service.")

	// TODO: Will need to set up support for various logfiles in the future.
	// 	runCmd.PersistentFlags().StringP("resultLog", "r", "", "The file to write structured result data logs to.")
}

type LogCloseFn func() error

// TODO: Am considering moving the log setup into the scheduler package
// and updating the WithLogger opt function to take the string, and
// set it so we can create the logger when initializing the
// service itself.
func setupLogger(logFile string) (LogCloseFn, *slog.Logger, error) {
	lFile, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	machineId := defaults.MachineName()

	logHandler := slog.NewJSONHandler(lFile, &slog.HandlerOptions{AddSource: true})
	logger := slog.New(logHandler).With(slog.String("machine_id", machineId))

	return lFile.Close, logger, nil
}

func runDefault(logFile string, hostStore *map[string]scheduler.Host) error {
	opts := make([]scheduler.OptionFn, 0)

	hostList := make([]scheduler.Host, len(*hostStore))
	for _, host := range *hostStore {
		hostList = append(hostList, host)
	}
	opts = append(opts, scheduler.WithHosts(hostList))

	var logCloseFn LogCloseFn
	var logger *slog.Logger
	var err error

	if logFile != "" {
		logCloseFn, logger, err = setupLogger(logFile)
		if err != nil {
			return err
		}
		opts = append(opts, scheduler.WithLogger(logger))
	}

	service, err := scheduler.New(opts...)
	if err != nil {
		return err
	}

	if service == nil {
		return errors.New("service is nil")
	}

	println("Starting service... Press Ctrl+C to stop.")
	if logFile == "" {
		logFile = "stdout"
	}
	println("Logging to: ", logFile)
	println()

	cancelSrv, err := service.Run()
	if err != nil {
		return err
	}

	// wait for Ctrl+C
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc

	// call our graceful shutdown
	cancelSrv()

	if logCloseFn != nil {
		if err := logCloseFn(); err != nil {
			return err
		}
	}

	return nil
}

func readHostsFile(path string, store *map[string]scheduler.Host) error {
	fInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	hosts := make([]scheduler.Host, 0)

	if fInfo.IsDir() {
		return errors.New("directory not supported yet")
	}

	hosts, err = processFile(path)
	if err != nil {
		return err
	}

	for _, host := range hosts {
		(*store)[host.Hostname] = host
	}

	return nil
}

func processFile(path string) ([]scheduler.Host, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var hosts []scheduler.Host
	if err := json.Unmarshal(data, &hosts); err != nil {
		return nil, err
	}

	return hosts, nil
}
