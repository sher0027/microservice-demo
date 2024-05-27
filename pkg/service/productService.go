package service

import (
	dto "go-microservice-demo/pkg/dto/product"
	"go-microservice-demo/pkg/model"
	"go-microservice-demo/pkg/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	repository *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{
		repository: repo,
	}
}

func (s *ProductService) CreateProduct(request *dto.ProductRequest) error {
	product := model.Product{
		ID:          primitive.NewObjectID(),
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
	}
	return s.repository.Save(&product)
}

func (s *ProductService) GetAllProducts() ([]dto.ProductResponse, error) {
	products, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.ProductResponse
	for _, product := range products {
		responses = append(responses, dto.ProductResponse{
			ID:          product.ID.Hex(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return responses, nil
}
