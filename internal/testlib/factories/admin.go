package factories

import (
	"fmt"
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/models"
	"github.com/xinchuantw/hoki-tabloid-backend/utils/database"
)

func CreateAdmin(db database.Queryer, with map[string]any) *models.Admin {
	admin := models.Admin{
		Name:          getOrDefault(with, "Name", faker.Name()).(string),
		Username:      getOrDefault(with, "Username", faker.Username()).(string),
		Provider:      getOrDefault(with, "Provider", "xinchuan-auth").(string),
		ProviderID:    getOrDefault(with, "ProviderID", fmt.Sprintf("%d", rand.Int())).(string),
		DeactivatedAt: getAsTime(with, "DeactivatedAt"),
	}
	if err := admin.Insert(db); err != nil {
		panic(err)
	}
	return &admin
}
