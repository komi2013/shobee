package i18n
// Sign
func Sign (lang string) map[string]string {
	en := make(map[string]string)
	en["1"] = "Google Sign In"
	en["2"] = "Facebook Sign In"
	en["3"] = "LINE Sign In"
	en["4"] = "email"
	en["5"] = "must be email format"
	en["6"] = "password"
	en["7"] = "include upper and lower case and number"
	en["8"] = "Sign In"
	en["9"] = "password confirm"
	en["10"] = "family name"
	en["11"] = "given name"
	en["12"] = "Birthday"
	en["13"] = "Register"
	en["14"] = "Password does not match"

	ja := make(map[string]string)
	ja["1"] = "Googleでログイン"
	ja["2"] = "Facebookでログイン"
	ja["3"] = "LINEでログイン"
	ja["4"] = "メール"
	ja["5"] = "emailのフォーマットを入力"
	ja["6"] = "パスワード"
	ja["7"] = "大文字、小文字、数字を入力"
	ja["8"] = "ログイン"
	ja["9"] = "パスワード確認"
	ja["10"] = "姓"
	ja["11"] = "名"
	ja["12"] = "生年月日"
	ja["13"] = "登録"
	ja["14"] = "パスワードが一致しません"

	language := make(map[string]map[string]string)
	language["en"] = en
	language["ja"] = ja
	return language[lang]
}
// 62pm7wS1