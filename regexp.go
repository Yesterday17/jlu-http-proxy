package main

import "regexp"

var PathMatchRegex = regexp.MustCompile("/(https?)(?:-\\d*)?/([0-9a-f]+)(.+)")

var RedirectLink = regexp.MustCompile("https://vpns.jlu.edu.cn/https?(?:-[0-9]*)?/([0-9a-f]+)([^\")]+)")

var HttpsToHttp = regexp.MustCompile("https://")

var VPNsLinkMatch = regexp.MustCompile("\"(?:https?://)?(?:vpns.jlu.edu.cn)?/https?(?:-[0-9]*)?/([0-9a-f]+)([^\")]+)")

var LinkUnescape = regexp.MustCompile("(href|link|src)=\"([^\"]+)\"")

var VPNEvalPrefix = regexp.MustCompile("vpn_eval\\(\\(function\\(\\){\r?\n")

var VPNEvalPostfix = regexp.MustCompile("}\r?\n\\).toString\\(\\).slice\\(12, -2\\),\"\"\\);")

var VPNRewritePrefix = regexp.MustCompile("var vpn_return;eval\\(vpn_rewrite_js\\(\\(function \\(\\) { ")

var VPNRewritePostfix = regexp.MustCompile(" }\\).toString\\(\\)\\.slice\\(14, -2\\), 2\\)\\);return vpn_return;")

var VPNParamRemoveFirst = regexp.MustCompile("vpn-\\d+&")

var VPNParamRemoveOther = regexp.MustCompile("\\?vpn-\\d+")

var VPNScriptInfo = regexp.MustCompile("<script>(?:\r?\nvar __vpn_[^ ]+ = [^;]+;)+\r?\n</script>")

var VPNScriptWEngine = regexp.MustCompile("<script src=\"/wengine-vpn/js/main.js[^>]+></script>")