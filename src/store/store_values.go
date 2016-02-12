package store

import (
	"models"
	"github.com/satori/go.uuid"
	"github.com/jackc/pgx"
	// "github.com/golang/glog"
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

func (s *ValueStore) FindByKeys(dto *models.ValueDTO, mode string) (res []models.Value, err error) {
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
	query += fmt.Sprintf(" WHERE keys %s ? AND is_removed = false", mode)
	
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

		// res = append(res, *model)
		_model := models.NewValue()
		_model.TransformFrom(model)
		res = append(res, *_model)

		// TODO: reset for maps
		model.Props = models.StringMap{}
	}

	return res, nil
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

func (s *ValueStore) GetOne(dto *models.ValueDTO) (models.Model, error) {
	model := models.NewValue()
	model.TransformFrom(dto)

	if uuid.Equal(model.ValueId, uuid.Nil) {
		return nil, models.ErrNotValid
	}

	if err := FindModel(model, s.db, dto.Tx, " AND is_removed = false"); err != nil {
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

	if err := UpdateModel(model, s.db, dto.Tx, " AND is_removed = false", dto.UpdateFields...); err != nil {
		return nil, err
	}
	
	return model, nil
}

func (s *ValueStore) Delete(dto *models.ValueDTO) (models.Model, error) {
	if err := models.V.StructPartial(dto, "ValueId"); err != nil {
		return nil, err
	}

	model := models.NewValue()
	model.TransformFrom(dto)

	if err := DeleteModel(model, s.db, dto.Tx, " AND is_removed = false"); err != nil {
		return nil, err
	}
	
	return model, nil
}
