package i18n
// Item
func Item (lang string) map[string]string {
	en := make(map[string]string)
	en["1"] = "genre"
	en["2"] = "exclude keyword"
	en["3"] = "minimum"
	en["4"] = "maximum"
	en["5"] = "Yen"
	en["6"] = "purchase confirm"

	ja := make(map[string]string)
	ja["1"] = "ジャンル"
	ja["2"] = "含まないキーワード"
	ja["3"] = "低い価格"
	ja["4"] = "高い価格"
	ja["5"] = "円"
	ja["6"] = "購入確認"

	language := make(map[string]map[string]string)
	language["en"] = en
	language["ja"] = ja
	return language[lang]
}
