package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/go-graylog"
	"github.com/suzuki-shunsuke/go-graylog/mockserver/logic"
	"github.com/suzuki-shunsuke/go-set"
)

// HandleGetStreams
func HandleGetStreams(
	user *graylog.User, ms *logic.Logic,
	w http.ResponseWriter, r *http.Request, _ httprouter.Params,
) (interface{}, int, error) {
	// GET /streams Get a list of all streams
	arr, err := ms.GetStreams()
	if err != nil {
		ms.Logger().WithFields(log.Fields{
			"error": err,
		}).Error("ms.StreamList() is failure")
		return nil, 500, err
	}

	return &graylog.StreamsBody{Streams: arr, Total: len(arr)}, 200, nil
}

// HandleCreateStream
func HandleCreateStream(
	user *graylog.User, ms *logic.Logic,
	w http.ResponseWriter, r *http.Request, _ httprouter.Params,
) (interface{}, int, error) {
	// POST /streams Create index set
	if sc, err := ms.Authorize(user, "streams:create"); err != nil {
		return nil, sc, err
	}
	requiredFields := set.NewStrSet("title", "index_set_id")
	allowedFields := set.NewStrSet(
		"rules", "description", "content_pack",
		"matching_type", "remove_matches_from_default_stream")
	body, sc, err := validateRequestBody(r.Body, requiredFields, allowedFields, nil)
	if err != nil {
		return nil, sc, err
	}

	stream := &graylog.Stream{}
	if err := msDecode(body, stream); err != nil {
		ms.Logger().WithFields(log.Fields{
			"body": body, "error": err,
		}).Info("Failed to parse request body as stream")
		return nil, 400, err
	}

	sc, err = ms.AddStream(stream)
	if err != nil {
		return nil, 400, err
	}
	return map[string]string{"stream_id": stream.ID}, 201, nil
}

// HandleGetEnabledStreams
func HandleGetEnabledStreams(
	user *graylog.User, ms *logic.Logic,
	w http.ResponseWriter, r *http.Request, _ httprouter.Params,
) (interface{}, int, error) {
	// GET /streams/enabled Get a list of all enabled streams
	arr, err := ms.EnabledStreamList()
	if err != nil {
		ms.Logger().WithFields(log.Fields{
			"error": err,
		}).Error("ms.EnabledStreamList() is failure")
		return nil, 500, err
	}
	return &graylog.StreamsBody{Streams: arr, Total: len(arr)}, 200, nil
}

// HandleGetStream
func HandleGetStream(
	user *graylog.User, ms *logic.Logic,
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
) (interface{}, int, error) {
	// GET /streams/{streamID} Get a single stream
	id := ps.ByName("streamID")
	if id == "enabled" {
		return HandleGetEnabledStreams(user, ms, w, r, ps)
	}
	if sc, err := ms.Authorize(user, "streams:read", id); err != nil {
		return nil, sc, err
	}
	stream, err := ms.GetStream(id)
	if err != nil {
		ms.Logger().WithFields(log.Fields{
			"error": err, "id": id,
		}).Error("ms.GetStream() is failure")
		return nil, 500, err
	}
	if stream == nil {
		return nil, 404, fmt.Errorf("no stream found with id <%s>", id)
	}
	return stream, 200, nil
}

// HandleUpdateStream
func HandleUpdateStream(
	user *graylog.User, ms *logic.Logic,
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
) (interface{}, int, error) {
	// PUT /streams/{streamID} Update a stream
	id := ps.ByName("streamID")
	if sc, err := ms.Authorize(user, "streams:edit", id); err != nil {
		return nil, sc, err
	}
	stream, err := ms.GetStream(id)
	if err != nil {
		ms.Logger().WithFields(log.Fields{
			"error": err, "id": id,
		}).Error("ms.GetStream() is failure")
		return nil, 500, err
	}
	if stream == nil {
		return nil, 404, fmt.Errorf("no stream found with id <%s>", id)
	}
	data := map[string]interface{}{}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&data); err != nil {
		return nil, 400, err
	}
	if title, ok := data["title"]; ok {
		t, ok := title.(string)
		if !ok {
			return nil, 400, fmt.Errorf("title must be string")
		}
		stream.Title = t
	}
	if description, ok := data["description"]; ok {
		d, ok := description.(string)
		if !ok {
			return nil, 400, fmt.Errorf("description must be string")
		}
		stream.Description = d
	}
	// TODO outputs
	if matchingType, ok := data["matching_type"]; ok {
		m, ok := matchingType.(string)
		if !ok {
			return nil, 400, fmt.Errorf("matching_type must be string")
		}
		stream.MatchingType = m
	}
	if removeMathcesFromDefaultStream, ok := data["remove_matches_from_default_stream"]; ok {
		m, ok := removeMathcesFromDefaultStream.(bool)
		if !ok {
			return nil, 400, fmt.Errorf("remove_matches_from_default_stream must be bool")
		}
		stream.RemoveMatchesFromDefaultStream = m
	}
	if indexSetID, ok := data["index_set_id"]; ok {
		m, ok := indexSetID.(string)
		if !ok {
			return nil, 400, fmt.Errorf("index_set_id must be string")
		}
		stream.IndexSetID = m
	}
	stream.ID = id
	if sc, err := ms.UpdateStream(stream); err != nil {
		return nil, sc, err
	}
	if err := ms.Save(); err != nil {
		return nil, 500, err
	}
	return stream, 200, nil
}

// HandleDeleteStream
func HandleDeleteStream(
	user *graylog.User, ms *logic.Logic,
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
) (interface{}, int, error) {
	// DELETE /streams/{streamID} Delete a stream
	id := ps.ByName("streamID")
	// TODO authorization
	sc, err := ms.DeleteStream(id)
	if err != nil {
		return nil, sc, err
	}
	if err := ms.Save(); err != nil {
		return nil, 500, err
	}
	return nil, sc, nil
}

// HandlePauseStream
func HandlePauseStream(
	user *graylog.User, ms *logic.Logic,
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
) (interface{}, int, error) {
	// POST /streams/{streamID}/pause Pause a stream
	id := ps.ByName("streamID")
	if sc, err := ms.Authorize(user, "streams:changestate", id); err != nil {
		return nil, sc, err
	}
	ok, err := ms.HasStream(id)
	if err != nil {
		ms.Logger().WithFields(log.Fields{
			"error": err, "id": id,
		}).Error("ms.HasStream() is failure")
		return nil, 500, err
	}
	if !ok {
		return nil, 404, fmt.Errorf("no stream found with id <%s>", id)
	}
	// TODO pause
	return nil, 200, nil
}

func HandleResumeStream(
	user *graylog.User, ms *logic.Logic,
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
) (interface{}, int, error) {
	id := ps.ByName("streamID")
	if sc, err := ms.Authorize(user, "streams:changestate", id); err != nil {
		return nil, sc, err
	}
	ok, err := ms.HasStream(id)
	if err != nil {
		ms.Logger().WithFields(log.Fields{
			"error": err, "id": id,
		}).Error("ms.HasStream() is failure")
		return nil, 500, err
	}
	if !ok {
		return nil, 404, fmt.Errorf("no stream found with id <%s>", id)
	}
	// TODO resume
	return nil, 200, nil
}
