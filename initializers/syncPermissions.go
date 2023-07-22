package initializers

import (
	"fmt"
	"strings"

	"github.com/guptaharsh13/balkanid-task/models"
)

var PermissibleTables = []string{"users", "tasks", "roles", "permissions"}
var Operations = []string{"CREATE", "READ", "UPDATE", "DELETE"}

func SyncPermissions() {

	var permissions []models.Permission

	for _, table := range PermissibleTables {
		for _, operation := range Operations {

			name := fmt.Sprintf("%s_%s", strings.ToLower(operation), table)
			if result := DB.Take(&models.Permission{}, "name = ?", name); result.RowsAffected > 0 {
				continue
			}
			permission := models.Permission{
				Name:        name,
				Description: fmt.Sprintf("This permission allows any user to perform %s operation on the %s table.", string(operation), table),
				Table:       table,
				Operation:   operation,
			}
			permissions = append(permissions, permission)
		}
	}

	if len(permissions) != 0 {
		if result := DB.Create(&permissions); result.Error != nil {
			panic(fmt.Sprintf("Couldn't create permissions: %s", result.Error))
		}
	}
	fmt.Println("✅ Permissions Synced")
}
