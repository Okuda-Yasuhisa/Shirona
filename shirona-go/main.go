package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

func fetch(url string) []byte {
	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	// mainがReturnしたら実行
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

// map[KeyType]ValueType
type Record map[string]any

// []Recordは、可変長のRecord
func asRecordSet(value any) ([]Record, bool) {
	// valueがanyのスライス型なら、arrayに値を代入する
	// anyのスライス型ではないなら、nilを返す
	array, ok := value.([]any)

	if !ok || len(array) == 0 {
		return nil, false
	}

	// Recordのスライスを作成
	// 初期の長さは0であるが、最大5まで拡張可能
	records := make([]Record, 0, len(array))

	for _, item := range array {
		// itemが string:any のmapなら、objectに値を代入する
		object, ok := item.(map[string]any)

		if !ok {
			return nil, false
		}

		records = append(records, Record(object))
	}

	return records, true
}

// ... は可変長引数
func getRecord(object map[string]any, keys ...string) (any, bool) {
	var current any = object

	// objectのキー分繰り返す
	for _, key := range keys {
		// currentが string:any のmap型かどうか?
		currentObject, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}

		// 取り出したオブジェクトに指定したキーがあるか?
		value, exists := currentObject[key]
		if !exists {
			return nil, false
		}

		current = value
	}

	return current, true
}

func createArgParser() (string, []string) {
	// 変数名 = flag.String("オプション名", "デフォルト値", "フラグの説明") で定義
	url := flag.String("url", "", "fetchを行うURL")

	flag.Parse()

	if *url == "" {
		log.Fatalln("--urlを指定してください")
	}

	keys := flag.Args()

	if len(keys) == 0 {
		log.Fatalln("JSONのキーを指定してください")
	}

	return *url, keys
}

func main() {
	url, keys := createArgParser()

	jsonStr := fetch(url)

	var data any
	err := json.Unmarshal(jsonStr, &data)
	if err != nil {
		log.Fatalln(err)
	}

	root, ok := data.(map[string]any)
	if !ok {
		log.Fatalln("ルートがObjectではない")
	}

	value, ok := getRecord(root, keys...)
	if !ok {
		log.Fatalln("パスが存在しない")
	}

	records, ok := asRecordSet(value)
	if !ok {
		log.Fatalln("RecordSetではない")
	}

	fmt.Println("レコード数: ", len(records))
	fmt.Println("最初のレコード: ", records[0])
	fmt.Println("タイトル: ", records[0]["title"])
}
