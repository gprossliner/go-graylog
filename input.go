package graylog

import (
	"encoding/json"

	"github.com/suzuki-shunsuke/go-graylog/util"
	"github.com/suzuki-shunsuke/go-ptr"
)

// Input represents Graylog Input.
type Input struct {
	// Select a name of your new input that describes it.
	Title string `json:"title,omitempty" v-create:"required"`
	// https://github.com/Graylog2/graylog2-server/issues/3480
	// update input overwrite attributes
	Attributes InputAttributes `json:"attributes,omitempty" v-create:"required"`

	ID string `json:"id,omitempty" v-create:"isdefault"`

	// Should this input start on all nodes
	Global bool `json:"global,omitempty"`
	// On which node should this input start
	// ex. "2ad6b340-3e5f-4a96-ae81-040cfb8b6024"
	Node string `json:"node,omitempty"`
	// ex. 2018-02-24T03:02:26.001Z
	CreatedAt string `json:"created_at,omitempty" v-create:"isdefault"`
	// ex. "admin"
	CreatorUserID string `json:"creator_user_id,omitempty" v-create:"isdefault"`
	// ContextPack `json:"context_pack,omitempty"`
	// StaticFields `json:"static_fields,omitempty"`
}

func (input Input) Type() string {
	if input.Attributes == nil {
		return ""
	}
	return input.Attributes.InputType()
}

// NewUpdateParams converts Input to InputUpdateParams.
func (input *Input) NewUpdateParams() *InputUpdateParams {
	return &InputUpdateParams{
		ID:         input.ID,
		Title:      input.Title,
		Type:       input.Type(),
		Attributes: input.Attributes,
		Node:       input.Node,
		Global:     ptr.PBool(input.Global),
	}
}

// InputUpdateParams represents Graylog Input update API's parameter.
type InputUpdateParams struct {
	ID         string          `json:"id,omitempty" v-update:"required,objectid"`
	Title      string          `json:"title,omitempty" v-update:"required"`
	Type       string          `json:"type,omitempty" v-update:"required"`
	Attributes InputAttributes `json:"attributes,omitempty" v-update:"required"`
	Global     *bool           `json:"global,omitempty"`
	Node       string          `json:"node,omitempty"`
}

func (input *Input) ToData() (*InputData, error) {
	d := &InputData{
		Title:         input.Title,
		Type:          input.Type(),
		ID:            input.ID,
		Global:        input.Global,
		Node:          input.Node,
		CreatedAt:     input.CreatedAt,
		CreatorUserID: input.CreatorUserID,
		Attributes:    map[string]interface{}{},
	}
	if input.Attributes == nil {
		return d, nil
	}
	return d, util.MSDecode(input.Attributes, &d.Attributes)
}

// UnmarshalJSON is the implementation of the json.Unmarshaler interface.
func (input *Input) UnmarshalJSON(b []byte) error {
	d, err := input.ToData()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, d); err != nil {
		return err
	}
	return d.ToInput(input)
}

// MarshalJSON is the implementation of the json.Marshaler interface.
func (input *Input) MarshalJSON() ([]byte, error) {
	d, err := input.ToData()
	if err != nil {
		return nil, err
	}
	return json.Marshal(d)
}

// InputsBody represents Get Inputs API's response body.
// Basically users don't use this struct, but this struct is public because some sub packages use this struct.
type InputsBody struct {
	Inputs []Input `json:"inputs"`
	Total  int     `json:"total"`
}
