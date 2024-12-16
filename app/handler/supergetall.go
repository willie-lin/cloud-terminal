package handler

//func GetAllTenantEntities(client *ent.Client) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		v := viewer.FromContext(c.Request().Context())
//		if v == nil || v.RoleName != "Tenant Admin" {
//			log.Printf("No viewer found in context or not authorized")
//			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
//		}
//
//		users, err := client.User.Query().Where(user.TenantIDEQ(v.TenantID)).All(context.Background())
//		if err != nil {
//			log.Printf("Error querying users: %v", err)
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying users"})
//		}
//
//		roles, err := client.Role.Query().Where(role.TenantIDEQ(v.TenantID)).All(context.Background())
//		if err != nil {
//			log.Printf("Error querying roles: %v", err)
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying roles"})
//		}
//
//		permissions, err := client.Permission.Query().Where(permission.TenantIDEQ(v.TenantID)).All(context.Background())
//		if err != nil {
//			log.Printf("Error querying permissions: %v", err)
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying permissions"})
//		}
//
//		resources, err := client.Resource.Query().Where(resource.TenantIDEQ(v.TenantID)).All(context.Background())
//		if err != nil {
//			log.Printf("Error querying resources: %v", err)
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying resources"})
//		}
//
//		response := map[string]interface{}{
//			"users":       users,
//			"roles":       roles,
//			"permissions": permissions,
//			"resources":   resources,
//		}
//
//		return c.JSON(http.StatusOK, response)
//	}
//}
//
