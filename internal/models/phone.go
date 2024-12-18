package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/xinchuantw/hoki-tabloid-backend/utils/database"
)

type Phone struct {
	ID             int        `db:"id" json:"id"`
	Name           string     `db:"name" json:"name"`
	BrandID        int        `db:"brand_id" json:"brand_id"`
	BrandName      string     `db:"brand_name" json:"brand_name"`
	Specifications string     `db:"specifications" json:"specifications"`
	Price          float64    `db:"price" json:"price"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at"`
	PublishedAt    *time.Time `db:"published_at" json:"published_at"`
	Tags           []Tag      `json:"tags"`
}

type Tag struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Installment struct {
	ID           int       `db:"id" json:"id"`
	PhoneID      int       `db:"phone_id" json:"phone_id"`
	ThreeMonths  float64   `db:"three_months" json:"three_months"`
	SixMonths    float64   `db:"six_months" json:"six_months"`
	TwelveMonths float64   `db:"twelve_months" json:"twelve_months"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type PriceHistory struct {
	ID        int       `db:"id" json:"id"`
	PhoneID   int       `db:"phone_id" json:"phone_id"`
	OldPrice  float64   `db:"old_price" json:"old_price"`
	NewPrice  float64   `db:"new_price" json:"new_price"`
	ChangedAt time.Time `db:"changed_at" json:"changed_at"`
}

type Brand struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (b *Brand) Bind(r *http.Request) error { return nil }

func (b *Brand) Insert(tx database.TxQueryer) error {
	query := `INSERT INTO brands (name) VALUES (:name);`
	_, err := tx.NamedExec(query, b)
	if err != nil {
		return fmt.Errorf("[Brand.Insert][NamedExec]%w", err)
	}
	var brandID int
	err = tx.QueryRow("SELECT LAST_INSERT_ID()").Scan(&brandID)
	if err != nil {
		return fmt.Errorf("[Brand.Insert][QueryRow]%w", err)
	}
	b.ID = brandID
	return nil
}

func (b *Brand) Update(tx database.TxQueryer) error {
	query := `UPDATE brands SET name = :name, updated_at = CURRENT_TIMESTAMP WHERE id = :id;`
	_, err := tx.NamedExec(query, b)
	if err != nil {
		return fmt.Errorf("[Brand.Update][NamedExec]%w", err)
	}
	return nil
}

func (b *Brand) Delete(tx database.TxQueryer) error {
	query := "DELETE FROM brands WHERE id = :id;"
	_, err := tx.NamedExec(query, b)
	if err != nil {
		return fmt.Errorf("[Brand.Delete][NamedExec]%w", err)
	}
	return nil
}

func GetBrands(db database.Queryer) ([]Brand, error) {
	brands := []Brand{}
	err := db.Select(&brands, "SELECT * FROM brands;")
	return brands, err
}

func GetBrand(db database.Queryer, id int) (Brand, error) {
	brand := Brand{}
	err := db.Get(&brand, "SELECT * FROM brands WHERE id=?;", id)
	return brand, err
}

func (p *Phone) Bind(r *http.Request) error {
	return nil
}

func (p *Phone) Insert(tx database.TxQueryer) error {
	query := `
    INSERT INTO phones (name, brand_id, specifications, price) 
    VALUES (?, ?, ?, ?)
    RETURNING id
    `
	err := tx.QueryRow(query, p.Name, p.BrandID, p.Specifications, p.Price).Scan(&p.ID)
	if err != nil {
		return fmt.Errorf("[Phone][Insert][QueryRow] %w", err)
	}

	for _, tag := range p.Tags {
		_, err := tx.Exec("INSERT INTO phone_tags (phone_id, tag_id) VALUES (?, ?)", p.ID, tag.ID)
		if err != nil {
			// Handle potential reconnection or other errors here
			return fmt.Errorf("[Phone][Insert][InsertTag] %w", err)
		}
	}

	return nil
}

func (p *Phone) Update(tx database.TxQueryer) error {
	// Get the old price
	var oldPrice float64
	err := tx.Get(&oldPrice, "SELECT price FROM phones WHERE id=?", p.ID)
	if err != nil {
		return fmt.Errorf("[Phone.Update][Get old price]%w", err)
	}

	// Insert price change into PriceHistory
	priceHistory := PriceHistory{
		PhoneID:   p.ID,
		OldPrice:  oldPrice,
		NewPrice:  p.Price,
		ChangedAt: time.Now(),
	}
	err = priceHistory.Insert(tx)
	if err != nil {
		return fmt.Errorf("[Phone.Update][PriceHistory.Insert]%w", err)
	}

	// Update the phone record
	query := `
    UPDATE phones SET name = :name, brand_id = :brand_id, specifications = :specifications, price = :price, published_at = :published_at, updated_at = CURRENT_TIMESTAMP
    WHERE id = :id;
  `
	_, err = tx.NamedExec(query, p)
	if err != nil {
		return fmt.Errorf("[Phone.Update][NamedExec]%w", err)
	}
	// Delete existing tags
	_, err = tx.Exec("DELETE FROM phone_tags WHERE phone_id = ?", p.ID)
	if err != nil {
		return fmt.Errorf("[Phone.Update][DeleteTags]%w", err)
	}
	// Insert updated tags
	for _, tag := range p.Tags {
		_, err := tx.Exec("INSERT INTO phone_tags (phone_id, tag_id) VALUES (?, ?)", p.ID, tag.ID)
		if err != nil {
			return fmt.Errorf("[Phone.Update][InsertTag]%w", err)
		}
	}
	return nil
}

func (p *Phone) Delete(tx database.TxQueryer) error {
	query := "UPDATE phones SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?;"
	_, err := tx.Exec(query, p.ID)
	if err != nil {
		return fmt.Errorf("[Phone.Delete][Exec]%w", err)
	}
	return nil
}
func GetPhones(db database.Queryer, limit, offset int, sortBy, order, filterBy, filterValue string) ([]Phone, error) {
	phones := []Phone{}
	baseQuery := `
    SELECT phones.id, phones.name, phones.brand_id, brands.name AS brand_name, phones.specifications, phones.price, phones.created_at, phones.updated_at, phones.deleted_at, phones.published_at 
    FROM phones 
    JOIN brands ON phones.brand_id = brands.id
    WHERE phones.deleted_at IS NULL
    `
	filterQuery := ""
	if filterBy != "" && filterValue != "" {
		filterQuery = fmt.Sprintf("AND %s LIKE '%%%s%%'", filterBy, filterValue)
	}
	sortQuery := ""
	if sortBy != "" && order != "" {
		sortQuery = fmt.Sprintf("ORDER BY %s %s", sortBy, order)
	}
	paginationQuery := fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	query := fmt.Sprintf("%s %s %s %s", baseQuery, filterQuery, sortQuery, paginationQuery)

	err := db.Select(&phones, query)
	if err != nil {
		return nil, fmt.Errorf("[GetPhones][Select]%w", err)
	}

	for i, phone := range phones {
		tags, err := GetTagsForPhone(db, phone.ID)
		if err != nil {
			return nil, fmt.Errorf("[GetPhones][GetTagsForPhone]%w", err)
		}
		phones[i].Tags = tags
	}

	return phones, nil
}

func GetPhone(db database.Queryer, id int) (Phone, error) {
	phone := Phone{}
	query := `
    SELECT phones.id, phones.name, phones.brand_id, brands.name AS brand_name, phones.specifications, phones.price, phones.created_at, phones.updated_at, phones.deleted_at, phones.published_at 
    FROM phones 
    JOIN brands ON phones.brand_id = brands.id 
    WHERE phones.id = ? AND phones.deleted_at IS NULL;
    `
	err := db.Get(&phone, query, id)
	if err != nil {
		return Phone{}, fmt.Errorf("[GetPhone][Get]%w", err)
	}

	tags, err := GetTagsForPhone(db, phone.ID)
	if err != nil {
		return Phone{}, fmt.Errorf("[GetPhone][GetTagsForPhone]%w", err)
	}
	phone.Tags = tags

	return phone, nil
}

func GetTagsForPhone(db database.Queryer, phoneID int) ([]Tag, error) {
	var tags []Tag

	query := `
    SELECT t.id, t.name
    FROM tags t
    JOIN phone_tags pt ON t.id = pt.tag_id
    WHERE pt.phone_id = ?
  `
	err := db.Select(&tags, query, phoneID)
	if err != nil {
		return nil, fmt.Errorf("[GetTagsForPhone][Select]%w", err)
	}

	return tags, nil
}

func (i *Installment) Bind(r *http.Request) error {
	return nil
}

func (i *Installment) Insert(tx database.TxQueryer) error {
	query := `
    INSERT INTO installments (phone_id, three_months, six_months, twelve_months) 
    VALUES (:phone_id, :three_months, :six_months, :twelve_months);
  `
	_, err := tx.NamedExec(query, i)
	if err != nil {
		return fmt.Errorf("[Installment.Insert][NamedExec]%w", err)
	}
	return nil
}

func CalculateInstallments(price float64) Installment {
	return Installment{
		ThreeMonths:  price / 3,
		SixMonths:    price / 6,
		TwelveMonths: price / 12,
	}
}

func (ph *PriceHistory) Bind(r *http.Request) error {
	return nil
}

func (ph *PriceHistory) Insert(tx database.TxQueryer) error {
	query := `
    INSERT INTO price_history (phone_id, old_price, new_price, changed_at) 
    VALUES (:phone_id, :old_price, :new_price, :changed_at);
  `
	_, err := tx.NamedExec(query, ph)
	if err != nil {
		return fmt.Errorf("[PriceHistory.Insert][NamedExec]%w", err)
	}
	return nil
}
