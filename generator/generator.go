package generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UniversityRadioYork/alias-go/utils"
	"github.com/UniversityRadioYork/myradio-go"
	"log"
	"sort"
	"strings"
	"time"
)

// Aliases holds aliases.
// The key is the source and the value is an array of destinations.
type Aliases map[string][]string

// GenerateAliases creates the aliases string using a config.
// It returns errors at the earliest opportunity.
func GenerateAliases(ury utils.URYFetcher, c utils.Configurer) (string, error) {
	err := checkConfig(c)
	if err != nil {
		return "", err
	}
	// Mailing List Aliases
	mailing, err := generateMailingListAliases(ury)
	if err != nil {
		return "", err
	}
	// Misc Aliases
	misc, err := generateMiscAliases(ury)
	if err != nil {
		return "", err
	}
	// Officer Aliases
	officer, err := generateOfficerAliases(ury, c)
	if err != nil {
		return "", err
	}
	// User Aliases
	user, err := generateUserAliases(ury)
	if err != nil {
		return "", err
	}
	aliases := mergeAliases(mailing, misc, officer, user)
	addManagementFallback(&aliases, c)
	addNonDottedAliases(&aliases)
	removeDuplicatesAndBlanks(&aliases)
	return aliasesToString(aliases), nil
}

func generateMailingListAliases(ury utils.URYFetcher) (Aliases, error) {
	lists, err := ury.GetMailingLists()
	if err != nil {
		return nil, err
	}
	var aliases = make(Aliases)
	for _, list := range lists {
		if len(list.Address) == 0 {
			log.Printf("Skipping list '%s' with id: %d, it has no address", list.Name, list.Listid)
			continue
		}
		members, err := ury.GetMailingListMembers(list)
		if err != nil {
			return nil, err
		}
		if len(members) > 0 {
			if _, exists := aliases[list.Address]; !exists {
				aliases[list.Address] = make([]string, 0, list.Recipients)
			}
			for _, member := range members {
				if member.Receiveemail {
					if member.Email == "" {
						log.Printf("Member with id: %d has receive_email set to true but has "+
							"no email set", member.MemberID)
					} else {
						aliases[list.Address] = append(aliases[list.Address], member.Email)
					}
				}
			}
		} else {
			log.Printf("Skipping list '%s' with id: %d, it has no members", list.Name, list.Listid)
		}
	}
	return aliases, nil
}

func generateMiscAliases(ury utils.URYFetcher) (Aliases, error) {
	raws, err := ury.GetMiscAliases()
	if err != nil {
		return nil, err
	}
	var aliases = make(Aliases)
	for _, raw := range raws {
		if len(raw.Source) == 0 {
			log.Printf("Skipping due to blank source for misc with id: %d", raw.Id)
			continue
		}
		if _, exists := aliases[raw.Source]; !exists {
			aliases[raw.Source] = make([]string, 0)
		}
		for _, dest := range raw.Destinations {
			var deststr string
			var err error
			switch dest.Atype {
			case "member":
				deststr, err = parseMemberAlias(dest.Value)
			case "text":
				deststr, err = parseTextAlias(dest.Value)
			case "officer":
				deststr, err = parseOfficerAlias(dest.Value)
			case "list":
				deststr, err = parseListAlias(dest.Value)
			default:
				err = errors.New(fmt.Sprintf("Invalid value for switch, '%s'", dest.Atype))
			}
			if err != nil {
				return nil, err
			}
			if deststr != "" {
				aliases[raw.Source] = append(aliases[raw.Source], deststr)
			}
		}
	}
	return aliases, nil
}

func generateOfficerAliases(ury utils.URYFetcher, c utils.Configurer) (Aliases, error) {
	officers, err := ury.GetOfficerAliases()
	if err != nil {
		return nil, err
	}
	var aliases = make(Aliases)
	for _, officer := range officers {
		if len(officer.Alias) == 0 {
			log.Printf("Skipping officer '%s' with id: %d as it has no alias", officer.Name, officer.OfficerID)
			continue
		}
		if _, exists := aliases[officer.Alias]; !exists {
			aliases[officer.Alias] = make([]string, 0)
		}
		err = addCurrentOfficers(&aliases, officer, ury, c)
		if err != nil {
			return nil, err
		}
		err = addHistoricalOfficers(&aliases, officer, c)
		if err != nil {
			return nil, err
		}
	}
	return aliases, nil
}

func generateUserAliases(ury utils.URYFetcher) (Aliases, error) {
	var userAliases, err = ury.GetMemberAliases()
	var aliases = make(Aliases)
	if err != nil {
		return nil, err
	}
	for _, v := range userAliases {
		if v.Source == "" || v.Destination == "" {
			log.Printf("Blank source or destination for member '%s' => '%s'", v.Source, v.Destination)
			continue
		}
		if _, exists := aliases[v.Source]; exists {
			aliases[v.Source] = append(aliases[v.Source], v.Destination)
		} else {
			aliases[v.Source] = []string{v.Destination}
		}
	}
	return aliases, nil
}

// addNonDottedAliases adds an alias for emails that have '.' in them
// eg 'head.of.computing' would need an alias from 'headofcomputing'
func addNonDottedAliases(a *Aliases) {
	n := make(Aliases)
	for s := range *a {
		if strings.Contains(s, ".") {
			nd := strings.Replace(s, ".", "", -1)
			if _, exists := n[nd]; exists {
				n[nd] = append(n[nd], s)
			} else {
				n[nd] = []string{s}
			}
		}
	}
	(*a) = mergeAliases(*a, n)
	return
}

