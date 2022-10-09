package grpcproxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func newResource(md metadata.MD, clientConn *grpc.ClientConn, headers []string) *Resource {
	return &Resource{md: md, clientConn: clientConn, headers: headers}
}

// Resource - hold 3 main function (List, Describe, and Invoke)
type Resource struct {
	clientConn *grpc.ClientConn
	descSource grpcurl.DescriptorSource
	refClient  *grpcreflect.Client

	headers []string
	md      metadata.MD
}

//openDescriptor - use it to reflect server descriptor
func (r *Resource) openDescriptor() error {
	ctx := context.Background()
	refCtx := metadata.NewOutgoingContext(ctx, r.md)
	r.refClient = grpcreflect.NewClient(refCtx, reflectpb.NewServerReflectionClient(r.clientConn))

	r.descSource = grpcurl.DescriptorSourceFromServer(ctx, r.refClient)
	return nil
}

//closeDescriptor - please ensure to always close after open in the same flow
func (r *Resource) closeDescriptor() {
	done := make(chan int)
	go func() {
		if r.refClient != nil {
			r.refClient.Reset()
		}
		done <- 1
	}()

	select {
	case <-done:
		return
	case <-time.After(3 * time.Second):
		log.Printf("Reflection %s failed to close\n", r.clientConn.Target())
		return
	}
}

func (r *Resource) ListAll() ([]string, error) {
	err := r.openDescriptor()
	if err != nil {
		return nil, err
	}
	defer r.closeDescriptor()

	var result []string
	svcs, err := grpcurl.ListServices(r.descSource)
	if err != nil {
		return result, err
	}
	if len(svcs) == 0 {
		return result, fmt.Errorf("No Services")
	}

	for _, svc := range svcs {
		result = append(result, fmt.Sprintf("%s\n", svc))
	}
	return result, nil
}

// List - To list all services exposed by a server
// symbol can be "" to list all available services
// symbol also can be service name to list all available method
func (r *Resource) List(symbol string) ([]string, error) {
	err := r.openDescriptor()
	if err != nil {
		return nil, err
	}
	defer r.closeDescriptor()

	var result []string
	methods, err := grpcurl.ListMethods(r.descSource, symbol)
	if err != nil {
		return result, err
	}

	if len(methods) == 0 {
		return result, fmt.Errorf("No Function") // probably unlikely
	}

	for _, m := range methods {
		result = append(result, fmt.Sprintf("%s\n", m))
	}

	return result, nil
}

// Describe - The "describe" verb will print the type of any symbol that the server knows about
// or that is found in a given protoset file.
// It also prints a description of that symbol, in the form of snippets of proto source.
// It won't necessarily be the original source that defined the element, but it will be equivalent.
func (r *Resource) Describe(symbol string) (string, string, error) {
	err := r.openDescriptor()
	if err != nil {
		return "", "", err
	}
	defer r.closeDescriptor()

	var result, template string

	var symbols []string
	if symbol != "" {
		symbols = []string{symbol}
	} else {
		// if no symbol given, describe all exposed services
		svcs, err := r.descSource.ListServices()
		if err != nil {
			return "", "", err
		}
		if len(svcs) == 0 {
			log.Println("Server returned an empty list of exposed services")
		}
		symbols = svcs
	}

	for _, s := range symbols {
		if s[0] == '.' {
			s = s[1:]
		}

		dsc, err := r.descSource.FindSymbol(s)
		if err != nil {
			return "", "", err
		}

		txt, err := grpcurl.GetDescriptorText(dsc, r.descSource)
		if err != nil {
			return "", "", err
		}
		result = txt

		if dsc, ok := dsc.(*desc.MessageDescriptor); ok {
			// for messages, also show a template in JSON, to make it easier to
			// create a request to invoke an RPC
			tmpl := grpcurl.MakeTemplate(dsc)

			_, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, r.descSource, nil, grpcurl.FormatOptions{
				EmitJSONDefaultFields: true,
				IncludeTextSeparator:  false,
			})
			if err != nil {
				return "", "", err
			}
			str, err := formatter(tmpl)
			if err != nil {
				return "", "", err
			}
			template = str
		}
	}

	return result, template, nil
}

// Invoke - invoking gRPC function
func (r *Resource) Invoke(ctx context.Context, metadata []string, symbol string, in io.Reader) (string, time.Duration, error) {
	err := r.openDescriptor()
	if err != nil {
		return "", 0, err
	}
	defer r.closeDescriptor()

	var resultBuffer bytes.Buffer

	rf, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, r.descSource, in, grpcurl.FormatOptions{
		EmitJSONDefaultFields: true,
		IncludeTextSeparator:  false,
	})
	if err != nil {
		return "", 0, err
	}
	h := &grpcurl.DefaultEventHandler{Out: &resultBuffer, Formatter: formatter}

	var headers = r.headers
	if len(metadata) != 0 {
		headers = metadata
	}

	start := time.Now()
	err = grpcurl.InvokeRPC(ctx, r.descSource, r.clientConn, symbol, headers, h, rf.Next)
	end := time.Now().Sub(start) / time.Millisecond
	if err != nil {
		return "", end, err
	}

	if h.Status.Code() != codes.OK {
		return "", end, fmt.Errorf(h.Status.Message())
	}

	return resultBuffer.String(), end, nil
}

// Close - to close all resources that was opened before
func (r *Resource) Close() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if r.clientConn != nil {
			r.clientConn.Close()
			r.clientConn = nil
		}
	}()

	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return
	case <-time.After(3 * time.Second):
		log.Printf("Connection %s failed to close\n", r.clientConn.Target())
		return
	}
}

func (r *Resource) isValid() bool {
	return r.refClient != nil && r.clientConn != nil
}

func (r *Resource) exit(code int) {
	// to force reset before os exit
	r.Close()
	os.Exit(code)
}
