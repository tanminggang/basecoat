package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/basecoat/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var cmdFormulasCreate = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a single formula",
	Long: `A formula is a combination of colors resulting in another color.
Formulas usually consist of one base with one or more colorants`,
	Args: cobra.MinimumNArgs(1),
	Run:  runFormulasCreateCmd,
}

func runFormulasCreateCmd(cmd *cobra.Command, args []string) {
	name := args[0]
	number, _ := cmd.Flags().GetString("number")
	notes, _ := cmd.Flags().GetString("notes")
	jobsraw, _ := cmd.Flags().GetStringSlice("jobs")

	var jobs []string
	for _, id := range jobsraw {
		jobs = append(jobs, id)
	}

	config, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed to read configuration")
	}

	creds, err := credentials.NewClientTLSFromFile(config.TLSCertPath, "")
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "failed to get certificates", err)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(creds))

	hostPortTuple := strings.Split(config.Backend.GRPCURL, ":")

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", hostPortTuple[0], hostPortTuple[1]), opts...)
	if err != nil {
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+config.CommandLine.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = basecoatClient.CreateFormula(ctx, &api.CreateFormulaRequest{
		Name:   name,
		Number: number,
		Notes:  notes,
		Jobs:   jobs,
	})
	if err != nil {
		log.Fatalf("could not create formula: %v", err)
	}
}

func init() {
	var number string
	cmdFormulasCreate.Flags().StringVarP(&number, "number", "u", "", "formula number used internally")

	var notes string
	cmdFormulasCreate.Flags().StringVarP(&notes, "notes", "o", "", "any additional notes for the formula")

	var jobs []string
	cmdFormulasCreate.Flags().StringSliceVarP(&jobs, "jobs", "j", []string{}, "comma seperated list of jobs by id in which this formula has been used")

	cmdFormulas.AddCommand(cmdFormulasCreate)
}