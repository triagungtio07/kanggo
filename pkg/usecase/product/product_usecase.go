package product

import (
	"context"
	"fmt"
	"kanggo/pkg/entity/model"
	"kanggo/pkg/entity/schema"
	storage "kanggo/pkg/storage/product"
)

//go:generate mockery --name ProductUsecase --case snake --output ../../mocks --disable-version-string

type (
	ProductUsecase interface {
		Insert(ctx context.Context, data model.ProductRequest) error
		Update(ctx context.Context, id uint, data model.ProductRequest) error
		GetAll(ctx context.Context) ([]model.ProductResponse, error)
		GetById(ctx context.Context, id int64) (*model.ProductResponse, error)
		Delete(ctx context.Context, id int64) error
	}

	productUsecase struct {
		productStorage storage.ProductStorage
	}
)

func NewProductUsecase(productStorage storage.ProductStorage) ProductUsecase {
	return &productUsecase{
		productStorage: productStorage,
	}
}

func (p *productUsecase) Insert(ctx context.Context, data model.ProductRequest) error {
	request := schema.Product{
		Name:  data.Name,
		Price: data.Price,
		Qty:   int64(data.Qty),
	}

	if err := p.productStorage.Insert(ctx, request); err != nil {
		return err
	}

	return nil
}

func (p *productUsecase) Update(ctx context.Context, id uint, data model.ProductRequest) error {
	request := schema.Product{
		Base:  schema.Base{Id: id},
		Name:  data.Name,
		Price: data.Price,
		Qty:   int64(data.Qty),
	}

	if err := p.productStorage.Update(ctx, request); err != nil {
		return err
	}

	return nil
}

func (p *productUsecase) GetAll(ctx context.Context) ([]model.ProductResponse, error) {
	res, err := p.productStorage.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	results := []model.ProductResponse{}
	for i := range res {
		rest := model.ProductResponse{
			Id:        int(res[i].Id),
			Name:      res[i].Name,
			Price:     float64(res[i].Qty),
			Qty:       int(res[i].Qty),
			CreatedAt: fmt.Sprintf("%v", res[i].CreatedAt),
		}

		results = append(results, rest)
	}

	return results, nil
}

func (p *productUsecase) GetById(ctx context.Context, id int64) (*model.ProductResponse, error) {

	res, err := p.productStorage.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	product := model.ProductResponse{
		Id:        int(res.Id),
		Name:      res.Name,
		Price:     float64(res.Qty),
		Qty:       int(res.Qty),
		CreatedAt: fmt.Sprintf("%v", res.CreatedAt),
	}

	return &product, nil
}

func (p *productUsecase) Delete(ctx context.Context, id int64) error {
	err := p.productStorage.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
