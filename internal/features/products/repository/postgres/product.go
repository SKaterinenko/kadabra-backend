package products_postgres

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	reviews_model "kadabra/internal/features/reviews/model"
	"time"

	"kadabra/internal/core"
	"kadabra/internal/core/config"
	products_model "kadabra/internal/features/products/model"
	products_service "kadabra/internal/features/products/service"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	db *pgxpool.Pool
}

func NewProductPostgres(db *pgxpool.Pool) *Product {
	return &Product{db: db}
}

func (c *Product) Create(ctx context.Context, req *products_service.CreateInput) (*products_model.ProductWithTranslations, error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query, args, err := config.Psql.
		Insert("products").
		Columns("product_type_id", "manufacturer_id").
		Values(req.ProductTypeId, req.ManufacturerId).
		Suffix("RETURNING id, product_type_id, manufacturer_id, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	productWT, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[products_model.ProductWithoutTranslations])
	if err != nil {
		return nil, core.QueryError(err)
	}

	product := &products_model.ProductWithTranslations{
		Id:             productWT.Id,
		ProductTypeId:  productWT.ProductTypeId,
		ManufacturerId: productWT.ManufacturerId,
		Translations:   make([]*products_model.ProductTranslate, 0),
		CreatedAt:      productWT.CreatedAt,
		UpdatedAt:      productWT.UpdatedAt,
	}

	query, args, err = config.Psql.
		Select("name").
		From("manufacturers").
		Where(sq.Eq{"id": req.ManufacturerId}).
		ToSql()

	rows, err = tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brand string

	for rows.Next() {
		err := rows.Scan(&brand)
		if err != nil {
			return nil, err
		}
	}

	for _, v := range req.Translations {
		slugName := slug.Make(v.Name)
		slugText := brand + "-" + slugName

		query, args, err := config.Psql.
			Insert("product_translations").
			Columns("product_id", "language_code", "name", "slug", "short_description", "description").
			Values(productWT.Id, v.LanguageCode, v.Name, slugText, v.ShortDescription, v.Description).
			Suffix("RETURNING id, product_id, language_code, name, slug, short_description, description, created_at, updated_at").
			ToSql()
		if err != nil {
			return nil, core.BuildSQLError(err)
		}

		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return nil, core.QueryError(err)
		}

		translate, err := pgx.CollectOneRow(
			rows,
			pgx.RowToAddrOfStructByName[products_model.ProductTranslate],
		)
		rows.Close()

		if err != nil {
			return nil, core.QueryError(err)
		}

		product.Translations = append(product.Translations, translate)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (c *Product) loadRatings(ctx context.Context, productIDs []int) (map[int]reviews_model.Rating, error) {
	if len(productIDs) == 0 {
		return make(map[int]reviews_model.Rating), nil
	}

	sql, args, err := config.Psql.
		Select(
			"product_id",
			"COUNT(*) as total_count",
			"COUNT(*) FILTER (WHERE rating = 5) as rating_5",
			"COUNT(*) FILTER (WHERE rating = 4) as rating_4",
			"COUNT(*) FILTER (WHERE rating = 3) as rating_3",
			"COUNT(*) FILTER (WHERE rating = 2) as rating_2",
			"COUNT(*) FILTER (WHERE rating = 1) as rating_1",
		).
		From("reviews").
		Where(sq.Eq{"product_id": productIDs}).
		GroupBy("product_id").
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	type RatingRow struct {
		ProductID  int                  `db:"product_id"`
		Rating     reviews_model.Rating `db:"-"`
		TotalCount int                  `db:"total_count"`
		Rating5    int                  `db:"rating_5"`
		Rating4    int                  `db:"rating_4"`
		Rating3    int                  `db:"rating_3"`
		Rating2    int                  `db:"rating_2"`
		Rating1    int                  `db:"rating_1"`
	}

	ratingsMap := make(map[int]reviews_model.Rating)
	for rows.Next() {
		var row RatingRow
		err := rows.Scan(
			&row.ProductID,
			&row.TotalCount,
			&row.Rating5,
			&row.Rating4,
			&row.Rating3,
			&row.Rating2,
			&row.Rating1,
		)
		if err != nil {
			return nil, core.ScanError(err)
		}

		ratingsMap[row.ProductID] = reviews_model.Rating{
			TotalCount: row.TotalCount,
			Rating5:    row.Rating5,
			Rating4:    row.Rating4,
			Rating3:    row.Rating3,
			Rating2:    row.Rating2,
			Rating1:    row.Rating1,
		}
	}

	return ratingsMap, nil
}

func (c *Product) GetAll(
	ctx context.Context,
	lang string,
	categories, types, manufacturers []int,
	limit, offset int,
) (*products_model.Products, error) {

	// Функция для применения общих условий
	applyFilters := func(q sq.SelectBuilder) sq.SelectBuilder {
		q = q.
			From("products p").
			Join("product_translations pt ON p.id = pt.product_id").
			Join("products_type t ON p.product_type_id = t.id").
			Where(sq.Eq{"pt.language_code": lang})

		if len(types) > 0 {
			q = q.Where(sq.Eq{"p.product_type_id": types})
		}

		if len(manufacturers) > 0 {
			q = q.Where(sq.Eq{"p.manufacturer_id": manufacturers})
		}

		if len(categories) > 0 {
			q = q.
				Join("sub_categories sc ON t.sub_category_id = sc.id").
				Join("categories c ON sc.category_id = c.id").
				Where(sq.Eq{"c.id": categories})
		}

		return q
	}

	dataQuery := applyFilters(
		config.Psql.Select(
			"p.id",
			"pt.name",
			"pt.slug",
			"p.product_type_id",
			"p.manufacturer_id",
			"pt.short_description",
			"pt.description",
			"p.created_at",
			"p.updated_at",
		)).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		OrderBy("pt.name ASC")

	query, args, err := dataQuery.ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	var products []*products_model.Product
	productIDs := make([]int, 0)

	for rows.Next() {
		var emptyProduct products_model.Product
		rows.Scan(
			&emptyProduct.Id,
			&emptyProduct.Name,
			&emptyProduct.Slug,
			&emptyProduct.ProductTypeId,
			&emptyProduct.ManufacturerId,
			&emptyProduct.ShortDescription,
			&emptyProduct.Description,
			&emptyProduct.CreatedAt,
			&emptyProduct.UpdatedAt,
		)

		// Загрузка вариаций
		sql, args, err := config.Psql.
			Select("id", "product_id", "image", "price", "created_at", "updated_at").
			From("product_variations").
			Where(sq.Eq{"product_id": emptyProduct.Id}).
			ToSql()
		if err != nil {
			return nil, core.BuildSQLError(err)
		}

		variationRows, err := c.db.Query(ctx, sql, args...)
		if err != nil {
			return nil, core.QueryError(err)
		}

		variations, err := pgx.CollectRows(variationRows, pgx.RowToAddrOfStructByName[products_model.ProductVariation])
		variationRows.Close()
		if err != nil {
			return nil, core.RowsError(err)
		}

		emptyProduct.Variations = variations
		products = append(products, &emptyProduct)
		productIDs = append(productIDs, emptyProduct.Id)
	}

	// Загрузка рейтингов
	ratingsMap, err := c.loadRatings(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	// Присваивание рейтингов продуктам
	for _, product := range products {
		if rating, ok := ratingsMap[product.Id]; ok {
			product.Rating = rating
		}
	}

	// Запрос для подсчета общего количества
	countQuery := applyFilters(config.Psql.Select("COUNT(*)"))

	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	var totalCount int
	err = c.db.QueryRow(ctx, countSQL, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, core.QueryError(err)
	}

	return &products_model.Products{
		Data:       products,
		TotalCount: totalCount,
	}, nil
}

func (c *Product) GetById(ctx context.Context, id int, lang string) (*products_model.Product, error) {
	query, args, err := config.Psql.
		Select(
			"p.id",
			"pt.name",
			"pt.slug",
			"p.product_type_id",
			"p.manufacturer_id",
			"pt.short_description",
			"pt.description",
			"p.created_at",
			"p.updated_at",
		).
		From("products p").
		Join("product_translations pt on p.id = pt.product_id").
		Where(sq.Eq{"pt.language_code": lang, "p.id": id}).
		OrderBy("pt.name ASC").
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	product, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[products_model.Product])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NoGetRecord
		}
		return nil, core.QueryError(err)
	}

	return product, nil
}

