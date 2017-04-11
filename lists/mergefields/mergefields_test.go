package mergefields

import (
	"fmt"
	"os"
	"testing"

	mailchimp "github.com/beeker1121/mailchimp-go"
	"github.com/beeker1121/mailchimp-go/lists"
)

var listID string

func createMergeField(tag, name string, fieldType FieldType) (*MergeField, error) {
	params := &NewParams{
		Tag:      tag,
		Name:     name,
		Type:     fieldType.String(),
		Required: false,
		Public:   true,
		Options: &Options{
			DefaultCountry: 164,
		},
	}

	return New(listID, params)
}

func TestNew(t *testing.T) {
	mergefield, err := createMergeField("ADDRESS", "Address", TypeAddress)
	if err != nil {
		t.Error(err)
	}

	if mergefield.Tag != "ADDRESS" {
		t.Errorf("Expected mergefield.Tag to equal \"ADDRESS\", got %s", mergefield.Tag)
	}

	if mergefield.Name != "Address" {
		t.Errorf("Expected mergefield.Name to equal \"Address\", got %s", mergefield.Tag)
	}

	if err = Delete(listID, mergefield.MergeID); err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	tag1 := "ADDRESS"
	name1 := "Address"
	mergefield1, err := createMergeField(tag1, name1, TypeAddress)
	if err != nil {
		t.Error(err)
	}

	tag2 := "RANDOM"
	name2 := "Random value"
	mergefield2, err := createMergeField(tag2, name2, TypeText)
	if err != nil {
		t.Error(err)
	}

	fields, err := Get(listID, nil)
	if err != nil {
		t.Error(err)
	}

	// Field ordering is not deterministic, check based on expected mergeID
	// each list defaults with FNAME/LNAME fields
	var errStr *string
	errStr = containsMergeID(fields.MergeFields, 3, tag1)
	if errStr != nil {
		t.Errorf(*errStr)
	}
	errStr = containsMergeID(fields.MergeFields, 4, tag2)
	if errStr != nil {
		t.Errorf(*errStr)
	}

	if err = Delete(listID, mergefield1.MergeID); err != nil {
		t.Error(err)
	}
	if err = Delete(listID, mergefield2.MergeID); err != nil {
		t.Error(err)
	}
}

func containsMergeID(fields []MergeField, ID int, tag string) *string {
	var field *MergeField
	for _, f := range fields {
		if ID == f.MergeID {
			field = &f
		}
	}
	if field == nil {
		err := fmt.Sprintf("Expected mergeID %d, with tag %s", ID, tag)
		return &err
	}
	return nil
}

func TestGetMergeField(t *testing.T) {
	tag := "ADDRESS"
	name := "Address"
	mergefield, err := createMergeField(tag, name, TypeAddress)
	if err != nil {
		t.Error(err)
	}

	gotField, err := GetMergeField(listID, mergefield.MergeID, nil)
	if err != nil {
		t.Error(err)
	}

	if gotField.Tag != tag {
		t.Errorf("Expected gotMember.Tag to equal %s, got %s", tag, gotField.Tag)
	}

	if gotField.Name != name {
		t.Errorf("Expected gotMember.Name to equal %s, got %s", name, gotField.Name)
	}

	if err = Delete(listID, mergefield.MergeID); err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	tag := "ADDRESS"
	name := "Address"
	mergefield, err := createMergeField(tag, name, TypeAddress)
	if err != nil {
		t.Error(err)
	}

	updateParams := &UpdateParams{
		Name: "A different label",
	}

	updatedField, err := Update(listID, mergefield.MergeID, updateParams)
	if err != nil {
		t.Error(err)
	}

	gotField, err := GetMergeField(listID, mergefield.MergeID, nil)
	if err != nil {
		t.Error(err)
	}

	if gotField.Name != updateParams.Name {
		t.Errorf("Expected gotMember.Name to equal %s, got %s", name, gotField.Name)
	}

	if gotField.Name != updatedField.Name {
		t.Errorf("Expected gotMember.Name to equal updatedField.Name")
	}

	if err = Delete(listID, mergefield.MergeID); err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	tag := "ADDRESS"
	name := "Address"
	mergefield, err := createMergeField(tag, name, TypeAddress)
	if err != nil {
		t.Error(err)
	}

	gotField, err := GetMergeField(listID, mergefield.MergeID, nil)
	if err != nil {
		t.Error(err)
	}

	if gotField.MergeID != mergefield.MergeID {
		t.Error("Expected gotField.MergeID to equal mergefield.MergeID")
	}

	if err = Delete(listID, mergefield.MergeID); err != nil {
		t.Error(err)
	}

	_, err = GetMergeField(listID, mergefield.MergeID, nil)

	apiErr := err.(*mailchimp.APIError)

	if apiErr.Status != 404 {
		t.Errorf("Expected err.Status to be 404, got %d", apiErr.Status)
	}
}

func createList() (*lists.List, error) {
	listParams := &lists.NewParams{
		Name: "mailchimp-go Test List",
		Contact: &lists.Contact{
			Company:  "Acme Corp",
			Address1: "123 Main St",
			City:     "Chicago",
			State:    "IL",
			Zip:      "60613",
			Country:  "United States",
		},
		PermissionReminder: "You opted to receive updates on Acme Corp",
		CampaignDefaults: &lists.CampaignDefaults{
			FromName:  "John Doe",
			FromEmail: "newsletter@acmecorp.com",
			Subject:   "Newsletter",
			Language:  "en",
		},
		EmailTypeOption: false,
		Visibility:      lists.VisibilityPublic,
	}

	return lists.New(listParams)
}

func TestMain(m *testing.M) {
	if err := mailchimp.SetKey(os.Getenv("MAILCHIMP_API_KEY")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	list, err := createList()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	listID = list.ID

	code := m.Run()

	if err = lists.Delete(list.ID); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(code)
}
