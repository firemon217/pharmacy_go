package transport

import (
	"context"
	"errors"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"pharmacy/medicine"
	"pharmacy/pb"
)

// GRPCServer адаптирует protobuf-интерфейс к Go kit endpoints.
type GRPCServer struct {
	pb.UnimplementedPharmacyServer
	create  kitgrpc.Handler
	getByID kitgrpc.Handler
	sell    kitgrpc.Handler
}

func NewGRPCServer(e Endpoints) *GRPCServer {
	return &GRPCServer{
		create:  kitgrpc.NewServer(e.Create, decodeGRPCCreate, encodeGRPCCreate),
		getByID: kitgrpc.NewServer(e.GetByID, decodeGRPCGetByID, encodeGRPCGetByID),
		sell:    kitgrpc.NewServer(e.Sell, decodeGRPCSell, encodeGRPCSell),
	}
}

func (s *GRPCServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	_, response, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return response.(*pb.CreateResponse), nil
}
func (s *GRPCServer) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.MedicineResponse, error) {
	_, response, err := s.getByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return response.(*pb.MedicineResponse), nil
}
func (s *GRPCServer) Sell(ctx context.Context, req *pb.SellRequest) (*pb.Empty, error) {
	_, response, err := s.sell.ServeGRPC(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return response.(*pb.Empty), nil
}

func decodeGRPCCreate(_ context.Context, request interface{}) (interface{}, error) {
	r := request.(*pb.CreateRequest)
	return CreateRequest{Medicine: medicine.Medicine{
		Name: r.Name, Manufacturer: r.Manufacturer, Price: r.Price,
		Prescription: r.Prescription, Quantity: int(r.Quantity),
	}}, nil
}
func encodeGRPCCreate(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(CreateResponse)
	if r.Err != nil {
		return nil, r.Err
	}
	return &pb.CreateResponse{Id: r.ID}, nil
}
func decodeGRPCGetByID(_ context.Context, request interface{}) (interface{}, error) {
	return GetByIDRequest{ID: request.(*pb.GetByIDRequest).Id}, nil
}
func encodeGRPCGetByID(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(GetByIDResponse)
	if r.Err != nil {
		return nil, r.Err
	}
	m := r.Medicine
	return &pb.MedicineResponse{Medicine: &pb.Medicine{
		Id: m.ID, Name: m.Name, Manufacturer: m.Manufacturer, Price: m.Price,
		Prescription: m.Prescription, Quantity: int32(m.Quantity),
	}}, nil
}
func decodeGRPCSell(_ context.Context, request interface{}) (interface{}, error) {
	r := request.(*pb.SellRequest)
	return SellRequest{ID: r.Id, Count: int(r.Count)}, nil
}
func encodeGRPCSell(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(SellResponse)
	if r.Err != nil {
		return nil, r.Err
	}
	return &pb.Empty{}, nil
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, medicine.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, medicine.ErrInvalidMedicine), errors.Is(err, medicine.ErrInvalidCount):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, medicine.ErrOutOfStock):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