func (c *Product) Delete(ctx context.Context, id int) error {
	query, args, err := config.Psql.Delete("products").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return core.BuildSQLError(err)
	}

	cmdTag, err := c.db.Exec(ctx, query, args...)
	if err != nil {
		return core.QueryError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return core.NoRecordDelete
	}
	return nil
}

func (c *Product) Patch(ctx context.Context, id int, update *products_model.ProductPatch) (*products_model.Product, error) {
	q := config.Psql.
		Update("products").
		Where(sq.Eq{"id": id})

	hasUpdates := false

	if update.Name != nil {
		q = q.Set("name", *update.Name)
		hasUpdates = true
	}

	if !hasUpdates {
		return nil, core.NoUpdate
	}

	query, args, err := q.
		Suffix(`
			RETURNING id,
			name, 
			product_type_id, 
			manufacturer_id, 
			short_description, 
			description, 
			created_at, 
			updated_at`).
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	product, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[products_model.Product])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NoGetRecord
		}
		return nil, core.QueryError(err)
	}

	return product, nil
}

func (c *Product) GetByCategoryIds(ctx context.Context, categoryIds []int, lang string) ([]*products_model.Product, error) {
	query, args, err := config.Psql.
		Select(
			"p.id",
			"pt.name",
			"pt.slug",
			"p.product_type_id",
			"p.manufacturer_id",
			"pt.short_description",
			"pt.description",
			"p.created_at",
			"p.updated_at",
		).
		From("products p").
		Join("product_translations pt on p.id = pt.product_id").
		Join("products_type pty ON p.product_type_id = pty.id").
		Join("sub_categories sc ON pty.sub_category_id = sc.id").
		Where(sq.Eq{"sc.category_id": categoryIds, "pt.language_code": lang}).
		Limit(30).
		OrderBy("p.created_at ASC").
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[products_model.Product])
	if err != nil {
		return nil, core.ScanError(err)
	}

	return products, nil
}

