package expressiontree

type IExecutionChecker interface {
	IHttpDataChecker
	IOptionsChecker
	IClientChecker
	IRequestChecker
	IResponseChecker
}

type IHttpDataChecker interface {
	RecursiveCheckHttpData(predicate Predicate, data *HttpData) (bool, error)
	CheckHttpDataHost(predicate Predicate, data *HttpData) (bool, error)
	CheckHttpDataProtocol(predicate Predicate, data *HttpData) (bool, error)
	CheckHttpDataPort(predicate Predicate, data *HttpData) (bool, error)
	CheckHttpDataHttpVersion(predicate Predicate, data *HttpData) (bool, error)
	CheckHttpDataTimestamp(predicate Predicate, data *HttpData) (bool, error)
}

type IOptionsChecker interface {
	RecursiveCheckOptions(predicate Predicate, data *HttpData) (bool, error)
	CheckOptionExistence(optionName string, data *HttpData) (bool, error)
	CheckOption(predicate Predicate, optionName string, data *HttpData) (bool, error)
}

type IGeoIpChecker interface {
	RecursiveCheckGeoIp(predicate Predicate, data *HttpData) (bool, error)
	CheckGeoIpCountry(predicate Predicate, data *HttpData) (bool, error)
	CheckGeoIpCountryCode(predicate Predicate, data *HttpData) (bool, error)
	CheckGeoIpCity(predicate Predicate, data *HttpData) (bool, error)
	CheckGeoIpLat(predicate Predicate, data *HttpData) (bool, error)
	CheckGeoIpLon(predicate Predicate, data *HttpData) (bool, error)
	CheckGeoIpAccuracyRadius(predicate Predicate, data *HttpData) (bool, error)
}

type IOsChecker interface {
	RecursiveCheckOs(predicate Predicate, data *HttpData) (bool, error)
	CheckOsName(predicate Predicate, data *HttpData) (bool, error)
	CheckOsVersion(predicate Predicate, data *HttpData) (bool, error)
}

type IBrowserChecker interface {
	RecursiveCheckBrowser(predicate Predicate, data *HttpData) (bool, error)
	CheckBrowserName(predicate Predicate, data *HttpData) (bool, error)
	CheckBrowserVersion(predicate Predicate, data *HttpData) (bool, error)
}

type IBasicAuthChecker interface {
	RecursiveCheckBasicAuth(predicate Predicate, data *HttpData) (bool, error)
	CheckBasicAuthUsername(predicate Predicate, data *HttpData) (bool, error)
	CheckBasicAuthPassword(predicate Predicate, data *HttpData) (bool, error)
}

type IClientChecker interface {
	IGeoIpChecker
	IOsChecker
	IBrowserChecker
	IBasicAuthChecker
	RecursiveCheckClient(predicate Predicate, data *HttpData) (bool, error)
	CheckClientId(predicate Predicate, data *HttpData) (bool, error)
	CheckClientIp(predicate Predicate, data *HttpData) (bool, error)
}

type IRequestGetChecker interface {
	RecursiveCheckRequestGet(predicate Predicate, data *HttpData) (bool, error)
	CheckRequestGetValueExistence(path ContentPath, data *HttpData) (bool, error)
	RecursiveCheckRequestGetValue(predicate Predicate, path ContentPath, data *HttpData) (bool, error)
}

type IRequestPostChecker interface {
	RecursiveCheckRequestPost(predicate Predicate, data *HttpData) (bool, error)
	CheckRequestPostValueExistence(path ContentPath, data *HttpData) (bool, error)
	RecursiveCheckRequestPostValue(predicate Predicate, path ContentPath, data *HttpData) (bool, error)
}

type IRequestHeadersChecker interface {
	RecursiveCheckRequestHeaders(predicate Predicate, data *HttpData) (bool, error)
	CheckRequestHeaderValueExistence(path ContentPath, data *HttpData) (bool, error)
	RecursiveCheckRequestHeaderValue(predicate Predicate, path ContentPath, data *HttpData) (bool, error)
}

type IRequestCookiesChecker interface {
	RecursiveCheckRequestCookies(predicate Predicate, data *HttpData) (bool, error)
	CheckRequestCookieValueExistence(path ContentPath, data *HttpData) (bool, error)
	RecursiveCheckRequestCookieValue(predicate Predicate, path ContentPath, data *HttpData) (bool, error)
}

type IRequestChecker interface {
	IRequestGetChecker
	IRequestPostChecker
	IRequestHeadersChecker
	IRequestCookiesChecker
	RecursiveCheckRequest(predicate Predicate, data *HttpData) (bool, error)
	CheckRequestId(predicate Predicate, data *HttpData) (bool, error)
	CheckRequestPath(predicate Predicate, data *HttpData) (bool, error)
	CheckRequestPaths(predicate Predicate, data *HttpData) (bool, error)
	CheckRequestPathsElement(predicate Predicate, index int, path ContentPath, data *HttpData) (bool, error)
	CheckRequestQuery(predicate Predicate, data *HttpData) (bool, error)
	CheckRequestMethod(predicate Predicate, data *HttpData) (bool, error)
	RecursiveCheckRequestBody(predicate Predicate, path ContentPath, data *HttpData) (bool, error)
	CheckRequestTime(predicate Predicate, data *HttpData) (bool, error)
	CheckRequestLength(predicate Predicate, data *HttpData) (bool, error)
}

type IResponseHeadersChecker interface {
	RecursiveCheckResponseHeaders(predicate Predicate, data *HttpData) (bool, error)
	CheckResponseHeaderValueExistence(path ContentPath, data *HttpData) (bool, error)
	RecursiveCheckResponseHeaderValue(predicate Predicate, path ContentPath, data *HttpData) (bool, error)
}

type IResponseChecker interface {
	IResponseHeadersChecker
	RecursiveCheckResponse(predicate Predicate, data *HttpData) (bool, error)
	RecursiveCheckResponseBody(predicate Predicate, path ContentPath, data *HttpData) (bool, error)
	CheckResponseCode(predicate Predicate, data *HttpData) (bool, error)
	CheckResponseSource(predicate Predicate, data *HttpData) (bool, error)
	CheckResponseLength(predicate Predicate, data *HttpData) (bool, error)
}
