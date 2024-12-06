package expressiontree

// main path
type TDataKey int

const (
	HttpDataKey TDataKey = iota + 1
	HttpDataHostKey
	HttpDataProtocolKey
	HttpDataPortKey
	HttpDataHttpVersionKey
	HttpDataTimestampKey
	OptionsKey
	ClientKey
	ClientIdKey
	ClientIpKey
	GeoIpKey
	GeoIpCountryKey
	GeoIpCountryCodeKey
	GeoIpCityKey
	GeoIpLatKey
	GeoIpLonKey
	GeoIpAccuracyRadiusKey
	OsKey
	OsNameKey
	OsVersionKey
	BrowserKey
	BrowserNameKey
	BrowserVersionKey
	BasicAuthKey
	BasicAuthUsernameKey
	BasicAuthPasswordKey
	RequestKey
	RequestIdKey
	RequestPathKey
	RequestPathsKey
	RequestQueryKey
	RequestMethodKey
	RequestBodyKey
	RequestGetKey
	RequestPostKey
	RequestHeadersKey
	RequestTimeKey
	RequestCookiesKey
	RequestLengthKey
	ResponseKey
	ResponseBodyKey
	ResponseCodeKey
	ResponseSourceKey
	ResponseHeadersKey
	ResponseLengthKey
)
