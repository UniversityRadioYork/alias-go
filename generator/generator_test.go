package generator

import (
	"encoding/json"
	"errors"
	"github.com/UniversityRadioYork/alias-go/utils"
	"github.com/UniversityRadioYork/myradio-go"
	"reflect"
	"sort"
	"testing"
	"time"
)

type uryTest struct {
	utils.URYFetcher
}

func (ury uryTest) GetMemberAliases() ([]myradio.UserAlias, error) {
	return []myradio.UserAlias{
		{
			Source:      "chris.taylor",
			Destination: "christaylor",
		},
		{
			Source:      "",
			Destination: "asdasd",
		},
		{
			Source:      "ghkhkfk",
			Destination: "",
		},
		{
			Source:      "qwer",
			Destination: "foop",
		},
		{
			Source:      "qwer",
			Destination: "baz",
		},
	}, nil
}

func (ury uryTest) GetMailingLists() ([]myradio.List, error) {

	return []myradio.List{
		{
			Recipients: 3,
			Address:    "test.list1",
			Listid:     1,
			Name:       "Test List 1",
		},
		{
			Recipients: 2,
			Address:    "test.list2",
			Listid:     2,
			Name:       "Test List 2",
		},
		{
			Recipients: 0,
			Address:    "Empty test.list3",
			Listid:     3,
			Name:       "Empty Test List 3",
		},
	}, nil

}

func (ury uryTest) GetMailingListMembers(list myradio.List) ([]myradio.Member, error) {

	switch list.Listid {
	case 1:
		return []myradio.Member{
			{
				Email:        "test.member1",
				Receiveemail: true,
				Memberid:     1,
			},
			{
				Email:        "test.member2",
				Receiveemail: false,
				Memberid:     2,
			},
			{
				Email:        "test.member3",
				Receiveemail: true,
				Memberid:     3,
			},
		}, nil
	case 2:
		return []myradio.Member{
			{
				Email:        "test.member4",
				Receiveemail: true,
				Memberid:     4,
			},
			{
				Email:        "",
				Receiveemail: true,
				Memberid:     5,
			},
		}, nil
	case 3:
		return []myradio.Member{}, nil
	default:
		return nil, errors.New("Bad switch in uryTest.GetMailingListMembers()")
	}

}

func (ury uryTest) GetMiscAliases() ([]myradio.Alias, error) {

	var temp1 = json.RawMessage(`"test1.dest"`)
	var temp2 = json.RawMessage(`{
    "officerid": 26,
    "name": "Test Officer Role",
    "alias": "test.officer.role",
    "team": {
      "teamid": 7,
      "name": "Test Team",
      "alias": "test.team",
      "ordering": 40,
      "description": "",
      "status": "c"
    },
    "ordering": 15,
    "description": "",
    "status": "h",
    "type": "a"
    }`)
	var temp3 = json.RawMessage(`{
    "memberid": 123345,
    "fname": "Test",
    "sname": "McTesterson",
    "sex": "m",
    "public_email": "test.mctesterson@ury.org.uk",
    "url": "//ury.org.uk/myradio/Profile/view/?memberid=123345",
    "photo": "/static/img/default_show_player.png",
    "bio": null,
    "receive_email": true
  }`)
	var temp4 = json.RawMessage(`{
    "memberid": 12341234,
    "fname": "Test",
    "sname": "McTesterson",
    "sex": "m",
    "public_email": "test.mctesterson@ury.org.uk",
    "url": "//ury.org.uk/myradio/Profile/view/?memberid=12341234",
    "photo": "/static/img/default_show_player.png",
    "bio": null,
    "receive_email": false
  }`)
	var temp5 = json.RawMessage(`{
    "listid": 123,
    "subscribed": false,
    "name": "Computing Team",
    "address": "computing",
    "recipient_count": 16
  }`)

	return []myradio.Alias{
		{
			Id:     1,
			Source: "test1.source",
			Destinations: []struct {
				Atype string `json:"type"`
				Value *json.RawMessage
			}{
				{
					Atype: "text",
					Value: &temp1,
				},
			},
		},
		{
			Id:     1,
			Source: "test2.source",
			Destinations: []struct {
				Atype string `json:"type"`
				Value *json.RawMessage
			}{
				{
					Atype: "officer",
					Value: &temp2,
				},
			},
		},
		{
			Id:     1,
			Source: "test3.source",
			Destinations: []struct {
				Atype string `json:"type"`
				Value *json.RawMessage
			}{
				{
					Atype: "member",
					Value: &temp3,
				},
				{
					Atype: "member",
					Value: &temp4,
				},
			},
		},
		{
			Id:     1,
			Source: "test4.source",
			Destinations: []struct {
				Atype string `json:"type"`
				Value *json.RawMessage
			}{
				{
					Atype: "list",
					Value: &temp5,
				},
			},
		},
	}, nil

}

