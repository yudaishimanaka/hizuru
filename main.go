package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

// Github ActionsでGorelaserを実行しリリースタグと最新のショートコミットをアプリケーションに埋め込む
var (
	Version  = "unset"
	Revision = "unset"
)

// 使用しているWindows Terminalがプレビュー版であるかチェックするフラグ
var isPreview bool

const (
	windowsTerminalDefault = "\\Packages\\Microsoft.WindowsTerminal_8wekyb3d8bbwe\\LocalState\\settings.json"
	windowsTerminalPreview = "\\Packages\\Microsoft.WindowsTerminalPreview_8wekyb3d8bbwe\\LocalState\\settings.json"
)

type image struct {
	Name string
	Base string
}

// 画像であるかチェックする関数
func isImage(imgName string) bool {
	ext := filepath.Ext(imgName)
	mime := mime.TypeByExtension(ext)

	isImage := false
	if mime != "" && (len(mime) > 5 && mime[0:5] == "image") {
		isImage = true
	}

	if isImage {
		return true
	} else {
		return false
	}
}

func getImageList(path string) ([]image, error) {
	var imgList []image

	// HIZURU_IMAGE_PATH配下のファイルとサブディレクトリをリスト化
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// ディレクトリでないかつ画像である時、画像情報をリストに追加
	for _, file := range files {
		if !file.IsDir() && isImage(file.Name()) {
			// 画像から拡張子を取り除く
			base := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			imgList = append(imgList, image{Name: file.Name(), Base: base})
		}
	}

	return imgList, nil
}

func selectImage(imgList []image) (image, error) {
	// プロンプトカスタマイズ用のテンプレート
	selectImageTemplates := &promptui.SelectTemplates{
		Active:   promptui.IconSelect + " {{ .Base | underline }}",
		Inactive: "  {{ .Base }}",
		Selected: promptui.IconGood + " {{ .Base | faint }}",
	}

	// 画像選択プロンプト
	selectImagePrompt := promptui.Select{
		Label:     "Select Background Image",
		Items:     imgList,
		Templates: selectImageTemplates,
	}

	// プロンプト実行
	i, _, err := selectImagePrompt.Run()
	if err != nil {
		return image{}, err
	}

	return imgList[i], nil
}

func saveJSON(jsonObj interface{}, path string) error {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer file.Close()
	err = json.NewEncoder(file).Encode(jsonObj)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func main() {
	// hizuruで設定するアプリケーション環境変数をロード
	hizuruPath := os.Getenv("HIZURU_IMAGE_PATH")
	if hizuruPath == "" {
		fmt.Println("Please set the HIZURU_IMAGE_PATH environment variable and run again")
		return
	}

	// LOCALAPPDATAをロード
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		fmt.Println("LOCALAPPDATA environment variable does not exist")
		return
	}

	app := &cli.App{
		Name:    "hizuru",
		Version: Version + " " + Revision,
		Usage:   "Change the background image of the Windows Terminal",
		// プレビュー版用のフラグ
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "preview",
				Usage:       "If you are using a preview version of Windows Treminal, please add the --preview flag",
				Aliases:     []string{"p"},
				Required:    false,
				Value:       false,
				Destination: &isPreview,
				DefaultText: "false",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "change",
				Aliases: []string{"c"},
				Usage:   "Change the background image by selecting from the list",
				Action: func(ctx *cli.Context) error {
					// 背景画像リストを取得
					imgList, err := getImageList(hizuruPath)
					if err != nil {
						return err
					}

					// 配列が空の場合は終了する
					if len(imgList) == 0 {
						fmt.Println("Image does not exist, please place the image under HIZURU_IMAGE_PATH")
						return nil
					}

					// 画像を選択するプロンプトを実行、選択された画像
					img, err := selectImage(imgList)
					if err != nil {
						return err
					}

					// プレビューフラグを用いてsettings.jsonのPATHを生成
					var settingFilePath string
					if isPreview {
						settingFilePath = localAppData + windowsTerminalPreview
					} else {
						settingFilePath = localAppData + windowsTerminalDefault
					}

					// settings.jsonをロード
					byteArray, err := ioutil.ReadFile(settingFilePath)
					if err != nil {
						fmt.Println("Configuration file could not be read, please check if the --preview flag is required")
						return err
					}

					var jsonObj interface{}
					_ = json.Unmarshal(byteArray, &jsonObj)

					jsonObj.(map[string]interface{})["profiles"].(map[string]interface{})["defaults"].(map[string]interface{})["backgroundImage"] = hizuruPath + "\\" + img.Name

					err = saveJSON(jsonObj, settingFilePath)
					if err != nil {
						fmt.Println(err)
						return err
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
