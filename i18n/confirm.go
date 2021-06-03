package i18n
// Confirm
func Confirm (lang string) map[string]string {
	en := make(map[string]string)
	en["yen"] = "Yen"
	en["pay"] = "Purchase"

	ja := make(map[string]string)
	ja["yen"] = "円"
	ja["pay"] = "購入"

	language := make(map[string]map[string]string)
	language["en"] = en
	language["ja"] = ja
	return language[lang]
}
// 62pm7wS1