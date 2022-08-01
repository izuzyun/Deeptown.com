package postgresql

import (
	"database/sql"
	_ "errors"
	_ "log"

	"deeptown.com/deepsearch/pkg/models"
)

// Обертка пула подключения sql.DB
type ProductModel struct {
	DB *sql.DB
}

/*func (m *ProductModel) Search(name string) ([]*models.Product, error) {
	request := `SELECT id, vendor, trust_level, name, price, sold, g_from, g_to, g_way, left_, category, subcategory
			 FROM products WHERE to_tsvector(name) @@ to_tsquery($1)
			 ORDER BY sold
			 DESC LIMIT 100`
	rows, err := m.DB.Query(request, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prods []*models.Product
	for rows.Next() {
		p := &models.Product{}
		err = rows.Scan(&p.ID, &p.Vendor, &p.TrLev, &p.Name, &p.Price, &p.Sold, &p.G_from, &p.G_to, &p.G_way, &p.Left_, &p.Category, &p.Subcategory)
		if err != nil {
			return nil, err
		}
		prods = append(prods, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return prods, nil
}
*/

func (m *ProductModel) SearchProduct(name string) ([]*models.Product, error) {
	request := `SELECT uuid, name, price, sold, category, subcategory, origin_country, dest_country, details, vendor_id, encode(decode(encode(image,'escape'),'base64'),'escape')
			 FROM public."Products" WHERE to_tsvector(name) @@ to_tsquery($1)
			 ORDER BY sold
			 DESC LIMIT 30`
	rows, err := m.DB.Query(request, name)
	if err != nil {
		return nil, err
	}

	var prods []*models.Product
	for rows.Next() {
		p := &models.Product{}
		err = rows.Scan(&p.UUID, &p.Name, &p.Price, &p.Sold, &p.Category, &p.Subcategory, &p.G_from, &p.G_to, &p.G_details, &p.Vendor_id, &p.Image)
		if err != nil {
			return nil, err
		}
		prods = append(prods, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	defer rows.Close()
	return prods, nil
}

func (m *ProductModel) GetVendor(products []*models.Product) ([]*models.Information, error) {
	request := `SELECT uuid, name, all_sales, link, trust_level, market
				FROM public."Vendor" WHERE id = $1`
	var information []*models.Information
	for _, element := range products {
		d := element.Vendor_id
		rows, err := m.DB.Query(request, d)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			s := &models.Vendor{}
			err = rows.Scan(&s.UUID, &s.Name, &s.All_sales, &s.Link, &s.Trust_level, &s.Market)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			m := &models.Information{element.UUID, element.Name, element.Price, element.Sold, element.Category, element.Subcategory, element.G_from, element.G_to, element.G_details, element.Image, s.UUID, s.Name, s.All_sales, s.Link, s.Trust_level, s.Market}
			information = append(information, m)
		}
	}
	return information, nil
}

func (m *ProductModel) Stat(name string) ([]*models.Statistics, error) {
	request := `SELECT date, category, subcategory, revenue, vendor_amount, products_amount, sold, avprice
				FROM public."Markets_statistics" WHERE market = 'Darkfox' AND date = (SELECT DISTINCT(MAX(date)) FROM public."Markets_statistics")`
	rows, err := m.DB.Query(request)
	if err != nil {
		return nil, err
	}
	var statistics []*models.Statistics
	for rows.Next() {
		p := &models.Statistics{}
		err = rows.Scan(&p.Date, &p.Category, &p.Subcategory, &p.Revenue, &p.Vendor_amount, &p.Products_amount, &p.Sold, &p.Avprice)
		if err != nil {
			return nil, err
		}
		statistics = append(statistics, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()
	return statistics, nil
}

/*func (m *ProductModel) ShowTen() ([]*models.Product, error) {
	request := `SELECT id, vendor, trust_level, name, price, sold, g_from, g_to, g_way, left_, category, subcategory
			 FROM products ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(request)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prods []*models.Product
	for rows.Next() {
		p := &models.Product{}
		err = rows.Scan(&p.ID, &p.Vendor, &p.TrLev, &p.Name, &p.Price, &p.Sold, &p.G_from, &p.G_to, &p.G_way, &p.Left_, &p.Category, &p.Subcategory)
		if err != nil {
			return nil, err
		}
		prods = append(prods, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return prods, nil
}
*/

/*func (m *ProductModel) Get(id int) (*models.Product, error) {
	request := `SELECT id, vendor, trust_level, name, price, sold, g_from, g_to, g_way, left_, category, subcategory
			 	FROM products WHERE id = $1`
	row := m.DB.QueryRow(request, id)
	p := &models.Product{}
	err := row.Scan(&p.ID, &p.Vendor, &p.TrLev, &p.Name, &p.Price, &p.Sold, &p.G_from, &p.G_to, &p.G_way, &p.Left_, &p.Category, &p.Subcategory)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return p, nil
}*/
