package ocr

import (
	"encoding/xml"
	"fmt"
	"strings"
	"worker/cfg"

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
	data, e := readTmp(fmt.Sprintf("%v.xml", f))
	if e != nil {
		log.Errorf("tess xml parse err: %v", e.Error())
	}
	xml.Unmarshal(data, &alt)
	return alt
}

func (a Alto) parse() []AlmoResult {
	var res []AlmoResult
	res = make([]AlmoResult, 0)
	tl := a.Layout.Page.PrintSpace.ComposedBlock.TextBlock.TextLine
	for i, line := range tl {
		for _, v := range line.String {
			if len(v.CONTENT) > 3 || slices.Contains(user.Exceptions, v.CONTENT) || strings.ContainsAny(v.CONTENT, "0123456789") {
				res = append(res, AlmoResult{Linechars: lowertrim(v.CONTENT), X: cfg.ToInt(v.HPOS), Y: cfg.ToInt(v.VPOS), LineNo: i})
			}
		}
	}
	return res
}

func lowertrim(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}
