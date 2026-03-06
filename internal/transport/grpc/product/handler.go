package product

import (
	"CatalogService/internal/app/product/queries/get_product"
	"CatalogService/internal/app/product/queries/list_products"
	"CatalogService/internal/app/product/usecases/activate_product"
	"CatalogService/internal/app/product/usecases/apply_discount"
	"CatalogService/internal/app/product/usecases/create_product"
	"CatalogService/internal/app/product/usecases/update_product"
	pb "CatalogService/proto/product/v1"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Commands struct {
	CreateProduct  *create_product.Interactor
	UpdateProduct  *update_product.Interactor
	Activate       *activate_product.ActivateInteractor
	Deactivate     *activate_product.DeactivateInteractor
	ApplyDiscount  *apply_discount.ApplyInteractor
	RemoveDiscount *apply_discount.RemoveInteractor
}

type Queries struct {
	GetProduct   *get_product.Query
	ListProducts *list_products.Query
}

type Handler struct {
	pb.UnimplementedProductServiceServer
	commands Commands
	queries  Queries
}

func NewHandler(commands Commands, queries Queries) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

func (h *Handler) Register(server *grpc.Server) {
	pb.RegisterProductServiceServer(server, h)
}

func ListGRPCMethods(server *grpc.Server) {
	log.Println("Registered gRPC methods:")
	for serviceName, info := range server.GetServiceInfo() {
		for _, method := range info.Methods {
			log.Printf("  / %s / %s", serviceName, method.Name)
		}
	}
}

func (h *Handler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductReply, error) {
	appReq := create_product.Request{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Category:    req.GetCategory(),
		PriceNum:    req.GetBasePriceNumerator(),
		PriceDenom:  req.GetBasePriceDenominator(),
	}

	id, err := h.commands.CreateProduct.Execute(ctx, appReq)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.CreateProductReply{ProductId: id}, nil
}

func (h *Handler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductReply, error) {
	appReq := update_product.Request{
		ProductID:   req.GetProductId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Category:    req.GetCategory(),
	}

	err := h.commands.UpdateProduct.Execute(ctx, appReq)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.UpdateProductReply{}, nil
}

func (h *Handler) ActivateProduct(ctx context.Context, req *pb.ActivateProductRequest) (*pb.ActivateProductReply, error) {
	err := h.commands.Activate.Execute(ctx, activate_product.Request{ProductID: req.ProductId})
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}
	return &pb.ActivateProductReply{}, nil
}

func (h *Handler) DeactivateProduct(ctx context.Context, req *pb.DeactivateProductRequest) (*pb.DeactivateProductReply, error) {
	err := h.commands.Deactivate.Execute(ctx, activate_product.Request{ProductID: req.ProductId})
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}
	return &pb.DeactivateProductReply{}, nil
}

func (h *Handler) ApplyDiscount(ctx context.Context, req *pb.ApplyDiscountRequest) (*pb.ApplyDiscountReply, error) {
	pct, ok := stringToRat(req.DiscountPercent)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid discount percent")
	}

	appReq := apply_discount.ApplyRequest{
		ProductID:  req.ProductId,
		Percentage: pct,
		StartDate:  req.StartDate.AsTime(),
		EndDate:    req.EndDate.AsTime(),
	}

	err := h.commands.ApplyDiscount.Execute(ctx, appReq)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.ApplyDiscountReply{}, nil
}

func (h *Handler) RemoveDiscount(ctx context.Context, req *pb.RemoveDiscountRequest) (*pb.RemoveDiscountReply, error) {
	err := h.commands.RemoveDiscount.Execute(ctx, apply_discount.RemoveRequest{ProductID: req.ProductId})
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}
	return &pb.RemoveDiscountReply{}, nil
}

func (h *Handler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductReply, error) {
	dto, err := h.queries.GetProduct.Execute(ctx, req.ProductId)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.GetProductReply{
		Product: productDetailToProto(dto),
	}, nil
}

func (h *Handler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsReply, error) {
	res, err := h.queries.ListProducts.Execute(ctx, list_products.Request{
		Category:  req.Category,
		PageSize:  req.PageSize,
		PageToken: req.PageToken,
	})
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	items := make([]*pb.ProductItem, len(res.Products))
	for i, p := range res.Products {
		items[i] = productItemToProto(p)
	}

	return &pb.ListProductsReply{
		Products:      items,
		NextPageToken: res.NextPageToken,
	}, nil
}
