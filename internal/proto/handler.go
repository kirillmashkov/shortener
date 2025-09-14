package pb

import (
	context "context"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *GRPCServer) GetURL(ctx context.Context, r *GetURLRequest) (*GetURLResponse, error) {
	urlID := r.GetUrlId()
	if urlID == "" {
		return nil, status.Error(codes.InvalidArgument, "url_id required")
	}

	url, deleted, exist := s.service.GetShortURL(ctx, urlID)

	if !exist {
		return nil, status.Error(codes.NotFound, "url not found")
	}

	if deleted {
		return nil, status.Error(codes.NotFound, "url was deleted")
	}

	return &GetURLResponse{
		FullUrl: url,
	}, nil

}
