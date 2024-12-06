package expressiontree

type IExecutionMatcher interface {
	IHttpDataMatcher
	IOptionsMatcher
	IClientMatcher
	IRequestMatcher
	IResponseMatcher
}

type IHttpDataMatcher interface {
	RecursiveMatchHttpData(patternId uint, data *HttpData) (bool, error)
	MatchHttpDataHost(patternId uint, data *HttpData) (bool, error)
	MatchHttpDataProtocol(patternId uint, data *HttpData) (bool, error)
	MatchHttpDataPort(patternId uint, data *HttpData) (bool, error)
	MatchHttpDataHttpVersion(patternId uint, data *HttpData) (bool, error)
	MatchHttpDataTimestamp(patternId uint, data *HttpData) (bool, error)
}

type IOptionsMatcher interface {
	RecursiveMatchOptions(patternId uint, data *HttpData) (bool, error)
	MatchOption(patternId uint, name string, data *HttpData) (bool, error)
}

type IGeoIpMatcher interface {
	RecursiveMatchGeoIp(patternId uint, data *HttpData) (bool, error)
	MatchGeoIpCountry(patternId uint, data *HttpData) (bool, error)
	MatchGeoIpCountryCode(patternId uint, data *HttpData) (bool, error)
	MatchGeoIpCity(patternId uint, data *HttpData) (bool, error)
	MatchGeoIpLat(patternId uint, data *HttpData) (bool, error)
	MatchGeoIpLon(patternId uint, data *HttpData) (bool, error)
	MatchGeoIpAccuracyRadius(patternId uint, data *HttpData) (bool, error)
}

type IOsMatcher interface {
	RecursiveMatchOs(patternId uint, data *HttpData) (bool, error)
	MatchOsName(patternId uint, data *HttpData) (bool, error)
	MatchOsVersion(patternId uint, data *HttpData) (bool, error)
}

type IBrowserMatcher interface {
	RecursiveMatchBrowser(patternId uint, data *HttpData) (bool, error)
	MatchBrowserName(patternId uint, data *HttpData) (bool, error)
	MatchBrowserVersion(patternId uint, data *HttpData) (bool, error)
}

type IBasicAuthMatcher interface {
	RecursiveMatchBasicAuth(patternId uint, data *HttpData) (bool, error)
	MatchBasicAuthUsername(patternId uint, data *HttpData) (bool, error)
	MatchBasicAuthPassword(patternId uint, data *HttpData) (bool, error)
}

type IClientMatcher interface {
	IGeoIpMatcher
	IOsMatcher
	IBrowserMatcher
	IBasicAuthMatcher
	RecursiveMatchClient(patternId uint, data *HttpData) (bool, error)
	MatchClientId(patternId uint, data *HttpData) (bool, error)
	MatchClientIp(patternId uint, data *HttpData) (bool, error)
}

type IRequestGetMatcher interface {
	RecursiveMatchRequestGet(patternId uint, data *HttpData) (bool, error)
	RecursiveMatchRequestGetValue(patternId uint, path ContentPath, data *HttpData) (bool, error)
}

type IRequestPostMatcher interface {
	RecursiveMatchRequestPost(patternId uint, data *HttpData) (bool, error)
	RecursiveMatchRequestPostValue(patternId uint, path ContentPath, data *HttpData) (bool, error)
}

type IRequestHeadersMatcher interface {
	RecursiveMatchRequestHeaders(patternId uint, data *HttpData) (bool, error)
	RecursiveMatchRequestHeaderValue(patternId uint, path ContentPath, data *HttpData) (bool, error)
}

type IRequestCookiesMatcher interface {
	RecursiveMatchRequestCookies(patternId uint, data *HttpData) (bool, error)
	RecursiveMatchRequestCookieValue(patternId uint, path ContentPath, data *HttpData) (bool, error)
}

type IRequestMatcher interface {
	IRequestGetMatcher
	IRequestPostMatcher
	IRequestHeadersMatcher
	IRequestCookiesMatcher
	RecursiveMatchRequest(patternId uint, data *HttpData) (bool, error)
	MatchRequestId(patternId uint, data *HttpData) (bool, error)
	MatchRequestPath(patternId uint, data *HttpData) (bool, error)
	MatchRequestPaths(patternId uint, data *HttpData) (bool, error)
	MatchRequestPathsElement(patternId uint, index int, path ContentPath, data *HttpData) (bool, error)
	MatchRequestQuery(patternId uint, data *HttpData) (bool, error)
	MatchRequestMethod(patternId uint, data *HttpData) (bool, error)
	RecursiveMatchRequestBody(patternId uint, path ContentPath, data *HttpData) (bool, error)
	MatchRequestTime(patternId uint, data *HttpData) (bool, error)
	MatchRequestLength(patternId uint, data *HttpData) (bool, error)
}

type IResponseHeadersMatcher interface {
	RecursiveMatchResponseHeaders(patternId uint, data *HttpData) (bool, error)
	RecursiveMatchResponseHeaderValue(patternId uint, path ContentPath, data *HttpData) (bool, error)
}

type IResponseMatcher interface {
	IResponseHeadersMatcher
	RecursiveMatchResponse(patternId uint, data *HttpData) (bool, error)
	RecursiveMatchResponseBody(patternId uint, path ContentPath, data *HttpData) (bool, error)
	MatchResponseCode(patternId uint, data *HttpData) (bool, error)
	MatchResponseSource(patternId uint, data *HttpData) (bool, error)
	MatchResponseLength(patternId uint, data *HttpData) (bool, error)
}
