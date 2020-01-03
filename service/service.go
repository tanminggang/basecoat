package service

import (
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/basecoat/search"
	"github.com/clintjedwards/basecoat/storage"
	"github.com/clintjedwards/toolkit/logger"
	"go.uber.org/zap"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// API represents a basecoat grpc backend service
type API struct {
	storage storage.BoltDB
	config  *config.Config
	search  *search.Search
	log     *zap.SugaredLogger
}

// NewBasecoatAPI inits a grpc basecoat api service
func NewBasecoatAPI(config *config.Config) *API {
	basecoatAPI := API{}
	basecoatAPI.log = logger.Log()

	storage, err := storage.NewBoltDB(config.Database.Path, config.Database.IDLength)
	if err != nil {
		basecoatAPI.log.Fatalw("failed to initialize storage", "error", err)
	}

	searchIndex, err := search.InitSearch(storage)
	if err != nil {
		basecoatAPI.log.Fatalw("failed to initialize search indexes", "error", err)
	}

	rebuildTime := time.Duration(config.Backend.SearchIndexRebuildTime) * time.Second
	go func() {
		searchIndex.BuildIndex(storage)
		time.Sleep(rebuildTime)
	}()

	basecoatAPI.config = config
	basecoatAPI.storage = storage
	basecoatAPI.search = searchIndex

	return &basecoatAPI
}

// CreateGRPCServer creates a grpc server with all the proper settings; TLS enabled
func CreateGRPCServer(basecoatAPI *API) *grpc.Server {

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(basecoatAPI.authenticate),
			grpc_prometheus.UnaryServerInterceptor,
		)),
	)

	grpc_prometheus.EnableHandlingTimeHistogram()

	reflection.Register(grpcServer)
	grpc_prometheus.Register(grpcServer)
	api.RegisterBasecoatServer(grpcServer, basecoatAPI)

	return grpcServer
}
