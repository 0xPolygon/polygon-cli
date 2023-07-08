package dashboard

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type (
	DashboardOptions struct {
		File                string
		Prefix              string
		Title               string
		Description         string
		WidgetWidth         int
		WidgetHeight        int
		TemplateVars        []string
		TemplateVarDefaults []string
		StripPrefixes       []string
		Pretty              bool
		ShowHelp            bool
	}

	DataDogTemplateVariables struct {
		Name            string   `json:"name,omitempty"`
		Prefix          string   `json:"prefix,omitempty"`
		AvailableValues []string `json:"available_values,omitempty"`
		Default         string   `json:"default,omitempty"`
	}

	DataDogFormula struct {
		Alias              string        `json:"alias,omitempty"`
		ConditionalFormats []interface{} `json:"conditional_formats,omitempty"`
		Limit              struct {
			Count int    `json:"count,omitempty"`
			Order string `json:"order,omitempty"`
		} `json:"limit,omitempty"`
		CellDisplayMode string `json:"cell_display_mode,omitempty"`
		Formula         string `json:"formula,omitempty"`
	}
	DataDogQuery struct {
		Query      string `json:"query,omitempty"`
		DataSource string `json:"data_source,omitempty"`
		Name       string `json:"name,omitempty"`
		Aggregator string `json:"aggregator,omitempty"`
	}
	DataDogRequest struct {
		Formulas       []DataDogFormula `json:"formulas,omitempty"`
		ResponseFormat string           `json:"response_format,omitempty"`
		Queries        []DataDogQuery   `json:"queries,omitempty"`
		Style          struct {
			Palette   string `json:"palette"`
			LineType  string `json:"line_type"`
			LineWidth string `json:"line_width"`
		} `json:"style"`
	}
	DataDogWidget struct {
		ID         int64 `json:"id"`
		Definition struct {
			helpText        string
			Title           string           `json:"title,omitempty"`
			TitleSize       string           `json:"title_size,omitempty"`
			TitleAlign      string           `json:"title_align,omitempty"`
			Type            string           `json:"type,omitempty"`
			Requests        []DataDogRequest `json:"requests,omitempty"`
			HasSearchBar    string           `json:"has_search_bar,omitempty"`
			ShowLegend      bool             `json:"show_legend,omitempty"`
			LegendLayout    string           `json:"legend_layout,omitempty"`
			LegendColumns   []string         `json:"legend_columns,omitempty"`
			Text            string           `json:"text,omitempty"`
			FontSize        string           `json:"font_size,omitempty"`
			Content         string           `json:"content,omitempty"`
			Widgets         DataDogWidgets   `json:"widgets,omitempty"`
			BackgroundColor string           `json:"background_color,omitempty"`
			LayoutType      string           `json:"layout_type,omitempty"`
		} `json:"definition"`
		Layout *DataDogLayout `json:"layout,omitempty"`
	}
	DataDogWidgets []DataDogWidget
	DataDogLayout  struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width,omitempty"`
		Height int `json:"height,omitempty"`
	}

	DataDogDashboard struct {
		Title             string                     `json:"title,omitempty"`
		Description       string                     `json:"description,omitempty"`
		Widgets           []DataDogWidget            `json:"widgets,omitempty"`
		TemplateVariables []DataDogTemplateVariables `json:"template_variables,omitempty"`
		LayoutType        string                     `json:"layout_type,omitempty"`
		IsReadOnly        bool                       `json:"is_read_only,omitempty"`
		NotifyList        []interface{}              `json:"notify_list,omitempty"`
		ReflowType        string                     `json:"reflow_type,omitempty"`
		ID                string                     `json:"id,omitempty"`
	}
)

func (a DataDogWidgets) Len() int           { return len(a) }
func (a DataDogWidgets) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DataDogWidgets) Less(i, j int) bool { return a[i].Definition.Title < a[j].Definition.Title }

