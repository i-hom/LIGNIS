package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"
)

func (a App) GetDefects(ctx context.Context, params api.GetDefectsParams) (*api.GetDefectsOK, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "manager" && user.Role != "admin" {
		return nil, errors.New("access denied")
	}

	defects, total, err := a.defectRepo.GetByPatter(params.Pattern.Value, int64(params.Limit.Value), int64(params.Page.Value))
	if err != nil {
		return nil, err
	}
	response := make([]api.Defect, 0)
	for i := range defects {

		defectProducts := make([]api.DefectProduct, 0)
		for i := range defects[i].Defects {
			defectProducts = append(defectProducts, api.DefectProduct{
				ProductID: defects[i].Defects[i].ProductID.Hex(),
				Quantity:  int(defects[i].Defects[i].Quantity),
				Remark:    api.OptString{Set: true, Value: defects[i].Defects[i].Remark},
			})
		}

		response = append(response, api.Defect{
			ID:        defects[i].ID.Hex(),
			Defects:   defectProducts,
			CreatedBy: defects[i].CreatedBy.Hex(),
		})
	}
	return &api.GetDefectsOK{Total: int(total), Defects: response}, nil
}
