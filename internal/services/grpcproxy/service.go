package grpcproxy

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/fullstorydev/grpcurl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// New constructor
func New() Service {
	maxLife, tick := 10, 3

	if val, err := strconv.Atoi(os.Getenv("MAX_LIFE_CONN")); err == nil {
		maxLife = val
	}

	if val, err := strconv.Atoi(os.Getenv("TICK_CLOSE_CONN")); err == nil {
		tick = val
	}

	c := NewConnectionStore()
	g := &serviceImpl{
		activeConn: c,
	}

	if maxLife > 0 && tick > 0 {
		g.maxLifeConn = time.Duration(maxLife) * time.Minute
		c.StartGC(time.Duration(tick) * time.Second)
	}

	return g
}

// serviceImpl - main object
type serviceImpl struct {
	KeepAlive float64

	activeConn  *ConnStore
	maxLifeConn time.Duration

	// TODO : utilize below args
	headers        []string
	reflectHeaders []string
	authority      string
	insecure       bool
	cacert         string
	cert           string
	key            string
	serverName     string
	isUnixSocket   func() bool
}

// GetResource - open resource to targeted grpc server
func (g *serviceImpl) GetResource(ctx context.Context, target string, plainText, isRestartConn bool) (*Resource, error) {
	if conn, ok := g.activeConn.getConnection(target); ok {
		if !isRestartConn && conn.isValid() {
			return conn, nil
		}
		g.CloseActiveConns(target)
	}

	h := append(g.headers, g.reflectHeaders...)
	md := grpcurl.MetadataFromHeaders(h)
	clientConn, err := g.dial(ctx, target, plainText)
	if err != nil {
		return nil, err
	}

	r := newResource(md, clientConn, h)
	g.activeConn.addConnection(target, r, g.maxLifeConn)
	return r, nil
}

// GetActiveConns - get all saved active connection
func (g *serviceImpl) GetActiveConns(ctx context.Context) []string {
	active := g.activeConn.getAllConn()
	result := make([]string, len(active))
	i := 0
	for k := range active {
		result[i] = k
		i++
	}
	return result
}

// CloseActiveConns - close conn by host or all
func (g *serviceImpl) CloseActiveConns(host string) error {
	if host == "all" {
		for k := range g.activeConn.getAllConn() {
			g.activeConn.delete(k)
		}
		return nil
	}

	g.activeConn.delete(host)
	return nil
}

// Extend extend connection based on setting max life
func (g *serviceImpl) Extend(host string) {
	g.activeConn.extend(host, g.maxLifeConn)
}

func (g *serviceImpl) dial(ctx context.Context, target string, plainText bool) (*grpc.ClientConn, error) {
	dialTime := 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, dialTime)
	defer cancel()
	var opts []grpc.DialOption

	// keep alive
	if g.KeepAlive > 0 {
		timeout := time.Duration(g.KeepAlive * float64(time.Second))
		opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    timeout,
			Timeout: timeout,
		}))
	}

	if g.authority != "" {
		opts = append(opts, grpc.WithAuthority(g.authority))
	}

	var creds credentials.TransportCredentials
	if !plainText {
		var err error
		tlsConf, err := grpcurl.ClientTLSConfig(g.insecure, g.cacert, g.cert, g.key)
		if err != nil {
			return nil, err
		}
		creds = credentials.NewTLS(tlsConf)
	}

	network := "tcp"
	if g.isUnixSocket != nil && g.isUnixSocket() {
		network = "unix"
	}

	cc, err := grpcurl.BlockingDial(ctx, network, target, creds, opts...)
	if err != nil {
		return nil, err
	}

	return cc, nil
}
