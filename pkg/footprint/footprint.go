package footprint

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"pairs/pkg/logging"
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
	logging.Info("検索条件を設定します")

	page := selenium.Page
	if err := page.Navigate(filteringURL); err != nil {
		return err
	}
	time.Sleep(time.Second * 2)

	// 検索条件のクリア
	if err := page.FindByButton("すべてリセット").Click(); err != nil {
		return err
	}
	time.Sleep(time.Second * 1)

	selects := page.AllByClass("css-1bq0nkw")

	// 最小年齢の設定
	min := selects.At(0)
	if err := min.Select(fmt.Sprintf("%d歳", f.MinAge)); err != nil {
		return err
	}
	time.Sleep(time.Second * 1)

	// 最大年齢の設定
	max := selects.At(1)
	if err := max.Select(fmt.Sprintf("%d歳", f.MaxAge)); err != nil {
		return err
	}
	time.Sleep(time.Second * 1)

	// 最終ログインの設定
	last := selects.At(4)
	if err := last.Select(f.LastLogin); err != nil {
		return err
	}
	time.Sleep(time.Second * 1)

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

	areas := page.AllByClass("css-17tl92q")
	for _, id := range ids {
		if err := areas.At(id).Click(); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}

	// 保存
	if err := page.FindByButton("決定").Click(); err != nil {
		return err
	}
	time.Sleep(time.Second * 1)
	if err := page.FindByButton("この条件で検索").Click(); err != nil {
		return err
	}

	return nil
}

// StepOn 足跡を残す
func (f *Footprint) StepOn(selenium *selenium.Selenium) error {
	page := selenium.Page
	if err := page.Navigate(homeURL); err != nil {
		return err
	}
	time.Sleep(time.Second * 2)

	// 1番目のユーザーをクリック
	if err := page.AllByClass("css-opde7s").At(0).Click(); err != nil {
		return err
	}
	time.Sleep(time.Second * 1)

	// 右矢印をクリック
	if err := page.AllByClass("css-1d94zew").At(0).Click(); err != nil {
		return err
	}

	// 乱数の初期化
	rand.Seed(time.Now().UnixNano())

	// 足跡を残す
	var c int
	for {
		if err := page.AllByClass("css-1d94zew").At(1).Click(); err != nil {
			// ユーザーの詳細画面を閉じる
			time.Sleep(time.Second * 2)
			if err := page.AllByClass("css-1ureyjg").At(0).Click(); err != nil {
				return err
			}
			time.Sleep(time.Second * 1)

			// ページの最下部まで移動し、次のユーザーを読み込む
			page.RunScript("window.scroll(0, document.documentElement.scrollHeight - document.documentElement.clientHeight);", nil, nil)
			time.Sleep(time.Second * 1)

			// 次のユーザーをクリックし、処理を再開
			if err := page.AllByClass("css-opde7s").At(c).Click(); err != nil {
				return err
			}
			time.Sleep(time.Second * 3)
		}

		c++
		println(c)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(8000)))
	}
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
