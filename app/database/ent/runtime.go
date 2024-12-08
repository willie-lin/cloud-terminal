// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/willie-lin/cloud-terminal/app/database/ent/asset"
	"github.com/willie-lin/cloud-terminal/app/database/ent/permission"
	"github.com/willie-lin/cloud-terminal/app/database/ent/resource"
	"github.com/willie-lin/cloud-terminal/app/database/ent/role"
	"github.com/willie-lin/cloud-terminal/app/database/ent/schema"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	assetMixin := schema.Asset{}.Mixin()
	assetMixinFields0 := assetMixin[0].Fields()
	_ = assetMixinFields0
	assetFields := schema.Asset{}.Fields()
	_ = assetFields
	// assetDescCreatedAt is the schema descriptor for created_at field.
	assetDescCreatedAt := assetMixinFields0[0].Descriptor()
	// asset.DefaultCreatedAt holds the default value on creation for the created_at field.
	asset.DefaultCreatedAt = assetDescCreatedAt.Default.(func() time.Time)
	// assetDescUpdatedAt is the schema descriptor for updated_at field.
	assetDescUpdatedAt := assetMixinFields0[1].Descriptor()
	// asset.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	asset.DefaultUpdatedAt = assetDescUpdatedAt.Default.(func() time.Time)
	// asset.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	asset.UpdateDefaultUpdatedAt = assetDescUpdatedAt.UpdateDefault.(func() time.Time)
	// assetDescID is the schema descriptor for id field.
	assetDescID := assetFields[0].Descriptor()
	// asset.DefaultID holds the default value on creation for the id field.
	asset.DefaultID = assetDescID.Default.(func() uuid.UUID)
	permissionMixin := schema.Permission{}.Mixin()
	permissionMixinFields0 := permissionMixin[0].Fields()
	_ = permissionMixinFields0
	permissionFields := schema.Permission{}.Fields()
	_ = permissionFields
	// permissionDescCreatedAt is the schema descriptor for created_at field.
	permissionDescCreatedAt := permissionMixinFields0[0].Descriptor()
	// permission.DefaultCreatedAt holds the default value on creation for the created_at field.
	permission.DefaultCreatedAt = permissionDescCreatedAt.Default.(func() time.Time)
	// permissionDescUpdatedAt is the schema descriptor for updated_at field.
	permissionDescUpdatedAt := permissionMixinFields0[1].Descriptor()
	// permission.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	permission.DefaultUpdatedAt = permissionDescUpdatedAt.Default.(func() time.Time)
	// permission.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	permission.UpdateDefaultUpdatedAt = permissionDescUpdatedAt.UpdateDefault.(func() time.Time)
	// permissionDescID is the schema descriptor for id field.
	permissionDescID := permissionFields[0].Descriptor()
	// permission.DefaultID holds the default value on creation for the id field.
	permission.DefaultID = permissionDescID.Default.(func() uuid.UUID)
	resourceMixin := schema.Resource{}.Mixin()
	resourceMixinFields0 := resourceMixin[0].Fields()
	_ = resourceMixinFields0
	resourceFields := schema.Resource{}.Fields()
	_ = resourceFields
	// resourceDescCreatedAt is the schema descriptor for created_at field.
	resourceDescCreatedAt := resourceMixinFields0[0].Descriptor()
	// resource.DefaultCreatedAt holds the default value on creation for the created_at field.
	resource.DefaultCreatedAt = resourceDescCreatedAt.Default.(func() time.Time)
	// resourceDescUpdatedAt is the schema descriptor for updated_at field.
	resourceDescUpdatedAt := resourceMixinFields0[1].Descriptor()
	// resource.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	resource.DefaultUpdatedAt = resourceDescUpdatedAt.Default.(func() time.Time)
	// resource.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	resource.UpdateDefaultUpdatedAt = resourceDescUpdatedAt.UpdateDefault.(func() time.Time)
	// resourceDescID is the schema descriptor for id field.
	resourceDescID := resourceFields[0].Descriptor()
	// resource.DefaultID holds the default value on creation for the id field.
	resource.DefaultID = resourceDescID.Default.(func() uuid.UUID)
	roleMixin := schema.Role{}.Mixin()
	roleMixinFields0 := roleMixin[0].Fields()
	_ = roleMixinFields0
	roleFields := schema.Role{}.Fields()
	_ = roleFields
	// roleDescCreatedAt is the schema descriptor for created_at field.
	roleDescCreatedAt := roleMixinFields0[0].Descriptor()
	// role.DefaultCreatedAt holds the default value on creation for the created_at field.
	role.DefaultCreatedAt = roleDescCreatedAt.Default.(func() time.Time)
	// roleDescUpdatedAt is the schema descriptor for updated_at field.
	roleDescUpdatedAt := roleMixinFields0[1].Descriptor()
	// role.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	role.DefaultUpdatedAt = roleDescUpdatedAt.Default.(func() time.Time)
	// role.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	role.UpdateDefaultUpdatedAt = roleDescUpdatedAt.UpdateDefault.(func() time.Time)
	// roleDescID is the schema descriptor for id field.
	roleDescID := roleFields[0].Descriptor()
	// role.DefaultID holds the default value on creation for the id field.
	role.DefaultID = roleDescID.Default.(func() uuid.UUID)
	tenantMixin := schema.Tenant{}.Mixin()
	tenantMixinFields0 := tenantMixin[0].Fields()
	_ = tenantMixinFields0
	tenantFields := schema.Tenant{}.Fields()
	_ = tenantFields
	// tenantDescCreatedAt is the schema descriptor for created_at field.
	tenantDescCreatedAt := tenantMixinFields0[0].Descriptor()
	// tenant.DefaultCreatedAt holds the default value on creation for the created_at field.
	tenant.DefaultCreatedAt = tenantDescCreatedAt.Default.(func() time.Time)
	// tenantDescUpdatedAt is the schema descriptor for updated_at field.
	tenantDescUpdatedAt := tenantMixinFields0[1].Descriptor()
	// tenant.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	tenant.DefaultUpdatedAt = tenantDescUpdatedAt.Default.(func() time.Time)
	// tenant.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	tenant.UpdateDefaultUpdatedAt = tenantDescUpdatedAt.UpdateDefault.(func() time.Time)
	// tenantDescID is the schema descriptor for id field.
	tenantDescID := tenantFields[0].Descriptor()
	// tenant.DefaultID holds the default value on creation for the id field.
	tenant.DefaultID = tenantDescID.Default.(func() uuid.UUID)
	userMixin := schema.User{}.Mixin()
	userMixinFields0 := userMixin[0].Fields()
	_ = userMixinFields0
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userMixinFields0[0].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(func() time.Time)
	// userDescUpdatedAt is the schema descriptor for updated_at field.
	userDescUpdatedAt := userMixinFields0[1].Descriptor()
	// user.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	user.DefaultUpdatedAt = userDescUpdatedAt.Default.(func() time.Time)
	// user.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	user.UpdateDefaultUpdatedAt = userDescUpdatedAt.UpdateDefault.(func() time.Time)
	// userDescNickname is the schema descriptor for nickname field.
	userDescNickname := userFields[2].Descriptor()
	// user.NicknameValidator is a validator for the "nickname" field. It is called by the builders before save.
	user.NicknameValidator = func() func(string) error {
		validators := userDescNickname.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(nickname string) error {
			for _, fn := range fns {
				if err := fn(nickname); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescBio is the schema descriptor for bio field.
	userDescBio := userFields[3].Descriptor()
	// user.BioValidator is a validator for the "bio" field. It is called by the builders before save.
	user.BioValidator = userDescBio.Validators[0].(func(string) error)
	// userDescUsername is the schema descriptor for username field.
	userDescUsername := userFields[4].Descriptor()
	// user.UsernameValidator is a validator for the "username" field. It is called by the builders before save.
	user.UsernameValidator = func() func(string) error {
		validators := userDescUsername.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
			validators[2].(func(string) error),
		}
		return func(username string) error {
			for _, fn := range fns {
				if err := fn(username); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescPassword is the schema descriptor for password field.
	userDescPassword := userFields[5].Descriptor()
	// user.PasswordValidator is a validator for the "password" field. It is called by the builders before save.
	user.PasswordValidator = func() func(string) error {
		validators := userDescPassword.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
			validators[2].(func(string) error),
		}
		return func(password string) error {
			for _, fn := range fns {
				if err := fn(password); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescEmail is the schema descriptor for email field.
	userDescEmail := userFields[6].Descriptor()
	// user.EmailValidator is a validator for the "email" field. It is called by the builders before save.
	user.EmailValidator = func() func(string) error {
		validators := userDescEmail.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(email string) error {
			for _, fn := range fns {
				if err := fn(email); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescOnline is the schema descriptor for online field.
	userDescOnline := userFields[9].Descriptor()
	// user.DefaultOnline holds the default value on creation for the online field.
	user.DefaultOnline = userDescOnline.Default.(bool)
	// userDescEnableType is the schema descriptor for enable_type field.
	userDescEnableType := userFields[10].Descriptor()
	// user.DefaultEnableType holds the default value on creation for the enable_type field.
	user.DefaultEnableType = userDescEnableType.Default.(bool)
	// userDescLastLoginTime is the schema descriptor for last_login_time field.
	userDescLastLoginTime := userFields[11].Descriptor()
	// user.DefaultLastLoginTime holds the default value on creation for the last_login_time field.
	user.DefaultLastLoginTime = userDescLastLoginTime.Default.(func() time.Time)
	// userDescID is the schema descriptor for id field.
	userDescID := userFields[0].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() uuid.UUID)
}
