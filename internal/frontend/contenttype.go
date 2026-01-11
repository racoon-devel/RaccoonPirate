package frontend

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
)

var contentTypeOrder = []media.ContentType{
	media.Movies,
	media.Music,
	media.Other,
}

var mapTypeToRu = map[media.ContentType]string{
	media.Movies: "Фильмы/Сериалы",
	media.Music:  "Музыка",
	media.Other:  "Другое",
}

var mapRuToType = map[string]media.ContentType{
	mapTypeToRu[media.Movies]: media.Movies,
	mapTypeToRu[media.Music]:  media.Music,
	mapTypeToRu[media.Other]:  media.Other,
}

var mapTypeToID = map[media.ContentType]string{
	media.Movies: "movies",
	media.Music:  "music",
	media.Other:  "other",
}

var mapIdToType = map[string]media.ContentType{
	mapTypeToID[media.Movies]: media.Movies,
	mapTypeToID[media.Music]:  media.Music,
	mapTypeToID[media.Other]:  media.Other,
}

func GetContentTypesRu() []string {
	result := make([]string, len(contentTypeOrder))
	for i, t := range contentTypeOrder {
		result[i] = mapTypeToRu[t]
	}
	return result
}

func GetContentTypesButtonsRu() []*communication.Button {
	titles := GetContentTypesRu()
	result := make([]*communication.Button, len(titles))
	for i, t := range titles {
		result[i] = &communication.Button{Title: t, Command: t}
	}
	return result
}

func DetermineContentType(t string) (contentType media.ContentType, ok bool) {
	contentType, ok = mapRuToType[t]
	if ok {
		return
	}

	contentType, ok = mapIdToType[t]
	return
}

func GetContentTypeID(contentType media.ContentType) string {
	return mapTypeToID[contentType]
}
