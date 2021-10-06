package add_test

import (
	"context"
	"log"
	"net"
	"os"
	"runtime/debug"
	"testing"

	pb "github.com/weaveworks/weave-gitops/pkg/api/applications"
	"github.com/weaveworks/weave-gitops/pkg/flux"
	"github.com/weaveworks/weave-gitops/pkg/middleware"
	"github.com/weaveworks/weave-gitops/pkg/osys"
	"github.com/weaveworks/weave-gitops/pkg/runner"
	"github.com/weaveworks/weave-gitops/pkg/server"
	"github.com/weaveworks/weave-gitops/pkg/testutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener
var addr = "0.0.0.0:9090"

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func failIfError(t *testing.T, err error) {
	if err != nil {
		debug.PrintStack()
		t.Fatal(err.Error())
	}
}

func TestAppAdd_ConfigRepoNone(t *testing.T) {
	ctx := context.Background()
	k8s, err := testutils.StartK8sTestEnvironment()
	failIfError(t, err)

	lis = bufconn.Listen(bufSize)

	defer k8s.Stop()
	flux.New(osys.New(), &runner.CLIRunner{}).SetupBin()

	cfg, err := server.DefaultConfig()
	if err != nil {
		t.Error(err)
	}

	cfg.KubeClient = k8s.Client

	s := grpc.NewServer()
	apps := server.NewApplicationsServer(cfg)
	pb.RegisterApplicationsServer(s, apps)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf(err.Error())
		}
	}()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	failIfError(t, err)

	defer conn.Close()

	client := pb.NewApplicationsClient(conn)

	secret := os.Getenv("GITHUB_TOKEN")

	// token, err := cfg.JwtClient.GenerateJWT(time.Duration(60*time.Second), gitproviders.GitProviderGitHub, secret)
	// failIfError(t, err)

	// fmt.Println(token)

	ctx = middleware.ContextWithGRPCAuth(ctx, secret)

	res, err := client.AddApplication(ctx, &pb.AddApplicationRequest{
		Name:      "my-app",
		Namespace: "wego-system",
		Url:       "git@github.com:jpellizzari/stringly.git",
		Branch:    "main",
		Path:      "k8s/overlays/development",
		AutoMerge: true,
	})
	failIfError(t, err)

	if !res.Success {
		t.Fail()
	}
}
