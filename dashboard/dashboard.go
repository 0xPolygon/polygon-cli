package dashboard

import (
	"encoding/json"
	"fmt"
	"os"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

type (
	DashboardOptions struct {
		File         string
		Prefix       string
		Title        string
		Description  string
		WidgetWidth  int
		WidgetHeight int
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
		// Limit              struct {
		// 	Count int    `json:"count,omitempty"`
		// 	Order string `json:"order,omitempty"`
		// } `json:"limit,omitempty"`
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
	}
	DataDogWidget struct {
		//		ID         int64 `json:"id"`
		Definition struct {
			Title         string           `json:"title,omitempty"`
			TitleSize     string           `json:"title_size,omitempty"`
			TitleAlign    string           `json:"title_align,omitempty"`
			Type          string           `json:"type,omitempty"`
			Requests      []DataDogRequest `json:"requests,omitempty"`
			HasSearchBar  string           `json:"has_search_bar,omitempty"`
			ShowLegend    bool             `json:"show_legend,omitempty"`
			LegendLayout  string           `json:"legend_layout,omitempty"`
			LegendColumns []string         `json:"legend_columns,omitempty"`
		} `json:"definition"`
		Layout DataDogLayout `json:"layout,omitempty"`
	}
	DataDogLayout struct {
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
		//		ID                string                     `json:"id,omitempty"`
	}
)

func ConvertMetricsToDashboard(input *DashboardOptions) ([]byte, error) {
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
			// Monotonically increasing counter. By default we probalby want to show the derivative
			w = NewDataDogCounterWidget(dopts, index, v)
		case dto.MetricType_GAUGE:
			// Numerical value that can go up and down
			continue
		case dto.MetricType_SUMMARY:
			// Samples of obervations
			continue
		case dto.MetricType_UNTYPED:
			continue
		case dto.MetricType_HISTOGRAM:
			continue
			// Saimples of observations
		default:
			continue
		}

		widgets = append(widgets, *w)
		index += 1

	}
	dash.Widgets = widgets
	dash.LayoutType = "ordered"
	dash.ReflowType = "fixed"
	dash.TemplateVariables = []DataDogTemplateVariables{
		DataDogTemplateVariables{
			Name:            "basedn",
			Prefix:          "basedn",
			AvailableValues: []string{},
			Default:         "",
		},
		DataDogTemplateVariables{
			Name:            "host",
			Prefix:          "host",
			AvailableValues: []string{},
			Default:         "",
		},
	}
	return dash, nil
}

func NewDataDogCounterWidget(dopts *DashboardOptions, index int, mf *dto.MetricFamily) *DataDogWidget {
	w := new(DataDogWidget)
	w.Definition.Title = *mf.Name
	w.Definition.Type = "timeseries"
	w.Definition.TitleSize = "16"
	w.Definition.TitleAlign = "left"

	w.Layout = DataDogLayout{
		X:      (index % (12 / dopts.WidgetWidth)) * dopts.WidgetWidth,
		Y:      0,
		Width:  dopts.WidgetWidth,
		Height: dopts.WidgetHeight,
	}

	f := DataDogFormula{}
	f.Formula = "autoquery"
	f.Alias = "fff"

	q := DataDogQuery{}
	q.Query = fmt.Sprintf("sum:%s%s.count{$basedn,$host}.as_count()", dopts.Prefix, *mf.Name)
	q.DataSource = "metrics"
	q.Name = "autoquery"

	r := DataDogRequest{}
	r.Formulas = []DataDogFormula{f}
	r.Queries = []DataDogQuery{q}
	r.ResponseFormat = "timeseries"

	w.Definition.Requests = []DataDogRequest{r}

	return w
}
