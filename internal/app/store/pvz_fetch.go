package store

import (
	"database/sql"
	"pvz_server/internal/app/model"
)

func aggregatePVZResults(rows *sql.Rows) ([]*model.PVZWithReceptions, error) {
	type pvzKey = string

	pvzMap := make(map[pvzKey]*model.PVZWithReceptions)

	for rows.Next() {
		var (
			pvzID, receptionID                  string
			productID                           sql.NullString
			pvzDate, receptionDate, productDate sql.NullTime
			pvzCity                             model.City
			receptionStatus                     model.ReceptionStatus
			productType                         sql.NullString
		)

		err := rows.Scan(
			&pvzID,
			&pvzDate,
			&pvzCity,
			&receptionID,
			&receptionDate,
			&receptionStatus,
			&productID,
			&productDate,
			&productType,
		)

		if err != nil {
			return nil, err
		}

		if _, exists := pvzMap[pvzID]; !exists {
			pvzMap[pvzID] = &model.PVZWithReceptions{
				PVZ: model.PVZ{
					ID:               pvzID,
					RegistrationDate: pvzDate.Time,
					City:             pvzCity,
				},
			}
		}

		if receptionID == "" {
			continue
		}

		pvz := pvzMap[pvzID]
		var currentReception *model.ReceptionWithProducts
		for i := range pvz.Receptions {
			if pvz.Receptions[i].Reception.ID == receptionID {
				currentReception = &pvz.Receptions[i]
				break
			}
		}

		if currentReception == nil {
			currentReception = &model.ReceptionWithProducts{
				Reception: model.Reception{
					ID:       receptionID,
					DateTime: receptionDate.Time,
					PvzID:    pvzID,
					Status:   receptionStatus,
				},
			}
			pvz.Receptions = append(pvz.Receptions, *currentReception)
			currentReception = &pvz.Receptions[len(pvz.Receptions)-1]
		}

		if productID.Valid && productType.Valid {
			currentReception.Products = append(currentReception.Products, model.Product{
				ID:          productID.String,
				DateTime:    productDate.Time,
				Type:        model.ProductType(productType.String),
				ReceptionID: receptionID,
			})
		}
	}

	var result []*model.PVZWithReceptions
	for _, pvz := range pvzMap {
		result = append(result, pvz)
	}
	return result, nil
}