func ConvertMetricsToDashboard(input *DashboardOptions) ([]byte, error) {
	if len(input.TemplateVars) != len(input.TemplateVarDefaults) && len(input.TemplateVarDefaults) > 0 {
		return nil, fmt.Errorf("the length of the template vars and template var defaults arguments do not match")
	}
	metrics, err := ParseMetricsFile(input.File)
	if err != nil {
		return nil, err
	}
	dd, err := MetricsToDataDog(input, metrics)
	if err != nil {
		return nil, err
	}
	jsBytes, err := json.Marshal(dd)
	if err != nil {
		return nil, err
	}
	return jsBytes, nil
}

func ParseMetricsFile(filePath string) (map[string]*dto.MetricFamily, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(file)
	if err != nil {
		return nil, err
	}
	return mf, nil
}

func MetricsToDataDog(dopts *DashboardOptions, metrics map[string]*dto.MetricFamily) (*DataDogDashboard, error) {
	dash := new(DataDogDashboard)
	dash.Title = dopts.Title
	dash.Description = dopts.Description
	widgets := make([]DataDogWidget, 0)
	index := 0

	for _, v := range metrics {
		var w *DataDogWidget
		switch *v.Type {

		case dto.MetricType_COUNTER:
			// Monotonically increasing counter. By default we probably want to show the derivative
			w = NewDataDogCounterWidget(dopts, v)
		case dto.MetricType_GAUGE:
			// Numerical value that can go up and down
			w = NewDataDogGaugeWidget(dopts, v)
		case dto.MetricType_SUMMARY:
			// Samples of observations
			continue
		case dto.MetricType_UNTYPED:
			continue
		case dto.MetricType_HISTOGRAM:
			w = NewDataDogHistogramWidget(dopts, v)
		default:
			continue
		}

		widgets = append(widgets, *w)
		index += 1

	}

	sort.Sort(DataDogWidgets(widgets))

	currentX := 0
	currentY := 0

	helpWidgets := make([]DataDogWidget, 0)
	for k := range widgets {
		currentX = (k % (12 / dopts.WidgetWidth)) * dopts.WidgetWidth
		if (k % (12 / dopts.WidgetWidth)) == 0 {
			currentY += dopts.WidgetHeight
		}

		adjustedHeight := dopts.WidgetHeight
		adjustedY := currentY
		if dopts.ShowHelp {
			hw := NewDataDogNoteWidget(dopts, widgets[k].Definition.helpText)
			hw.Layout = &DataDogLayout{
				X:      currentX,
				Y:      currentY,
				Width:  dopts.WidgetWidth,
				Height: 1,
			}
			helpWidgets = append(helpWidgets, *hw)

			adjustedY = currentY + 1
			adjustedHeight = dopts.WidgetHeight - 1
		}

		l := DataDogLayout{
			// X:      (k % (12 / dopts.WidgetWidth)) * dopts.WidgetWidth,
			// Y:      int(math.Floor(float64(k/(12/dopts.WidgetWidth)))) * dopts.WidgetHeight,
			X:      currentX,
			Y:      adjustedY,
			Width:  dopts.WidgetWidth,
			Height: adjustedHeight,
		}
		widgets[k].Layout = &l
	}

	if dopts.ShowHelp {
		finalWidgets := make([]DataDogWidget, 0)
		itemsPerRow := 12 / dopts.WidgetWidth
		for i := 0; i < len(widgets); i = i + itemsPerRow {
			for j := i; j-i < itemsPerRow && j < len(helpWidgets); j = j + 1 {
				finalWidgets = append(finalWidgets, helpWidgets[j])
			}
			for j := i; j-i < itemsPerRow && j < len(widgets); j = j + 1 {
				finalWidgets = append(finalWidgets, widgets[j])
			}
		}
		dash.Widgets = finalWidgets
	} else {
		dash.Widgets = widgets
	}

	g := NewDataDogGroupWidget(dopts, "Auto Widgets", "vivid_purple")
	g.Definition.Widgets = dash.Widgets
	dash.Widgets = []DataDogWidget{*g}

	// dash.LayoutType = "free"
	dash.LayoutType = "ordered"
	dash.ReflowType = "fixed"
	// dash.ReflowType = "auto"
	dash.IsReadOnly = true
	tv := make([]DataDogTemplateVariables, 0)
	for k, v := range dopts.TemplateVars {
		dv := DataDogTemplateVariables{
			Name:            v,
			Prefix:          v,
			AvailableValues: []string{},
			Default:         dopts.TemplateVarDefaults[k],
		}
		tv = append(tv, dv)
	}

	dash.TemplateVariables = tv
	return dash, nil
}

