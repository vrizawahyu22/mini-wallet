package app

import (
	"fmt"
	"mini-wallet/handler"
	"mini-wallet/utils"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/soheilhy/cmux"
)

type App interface {
	Start()
}

type AppImpl struct {
	route         *chi.Mux
	walletHandler handler.WalletHandler
	config        *utils.BaseConfig
}

func NewApp(
	route *chi.Mux,
	walletHandler handler.WalletHandler,
	config *utils.BaseConfig,
) App {
	return &AppImpl{
		route:         route,
		walletHandler: walletHandler,
		config:        config,
	}
}

func (s *AppImpl) Start() {
	utils.SetupMiddleware(s.route, s.config)
	s.walletHandler.SetupWalletRoutes(s.route)

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", "80"))
	if err != nil {
		utils.LogAndPanicIfError(err, "failed when start listening")
	}
	m := cmux.New(l)

	go func() {
		http1 := m.Match(cmux.HTTP1Fast())
		err = http.Serve(http1, s.route)
		utils.LogIfError(err)
	}()

	utils.LogInfo("server started")
	err = m.Serve()
	utils.LogIfError(err)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
}
