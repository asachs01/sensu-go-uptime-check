package main

//Import the packages we need
import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hako/durafmt"
	"github.com/sensu/sensu-go/types"
	"github.com/shirou/gopsutil/host"
	"github.com/spf13/cobra"
)

//Set up some variables. Most notably, warning and critical as time durations
var (
	warning, critical time.Duration
	greaterThan       bool
	stdin             *os.File
)

//Start our main function
func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

//Set up our flags for the command. Note that we have time duration defaults for warning & critical
func configureRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sensu-go-uptime-status",
		Short: "The Sensu Go check for system uptime",
		RunE:  run,
	}

	cmd.Flags().BoolVarP(&greaterThan,
		"greater-than",
		"g",
		false,
		"This compare uptime > threshold. Default behavior uptime < threshold")

	cmd.Flags().DurationVarP(&warning,
		"warning",
		"w",
		time.Duration(72*time.Hour),
		"Warning value in seconds, minutes, or hours, default is 72 hours (72h)")

	cmd.Flags().DurationVarP(&critical,
		"critical",
		"c",
		time.Duration(168*time.Hour),
		"Critical value in seconds, minutes, or hours default is 1 week (168h)")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {

	if len(args) != 0 {
		_ = cmd.Help()
		return fmt.Errorf("invalid argument(s) received")
	}

	if stdin == nil {
		stdin = os.Stdin
	}

	event := &types.Event{}

	return checkUptime(event)
}

//Here we start the meat of what we do.
func checkUptime(event *types.Event) error {
	// Setting uptime as the value retrieved from gopsutil
	uptime, err := host.Uptime()

	// Let's set up some error handling
	if err != nil {
		msg := fmt.Sprintf("Failed to determine uptime %s", err.Error())
		io.WriteString(os.Stdout, msg)
		os.Exit(3)
	}

	// Add a variable for uptimeSecs, which converts uptime to a duration
	uptimeSecs := time.Duration(uptime) * time.Second

	if greaterThan {
		return checkUptimeGreaterThan(uptimeSecs)
	}
	return checkUptimeLessThan(uptimeSecs)
}

func report(res string, uptimeSecs time.Duration) {
	// Setting "CheckUptime" as a constant
	const checkName = "CheckUptime"
	const metricName = "current_system_uptime"

	msg := fmt.Sprintf("%s %s - value = %v | %s=%d\n", checkName, res, durafmt.Parse(uptimeSecs), metricName, int64(uptimeSecs.Seconds()))
	io.WriteString(os.Stdout, msg)
}

func checkUptimeLessThan(uptimeSecs time.Duration) error {
	// Sets up conditionss for a comparison
	if uptimeSecs > critical {
		report("CRITICAL", uptimeSecs)
		os.Exit(2)
	} else if uptimeSecs >= warning && uptimeSecs <= critical {
		report("WARNING", uptimeSecs)
		os.Exit(1)
	} else {
		report("OK", uptimeSecs)
		os.Exit(0)
	}

	return nil
}

func checkUptimeGreaterThan(uptimeSecs time.Duration) error {
	// Sets up conditionss for a comparison
	if uptimeSecs < critical {
		report("CRITICAL", uptimeSecs)
		os.Exit(2)
	} else if uptimeSecs < warning {
		report("WARNING", uptimeSecs)
		os.Exit(1)
	} else {
		report("OK", uptimeSecs)
		os.Exit(0)
	}

	return nil
}