func NewDataDogCounterWidget(dopts *DashboardOptions, mf *dto.MetricFamily) *DataDogWidget {
	w := newDataDogWidget(dopts, mf)

	w.Definition.Requests[0].Queries[0].Query = fmt.Sprintf("sum:%s%s.count{$basedn,$host}.as_count()", dopts.Prefix, *mf.Name)
	w.Definition.Requests[0].Queries[0].DataSource = "metrics"
	w.Definition.Requests[0].Queries[0].Name = "autoquery"
	return w
}

func NewDataDogGaugeWidget(dopts *DashboardOptions, mf *dto.MetricFamily) *DataDogWidget {
	w := newDataDogWidget(dopts, mf)

	w.Definition.Requests[0].Queries[0].Query = fmt.Sprintf("avg:%s%s{$basedn,$host}", dopts.Prefix, *mf.Name)
	w.Definition.Requests[0].Queries[0].DataSource = "metrics"
	w.Definition.Requests[0].Queries[0].Name = "autoquery"
	return w
}

func NewDataDogHistogramWidget(dopts *DashboardOptions, mf *dto.MetricFamily) *DataDogWidget {
	w := newDataDogWidget(dopts, mf)

	// Datadog ignores buckets by default, so we'll just use the sum field
	w.Definition.Requests[0].Queries[0].Query = fmt.Sprintf("avg:%s%s.sum{$basedn,$host}.as_count()", dopts.Prefix, *mf.Name)
	w.Definition.Requests[0].Queries[0].DataSource = "metrics"
	w.Definition.Requests[0].Queries[0].Name = "autoquery"
	return w
}

// newDataDogWidget will initialize a basic object with arrays with one item in them
func newDataDogWidget(dopts *DashboardOptions, mf *dto.MetricFamily) *DataDogWidget {
	w := new(DataDogWidget)
	name := *mf.Name
	for _, strip := range dopts.StripPrefixes {
		name = strings.TrimPrefix(name, strip)
	}
	if dopts.Pretty {
		name = strings.ReplaceAll(name, "_", " ")
		name = cases.Title(language.English, cases.Compact).String(name)
	}
	w.Definition.Title = name
	w.Definition.helpText = *mf.Help
	w.Definition.Type = "timeseries"
	w.Definition.TitleSize = "16"
	w.Definition.TitleAlign = "left"

	f := DataDogFormula{}
	f.Formula = "autoquery"

	q := DataDogQuery{}

	r := DataDogRequest{}
	r.Formulas = []DataDogFormula{f}
	r.Queries = []DataDogQuery{q}
	r.ResponseFormat = "timeseries"
	r.Style.Palette = "dog_classic"
	r.Style.LineType = "solid"
	r.Style.LineWidth = "normal"

	w.Definition.Requests = []DataDogRequest{r}

	return w
}

func NewDataDogTextWidget(dopts *DashboardOptions, text string) *DataDogWidget {
	w := new(DataDogWidget)

	w.Definition.Type = "free_text"
	w.Definition.Text = text
	w.Definition.FontSize = "auto"
	return w
}
func NewDataDogNoteWidget(dopts *DashboardOptions, text string) *DataDogWidget {
	w := new(DataDogWidget)

	w.Definition.Type = "note"
	w.Definition.Content = text
	w.Definition.FontSize = "auto"
	return w
}

func NewDataDogGroupWidget(dopts *DashboardOptions, title, color string) *DataDogWidget {
	w := new(DataDogWidget)

	w.Definition.Type = "group"
	w.Definition.Title = title
	w.Definition.BackgroundColor = color
	w.Definition.LayoutType = "ordered"
	w.Definition.Widgets = make([]DataDogWidget, 0)

	return w
}
