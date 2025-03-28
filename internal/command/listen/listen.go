package listen

import (
	"fmt"
	"strings"

	"github.com/panjf2000/gnet/v2"
	"github.com/rs/zerolog/log"

	"github.com/moyu-x/level-5/pkg/logger"
)

func init() {
	logger.NewLogger(nil)
}

type gnetServer struct {
	*gnet.BuiltinEventEngine
}

func (s *gnetServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	buf, _ := c.Next(-1)
	log.Info().Msgf("Received: %s", strings.TrimSuffix(string(buf), "\n"))
	return
}

func (s *gnetServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Info().Msgf("New connection established: %s", c.RemoteAddr())
	return
}

func (s *gnetServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	log.Info().Msgf("Connection disconnected: %s", c.RemoteAddr())
	return
}

func Run(port int, protocol string) {
	options := []gnet.Option{
		gnet.WithMulticore(true),
		gnet.WithReusePort(true),
		gnet.WithTCPKeepAlive(0),
	}

	addr := fmt.Sprintf("tcp://:%d", port)
	if "udp" == protocol {
		addr = fmt.Sprintf("udp://:%d", port)
	}
	server := &gnetServer{}

	err := gnet.Run(server, addr, options...)
	if err != nil {
		log.Fatal().Err(err).Msg("gnet server run failed")
	}
}
