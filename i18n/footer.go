package i18n
// Footer
func Footer (lang string) map[string]string {
	en := make(map[string]string)
	en["1"] = "guide"
	en["2"] = "contact"
	en["3"] = "company"
	en["4"] = "rule"
	en["5"] = "privacy policy"
	en["6"] = "pecified Commercial Transaction Act"
	en["7"] = "top"

	ja := make(map[string]string)
	ja["1"] = "ご利用ガイド"
	ja["2"] = "お問い合わせ"
	ja["3"] = "会社概要"
	ja["4"] = "ご利用規約"
	ja["5"] = "プライバシーポリシー"
	ja["6"] = "特定商取引法に基づく表記"
	ja["7"] = "トップ"

	language := make(map[string]map[string]string)
	language["en"] = en
	language["ja"] = ja
	return language[lang]
}