// Package starlighttest contains agent-level integration tests for starlight

package starlighttest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/interstellar/starlight/starlight"
	"github.com/interstellar/starlight/starlight/walletrpc"
	"github.com/interstellar/starlight/starlight/xlm"
	"github.com/interstellar/starlight/worizon"
)

type Starlightd struct {
	g             *starlight.Agent
	wclient       *worizon.Client
	handler       http.Handler
	server        *httptest.Server
	cookie        string
	address       string
	accountID     string
	nextUpdateNum uint64
	balance       xlm.Amount
}

func StartServer(ctx context.Context, testdir, name string) *Starlightd {
	return start(nil, ctx, testdir, name)
}

func (g *Starlightd) Address() string {
	return g.address
}

func start(t *testing.T, ctx context.Context, testdir, name string) *Starlightd {
	g, wclient := starlight.StartTestnetAgent(ctx, t, fmt.Sprintf("%s/testdb_%s", testdir, name))
	s := &Starlightd{
		g:             g,
		wclient:       wclient,
		nextUpdateNum: 1,
	}
	s.handler = logWrapper(walletrpc.Handler(s.g), name)
	s.server = httptest.NewTLSServer(s.handler)
	s.address = strings.TrimPrefix(s.server.URL, "https://")
	return s
}