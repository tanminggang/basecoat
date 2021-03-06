package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/app"
	"github.com/clintjedwards/toolkit/random"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type testHarness struct {
	client       api.BasecoatClient
	apikey       string
	adminkey     string
	databasePath string
}

func (info *testHarness) setup() {
	databasePath := fmt.Sprintf("/tmp/basecoat%s.db", random.GenerateRandString(4))

	os.Setenv("TLS_CERT_PATH", "../localhost.crt")
	os.Setenv("TLS_KEY_PATH", "../localhost.key")
	os.Setenv("DATABASE_PATH", databasePath)
	os.Setenv("BACKEND_ADMIN_TOKEN", "admin")
	os.Setenv("LOGLEVEL", "error")

	go app.StartServices()
	time.Sleep(time.Second)

	var opts []grpc.DialOption

	creds, err := credentials.NewClientTLSFromFile("../localhost.crt", "")
	if err != nil {
		zap.S().Fatalf("failed to get certificates: %v", err)
	}

	opts = append(opts, grpc.WithTransportCredentials(creds))

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8080"), opts...)
	if err != nil {
		zap.S().Fatalf("could not connect to basecoat: %v", err)
	}

	basecoatClient := api.NewBasecoatClient(conn)
	info.client = basecoatClient
	info.adminkey = "admin"
	info.databasePath = databasePath
}

func (info *testHarness) cleanup() {
	os.Unsetenv("TLS_CERT_PATH")
	os.Unsetenv("TLS_KEY_PATH")
	os.Unsetenv("DATABASE_PATH")
	os.Unsetenv("BACKEND_ADMIN_TOKEN")
	os.Unsetenv("LOGLEVEL")
	os.Remove(info.databasePath)
}

func TestFullApplication(t *testing.T) {
	info := testHarness{}
	info.setup()

	// Test accounts
	info.TestCreateAccount(t)
	info.TestGetAccount(t)
	info.TestListAccounts(t)
	info.TestDisableAccount(t)
	info.TestUpdateAccount(t)
	// Test auth
	info.TestCreateAPIToken(t)
	// Test contractors
	info.TestCreateContractor(t)
	info.TestGetContractor(t)
	info.TestListContractors(t)
	info.TestUpdateContractor(t)
	info.TestDeleteContractor(t)
	info.TestCreateContractorAddJob(t)
	// Test jobs
	info.TestCreateJob(t)
	info.TestGetJob(t)
	info.TestListJobs(t)
	info.TestUpdateJob(t)
	info.TestDeleteJob(t)
	info.TestCreateJobAddContractor(t)
	info.TestUpdateJobRemoveContractor(t)

	info.cleanup()
}
