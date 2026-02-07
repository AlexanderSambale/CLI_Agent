package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"cli_agent/internal/openai"

	"github.com/spf13/pflag"
)

var (
	modelsList bool
	modelsGet  string
)

// ModelsCmd represents the models command
func ModelsCmd() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("models", pflag.ExitOnError)

	flagSet.BoolVarP(&modelsList, "list", "l", false, "List all available models")
	flagSet.StringVarP(&modelsGet, "get", "g", "", "Get details for a specific model")

	return flagSet
}

// ExecuteModels runs the models command
func ExecuteModels(client openai.CLIClient, args []string) error {
	flagSet := ModelsCmd()
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	ctx := context.Background()

	if modelsGet != "" {
		return getModelDetails(client, ctx, modelsGet)
	}

	return listModels(client, ctx)
}

func listModels(client openai.CLIClient, ctx context.Context) error {
	models, err := client.ListModels(ctx)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tOwned By\tCreated")
	for _, model := range models {
		fmt.Fprintf(w, "%s\t%s\t%d\n", model.ID, model.OwnedBy, model.Created)
	}
	w.Flush()

	return nil
}

func getModelDetails(client openai.CLIClient, ctx context.Context, modelID string) error {
	model, err := client.GetModel(ctx, modelID)
	if err != nil {
		return err
	}

	fmt.Printf("ID: %s\n", model.ID)
	fmt.Printf("Object: %s\n", model.Object)
	fmt.Printf("Owned By: %s\n", model.OwnedBy)
	fmt.Printf("Created: %d\n", model.Created)

	return nil
}
