package store

import (
	"models"
	"github.com/satori/go.uuid"
	"github.com/jackc/pgx"
	"fmt"
)

func init() {
	registrationOfStoreBuilder("value", func(sm *StoreManager) Store {
		return NewValueStore(sm)
	})
}

type ValueStore struct {
	*StoreManager
}

func NewValueStore(_store *StoreManager) *ValueStore {

	return &ValueStore{_store}
}

func (_manager ValueStore) ErrorLog(args ...interface{}) {
	_manager.StoreManager.ErrorLog(_manager.Name(), args...)
}

func (_manager ValueStore) Name() string {
	return "value"
}

func (s *ValueStore) FindByKeys(dto *models.ValueDTO, mode string) (res []*models.Value, err error) {
	if err := models.V.StructPartial(dto, "Keys"); err != nil {
		return nil, err
	}

	if _, exist := map[string]bool {
		"&&": true,
		"@>": true,
		"=": true,
	}[mode]; !exist {
		mode = "="
	}

	model := models.NewValue()
	model.TransformFrom(dto)

	fields, modelValues := model.Fields()

	query := SqlSelect(model.TableName(), fields)
	// any - "&&", contains - "@>"", equal - "="
	query += fmt.Sprintf(" WHERE keys %s ?", mode)
	
	query = FormateToPQuery(query)

	var rows *pgx.Rows

	if dto.Tx != nil {
		rows, err = dto.Tx.Query(query, model.Keys)
	} else {
		rows, err = s.db.Query(query, model.Keys)
	}

	defer rows.Close()

	if err != nil {
		s.ErrorLog("action", "поиск записей по keys", "err", err, "keys", model.Keys)
		return
	}

	for rows.Next() {
		if err = rows.Scan(modelValues...); err != nil {
			s.ErrorLog("action", "поиск записей по keys", "subaction", "сканирование строки", "err", err, "keys", model.Keys)
			return
		}

		res = append(res, model)
	}

	return 
}

func (s *ValueStore) Create(dto *models.ValueDTO) (models.Model, error) {
	if err := models.V.StructPartial(dto, "Keys", "Value"); err != nil {
		return nil, err
	}

	model := models.NewValue()
	model.TransformFrom(dto)
	model.ValueId = uuid.NewV1()

	fields, _ := model.Fields()

	if err := CreateModel(model, s.db, dto.Tx, fields...); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *ValueStore) Delete(dto *models.ValueDTO) (error) {
	model, err := s.GetOne(dto)

	if err != nil {

		return err
	}

	return UpdateModel(model, s.db, nil, "update_at", "is_removed")
}

func (s *ValueStore) GetOne(dto *models.ValueDTO) (models.Model, error) {
	model := models.NewValue()
	model.TransformFrom(dto)

	if len(model.Keys) == 0 && uuid.Equal(model.ValueId, uuid.Nil) {
		return nil, models.ErrNotValid
	}

	fields, modelValues := model.Fields()

	query := SqlSelect(model.TableName(), fields)

	args := []interface{}{}

	if uuid.Equal(model.ValueId, uuid.Nil) {
		query += fmt.Sprintf(" WHERE %[1]s = ? LIMIT 1", model.PrimaryName())
		args = append(args, model.PrimaryValue())
	} else {
		query += " WHERE sort_text_array(keys) = sort_text_array(?) LIMIT 1"
		args = append(args, model.Keys)
	}

	query = FormateToPQuery(query)

	var err error

	if dto.Tx != nil {
		err = dto.Tx.QueryRow(query, args...).Scan(modelValues...)
	} else {
		err = s.db.QueryRow(query, args...).Scan(modelValues...)
	}

	if err != nil {
		return nil, err
	}
	
	return model, nil
}

func (s *ValueStore) Update(dto *models.ValueDTO) (models.Model, error) {
	if err := models.V.StructPartial(dto, "ValueId", "Value"); err != nil {
		return nil, err
	}

	model := models.NewValue()
	model.TransformFrom(dto)
	model.BeforeSave()

	if err := UpdateModel(model, s.db, dto.Tx, "value", "props", "flags", "update_at", "is_enabled"); err != nil {
		return nil, err
	}
	
	return model, nil
}
