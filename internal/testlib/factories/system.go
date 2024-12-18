package factories

import (
	"github.com/go-faker/faker/v4"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/models"
	"github.com/xinchuantw/hoki-tabloid-backend/utils/database"
	"gopkg.in/guregu/null.v4"
)

func CreateSystem(db database.Queryer, with map[string]any) *models.System {
	s := models.System{
		Name: getOrDefault(with, "Name", faker.Name()).(string),
		URL:  getOrDefault(with, "URL", faker.URL()).(string),
	}

	if mapContains(with, "SecretKey") {
		s.PlainSecretKey = null.StringFrom(with["SecretKey"].(string))
	}

	if err := s.Insert(db); err != nil {
		panic(err)
	}
	return &s
}
