package app

import (
	"context"

	"lignis/internal/generated/api"
)

func (a *App) GetImage(ctx context.Context, params api.GetImageParams) (api.GetImageOK, error) {
	obj, err := a.minio.Download(ctx, params.ID)
	return api.GetImageOK{Data: obj}, err
}
