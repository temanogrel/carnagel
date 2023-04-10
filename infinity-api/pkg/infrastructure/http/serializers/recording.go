package serializers

import (
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/tuvistavie/structomap"
)

func hasPerformer(recording interface{}) bool {
	if rec, ok := recording.(*infinity.Recording); ok {
		return rec.Performer != nil
	}

	if rec, ok := recording.(infinity.Recording); ok {
		return rec.Performer != nil
	}

	return true
}

func hasUserCtx(recording interface{}) bool {
	if _, ok := recording.(infinity.RecordingWithUserData); ok {
		return true
	}

	if _, ok := recording.(*infinity.RecordingWithUserData); ok {
		return true
	}

	return false
}

type RecordingImageSerializer struct {
	*structomap.Base
}

func (s *RecordingImageSerializer) WithStorageInfo() *RecordingImageSerializer {

	return s
}

type RecordingSerializer struct {
	*structomap.Base

	performerSerializer *PerformerSerializer
}

func (s *RecordingSerializer) WithStorageInfo() *RecordingSerializer {
	return s
}

func NewRecordingSerializer() *RecordingSerializer {
	fields := []string{
		"Uuid",
		"Slug",
		"PerformerUuid",
		"VideoUuid",
		"CollageUuid",
		"Sprites",
		"Images",
		"StageName",
		"ViewCount",
		"LikeCount",
		"Duration",
		"CreatedAt",
		"UpdatedAt",
	}

	serializer := &RecordingSerializer{
		structomap.New(),
		NewPerformerSerializer().Minimal(),
	}

	serializer.UseCamelCase()
	serializer.Pick(fields...)

	// Embedded the performer if it exists
	serializer.PickFuncIf(hasPerformer, func(performer interface{}) interface{} {
		return serializer.performerSerializer.Transform(performer)
	}, "Performer")

	serializer.PickIf(hasUserCtx, "IsLiked", "IsFavorite")

	// No need to double expose the performer uuid
	serializer.OmitIf(hasPerformer, "PerformerUuid")

	return serializer
}
