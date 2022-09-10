package footprint

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"pairs/pkg/selenium"
	"strconv"
	"time"
)

var (
	errOutOfRangeMinAge              = errors.New("最小年齢が範囲外です")
	errOutOfRangeMaxAge              = errors.New("最大年齢が範囲外です")
	errEmptyResidenceArea            = errors.New("居住地が空です")
	errEmptyLastLogin                = errors.New("最終ログインが空です")
	errCannotFindIDFromResidenceArea = errors.New("居住地からIDを見つけることができません")
)

const (
	homeURL      = "https://pairs.lv/search"
	filteringURL = "https://pairs.lv/search/setup_condition/grid/index"
)

// Footprint 足跡
type Footprint struct {
	MinAge        uint8
	MaxAge        uint8
	ResidenceArea []string
	LastLogin     string
}

// NewFootprint 初期化
func NewFootprint(minAge, maxAge uint8, residenceArea []string, lastLogin string) (*Footprint, error) {
	f := &Footprint{
		MinAge:        minAge,
		MaxAge:        maxAge,
		ResidenceArea: residenceArea,
		LastLogin:     lastLogin,
	}
	if err := f.validate(); err != nil {
		return nil, err
	}
	return f, nil
}

// Filtering 条件でフィルタリング
func (f *Footprint) Filtering(selenium *selenium.Selenium) error {
	page := selenium.Page
	if err := page.Navigate(filteringURL); err != nil {
		return err
	}
	time.Sleep(time.Second * 2)

	selects := page.AllByClass("css-1bq0nkw")

	// 最小年齢の設定
	min := selects.At(0)
	if err := min.Select(fmt.Sprintf("%d歳", f.MinAge)); err != nil {
		return err
	}

	// 最大年齢の設定
	max := selects.At(1)
	if err := max.Select(fmt.Sprintf("%d歳", f.MaxAge)); err != nil {
		return err
	}

	// 最終ログインの設定
	last := selects.At(4)
	if err := last.Select(f.LastLogin); err != nil {
		return err
	}

	// 居住地から ID を取得
	ids, err := f.findIDForResidenceArea()
	if err != nil {
		return err
	}

	// 居住地の設定
	if err := page.AllByClass("css-a3zx38").At(0).Click(); err != nil {
		return err
	}
	time.Sleep(time.Second * 1)
	if err := page.AllByClass("css-17tl92q").At(0).Click(); err != nil {
		return err
	}
	if err := page.AllByClass("css-17tl92q").At(0).Click(); err != nil {
		return err
	}

	for _, id := range ids {
		if err := page.AllByClass("css-17tl92q").At(id).Click(); err != nil {
			return err
		}
	}

	// 保存
	if err := page.FindByButton("決定").Click(); err != nil {
		return err
	}
	time.Sleep(time.Second * 1)
	if err := page.FindByButton("この条件で検索").Click(); err != nil {
		return err
	}

	time.Sleep(time.Second * 4)

	return nil
}

// findIDForResidenceArea 居住地から ID を検索する
func (f *Footprint) findIDForResidenceArea() ([]int, error) {
	file, err := os.Open("./config/residence-area-master.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var ids []int
	for _, area := range f.ResidenceArea {
		for _, r := range rows {
			if area == r[1] {
				id, err := strconv.Atoi(r[0])
				if err != nil {
					return nil, err
				}
				ids = append(ids, id)
			}
		}
	}

	if len(ids) == 0 {
		return nil, errCannotFindIDFromResidenceArea
	}

	return ids, nil
}

func (f *Footprint) validate() error {
	if f.MinAge < 18 || f.MinAge > 65 {
		return errOutOfRangeMinAge
	}
	if f.MinAge < 18 || f.MinAge > 65 {
		return errOutOfRangeMinAge
	}
	for _, area := range f.ResidenceArea {
		if area == "" {
			return errEmptyResidenceArea
		}
	}
	if f.LastLogin == "" {
		return errEmptyLastLogin
	}
	return nil
}
