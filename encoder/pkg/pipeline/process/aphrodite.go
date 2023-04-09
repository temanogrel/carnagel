package process

import (
	"context"
	"git.misc.vee.bz/carnagel/encoder/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
)

func AphroditeUpdateRecording(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	if err := app.AphroditeClient.UpdateRecording(rec); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func AphroditeDeleteRecording(ctx context.Context, app *encoder.Application, rec *ecosystem.Recording) (context.Context, error) {
	if err := app.AphroditeClient.DeleteRecording(rec.Id); err != nil {
		return ctx, err
	}

	return ctx, nil
}