func (c *Product) GetByProductsTypeIds(ctx context.Context, categoryIds []int, lang string) ([]*products_model.Product, error) {
	sql, args, err := config.Psql.Select(
		"p.id",
		"pt.name",
		"pt.slug",
		"p.product_type_id",
		"p.manufacturer_id",
		"pt.short_description",
		"pt.description",
		"p.created_at",
		"p.updated_at").
		From("products p").
		Join("product_translations pt on p.id = pt.product_id").
		Where(sq.Eq{"p.product_type_id": categoryIds, "pt.language_code": lang}).
		Limit(30).
		OrderBy("p.created_at ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[products_model.Product])
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (c *Product) GetByCategorySlug(ctx context.Context, lang, slug string) ([]*products_model.Product, error) {
	sql, args, err := config.Psql.
		Select(
			"p.id",
			"pt.name",
			"pt.slug",
			"p.product_type_id",
			"p.manufacturer_id",
			"pt.short_description",
			"pt.description",
			"p.created_at",
			"p.updated_at",
		).From("products p").
		Join("product_translations pt on pt.product_id = p.id").
		Where(sq.Eq{"pt.language_code": lang}).
		Join("products_type pty on p.product_type_id = pty.id").
		Join("sub_categories sc on pty.sub_category_id = sc.id").
		Join("category_translations ct on sc.category_id = ct.category_id").
		Where(sq.Eq{"ct.slug": slug}).
		Limit(30).
		OrderBy("p.created_at ASC").
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}
	rows, err := c.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[products_model.Product])
	if err != nil {
		return nil, core.ScanError(err)
	}
	return products, nil
}

