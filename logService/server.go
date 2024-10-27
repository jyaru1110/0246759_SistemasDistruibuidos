package logService

import (
	"context"

	api "server/api/v1"

	log "server/log"
)

var _ api.LogServer = (*GrpcServer)(nil)

type GrpcServer struct {
	api.UnimplementedLogServer
	CommitLog *log.Log
}

func newgrpcServer(commitlog *log.Log) (srv *GrpcServer, err error) {
	srv = &GrpcServer{
		CommitLog: commitlog,
	}
	return srv, nil
}

func (s *GrpcServer) Produce(ctx context.Context, req *api.ProduceRequest) (*api.ProduceResponse, error) {
	offset, err := s.CommitLog.Append(req.Record)
	if err != nil {
		return nil, err
	}
	return &api.ProduceResponse{Offset: offset}, nil
}

func (s *GrpcServer) Consume(ctx context.Context, req *api.ConsumeRequest) (*api.ConsumeResponse, error) {
	record, err := s.CommitLog.Read(req.Offset)
	if err != nil {
		re, ok := err.(*api.ErrOffsetOutOfRange)
		if ok {
			return nil, re.GRPCStatus().Err()
		}
		return nil, err
	}
	return &api.ConsumeResponse{Record: record}, nil
}

func (s *GrpcServer) ProduceStream(stream api.Log_ProduceStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		res, err := s.Produce(stream.Context(), req)
		if err != nil {
			return err
		}
		if err = stream.Send(res); err != nil {
			return err
		}
	}
}

func (s *GrpcServer) ConsumeStream(req *api.ConsumeRequest, stream api.Log_ConsumeStreamServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			res, err := s.Consume(stream.Context(), req)
			switch err.(type) {
			case nil:
			case api.ErrOffsetOutOfRange:
				continue
			default:
				return err
			}
			if err = stream.Send(res); err != nil {
				return err
			}
			req.Offset++
		}
	}
}
