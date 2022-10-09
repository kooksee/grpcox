package grpcproxy

import "context"

type Service interface {
	GetResource(ctx context.Context, target string, plainText, isRestartConn bool) (*Resource, error)
	GetActiveConns(ctx context.Context) []string
	CloseActiveConns(host string) error
	Extend(host string)
}
