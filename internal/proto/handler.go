package pb

import (
	context "context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/kirillmashkov/shortener.git/internal/app"
	"github.com/kirillmashkov/shortener.git/internal/model"
	"go.uber.org/zap"
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

func (s *GRPCServer) CreateShort(ctx context.Context, r *CreateShortRequest) (*CreateShortResponse, error) {
	if r.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "url required")
	}

	var userID int
	var err error
	if r.UserId == "" {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		userID = r.Int()
	} else {
		userID, err = strconv.Atoi(r.UserId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "userId must be int")	
		}
	}

	shortURL, err := s.service.ProcessURL(ctx, r.Url, userID)

	if err != nil {
		if errors.Is(err, model.ErrDuplicateURL) {
			return &CreateShortResponse {ResultUrl: shortURL, UserId: "", UrlId: shortURL}, nil
		}
		app.Log.Error("Error process URL", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &CreateShortResponse {ResultUrl: shortURL, UserId: string(userID)}, nil

}