func (c *Product) GetByManufacturersIds(ctx context.Context, ids []int, lang string) ([]*products_model.Product, error) {
	sql, args, err := config.Psql.Select(
		"p.id",
		"pt.name",
		"pt.slug",
		"p.product_type_id",
		"p.manufacturer_id",
		"pt.short_description",
		"pt.description",
		"p.created_at",
		"p.updated_at",
	).From("products p").
		Join("product_translations pt on pt.product_id = p.id").
		Where(sq.Eq{"pt.language_code": lang}).
		Join("manufacturers m on m.id = p.manufacturer_id").
		Where(sq.Eq{"m.id": ids}).
		Limit(30).
		OrderBy("p.created_at ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()
	products, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[products_model.Product])
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (c *Product) GetBySlug(ctx context.Context, slug, lang string) (*products_model.ProductWithParents, error) {
	sql, args, err := config.Psql.Select(
		// Product fields
		"p.id",
		"pt.name AS name",
		"pt.slug AS slug",
		"p.product_type_id",
		"p.manufacturer_id",
		"pt.short_description AS short_description",
		"pt.description AS description",
		"p.created_at",
		"p.updated_at",

		// ProductType fields
		"ptt.name AS product_type_name",

		// SubCategory fields
		"sc.id AS sub_category_id",
		"sct.name AS sub_category_name",

		// Category fields
		"c.id AS category_id",
		"ct.name AS category_name",
		"ct.slug AS category_slug",

		//Manufacturer
		"m.name AS manufacturer_name",
		"m.slug AS manufacturer_slug",
	).From("products p").
		Join("product_translations pt ON pt.product_id = p.id").
		Join("products_type pty ON pty.id = p.product_type_id").
		Join("product_type_translations ptt ON ptt.product_type_id = pty.id AND ptt.language_code = ?", lang).
		Join("sub_categories sc ON sc.id = pty.sub_category_id").
		Join("sub_category_translations sct ON sct.sub_category_id = sc.id AND sct.language_code = ?", lang).
		Join("categories c ON c.id = sc.category_id").
		Join("category_translations ct ON ct.category_id = c.id AND ct.language_code = ?", lang).
		Join("manufacturers m ON m.id = p.manufacturer_id").
		Where(sq.Eq{"pt.language_code": lang, "pt.slug": slug}).
		ToSql()

	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	type FlatResult struct {
		// Product
		Id               int                                `db:"id"`
		Name             string                             `db:"name"`
		Slug             string                             `db:"slug"`
		ProductTypeId    int                                `db:"product_type_id"`
		ManufacturerId   int                                `db:"manufacturer_id"`
		ShortDescription string                             `db:"short_description"`
		Description      string                             `db:"description"`
		Variations       []*products_model.ProductVariation `db:"-"`
		CreatedAt        time.Time                          `db:"created_at"`
		UpdatedAt        time.Time                          `db:"updated_at"`

		// ProductType
		PTName string `db:"product_type_name"`

		// SubCategory
		SCId   int    `db:"sub_category_id"`
		SCName string `db:"sub_category_name"`

		// Category
		CatId   int    `db:"category_id"`
		CatName string `db:"category_name"`
		CatSlug string `db:"category_slug"`

		//Manufacturer
		MId   int    `db:"-"`
		MName string `db:"manufacturer_name"`
		MSlug string `db:"manufacturer_slug"`
	}

	flat, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[FlatResult])
	if err != nil {
		return nil, err
	}

	sql, args, err = config.Psql.
		Select("id", "product_id", "image", "price", "created_at", "updated_at").
		From("product_variations").
		Where(sq.Eq{"product_id": flat.Id}).
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err = c.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	variations, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[products_model.ProductVariation])
	if err != nil {
		return nil, core.RowsError(err)
	}

	flat.Variations = variations

	// Собираем вложенную структуру
	result := &products_model.ProductWithParents{
		Product: products_model.Product{
			Id:               flat.Id,
			Name:             flat.Name,
			Slug:             flat.Slug,
			ProductTypeId:    flat.ProductTypeId,
			ManufacturerId:   flat.ManufacturerId,
			ShortDescription: flat.ShortDescription,
			Description:      flat.Description,
			Variations:       flat.Variations,
			CreatedAt:        flat.CreatedAt,
			UpdatedAt:        flat.UpdatedAt,
		},
		ProductType: products_model.ProductType{
			Id:   flat.ProductTypeId,
			Name: flat.PTName,
			SubCategory: products_model.SubCategory{
				Id:   flat.SCId,
				Name: flat.SCName,
				Category: products_model.Category{
					Id:   flat.CatId,
					Name: flat.CatName,
					Slug: flat.CatSlug,
				},
			},
		},
		Manufacturer: products_model.Manufacturer{
			Id:   flat.ManufacturerId,
			Name: flat.MName,
			Slug: flat.MSlug,
		},
	}

	return result, nil
}

func (c *Product) CreateProductVariations(ctx context.Context, req *products_service.VariationReq) (*products_model.ProductVariation, error) {
	sql, args, err := config.Psql.Insert("product_variations").
		Columns("product_id", "image", "price").
		Values(req.ProductId, req.Image, req.Price).
		Suffix("RETURNING id, product_id, image, price, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, core.BuildSQLError(err)
	}

	rows, err := c.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, core.QueryError(err)
	}
	defer rows.Close()

	variation, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[products_model.ProductVariation])
	if err != nil {
		return nil, core.RowsError(err)
	}

	return variation, nil
}

func (c *Product) DeleteProductVariation(ctx context.Context, id int) error {
	sql, args, err := config.Psql.Delete("product_variations").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return core.BuildSQLError(err)
	}
	_, err = c.db.Exec(ctx, sql, args...)
	if err != nil {
		return core.QueryError(err)
	}
	return nil
}
