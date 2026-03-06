package product

import (
	"CatalogService/internal/app/product/queries/get_product"
	"CatalogService/internal/app/product/queries/list_products"
	pb "CatalogService/proto/product/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/big"
	"time"
)

func ratToString(r *big.Rat) string {
	if r == nil {
		return ""
	}
	return r.RatString()
}

func stringToRat(s string) (*big.Rat, bool) {
	r := new(big.Rat)
	_, ok := r.SetString(s)
	return r, ok
}

func toTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func productDetailToProto(dto *get_product.ProductDTO) *pb.ProductDetail {
	return &pb.ProductDetail{
		Id:              dto.ID,
		Name:            dto.Name,
		Description:     dto.Description,
		Category:        dto.Category,
		BasePrice:       ratToString(dto.BasePrice),
		EffectivePrice:  ratToString(dto.EffectivePrice),
		DiscountPercent: ratToString(dto.DiscountPercent),
		Status:          dto.Status,
		CreatedAt:       toTimestamp(dto.CreatedAt),
		UpdatedAt:       toTimestamp(dto.UpdatedAt),
	}
}

func productItemToProto(dto *list_products.ProductItemDTO) *pb.ProductItem {
	return &pb.ProductItem{
		Id:             dto.ID,
		Name:           dto.Name,
		Category:       dto.Category,
		BasePrice:      ratToString(dto.BasePrice),
		EffectivePrice: ratToString(dto.EffectivePrice),
		Status:         dto.Status,
		CreatedAt:      toTimestamp(dto.CreatedAt),
	}
}
