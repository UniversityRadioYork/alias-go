package utils

import "github.com/UniversityRadioYork/myradio-go"

type URYFetcher interface {
	GetMailingLists() ([]myradio.List, error)
	GetMailingListMembers(list myradio.List) ([]myradio.User, error)
	GetMiscAliases() ([]myradio.Alias, error)
	GetOfficerAliases() ([]myradio.OfficerPosition, error)
	GetMemberAliases() ([]myradio.UserAlias, error)
	GetHeadOfTeam(myradio.Team) ([]myradio.Officer, error)
}

type URY struct {
	URYFetcher
	session myradio.Session
}

func NewURY(apikey string) (u URY, err error) {
	s, err := myradio.NewSession(apikey)
	if err != nil {
		return
	}
	u = URY{session: *s}
	return
}

func (u URY) GetMailingLists() ([]myradio.List, error) {
	return u.session.GetAllLists()
}

func (u URY) GetMailingListMembers(list myradio.List) ([]myradio.User, error) {
	return u.session.GetUsers(&list)
}

func (u URY) GetMiscAliases() ([]myradio.Alias, error) {
	return u.session.GetAllAliases([]string{})
}

func (u URY) GetOfficerAliases() ([]myradio.OfficerPosition, error) {
	return u.session.GetAllOfficerPositions([]string{"history", "current"})
}

func (u URY) GetMemberAliases() ([]myradio.UserAlias, error) {
	return u.session.GetUserAliases()
}

func (u URY) GetHeadOfTeam(t myradio.Team) ([]myradio.Officer, error) {
	return u.session.GetTeamHeadPositions(int(t.TeamID), []string{})
}
