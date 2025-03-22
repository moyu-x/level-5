package listen

import (
	"bufio"
	"context"
	"strconv"
	"strings"

	"github.com/cloudwego/netpoll"
	"github.com/rs/zerolog/log"

	"github.com/moyu-x/level-5/pkg/logger"
)

func init() {
	logger.NewLogger(nil)
}

func Run(port int, protocol string) {
	listener, err := netpoll.CreateListener(protocol, "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal().Err(err).Msg("create listener failed")
		return
	}
	defer listener.Close()

	log.Info().Msg("start listen on: " + listener.Addr().String())
	// 创建一个事件循环
	eventLoop, err := netpoll.NewEventLoop(
		handleConnection,
		netpoll.WithOnPrepare(onPrepare),
		netpoll.WithOnConnect(onConnect),
		netpoll.WithOnDisconnect(onDisconnect),
	)
	if err != nil {
		log.Fatal().Msgf("Failed to create event loop: %v", err)
	}
	// 开始监听
	err = eventLoop.Serve(listener)
	if err != nil {
		log.Fatal().Msgf("Failed to serve: %v", err)
	}
}

func handleConnection(ctx context.Context, connection netpoll.Connection) error {
	reader := bufio.NewReader(connection)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		log.Info().Msgf("Received: %s", strings.TrimSuffix(line, "\n"))
	}
}

func onPrepare(connection netpoll.Connection) context.Context {
	log.Info().Msgf("New connection prepared: %s", connection.RemoteAddr())
	return context.Background()
}

func onConnect(ctx context.Context, connection netpoll.Connection) context.Context {
	log.Info().Msgf("New connection established: %s", connection.RemoteAddr())
	return ctx
}

func onDisconnect(ctx context.Context, connection netpoll.Connection) {
	log.Info().Msgf("Connection disconnected: %s", connection.RemoteAddr())
}
