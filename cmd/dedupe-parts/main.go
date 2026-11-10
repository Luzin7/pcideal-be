package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Luzin7/pcideal-be/infra/database"
	"github.com/joho/godotenv"
)

const (
	defaultRunTimeout = 5 * time.Minute
	disconnectTimeout = 10 * time.Second
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, out io.Writer, errOut io.Writer) int {
	flags := flag.NewFlagSet("dedupe-parts", flag.ContinueOnError)
	flags.SetOutput(errOut)

	apply := flags.Bool("apply", false, "delete duplicate parts and create the unique store/url index")
	confirmation := flags.String("confirm", "", "confirmation value required by --apply")
	check := flags.Bool("check", false, "exit 1 when duplicate groups or missing store/url documents exist")
	timeout := flags.Duration("timeout", defaultRunTimeout, "maximum runtime for the command")

	if err := flags.Parse(args); err != nil {
		return 2
	}

	if err := validateApplyMode(*apply, *confirmation); err != nil {
		fmt.Fprintf(errOut, "error: %v\n", err)
		return 2
	}

	_ = godotenv.Load(".env")

	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		fmt.Fprintf(errOut, "error: %v\n", &MissingConfigError{Variable: "DATABASE_URL"})
		return 2
	}

	databaseName := os.Getenv("PCIDEAL_DB_NAME")
	if databaseName == "" {
		fmt.Fprintf(errOut, "error: %v\n", &MissingConfigError{Variable: "PCIDEAL_DB_NAME"})
		return 2
	}

	db, err := database.MongoConnection(connectionString, databaseName)
	if err != nil {
		fmt.Fprintf(errOut, "error: connect to database: %v\n", err)
		return 2
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), disconnectTimeout)
		defer cancel()

		if err := db.Client().Disconnect(ctx); err != nil {
			fmt.Fprintf(errOut, "error: disconnect from database: %v\n", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	service := NewDedupeService(db.Collection("parts"), out)

	summary := auditSummary{}
	if *apply {
		err = service.Apply(ctx)
	} else {
		summary, err = service.DryRun(ctx)
	}

	if err != nil {
		fmt.Fprintf(errOut, "error: %v\n", err)
		if *apply {
			fmt.Fprintln(errOut, "the operation is idempotent: rerun --apply --confirm=DELETE-DUPLICATES to complete the job")
		}
		return 2
	}

	if *check && checkFailed(summary) {
		return 1
	}

	return 0
}

func checkFailed(summary auditSummary) bool {
	return summary.Groups > 0 || summary.MissingKeys > 0
}

func validateApplyMode(apply bool, confirmation string) error {
	if !apply {
		return nil
	}

	if confirmation != deleteConfirmation {
		return &MissingConfirmationError{Provided: confirmation}
	}

	return nil
}
