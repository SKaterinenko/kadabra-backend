package products_service

import (
	"context"
	"fmt"
	"kadabra/internal/core/config"
	products_model "kadabra/internal/features/products/model"
	"strings"
	"time"
)

type Service struct {
	repo        ProductRepository
	s3Client    *config.S3Client
	cacheClient *config.Cache
}

func NewService(repo ProductRepository, s3Client *config.S3Client, cacheClient *config.Cache) *Service {
	return &Service{repo: repo, s3Client: s3Client, cacheClient: cacheClient}
}

func (s *Service) Create(ctx context.Context, product *CreateInput) (*products_model.ProductWithTranslations, error) {
	out, err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) GetAll(ctx context.Context, lang string, categories, types, manufacturers []int, limit, offset int) (*products_model.Products, error) {
	if limit == 0 {
		limit = 20
	}

	cacheKey := fmt.Sprintf("products:all:%s:%v:%v:%v:%d:%d",
		lang, categories, types, manufacturers, limit, offset)

	// Пробуем кеш
	var cached products_model.Products
	err := s.cacheClient.Get(ctx, cacheKey, &cached)
	if err == nil {
		return &cached, nil
	}

	// Идём в БД
	out, err := s.repo.GetAll(ctx, lang, categories, types, manufacturers, limit, offset)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш на 5 минут
	err = s.cacheClient.Set(ctx, cacheKey, out, 5*time.Minute)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) GetById(ctx context.Context, id int, lang string) (*products_model.Product, error) {
	out, err := s.repo.GetById(ctx, id, lang)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Patch(ctx context.Context, id int, update *PatchInput) (*products_model.Product, error) {
	newPatch := products_model.NewProductPatch(
		*update.Name,
		*update.Description,
		*update.ShortDescription,
		*update.ProductTypeId,
		*update.ManufacturerId)
	out, err := s.repo.Patch(ctx, id, newPatch)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetByCategoryIds(ctx context.Context, categoryIds []int, lang string) ([]*products_model.Product, error) {
	if len(categoryIds) == 0 {
		return []*products_model.Product{}, nil
	}

	out, err := s.repo.GetByCategoryIds(ctx, categoryIds, lang)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetByProductsTypeIds(ctx context.Context, productsTypeIds []int, lang string) ([]*products_model.Product, error) {
	if len(productsTypeIds) == 0 {
		return []*products_model.Product{}, nil
	}

	out, err := s.repo.GetByProductsTypeIds(ctx, productsTypeIds, lang)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Service) GetByCategorySlug(ctx context.Context, lang, slug string) ([]*products_model.Product, error) {
	if slug == "" {
		return []*products_model.Product{}, nil
	}
	products, err := s.repo.GetByCategorySlug(ctx, lang, slug)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Service) GetByManufacturersIds(ctx context.Context, ids []int, lang string) ([]*products_model.Product, error) {
	products, err := s.repo.GetByManufacturersIds(ctx, ids, lang)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Service) GetBySlug(ctx context.Context, slug, lang string) (*products_model.ProductWithParents, error) {
	newSlug := strings.ToUpper(slug[:1]) + slug[1:]
	product, err := s.repo.GetBySlug(ctx, newSlug, lang)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Service) CreateProductVariations(ctx context.Context, input *VariationInput) (*products_model.ProductVariation, error) {
	imageURL, err := s.s3Client.UploadFile(ctx, input.Image)
	if err != nil {
		return nil, err
	}

	req := &VariationReq{
		ProductId: input.ProductId,
		Image:     imageURL,
		Price:     input.Price,
	}

	variation, err := s.repo.CreateProductVariations(ctx, req)
	if err != nil {
		return nil, err
	}

	return variation, nil
}

func (s *Service) DeleteProductVariation(ctx context.Context, id int) error {
	err := s.repo.DeleteProductVariation(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
