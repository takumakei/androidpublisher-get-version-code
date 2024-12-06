package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/goaux/stacktrace"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", stacktrace.Format(err))
		os.Exit(1)
	}
}

//go:embed usageCredentials.txt
var usageCredentials string

//go:embed usagePackageName.txt
var usagePackageName string

//go:embed usageOutputStyle.txt
var usageOutputStyle string

//go:embed usageTimeLimit.txt
var usageTimeLimit string

var printers = map[string]func(*androidpublisher.TracksListResponse) error{
	"":           printHighest,
	"highest":    printHighest,
	"production": newPrinter("production"),
	"beta":       newPrinter("beta"),
	"alpha":      newPrinter("alpha"),
	"internal":   newPrinter("internal"),
	"response":   printResponse,
}

var flags struct {
	Credentials string
	PackageName string
	OutputStyle string
	TimeLimit   time.Duration
}

func run() error {
	flag.StringVar(&flags.Credentials, "credentials", "", usageString(usageCredentials))
	flag.StringVar(&flags.PackageName, "package-name", os.Getenv("PACKAGE_NAME"), usageString(usagePackageName))
	flag.StringVar(&flags.OutputStyle, "output-style", os.Getenv("OUTPUT_STYLE"), usageString(usageOutputStyle))
	flag.DurationVar(&flags.TimeLimit, "time-limit", getenvDuration("TIME_LIMIT", 30*time.Second), usageString(usageTimeLimit))
	flag.Parse()

	if flags.PackageName == "" {
		return fmt.Errorf("Package name must be specified")
	}

	var printer func(*androidpublisher.TracksListResponse) error
	if v, ok := printers[flags.OutputStyle]; ok {
		printer = v
	} else {
		return fmt.Errorf("output-style must be one of highest, production, beta, alpha, internal or response")
	}

	if flags.Credentials == "" {
		// To avoid credentials appearing as a default value in the usage text,
		// the value is read here from an environment variable.
		flags.Credentials = os.Getenv("CREDENTIALS")
		if flags.Credentials == "" {
			return fmt.Errorf("Credentials must be specified")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), flags.TimeLimit)
	defer cancel()

	service, err := androidpublisher.NewService(ctx, withCredentials(flags.Credentials))
	if err != nil {
		return stacktrace.With(err)
	}

	appEdit, err := androidpublisher.
		NewEditsService(service).
		Insert(flags.PackageName, &androidpublisher.AppEdit{}).
		Context(ctx).
		Do()
	if err != nil {
		return stacktrace.With(err)
	}

	response, err := androidpublisher.
		NewEditsTracksService(service).
		List(flags.PackageName, appEdit.Id).
		Context(ctx).
		Do()
	if err != nil {
		return stacktrace.With(err)
	}

	return printer(response)
}

func usageString(s string) string {
	return strings.TrimSpace(s) + "\n"
}

func withCredentials(credentials string) option.ClientOption {
	switch {
	case strings.HasPrefix(credentials, "@env:"):
		return option.WithCredentialsJSON([]byte(os.Getenv(credentials[5:])))
	case strings.HasPrefix(credentials, "@file:"):
		return option.WithCredentialsFile(credentials[6:])
	}
	return option.WithCredentialsJSON([]byte(credentials))
}

func printHighest(response *androidpublisher.TracksListResponse) error {
	return printCode(
		response,
		func(track string) bool { return true },
	)
}

func newPrinter(targetTrack string) func(*androidpublisher.TracksListResponse) error {
	return func(response *androidpublisher.TracksListResponse) error {
		return printCode(
			response,
			func(track string) bool { return track == targetTrack },
		)
	}
}

func printCode(response *androidpublisher.TracksListResponse, selector func(string) bool) error {
	var result struct {
		Track string `json:"track,omitempty"`
		Name  string `json:"name,omitempty"`
		Code  int64  `json:"code,omitempty"`
	}
	for _, track := range response.Tracks {
		for _, release := range track.Releases {
			for _, code := range release.VersionCodes {
				if selector(track.Track) && result.Code < code {
					result.Track = track.Track
					result.Name = release.Name
					result.Code = code
				}
			}
		}
	}
	return jsonify(os.Stdout, &result)
}

func printResponse(response *androidpublisher.TracksListResponse) error {
	return jsonify(os.Stdout, response)
}

func jsonify(w io.Writer, v any) error {
	j := json.NewEncoder(w)
	j.SetEscapeHTML(false)
	j.SetIndent("", "  ")
	return stacktrace.With(j.Encode(v))
}

func getenvDuration(key string, fallback time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
