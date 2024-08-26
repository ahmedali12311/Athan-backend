package setting

import (
	"context"

	"github.com/Masterminds/squirrel"
)

func (m *Queries) GetDBTime() (string, error) {
	var now string
	if err := m.DB.GetContext(
		context.Background(),
		&now,
		`SELECT now() as now`,
	); err != nil {
		return "", err
	}
	return now, nil
}

func (m *Queries) GetBoolByKey(key string) (bool, error) {
	var val bool
	if err := m.DB.GetContext(
		context.Background(),
		&val,
		`SELECT value FROM settings WHERE key = $1`,
		key,
	); err != nil {
		return val, err
	}
	return val, nil
}

func (m *Queries) GetForMeta(settings *[]Model) error {
	if err := m.DB.SelectContext(
		context.Background(),
		settings,
		`
            SELECT key, value 
            FROM settings 
            WHERE is_disabled = false
        `,
	); err != nil {
		return err
	}
	return nil
}

func (m *Queries) GetByKeys(keys []string) ([]MinimalModel, error) {
	keyVals := []MinimalModel{}
	query, args, err := m.QB.
		Select("key", "value").
		From("settings").
		Where(squirrel.Eq{"key": keys}).
		ToSql()
	if err != nil {
		return keyVals, err
	}
	if err := m.DB.SelectContext(
		context.Background(),
		&keyVals,
		query,
		args...,
	); err != nil {
		return keyVals, err
	}
	return keyVals, nil
}

func (m *Queries) GetForCosting(settings *[]MinimalModel) error {
	query, args, err := m.QB.
		Select("key", "value").
		From("settings").
		Where(squirrel.Eq{"key": []string{
			KeyOrderMinimum,
			KeyFreeDelivery,
		}}).
		ToSql()
	if err != nil {
		return err
	}
	if err := m.DB.SelectContext(
		context.Background(),
		settings,
		query,
		args...,
	); err != nil {
		return err
	}
	return nil
}

func (m *Queries) GetSASv4() (*SASv4Values, error) {
	var settings []MinimalModel

	query, args, err := m.QB.
		Select("key", "value").
		From("settings").
		Where(squirrel.Eq{"key": []string{
			KeySas4Username,
			KeySas4Password,
			KeySas4Token,
		}}).
		ToSql()
	if err != nil {
		return nil, err
	}
	if err := m.DB.SelectContext(
		context.Background(),
		&settings,
		query,
		args...,
	); err != nil {
		return nil, err
	}
	var sasv4Values SASv4Values
	for _, v := range settings {
		switch v.Key {
		case "sas4_username":
			sasv4Values.Username = v.Value
		case "sas4_password":
			sasv4Values.Password = v.Value
		case "sas4_token":
			sasv4Values.Token = v.Value
		}
	}
	return &sasv4Values, nil
}

// func (m *Queries) GetForTlync(tlyncSettings *tlync.Settings) error {
// 	keyVals := []MinimalModel{}
//
// 	query, args, err := m.QB.
// 		Select("key", "value").
// 		From("settings").
// 		Where(squirrel.Eq{"key": []string{
// 			"tlync_endpoint",
// 			"tlync_token",
// 			"tlync_store_id",
// 			"tlync_front_url",
// 		}}).
// 		ToSql()
// 	if err != nil {
// 		return err
// 	}
// 	if err := m.DB.SelectContext(
//      context.Background(),
//      &keyVals,
//      query,
//      args...,
//  ); err != nil {
// 		return err
// 	}
// 	for _, v := range keyVals {
// 		switch v.Key {
// 		case "tlync_endpoint":
// 			tlyncSettings.Endpoint = v.Value
// 		case "tlync_token":
// 			tlyncSettings.Token = v.Value
// 		case "tlync_store_id":
// 			tlyncSettings.StoreID = v.Value
// 		case "tlync_front_url":
// 			tlyncSettings.FrontURL = v.Value
// 		}
// 	}
// 	return nil
// }
