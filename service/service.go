package service

import (
	"bytes"
	"encoding/xml"
	"log"
	"os"
	"strconv"
	"strings"
	"swim-tools/enum"
	"time"

	"github.com/konrad2002/dsvparser/model"
	"github.com/konrad2002/dsvparser/model/types"
	"github.com/konrad2002/dsvparser/parser"
	"github.com/konrad2002/lenexparser/model/elements"
	"github.com/konrad2002/lenexparser/model/enums"
	parser2 "github.com/konrad2002/lenexparser/parser"
)

func Convert(filename string, origin enum.FileType, target enum.FileType) error {

	dat, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(dat)
	r := parser.NewReader(buf)
	res, err := r.Read()
	if err != nil {
		panic(err)
	}
	def := res.(*model.Wettkampfergebnisliste)

	lef := elements.Lenex{
		Constructor: elements.Constructor{
			Contact: elements.Contact{
				City:     "Freising",
				Country:  "GER",
				Email:    "konrad@swimresults.de",
				Internet: "https://swimresults.de",
				Name:     "Konrad Weiß",
			},
			Name:    "SwimResults",
			Version: "0.0.1",
		},
		Meets:   &[]elements.Meet{},
		Created: parser2.DateTime{Time: time.Now()},
		Version: "3.0",
	}

	var timing enums.Timing

	switch def.Veranstaltung.Zeitmessung {
	case "AUTOMATISCH":
		timing = enums.TimingAutomatic
	case "HALBAUTOMATISCH":
		timing = enums.TimingSemiAutomatic
	case "HANDZEIT":
		timing = enums.TimingManual1
	}

	meet := elements.Meet{
		City:  def.Veranstaltung.Veranstaltungsort,
		Clubs: nil, // set later
		Contact: &elements.Contact{
			City:    def.Ausrichter.Ort,
			Country: def.Ausrichter.Land,
			Email:   def.Ausrichter.Email,
			Fax:     def.Ausrichter.Fax,
			Name:    def.Ausrichter.Name,
			Phone:   def.Ausrichter.Telefon,
			Street:  def.Ausrichter.Strasse,
			Zip:     def.Ausrichter.PLZ,
		},
		Facility: &elements.Facility{
			City: def.Veranstaltung.Veranstaltungsort,
		},
		HostClub: def.Ausrichter.NameDesAusrichters,
		Name:     def.Veranstaltung.Veranstaltungsbezeichnung,
		Sessions: nil, // set later
		Timing:   timing,
	}

	sessionMap := make(map[int]elements.Session)
	eventMap := make(map[int]elements.Event)
	clubMap := make(map[string]elements.Club)
	athleteMap := make(map[int]elements.Athlete)

	for _, dpart := range def.Abschnitte {
		sessionMap[dpart.Abschnittsnummer] = elements.Session{
			Date:            parser2.DateTime{Time: time.Date(dpart.Abschnittsdatum.Jahr, time.Month(dpart.Abschnittsdatum.Monat), dpart.Abschnittsdatum.Tag, dpart.Anfangszeit.Stunde, dpart.Anfangszeit.Minute, 0, 0, time.UTC)},
			Daytime:         parser2.DateTime{Time: time.Date(dpart.Abschnittsdatum.Jahr, time.Month(dpart.Abschnittsdatum.Monat), dpart.Abschnittsdatum.Tag, dpart.Anfangszeit.Stunde, dpart.Anfangszeit.Minute, 0, 0, time.UTC)},
			Events:          []elements.Event{}, // set later
			Number:          dpart.Abschnittsnummer,
			OfficialMeeting: parser2.DateTime{Time: time.Date(1970, time.January, 1, dpart.Kampfrichtersitzung.Stunde, dpart.Kampfrichtersitzung.Minute, 0, 0, time.UTC)},
		}
	}

	i := 1
	for _, devent := range def.Wettkaempfe {
		var gender enums.EventGender

		switch devent.Geschlecht {
		case types.MAENNLICH:
			gender = enums.EventGenderMale
		case types.WEIBLICH:
			gender = enums.EventGenderFemale
		case types.DIVERS:
			gender = enums.EventGenderMixed
		default:
			gender = enums.EventGenderAll
		}

		var stroke enums.Stroke

		switch devent.Technik {
		case 'F':
			stroke = enums.StrokeFree
		case 'R':
			stroke = enums.StrokeBack
		case 'B':
			stroke = enums.StrokeBreast
		case 'S':
			stroke = enums.StrokeFly
		case 'L':
			stroke = enums.StrokeMedley
		case 'X':
			stroke = enums.StrokeUnknown
		default:
			stroke = enums.StrokeUnknown
		}

		var tech enums.Technique
		switch devent.Ausuebung {
		case "GL":
			tech = ""
		case "BE":
			tech = enums.TechniqueKick
		case "AR":
			tech = enums.TechniquePull
		case "ST":
			tech = enums.TechniqueStart
		case "WE":
			tech = enums.TechniqueTurn
		case "GB":
			tech = enums.TechniqueGlide
		case "X":
			tech = ""
		default:
			tech = ""
		}

		e := elements.Event{
			AgeGroups:  nil, // TODO
			EventId:    devent.Wettkampfnummer,
			Gender:     gender,
			Heats:      nil,
			Number:     devent.Wettkampfnummer,
			Order:      i,
			PreEventId: devent.Qualifikationswettkampfnummer,
			SwimStyle: elements.SwimStyle{
				Distance:   devent.Einzelstrecke,
				RelayCount: devent.AnzahlStarter,
				Stroke:     stroke,
				Technique:  tech,
			},
		}

		i += 1

		eventMap[devent.Wettkampfnummer] = e

		if entry, ok := sessionMap[devent.Abschnittsnummer]; ok {
			entry.Events = append(entry.Events, e)
			sessionMap[devent.Abschnittsnummer] = entry
		}
	}

	for _, dclub := range def.Vereine {
		c := elements.Club{
			Athletes: &[]elements.Athlete{},
			Code:     strconv.Itoa(dclub.Vereinskennzahl),
			Contact:  nil,
			Name:     dclub.Vereinsbezeichnung,
			Nation:   enums.Nation(dclub.FinaNationenkuerzel),
			Region:   strconv.Itoa(dclub.Landesschwimmverband),
			Relays:   nil,
		}

		clubMap[c.Name] = c
	}

	for _, dergebnis := range def.PNErgebnisse {
		if _, ok := athleteMap[dergebnis.VeranstaltungsIdSchwimmer]; ok {
			continue
		}

		var gender enums.Gender

		switch dergebnis.Geschlecht {
		case types.MAENNLICH:
			gender = enums.GenderMale
		case types.WEIBLICH:
			gender = enums.GenderFemale
		case types.DIVERS:
			gender = enums.GenderNonBinary
		default:
			gender = enums.GenderNonBinary
		}

		name := strings.Split(dergebnis.Name, ",")

		a := elements.Athlete{
			AthleteId: dergebnis.VeranstaltungsIdSchwimmer,
			Birthdate: parser2.DateTime{Time: time.Date(dergebnis.Jahrgang, time.January, 1, 0, 0, 0, 0, time.UTC)},
			Entries:   nil, // TODO
			Firstname: name[1],
			Gender:    gender,
			Handicap:  nil,
			Lastname:  name[0],
			Level:     "", // TODO
			License:   strconv.Itoa(dergebnis.DsvId),
			Results:   nil, // TODO
		}

		athleteMap[a.AthleteId] = a

		if entry, ok := clubMap[dergebnis.Verein]; ok {
			*entry.Athletes = append(*entry.Athletes, a)
			clubMap[dergebnis.Verein] = entry
		}
	}

	var clubs []elements.Club
	for _, club := range clubMap {
		clubs = append(clubs, club)
	}
	meet.Clubs = &clubs

	var sessions []elements.Session
	for _, session := range sessionMap {
		sessions = append(sessions, session)
	}
	meet.Sessions = &sessions

	*lef.Meets = append(*lef.Meets, meet)

	lenex, err := xml.Marshal(lef)
	if err != nil {
		log.Fatal(err)
	}

	err1 := os.WriteFile("output.lef", lenex, 0644)
	if err1 != nil {
		return err1
	}

	return nil
}
