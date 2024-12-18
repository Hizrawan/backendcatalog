package factories

import (
	"log"

	"github.com/xinchuantw/hoki-tabloid-backend/internal/models"
	"github.com/xinchuantw/hoki-tabloid-backend/utils/database"
	"gopkg.in/guregu/null.v4"
)

func AssignRoleToAdmin(db database.Queryer, permissionIdentifier string, admin *models.Admin) {
	permission := models.Permission{
		Identifier:  permissionIdentifier,
		Module:      "Test",
		Name:        "Test",
		Description: "for testing purpose",
	}
	err := permission.Insert(db)
	if err != nil {
		log.Fatalln(err)
	}

	role := models.Role{
		Permissions: []models.Permission{permission},
		Name:        "Test",
	}
	err = role.Insert(db)
	if err != nil {
		log.Fatalln(err)
	}

	auth := models.Authorities{
		RoleID:       role.ID,
		PermissionID: permission.ID,
	}
	err = auth.Insert(db)
	if err != nil {
		log.Fatalln(err)
	}

	admin.RoleID = null.StringFrom(role.ID)
	err = admin.Update(db)
	if err != nil {
		log.Fatalln(err)
	}
}
