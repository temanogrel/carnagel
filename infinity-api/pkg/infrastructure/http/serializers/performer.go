package serializers

import (
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/tuvistavie/structomap"
)

type PerformerSerializer struct {
	*structomap.Base
}

func NewPerformerSerializer() *PerformerSerializer {

	serializer := &PerformerSerializer{structomap.New()}
	serializer.UseCamelCase()
	serializer.Pick("Uuid", "Slug", "StageName", "Aliases", "RecordingCount", "OriginService", "CreatedAt", "UpdatedAt")
	serializer.AddFunc("LatestRecording", func(performer interface{}) interface{} {
		if len(performer.(*infinity.Performer).Recordings) > 0 {
			return performer.(*infinity.Performer).Recordings[0]
		}

		return nil
	})
	serializer.AddFunc("Section", func(performer interface{}) interface{} {
		return performer.(*infinity.Performer).OriginSection.String
	})

	return serializer
}

// This is used when we nest the performer in the recording serializer
func (s *PerformerSerializer) Minimal() *PerformerSerializer {
	s.Omit("Aliases", "CreatedAt", "UpdatedAt")

	return s
}
