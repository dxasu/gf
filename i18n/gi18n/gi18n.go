// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/dxasu/gf.

// Package gi18n implements internationalization and localization.
package gi18n

// SetPath sets the directory path storing i18n files.
func SetPath(path string) error {
	return Instance().SetPath(path)
}

// SetLanguage sets the language for translator.
func SetLanguage(language string) {
	Instance().SetLanguage(language)
}

// SetDelimiters sets the delimiters for translator.
func SetDelimiters(left, right string) {
	Instance().SetDelimiters(left, right)
}

// T is alias of Translate for convenience.
func T(content string, language ...string) string {
	return Instance().T(content, language...)
}

// Tf is alias of TranslateFormat for convenience.
func Tf(format string, values ...interface{}) string {
	return Instance().TranslateFormat(format, values...)
}

// Tfl is alias of TranslateFormatLang for convenience.
func Tfl(language string, format string, values ...interface{}) string {
	return Instance().TranslateFormatLang(language, format, values...)
}

// TranslateFormat translates, formats and returns the <format> with configured language
// and given <values>.
func TranslateFormat(format string, values ...interface{}) string {
	return Instance().TranslateFormat(format, values...)
}

// TranslateFormatLang translates, formats and returns the <format> with configured language
// and given <values>. The parameter <language> specifies custom translation language ignoring
// configured language. If <language> is given empty string, it uses the default configured
// language for the translation.
func TranslateFormatLang(language string, format string, values ...interface{}) string {
	return Instance().TranslateFormatLang(language, format, values...)
}

// Translate translates <content> with configured language and returns the translated content.
// The parameter <language> specifies custom translation language ignoring configured language.
func Translate(content string, language ...string) string {
	return Instance().Translate(content, language...)
}

// GetValue retrieves and returns the configured content for given key and specified language.
// It returns an empty string if not found.
func GetContent(key string, language ...string) string {
	return Instance().GetContent(key, language...)
}
