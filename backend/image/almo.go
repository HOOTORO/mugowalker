package image

import (
	"encoding/xml"
	"fmt"
	"mugowalker/backend/cfg"
	"mugowalker/backend/localstore"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

type Alto struct {
	XMLName        xml.Name `xml:"alto"`
	Text           string   `xml:",chardata"`
	Xmlns          string   `xml:"xmlns,attr"`
	Xlink          string   `xml:"xlink,attr"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Description    struct {
		Text                   string `xml:",chardata"`
		MeasurementUnit        string `xml:"MeasurementUnit"`
		SourceImageInformation struct {
			Text     string `xml:",chardata"`
			FileName string `xml:"fileName"`
		} `xml:"sourceImageInformation"`
		OCRProcessing struct {
			Text              string `xml:",chardata"`
			ID                string `xml:"ID,attr"`
			OcrProcessingStep struct {
				Text               string `xml:",chardata"`
				ProcessingSoftware struct {
					Text         string `xml:",chardata"`
					SoftwareName string `xml:"softwareName"`
				} `xml:"processingSoftware"`
			} `xml:"ocrProcessingStep"`
		} `xml:"OCRProcessing"`
	} `xml:"Description"`
	Layout struct {
		Text string `xml:",chardata"`
		Page struct {
			Text          string `xml:",chardata"`
			WIDTH         string `xml:"WIDTH,attr"`
			HEIGHT        string `xml:"HEIGHT,attr"`
			PHYSICALIMGNR string `xml:"PHYSICAL_IMG_NR,attr"`
			ID            string `xml:"ID,attr"`
			PrintSpace    struct {
				Text          string `xml:",chardata"`
				HPOS          string `xml:"HPOS,attr"`
				VPOS          string `xml:"VPOS,attr"`
				WIDTH         string `xml:"WIDTH,attr"`
				HEIGHT        string `xml:"HEIGHT,attr"`
				ComposedBlock struct {
					Text      string `xml:",chardata"`
					ID        string `xml:"ID,attr"`
					HPOS      string `xml:"HPOS,attr"`
					VPOS      string `xml:"VPOS,attr"`
					WIDTH     string `xml:"WIDTH,attr"`
					HEIGHT    string `xml:"HEIGHT,attr"`
					TextBlock struct {
						Text     string `xml:",chardata"`
						ID       string `xml:"ID,attr"`
						HPOS     string `xml:"HPOS,attr"`
						VPOS     string `xml:"VPOS,attr"`
						WIDTH    string `xml:"WIDTH,attr"`
						HEIGHT   string `xml:"HEIGHT,attr"`
						TextLine []struct {
							Text   string `xml:",chardata"`
							ID     string `xml:"ID,attr"`
							HPOS   string `xml:"HPOS,attr"`
							VPOS   string `xml:"VPOS,attr"`
							WIDTH  string `xml:"WIDTH,attr"`
							HEIGHT string `xml:"HEIGHT,attr"`
							String []struct {
								Text    string `xml:",chardata"`
								ID      string `xml:"ID,attr"`
								HPOS    string `xml:"HPOS,attr"`
								VPOS    string `xml:"VPOS,attr"`
								WIDTH   string `xml:"WIDTH,attr"`
								HEIGHT  string `xml:"HEIGHT,attr"`
								WC      string `xml:"WC,attr"`
								CONTENT string `xml:"CONTENT,attr"`
							} `xml:"String"`
							SP []struct {
								Text  string `xml:",chardata"`
								WIDTH string `xml:"WIDTH,attr"`
								VPOS  string `xml:"VPOS,attr"`
								HPOS  string `xml:"HPOS,attr"`
							} `xml:"SP"`
						} `xml:"TextLine"`
					} `xml:"TextBlock"`
				} `xml:"ComposedBlock"`
			} `xml:"PrintSpace"`
		} `xml:"Page"`
	} `xml:"Layout"`
}

func UnmarshalAlto(f string) Alto {
	var alt Alto
	data, e := localstore.ReadTempFile(fmt.Sprintf("%v.xml", f))
	if e != nil {
		log("tess xml parse err: " + e.Error())
	}
	xml.Unmarshal(data, &alt)
	return alt
}

func (a Alto) parse(ex []string) []*ScreenWord {
	var res []*ScreenWord
	res = make([]*ScreenWord, 0)
	tl := a.Layout.Page.PrintSpace.ComposedBlock.TextBlock.TextLine
	for i, line := range tl {
		for _, v := range line.String {
			clean := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(v.CONTENT, "")
			if len(clean) > 3 || slices.Contains(ex, clean) {
				res = append(res, SW(lowertrim(clean), cfg.ToInt(v.HPOS), cfg.ToInt(v.VPOS), i))
			}
		}
	}
	return res
}

func lowertrim(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}

var almoResultsStringer = func(arr []*ScreenWord, psm int) string {
	var s string
	s = fmt.Sprintf("__________|> PSM %v <|__________", psm)
	for i, elem := range arr {
		if i%3 == 0 {
			s += "\n"
		}
		s += elem.String()
	}
	return s
}
