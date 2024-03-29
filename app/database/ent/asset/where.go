// Code generated by ent, DO NOT EDIT.

package asset

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/willie-lin/cloud-terminal/app/database/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.Asset {
	return predicate.Asset(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.Asset {
	return predicate.Asset(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.Asset {
	return predicate.Asset(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.Asset {
	return predicate.Asset(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.Asset {
	return predicate.Asset(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.Asset {
	return predicate.Asset(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.Asset {
	return predicate.Asset(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.Asset {
	return predicate.Asset(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.Asset {
	return predicate.Asset(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldEQ(FieldUpdatedAt, v))
}

// AssetName applies equality check predicate on the "asset_name" field. It's identical to AssetNameEQ.
func AssetName(v string) predicate.Asset {
	return predicate.Asset(sql.FieldEQ(FieldAssetName, v))
}

// AssetType applies equality check predicate on the "asset_type" field. It's identical to AssetTypeEQ.
func AssetType(v string) predicate.Asset {
	return predicate.Asset(sql.FieldEQ(FieldAssetType, v))
}

// AssetDetails applies equality check predicate on the "asset_details" field. It's identical to AssetDetailsEQ.
func AssetDetails(v string) predicate.Asset {
	return predicate.Asset(sql.FieldEQ(FieldAssetDetails, v))
}

// GroupID applies equality check predicate on the "group_id" field. It's identical to GroupIDEQ.
func GroupID(v int) predicate.Asset {
	return predicate.Asset(sql.FieldEQ(FieldGroupID, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Asset {
	return predicate.Asset(sql.FieldLTE(FieldUpdatedAt, v))
}

// AssetNameEQ applies the EQ predicate on the "asset_name" field.
func AssetNameEQ(v string) predicate.Asset {
	return predicate.Asset(sql.FieldEQ(FieldAssetName, v))
}

// AssetNameNEQ applies the NEQ predicate on the "asset_name" field.
func AssetNameNEQ(v string) predicate.Asset {
	return predicate.Asset(sql.FieldNEQ(FieldAssetName, v))
}

// AssetNameIn applies the In predicate on the "asset_name" field.
func AssetNameIn(vs ...string) predicate.Asset {
	return predicate.Asset(sql.FieldIn(FieldAssetName, vs...))
}

// AssetNameNotIn applies the NotIn predicate on the "asset_name" field.
func AssetNameNotIn(vs ...string) predicate.Asset {
	return predicate.Asset(sql.FieldNotIn(FieldAssetName, vs...))
}

// AssetNameGT applies the GT predicate on the "asset_name" field.
func AssetNameGT(v string) predicate.Asset {
	return predicate.Asset(sql.FieldGT(FieldAssetName, v))
}

// AssetNameGTE applies the GTE predicate on the "asset_name" field.
func AssetNameGTE(v string) predicate.Asset {
	return predicate.Asset(sql.FieldGTE(FieldAssetName, v))
}

// AssetNameLT applies the LT predicate on the "asset_name" field.
func AssetNameLT(v string) predicate.Asset {
	return predicate.Asset(sql.FieldLT(FieldAssetName, v))
}

// AssetNameLTE applies the LTE predicate on the "asset_name" field.
func AssetNameLTE(v string) predicate.Asset {
	return predicate.Asset(sql.FieldLTE(FieldAssetName, v))
}

// AssetNameContains applies the Contains predicate on the "asset_name" field.
func AssetNameContains(v string) predicate.Asset {
	return predicate.Asset(sql.FieldContains(FieldAssetName, v))
}

// AssetNameHasPrefix applies the HasPrefix predicate on the "asset_name" field.
func AssetNameHasPrefix(v string) predicate.Asset {
	return predicate.Asset(sql.FieldHasPrefix(FieldAssetName, v))
}

// AssetNameHasSuffix applies the HasSuffix predicate on the "asset_name" field.
func AssetNameHasSuffix(v string) predicate.Asset {
	return predicate.Asset(sql.FieldHasSuffix(FieldAssetName, v))
}

// AssetNameEqualFold applies the EqualFold predicate on the "asset_name" field.
func AssetNameEqualFold(v string) predicate.Asset {
	return predicate.Asset(sql.FieldEqualFold(FieldAssetName, v))
}

// AssetNameContainsFold applies the ContainsFold predicate on the "asset_name" field.
func AssetNameContainsFold(v string) predicate.Asset {
	return predicate.Asset(sql.FieldContainsFold(FieldAssetName, v))
}

// AssetTypeEQ applies the EQ predicate on the "asset_type" field.
func AssetTypeEQ(v string) predicate.Asset {
	return predicate.Asset(sql.FieldEQ(FieldAssetType, v))
}

// AssetTypeNEQ applies the NEQ predicate on the "asset_type" field.
func AssetTypeNEQ(v string) predicate.Asset {
	return predicate.Asset(sql.FieldNEQ(FieldAssetType, v))
}

// AssetTypeIn applies the In predicate on the "asset_type" field.
func AssetTypeIn(vs ...string) predicate.Asset {
	return predicate.Asset(sql.FieldIn(FieldAssetType, vs...))
}

// AssetTypeNotIn applies the NotIn predicate on the "asset_type" field.
func AssetTypeNotIn(vs ...string) predicate.Asset {
	return predicate.Asset(sql.FieldNotIn(FieldAssetType, vs...))
}

// AssetTypeGT applies the GT predicate on the "asset_type" field.
func AssetTypeGT(v string) predicate.Asset {
	return predicate.Asset(sql.FieldGT(FieldAssetType, v))
}

// AssetTypeGTE applies the GTE predicate on the "asset_type" field.
func AssetTypeGTE(v string) predicate.Asset {
	return predicate.Asset(sql.FieldGTE(FieldAssetType, v))
}

// AssetTypeLT applies the LT predicate on the "asset_type" field.
func AssetTypeLT(v string) predicate.Asset {
	return predicate.Asset(sql.FieldLT(FieldAssetType, v))
}

// AssetTypeLTE applies the LTE predicate on the "asset_type" field.
func AssetTypeLTE(v string) predicate.Asset {
	return predicate.Asset(sql.FieldLTE(FieldAssetType, v))
}

// AssetTypeContains applies the Contains predicate on the "asset_type" field.
func AssetTypeContains(v string) predicate.Asset {
	return predicate.Asset(sql.FieldContains(FieldAssetType, v))
}

// AssetTypeHasPrefix applies the HasPrefix predicate on the "asset_type" field.
func AssetTypeHasPrefix(v string) predicate.Asset {
	return predicate.Asset(sql.FieldHasPrefix(FieldAssetType, v))
}

// AssetTypeHasSuffix applies the HasSuffix predicate on the "asset_type" field.
func AssetTypeHasSuffix(v string) predicate.Asset {
	return predicate.Asset(sql.FieldHasSuffix(FieldAssetType, v))
}

// AssetTypeEqualFold applies the EqualFold predicate on the "asset_type" field.
func AssetTypeEqualFold(v string) predicate.Asset {
	return predicate.Asset(sql.FieldEqualFold(FieldAssetType, v))
}

// AssetTypeContainsFold applies the ContainsFold predicate on the "asset_type" field.
func AssetTypeContainsFold(v string) predicate.Asset {
	return predicate.Asset(sql.FieldContainsFold(FieldAssetType, v))
}

// AssetDetailsEQ applies the EQ predicate on the "asset_details" field.
func AssetDetailsEQ(v string) predicate.Asset {
	return predicate.Asset(sql.FieldEQ(FieldAssetDetails, v))
}

// AssetDetailsNEQ applies the NEQ predicate on the "asset_details" field.
func AssetDetailsNEQ(v string) predicate.Asset {
	return predicate.Asset(sql.FieldNEQ(FieldAssetDetails, v))
}

// AssetDetailsIn applies the In predicate on the "asset_details" field.
func AssetDetailsIn(vs ...string) predicate.Asset {
	return predicate.Asset(sql.FieldIn(FieldAssetDetails, vs...))
}

// AssetDetailsNotIn applies the NotIn predicate on the "asset_details" field.
func AssetDetailsNotIn(vs ...string) predicate.Asset {
	return predicate.Asset(sql.FieldNotIn(FieldAssetDetails, vs...))
}

// AssetDetailsGT applies the GT predicate on the "asset_details" field.
func AssetDetailsGT(v string) predicate.Asset {
	return predicate.Asset(sql.FieldGT(FieldAssetDetails, v))
}

// AssetDetailsGTE applies the GTE predicate on the "asset_details" field.
func AssetDetailsGTE(v string) predicate.Asset {
	return predicate.Asset(sql.FieldGTE(FieldAssetDetails, v))
}

// AssetDetailsLT applies the LT predicate on the "asset_details" field.
func AssetDetailsLT(v string) predicate.Asset {
	return predicate.Asset(sql.FieldLT(FieldAssetDetails, v))
}

// AssetDetailsLTE applies the LTE predicate on the "asset_details" field.
func AssetDetailsLTE(v string) predicate.Asset {
	return predicate.Asset(sql.FieldLTE(FieldAssetDetails, v))
}

// AssetDetailsContains applies the Contains predicate on the "asset_details" field.
func AssetDetailsContains(v string) predicate.Asset {
	return predicate.Asset(sql.FieldContains(FieldAssetDetails, v))
}

// AssetDetailsHasPrefix applies the HasPrefix predicate on the "asset_details" field.
func AssetDetailsHasPrefix(v string) predicate.Asset {
	return predicate.Asset(sql.FieldHasPrefix(FieldAssetDetails, v))
}

// AssetDetailsHasSuffix applies the HasSuffix predicate on the "asset_details" field.
func AssetDetailsHasSuffix(v string) predicate.Asset {
	return predicate.Asset(sql.FieldHasSuffix(FieldAssetDetails, v))
}

// AssetDetailsEqualFold applies the EqualFold predicate on the "asset_details" field.
func AssetDetailsEqualFold(v string) predicate.Asset {
	return predicate.Asset(sql.FieldEqualFold(FieldAssetDetails, v))
}

// AssetDetailsContainsFold applies the ContainsFold predicate on the "asset_details" field.
func AssetDetailsContainsFold(v string) predicate.Asset {
	return predicate.Asset(sql.FieldContainsFold(FieldAssetDetails, v))
}

// GroupIDEQ applies the EQ predicate on the "group_id" field.
func GroupIDEQ(v int) predicate.Asset {
	return predicate.Asset(sql.FieldEQ(FieldGroupID, v))
}

// GroupIDNEQ applies the NEQ predicate on the "group_id" field.
func GroupIDNEQ(v int) predicate.Asset {
	return predicate.Asset(sql.FieldNEQ(FieldGroupID, v))
}

// GroupIDIn applies the In predicate on the "group_id" field.
func GroupIDIn(vs ...int) predicate.Asset {
	return predicate.Asset(sql.FieldIn(FieldGroupID, vs...))
}

// GroupIDNotIn applies the NotIn predicate on the "group_id" field.
func GroupIDNotIn(vs ...int) predicate.Asset {
	return predicate.Asset(sql.FieldNotIn(FieldGroupID, vs...))
}

// GroupIDGT applies the GT predicate on the "group_id" field.
func GroupIDGT(v int) predicate.Asset {
	return predicate.Asset(sql.FieldGT(FieldGroupID, v))
}

// GroupIDGTE applies the GTE predicate on the "group_id" field.
func GroupIDGTE(v int) predicate.Asset {
	return predicate.Asset(sql.FieldGTE(FieldGroupID, v))
}

// GroupIDLT applies the LT predicate on the "group_id" field.
func GroupIDLT(v int) predicate.Asset {
	return predicate.Asset(sql.FieldLT(FieldGroupID, v))
}

// GroupIDLTE applies the LTE predicate on the "group_id" field.
func GroupIDLTE(v int) predicate.Asset {
	return predicate.Asset(sql.FieldLTE(FieldGroupID, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Asset) predicate.Asset {
	return predicate.Asset(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Asset) predicate.Asset {
	return predicate.Asset(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Asset) predicate.Asset {
	return predicate.Asset(sql.NotPredicates(p))
}
