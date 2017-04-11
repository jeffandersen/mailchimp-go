package mergefields

import (
	"fmt"
	"strconv"
	"strings"

	mailchimp "github.com/beeker1121/mailchimp-go"
	"github.com/beeker1121/mailchimp-go/query"
)

// FieldType defines the type of field asked for
type FieldType string

// String implements the string interface for FieldType
func (ft *FieldType) String() string {
	return string(*ft)
}

// The merge field type definitions
const (
	TypeText        FieldType = "text"
	TypeAddress     FieldType = "address"
	TypeBirthday    FieldType = "birthday"
	TypeDate        FieldType = "date"
	TypeDropdown    FieldType = "dropdown"
	TypeImage       FieldType = "imageurl"
	TypeNumber      FieldType = "number"
	TypePhoneNumber FieldType = "phone"
	TypeRadio       FieldType = "radio"
)

// ListMergeFields defines the list of merge fields
type ListMergeFields struct {
	MergeFields []MergeField `json:"merge_fields"`
	ListID      string       `json:"list_id"`
	TotalItems  int          `json:"total_items"`
}

// MergeField defines a single merge field within a list
type MergeField struct {
	MergeID      int      `json:"merge_id"`
	Tag          string   `json:"tag"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Required     bool     `json:"required"`
	DefaultValue string   `json:"default_value"`
	Public       bool     `json:"public"`
	DisplayOrder int      `json:"display_order"`
	Options      *Options `json:"options"`
	HelpText     string
	ListID       string `json:"list_id"`
}

// Options defines a merge field's options
type Options struct {
	DefaultCountry int      `json:"default_country,omitempty"`
	PhoneFormat    string   `json:"phone_format,omitempty"`
	DateFormat     string   `json:"date_format,omitempty"`
	Choices        []string `json:"choices,omitempty"`
	Size           int      `json:"size,omitempty"`
}

// NewParams defines the available parameters that can be used when
// adding a new merge field to a list via the New function
type NewParams struct {
	Tag          string   `json:"tag,omitempty"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Required     bool     `json:"required,omitempty"`
	DefaultValue string   `json"default_value,omitempty"`
	Public       bool     `json:"public,omitempty"`
	DisplayOrder int      `json:"display_order,omitempty"`
	Options      *Options `json:"options,omitempty"`
	HelpText     string   `json:"help_text,omitempty"`
}

// UpdateParams defines the available parameters that can be used when
// updating an existing merge field for a list via the Update function
type UpdateParams struct {
	Tag          string   `json:"tag,omitempty"`
	Name         string   `json:"name"`
	Required     bool     `json:"required,omitempty"`
	DefaultValue string   `json"default_value,omitempty"`
	Public       bool     `json:"public,omitempty"`
	DisplayOrder int      `json:"display_order,omitempty"`
	Options      *Options `json:"options,omitempty"`
	HelpText     string   `json:"help_text,omitempty"`
}

// GetParams defines the available parameters that can be used when
// retrieving a list of merge fields via the Get function
type GetParams struct {
	Fields        []string  `url:"fields,omitempty"`
	ExcludeFields []string  `url:"exclude_fields,omitempty"`
	Count         int       `url:"count,omitempty"`
	Offset        int       `url:"offset,omitempty"`
	FieldType     FieldType `url:"type,omitempty"`
	Required      bool      `url:"required,omitempty"`
}

// EncodeQueryString handles custom query string encoding for the
// GetParams object.
func (gp *GetParams) EncodeQueryString(v interface{}) (string, error) {
	return query.Encode(struct {
		Fields        string    `url:"fields,omitempty"`
		ExcludeFields string    `url:"exclude_fields,omitempty"`
		Count         int       `url:"count,omitempty"`
		Offset        int       `url:"offset,omitempty"`
		FieldType     FieldType `url:"field_type,omitempty"`
		Required      bool      `url:"required,omitempty"`
	}{
		Fields:        strings.Join(gp.Fields, ","),
		ExcludeFields: strings.Join(gp.ExcludeFields, ","),
		Count:         gp.Count,
		Offset:        gp.Offset,
		FieldType:     gp.FieldType,
		Required:      gp.Required,
	})
}

// GetMergeFieldParams defines the available parameters that can be used when
// retrieving a single merge field via the GetMergeField function
type GetMergeFieldParams struct {
	Fields        []string `url:"fields,omitempty"`
	ExcludeFields []string `url:"exclude_fields,omitempty"`
}

// EncodeQueryString handles custom query string encoding for the
// GetMergeFieldParams object.
func (gp *GetMergeFieldParams) EncodeQueryString(v interface{}) (string, error) {
	return query.Encode(struct {
		Fields        string `url:"fields,omitempty"`
		ExcludeFields string `url:"exclude_fields,omitempty"`
	}{
		Fields:        strings.Join(gp.Fields, ","),
		ExcludeFields: strings.Join(gp.ExcludeFields, ","),
	})
}

// New adds a new list merge field
func New(listID string, params *NewParams) (*MergeField, error) {
	res := &MergeField{}
	path := fmt.Sprintf("lists/%s/merge-fields", listID)

	if params == nil {
		if err := mailchimp.Call("POST", path, nil, nil, res); err != nil {
			return nil, err
		}
	}

	if err := mailchimp.Call("POST", path, nil, params, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Get retrieves a list of merge fields
func Get(listID string, params *GetParams) (*ListMergeFields, error) {
	res := &ListMergeFields{}
	path := fmt.Sprintf("lists/%s/merge-fields", listID)

	if params == nil {
		if err := mailchimp.Call("GET", path, nil, nil, res); err != nil {
			return nil, err
		}
		return res, nil
	}

	if err := mailchimp.Call("GET", path, params, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetMergeField retrieves a single merge field
func GetMergeField(listID string, mergeID int, params *GetParams) (*MergeField, error) {
	res := &MergeField{}
	path := fmt.Sprintf("lists/%s/merge-fields/%s", listID, strconv.Itoa(mergeID))

	if params == nil {
		if err := mailchimp.Call("GET", path, nil, nil, res); err != nil {
			return nil, err
		}
		return res, nil
	}

	if err := mailchimp.Call("GET", path, params, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Update updates a merge field in a list
func Update(listID string, mergeID int, params *UpdateParams) (*MergeField, error) {
	res := &MergeField{}
	path := fmt.Sprintf("lists/%s/merge-fields/%s", listID, strconv.Itoa(mergeID))

	if params == nil {
		if err := mailchimp.Call("PATCH", path, nil, nil, res); err != nil {
			return nil, err
		}
	}

	if err := mailchimp.Call("PATCH", path, nil, params, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Delete deletes a list merge field.
func Delete(listID string, mergeID int) error {
	path := fmt.Sprintf("lists/%s/merge-fields/%s", listID, strconv.Itoa(mergeID))
	return mailchimp.Call("DELETE", path, nil, nil, nil)
}