func (ury uryTest) GetOfficerAliases() ([]myradio.OfficerPosition, error) {
	return []myradio.OfficerPosition{
		{
			OfficerID: 1,
			Alias:     "",
			Name:      "This shouldn't appear",
			Team: myradio.Team{
				TeamID: 1,
			},
		},
		{
			OfficerID: 2,
			Alias:     "boop",
			Name:      "No current officer, should add head of team",
			Team: myradio.Team{
				TeamID: 2,
			},
			Current: []myradio.Member{},
			History: []struct {
				User            myradio.Member
				From            time.Time
				FromRaw         int64 `json:"from"`
				To              time.Time
				ToRaw           int64 `json:"to"`
				MemberOfficerID int
			}{},
		},
		{
			OfficerID: 3,
			Alias:     "foop",
			Name:      "No history",
			Team: myradio.Team{
				TeamID: 3,
			},
			Current: []myradio.Member{
				{
					Receiveemail: true,
					Email:        "boop",
				},
				{
					Receiveemail: true,
					Email:        "baz",
				},
				{
					Receiveemail: false,
					Email:        "foo",
				},
			},
			History: []struct {
				User            myradio.Member
				From            time.Time
				FromRaw         int64 `json:"from"`
				To              time.Time
				ToRaw           int64 `json:"to"`
				MemberOfficerID int
			}{},
		},
		{
			OfficerID: 3,
			Alias:     "asda",
			Name:      "Has history",
			Team: myradio.Team{
				TeamID: 4,
			},
			Current: []myradio.Member{
				{
					Receiveemail: true,
					Email:        "boop",
				},
			},
			History: []struct {
				User            myradio.Member
				From            time.Time
				FromRaw         int64 `json:"from"`
				To              time.Time
				ToRaw           int64 `json:"to"`
				MemberOfficerID int
			}{
				{
					User: myradio.Member{
						Email:        "123123",
						Receiveemail: true,
					},
				},
				{
					User: myradio.Member{
						Email:        "456678",
						Receiveemail: false,
					},
				},
			},
		},
	}, nil
}

func (ury uryTest) GetHeadOfTeam(t myradio.Team) ([]myradio.HeadPosition, error) {
	switch t.TeamID {
	case 1:
		return []myradio.HeadPosition{
			{
				User: myradio.Member{
					Email:        "foo@baz",
					Receiveemail: true,
					Memberid:     123,
				},
			},
		}, nil
	case 2:
		return []myradio.HeadPosition{
			{
				User: myradio.Member{
					Email:        "qwexgd@baz",
					Receiveemail: true,
					Memberid:     456,
				},
			},
		}, nil
	case 3:
		return []myradio.HeadPosition{
			{
				User: myradio.Member{
					Email:        "asdqweqwe@baz",
					Receiveemail: true,
					Memberid:     674,
				},
			},
		}, nil
	case 4:
		return []myradio.HeadPosition{
			{
				User: myradio.Member{
					Email:        "asdasd@baz",
					Receiveemail: false,
					Memberid:     234,
				},
			},
		}, nil
	default:
		return nil, errors.New("Base case statement")
	}
}

type configTest struct {
	utils.Configurer
	Valid bool
	SM    string
	ASM   string
	API   string
}

func (tc configTest) IsHistoricalOfficerValid(now, to time.Time) (bool, error) {
	return tc.Valid, nil
}

func (tc configTest) GetHeadOfStation() string {
	return tc.SM
}

func (tc configTest) GetAssistantHeadOfStation() string {
	return tc.ASM
}

func (tc configTest) GetApiKey() string {
	return tc.API
}

func TestGenerator_generateMailingListAliases(t *testing.T) {

	var ury uryTest

	expected := Aliases{
		"test.list1": {
			"test.member1",
			"test.member3",
		},
		"test.list2": {
			"test.member4",
		},
	}

	actual, err := generateMailingListAliases(ury)

	if err != nil {
		t.Error(err)
	}

	assertAliases(actual, expected, t)

}

func TestGenerator_generateMiscAliases(t *testing.T) {

	var ury uryTest

	expected := Aliases{
		"test1.source": {
			"test1.dest",
		},
		"test2.source": {
			"test.officer.role",
		},
		"test3.source": {
			"test.mctesterson@ury.org.uk",
		},
		"test4.source": {
			"computing",
		},
	}

	actual, err := generateMiscAliases(ury)

	if err != nil {
		t.Error(err)
	}

	assertAliases(actual, expected, t)
}

func TestGenerator_generateOfficerAliases1(t *testing.T) {

	var ury uryTest
	var config = configTest{
		Valid: true,
	}

	actual, err := generateOfficerAliases(ury, config)

	expected := Aliases{
		"boop": {
			"qwexgd@baz",
		},
		"foop": {
			"boop",
			"baz",
		},
		"asda": {
			"boop",
			"123123",
		},
	}

	if err != nil {
		t.Error(err)
	}

	assertAliases(actual, expected, t)

}

func TestGenerator_generateOfficerAliases2(t *testing.T) {

	var ury uryTest
	var config = configTest{
		Valid: false,
	}

	actual, err := generateOfficerAliases(ury, config)

	expected := Aliases{
		"boop": {
			"qwexgd@baz",
		},
		"foop": {
			"boop",
			"baz",
		},
		"asda": {
			"boop",
		},
	}

	if err != nil {
		t.Error(err)
	}

	assertAliases(actual, expected, t)

}

