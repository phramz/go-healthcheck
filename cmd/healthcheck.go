package main

import (
	"fmt"
	"os"

	"github.com/phramz/go-healthcheck/internal"
	"github.com/phramz/go-healthcheck/pkg/contract"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := newApp().Run(os.Args); err != nil {
		panic(err)
	}
}

func newApp() *cli.App {
	sc := internal.NewEcosystem()

	commands := []*cli.Command{
		{
			Name:  "probe",
			Usage: "Performs a HTTP connection check",
			Action: func(c *cli.Context) error {
				// set up probes
				probes := make([]contract.Probe, 0)
				for i := 0; i < c.NArg(); i++ {
					name := fmt.Sprintf(`p%d`, i)
					tpl := c.Args().Get(i)
					p, err := internal.NewProbe(name, tpl, sc.Logger().WithField(internal.LogFieldProbe, name), sc.Config())
					if err != nil {
						return cli.Exit(err, -1)
					}

					probes = append(probes, p)
				}

				// run probes
				numFail := 0
				numPass := 0
				for _, p := range probes {
					out, err := p.Run()
					if err != nil {
						return cli.Exit(err, -1)
					}

					if sc.Config().Bool(contract.ConfigKeyOutput, false) {
						_, err = fmt.Fprint(os.Stdout, out)

						if err != nil {
							sc.Logger().Error(err)
						}
					}

					if p.HasFailed() {
						numFail++
					}

					if p.HasPassed() {
						numPass++
					}
				}

				// finish
				if numFail > 0 {
					sc.Logger().Errorf(`%d of %d probes failed 😭`, numFail, numFail+numPass)
					if sc.Config().Bool(contract.ConfigKeyNoFail, false) || (sc.Config().Bool(contract.ConfigKeyNoFailAffirmative, false) && numPass > 0) {
						return nil
					}

					return cli.Exit("FAIL", 1)
				}

				sc.Logger().Infof(`%d of %d probes passed 👍`, numPass, numFail+numPass)
				return nil
			},
		},
	}

	return &cli.App{
		Name:     "healthcheck",
		Usage:    "HealthCheck Utility",
		Commands: commands,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    contract.ConfigKeyLogLevel,
				Usage:   "Set the log output verbosity",
				Value:   internal.LogLevelInfo,
				EnvVars: []string{"HC_DEBUG"},
			},
			&cli.BoolFlag{
				Name:    contract.ConfigKeyOutput,
				Usage:   "If set rendered output will be printed to stdout",
				Value:   false,
				EnvVars: []string{"HC_OUTPUT"},
			},
			&cli.BoolFlag{
				Name:    contract.ConfigKeyNoFail,
				Usage:   "If set failed checks won't cause the program to exit with an error code",
				Value:   false,
				EnvVars: []string{"HC_NO_FAIL"},
			},
			&cli.BoolFlag{
				Name:    contract.ConfigKeyNoFailAffirmative,
				Usage:   "If set the program won't exit with an error code as long as at least one probe passed",
				Value:   false,
				EnvVars: []string{"HC_NO_FAIL_AFFIRMATIVE"},
			},
			&cli.IntFlag{
				Name:    contract.ConfigKeyProbeTimeout,
				Usage:   "Maximum time in seconds after an unfinished probe will fail",
				Value:   30,
				EnvVars: []string{"HC_TIMEOUT"},
			},
		},
		Before: func(c *cli.Context) error {
			config := sc.Config().
				WithBool(contract.ConfigKeyNoFail, c.Bool(contract.ConfigKeyNoFail)).
				WithBool(contract.ConfigKeyNoFailAffirmative, c.Bool(contract.ConfigKeyNoFailAffirmative)).
				WithBool(contract.ConfigKeyOutput, c.Bool(contract.ConfigKeyOutput)).
				WithInt(contract.ConfigKeyProbeTimeout, c.Int(contract.ConfigKeyProbeTimeout)).
				WithString(contract.ConfigKeyLogLevel, c.String(contract.ConfigKeyLogLevel))

			sc = sc.WithConfig(config)

			return nil
		},
	}
}
