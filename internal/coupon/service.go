package coupon

import (
  "context"
  "crypto/rand"
  "encoding/hex"
  "time"

  "log"
  "github.com/bufbuild/connect-go"
  couponv1 "github.com/Nileshmaharjan/coupon-service/gen/coupon"
  "google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
  repo Store
}

func NewService(repo Store) *Service {
  return &Service{repo: repo}
}

// create a campaign
func (s *Service) CreateCampaign(
  ctx context.Context,
  req *connect.Request[couponv1.CreateCampaignRequest],
) (*connect.Response[couponv1.CreateCampaignResponse], error) {
  id := hex.EncodeToString(mustRandomBytes(4))
  c := &Campaign{
      ID:        id,
      Name:      req.Msg.Name,
      Total:     req.Msg.Total,
      StartTime: req.Msg.StartTime.AsTime(),
  }
  if err := s.repo.Create(c); err != nil {
      return nil, connect.NewError(connect.CodeAlreadyExists, err)
  }
  return connect.NewResponse(&couponv1.CreateCampaignResponse{CampaignId: id}), nil
}

// get a campaign
func (s *Service) GetCampaign(
  ctx context.Context,
  req *connect.Request[couponv1.GetCampaignRequest],
) (*connect.Response[couponv1.GetCampaignResponse], error) {
  c, err := s.repo.Get(req.Msg.CampaignId)
  if err != nil {
      return nil, connect.NewError(connect.CodeNotFound, err)
  }
  log.Printf("DEBUG GetCampaign: ID=%s, Codes=%#v", c.ID, c.Codes())

  resp := &couponv1.GetCampaignResponse{
      CampaignId:  c.ID,
      Name:        c.Name,
      Total:       c.Total,
      StartTime:   timestamppb.New(c.StartTime),
      IssuedCodes: c.Codes(),
  }
  return connect.NewResponse(resp), nil
}


// issue a coupon
func (s *Service) IssueCoupon(
  ctx context.Context,
  req *connect.Request[couponv1.IssueCouponRequest],
) (*connect.Response[couponv1.IssueCouponResponse], error) {
  code, err := s.repo.Issue(req.Msg.CampaignId, time.Now())
  if err != nil {
      return connect.NewResponse(&couponv1.IssueCouponResponse{
          Error: err.Error(),
      }), nil
  }
  return connect.NewResponse(&couponv1.IssueCouponResponse{
      Code: code,
  }), nil
}

func mustRandomBytes(n int) []byte {
  b := make([]byte, n)
  if _, err := rand.Read(b); err != nil {
    panic(err)
  }
  return b
}