func TestGenerator_generateUserAliases(t *testing.T) {

	var ury uryTest

	actual, err := generateUserAliases(ury)

	expected := Aliases{
		"chris.taylor": {
			"christaylor",
		},
		"qwer": {
			"foop",
			"baz",
		},
	}

	if err != nil {
		t.Error(err)
	}

	assertAliases(actual, expected, t)

}

func TestGenerator_removeDuplicatesAndBlanks(t *testing.T) {

	actual := Aliases{
		"root": {
			"heres1",
			"another.one",
			"duplicate",
			"duplicate",
			"woo",
			"yarrrrr",
			"testing",
			"is",
			"fun",
			"...",
			"not",
			"heres1",
		},
	}

	expected := Aliases{
		"root": {
			"heres1",
			"another.one",
			"duplicate",
			"woo",
			"yarrrrr",
			"testing",
			"is",
			"fun",
			"...",
			"not",
		},
	}

	removeDuplicatesAndBlanks(&actual)

	assertAliases(actual, expected, t)

}

func TestGenerator_aliasesToString(t *testing.T) {

	a := Aliases{
		"root": {
			"testing",
			"is",
			"fun",
		},
		"not.root": {
			"testing",
			"is",
			"exciting",
		},
		"empty": {},
	}

	expected := "not.root: exciting, is, testing, \nroot: fun, is, testing, \n"

	actual := aliasesToString(a)

	if eq := reflect.DeepEqual(expected, actual); !eq {
		t.Errorf("expected \n%s, got \n%s", expected, actual)
	}

}

func TestGenerator_addNonDottedAliases(t *testing.T) {

	actual := Aliases{
		"something.with.dots": {
			"testing",
			"is",
			"fun",
		},
		"nodots": {
			"testing",
			"is",
			"fun",
		},
		"dots.already.exists": {
			"root",
		},
		"dotsalreadyexists": {
			"notroot",
		},
	}

	expected := Aliases{
		"somethingwithdots": {
			"something.with.dots",
		},
		"something.with.dots": {
			"testing",
			"is",
			"fun",
		},
		"nodots": {
			"testing",
			"is",
			"fun",
		},
		"dots.already.exists": {
			"root",
		},
		"dotsalreadyexists": {
			"dots.already.exists",
			"notroot",
		},
	}

	addNonDottedAliases(&actual)

	assertAliases(actual, expected, t)

}

func TestGenerator_addManagementFallback1(t *testing.T) {

	config := configTest{
		SM:  "sm",
		ASM: "asm",
	}

	actual := Aliases{
		"sm": {
			"station.manager",
		},
	}

	expected := Aliases{
		"sm": {
			"station.manager",
		},
	}

	addManagementFallback(&actual, config)

	assertAliases(actual, expected, t)

}

func TestGenerator_addManagementFallback2(t *testing.T) {

	config := configTest{
		SM:  "sm",
		ASM: "asm",
	}

	actual := Aliases{}

	expected := Aliases{
		"sm": {
			"asm",
		},
	}

	addManagementFallback(&actual, config)

	assertAliases(actual, expected, t)

}

func TestGenerator_checkConfig_1(t *testing.T) {

	var tc = configTest{
		SM:  "123123",
		ASM: "123123",
	}

	err := checkConfig(tc)

	if err != nil {
		t.Errorf("Expected nil, got '%s'", err.Error())
	}

	tc.ASM = ""

	err = checkConfig(tc)

	assertErrorMessage(err, "No ASM set in config", t)

	tc.ASM = "123"
	tc.SM = ""

	err = checkConfig(tc)

	assertErrorMessage(err, "No SM set in config", t)

}

func assertErrorMessage(actual error, expected string, t *testing.T) {
	if actual == nil {
		t.Error("Value is not an error")
	}
	if eq := reflect.DeepEqual(actual.Error(), expected); !eq {
		t.Errorf("\n\nValues do not match:\n\nExpected:\n\n%s\n\nGot\n\n%s", expected, actual)
		return
	}
}

func assertAliases(actual, expected Aliases, t *testing.T) {

	if len(actual) != len(expected) {
		t.Errorf("\n\nDifferent number of keys:\n\nExpected:\n\n%s\n\nGot\n\n%s", expected, actual)
		return
	}

	for _, v := range actual {
		sort.Strings(v)
	}
	for _, v := range expected {
		sort.Strings(v)
	}

	for key, value := range expected {
		if _, exists := actual[key]; exists {
			if eq := reflect.DeepEqual(value, actual[key]); !eq {
				t.Errorf("\n\nValues for '%s' do not match:\n\nExpected:\n\n%s\n\nGot\n\n%s", key,
					expected,
					actual)
				return
			}
		} else {
			t.Errorf("\n\nKey '%s' does not exist in actual:\n\nExpected:\n\n%s\n\nGot\n\n%s", key,
				expected,
				actual)
			return
		}
	}

}
