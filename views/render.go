package views

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"html/template"
	"io/ioutil"
	"os"
)

var MixManifestMap = make(map[string]string)

func LoadAssetMap()  {
	// load mix manifest file
	if mixFileStat, err := os.Stat("mix-manifest.json"); err == nil && mixFileStat.IsDir() == false {
		if mixManifestFile, err := os.Open("mix-manifest.json"); err == nil {
			b, err := ioutil.ReadAll(mixManifestFile)
			if err == nil {
				_ = json.Unmarshal(b, &MixManifestMap)
			}
		}
	}
}

func assetFunc(path string) string {
	pathWithQuery, ok := MixManifestMap["/static/" + path]
	if ok && len(pathWithQuery) != 0 {
		return pathWithQuery
	}
	return "/static/" + path
}

func marshalFunc(v interface {})  template.JS {
	a, _ := json.Marshal(v)
	return template.JS(a)
}

func InitRenderFunc() {
	var err error
	LoadAssetMap()
	err = beego.AddFuncMap("asset", assetFunc)
	if err != nil {
		logs.Warn("add function fail")
	}

	err = beego.AddFuncMap("marshal", marshalFunc)
}