// mergeAliases takes an amount of aliases and combines them.
// It does *not* check for duplicates as this would take longer.
func mergeAliases(args ...Aliases) Aliases {
	merged := make(Aliases)
	for _, a := range args {
		for s, ds := range a {
			if _, exists := merged[s]; exists {
				merged[s] = append(merged[s], ds...)
			} else {
				merged[s] = ds
			}
		}
	}
	return merged
}

func aliasesToString(a Aliases) string {
	// Because it's nice to generate the aliases
	// in alphabetical order, and to make the tests
	// pass 100% of the time, we make an array of the
	// keys in the map, and use those instead
	keys := make([]string, 0, len(a))
	for key := range a {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	str := ""
	for _, key := range keys {
		if len(a[key]) > 0 {
			sort.Strings(a[key])
			str += key + ": " + strings.Join(a[key], ", ") + ", \n"
		} else {
			log.Printf("Skipping writing source '%s', as it has no destinations", key)
		}
	}
	return str
}

func removeDuplicatesAndBlanks(a *Aliases) {
	for s, ds := range *a {
		found := make(map[string]bool)
		n := make([]string, 0, len(ds))
		for _, d := range ds {
			if d != "" && !found[d] {
				found[d] = true
				n = append(n, d)
			}
		}
		(*a)[s] = n
	}
}

func parseTextAlias(raw *json.RawMessage) (string, error) {
	var result string
	err := json.Unmarshal(*raw, &result)
	if err != nil {
		return "", err
	}
	return result, nil
}

func parseOfficerAlias(raw *json.RawMessage) (string, error) {
	var result myradio.OfficerPosition
	err := json.Unmarshal(*raw, &result)
	if err != nil {
		return "", err
	}
	return result.Alias, nil
}

func parseMemberAlias(raw *json.RawMessage) (string, error) {
	var result myradio.User
	err := json.Unmarshal(*raw, &result)
	if err != nil {
		return "", err
	}
	if result.Receiveemail {
		return result.Email, nil
	} else {
		log.Printf("memberid %d has receive_email set to false", result.MemberID)
		return "", nil
	}
}

func parseListAlias(raw *json.RawMessage) (string, error) {
	var result myradio.List
	err := json.Unmarshal(*raw, &result)
	if err != nil {
		return "", err
	}
	return result.Address, nil
}

func addCurrentOfficers(a *Aliases, o myradio.OfficerPosition, ury utils.URYFetcher, c utils.Configurer) error {
	if len(o.Current) > 0 {
		for _, officer := range o.Current {
			if officer.Receiveemail {
				if officer.Email == "" {
					log.Printf("Member with id: %d has receive_email set to true but has "+
						"no email set", officer.MemberID)
				} else {
					(*a)[o.Alias] = append((*a)[o.Alias], officer.Email)
				}
			}
		}
		return nil
	} else {
		log.Printf("No current officer '%s' in team: '%d '%s', deferring to head of team",
			o.Name, o.Team.TeamID, o.Team.Name)
		return addHeadOfTeam(a, o, ury, c)
	}
}

func addHistoricalOfficers(a *Aliases, o myradio.OfficerPosition, c utils.Configurer) error {

	for _, officer := range o.History {
		v, err := c.IsHistoricalOfficerValid(time.Now(), officer.To)
		if err != nil {
			return err
		}
		if v {
			if officer.User.Receiveemail {
				if officer.User.Email == "" {
					log.Printf("Member with id: %d has receive_email set to true but has "+
						"no email set", officer.User.MemberID)
				} else {
					(*a)[o.Alias] = append((*a)[o.Alias], officer.User.Email)
				}
			}
		}
	}
	return nil
}

func addHeadOfTeam(a *Aliases, o myradio.OfficerPosition, ury utils.URYFetcher, c utils.Configurer) error {
	heads, err := ury.GetHeadOfTeam(o.Team)
	if err != nil {
		return err
	}
	if len(heads) > 0 {
		for _, head := range heads {
			if head.User.Receiveemail {
				if head.User.Email == "" {
					log.Printf("Member with id: %d has receive_email set to true but has "+
						"no email set", head.User.MemberID)
				} else {
					(*a)[o.Alias] = append((*a)[o.Alias], head.User.Email)
				}
			}
		}
	} else {
		log.Printf("Deferring head of team '%s' with id: %d to (assistant) station manager", o.Team.Name, o.OfficerID)
		if o.Alias == c.GetHeadOfStation() {
			(*a)[o.Alias] = append((*a)[o.Alias], c.GetAssistantHeadOfStation())
		} else {
			(*a)[o.Alias] = append((*a)[o.Alias], c.GetHeadOfStation())
		}
	}
	return nil
}

func addManagementFallback(a *Aliases, c utils.Configurer) {
	// Fall back to ASM if there is no SM
	if h, exists := (*a)[c.GetHeadOfStation()]; !exists || len(h) == 0 {
		if !exists {
			(*a)[c.GetHeadOfStation()] = make([]string, 0, 1)
		}
		(*a)[c.GetHeadOfStation()] = append((*a)[c.GetHeadOfStation()], c.GetAssistantHeadOfStation())
	}
}

func checkConfig(c utils.Configurer) error {
	if c.GetHeadOfStation() == "" {
		return errors.New("No SM set in config")
	}
	if c.GetAssistantHeadOfStation() == "" {
		return errors.New("No ASM set in config")
	}
	return nil
}
