package main

import "regexp"

var PathMatchRegex = regexp.MustCompile("/(https?)(?:-\\d*)?/([0-9a-f]{32,})(.+)")

var RedirectLink = regexp.MustCompile("https://vpns.jlu.edu.cn/(https?)(?:-[0-9]*)?/([0-9a-f]{32,})(.*)")

var VPNsLinkMatch = regexp.MustCompile("(?:https:)?(?://)?(?:vpns.jlu.edu.cn)?/(https?)(?:-[0-9]*)?/([0-9a-f]{32,})([^\")]+)")

var LinkUnescape = regexp.MustCompile("(href|link|src)=\"([^\"]+)\"")

var VPNEvalPrefix = regexp.MustCompile("vpn_eval\\(\\(function\\(\\){\r?\n")

var VPNEvalPostfix = regexp.MustCompile("}\r?\n\\).toString\\(\\).slice\\(12, -2\\),\"\"\\);")

var VPNRewritePrefix = regexp.MustCompile("var vpn_return;eval\\(vpn_rewrite_js\\(\\(function \\(\\) { ")

var VPNRewritePostfix = regexp.MustCompile(" }\\).toString\\(\\)\\.slice\\(14, -2\\), 2\\)\\);return vpn_return;")

var VPNInjectPrefix = regexp.MustCompile("javascript:this.top.vpn_inject_scripts_window\\(this\\);vpn_eval\\(\\(function \\(\\) { ")

var VPNInjectPostfix = regexp.MustCompile(" }\\)\\.toString\\(\\).slice\\(14, -2\\)\\)")

var VPNParamRemoveFirst = regexp.MustCompile("vpn-\\d+&")

var VPNParamRemoveOther = regexp.MustCompile("\\?vpn-\\d+")

var VPNScriptInfo = regexp.MustCompile("<script>(?:\r?\nvar __vpn_[^ ]+ = [^;]+;)+\r?\n</script>")

var VPNScriptWEngine = regexp.MustCompile("<script src=\"/wengine-vpn/js/main.js[^>]+></script>")
