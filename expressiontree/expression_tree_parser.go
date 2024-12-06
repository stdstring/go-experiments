package expressiontree

import (
	"errors"
	"strconv"
	"strings"
)

// AND(COND1,COND2,...)
// OR(COND1,COND2,...)
// NOT(COND)
// COND: CHECK(PATH,OP,ARG) | EXISTS(PATH) | MATCH(PATH,PATTERN)
// OP: INT (in table), PATH: INT (in table), ARG: INT (in table), PATTERN: INT (in table)

const (
	OperationEqual    = 0
	OperationNotEqual = 1
)

var parseError = errors.New("parse error")
var badArgsError = errors.New("bad args error")
var unknownExpressionError = errors.New("unknown expression path")
var unsupportedOperationError = errors.New("unsupported operation error")
var unknownMainPathError = errors.New("unknown main path")
var badContentPathError = errors.New("bad content path")
var badRequestPathIndexError = errors.New("bad request path index")

type PredicateWithError func(data *HttpData, manager IExecutionManager) (bool, error)

type sourceReader struct {
	//position int
	source string
}

func (r *sourceReader) readTo(end string) (string, error) {
	index := strings.IndexAny(r.source, end)
	if index == -1 {
		return "", parseError
	}
	result := string([]rune(r.source)[0 : index+1])
	rest := string([]rune(r.source)[index+1:])
	//r.position = index + 1
	r.source = rest
	return result, nil
}

func (r *sourceReader) readCurrent() (string, error) {
	if len(r.source) == 0 {
		return "", parseError
	}
	result := string([]rune(r.source)[0])
	rest := string([]rune(r.source)[1:])
	r.source = rest
	return result, nil
}

func (r *sourceReader) isEmpty() bool {
	return len(r.source) == 0
}

func newSourceReader(source string) *sourceReader {
	return &sourceReader{source}
}

type parseStorage struct {
	knownPath      []DataPath
	checkArguments []any
}

func parseExpressionTree(source string, storage *parseStorage) (PredicateWithError, error) {
	reader := newSourceReader(source)
	expression, expressionError := parseExpression(reader, storage)
	if expressionError != nil {
		return nil, expressionError
	}
	if !reader.isEmpty() {
		return nil, parseError
	}
	return expression, nil
}

func parseExpression(reader *sourceReader, storage *parseStorage) (PredicateWithError, error) {
	expressionHead, readError := reader.readTo("(")
	if readError != nil {
		return nil, readError
	}
	switch strings.TrimSuffix(expressionHead, "(") {
	case "AND":
		arguments, argumentsError := parseLogicalExpressionArgs(reader, storage)
		if argumentsError != nil {
			return nil, argumentsError
		}
		return createLogicalAnd(arguments...), nil
	case "OR":
		arguments, argumentsError := parseLogicalExpressionArgs(reader, storage)
		if argumentsError != nil {
			return nil, argumentsError
		}
		return createLogicalOr(arguments...), nil
	case "NOT":
		innerExpression, innerExpressionErr := parseNotArg(reader, storage)
		if innerExpressionErr != nil {
			return nil, innerExpressionErr
		}
		return createLogicalNot(innerExpression), nil
	case "CHECK":
		return parseCheck(reader, storage)
	case "EXISTS":
		return parseExists(reader, storage)
	case "MATCH":
		return parseMatch(reader, storage)
	default:
		return nil, unknownExpressionError
	}
}

func parseLogicalExpressionArgs(reader *sourceReader, storage *parseStorage) ([]PredicateWithError, error) {
	arguments := make([]PredicateWithError, 0)
	for {
		argument, argumentError := parseExpression(reader, storage)
		if argumentError != nil {
			return nil, argumentError
		}
		arguments = append(arguments, argument)
		currentRune, readError := reader.readCurrent()
		if readError != nil {
			return nil, readError
		}
		switch currentRune {
		case ",":
		case ")":
			if len(arguments) <= 1 {
				return nil, badArgsError
			}
			return arguments, nil
		default:
			return nil, parseError
		}
	}
}

func parseNotArg(reader *sourceReader, storage *parseStorage) (PredicateWithError, error) {
	innerExpression, innerExpressionErr := parseExpression(reader, storage)
	if innerExpressionErr != nil {
		return nil, innerExpressionErr
	}
	currentRune, readError := reader.readCurrent()
	if readError != nil {
		return nil, readError
	}
	if currentRune != ")" {
		return nil, parseError
	}
	return innerExpression, nil
}

func parseExists(reader *sourceReader, storage *parseStorage) (PredicateWithError, error) {
	value, readError := reader.readTo(")")
	if readError != nil {
		return nil, readError
	}
	arguments, argumentsError := parseArguments(strings.TrimSuffix(value, ")"), 1)
	if argumentsError != nil {
		return nil, argumentsError
	}
	path := storage.knownPath[arguments[0]]
	return createExists(path)
}

func parseMatch(reader *sourceReader, storage *parseStorage) (PredicateWithError, error) {
	value, readError := reader.readTo(")")
	if readError != nil {
		return nil, readError
	}
	arguments, argumentsError := parseArguments(strings.TrimSuffix(value, ")"), 2)
	if argumentsError != nil {
		return nil, argumentsError
	}
	path := storage.knownPath[arguments[0]]
	patternId := arguments[1]
	return createMatch(path, uint(patternId))
}

func parseCheck(reader *sourceReader, storage *parseStorage) (PredicateWithError, error) {
	value, readError := reader.readTo(")")
	if readError != nil {
		return nil, readError
	}
	arguments, argumentsError := parseArguments(strings.TrimSuffix(value, ")"), 3)
	if argumentsError != nil {
		return nil, argumentsError
	}
	path := storage.knownPath[arguments[0]]
	operation := arguments[1]
	checkArg := storage.checkArguments[arguments[2]]
	predicate, predicateError := parsePredicate(operation, checkArg)
	if predicateError != nil {
		return nil, predicateError
	}
	return createCheck(path, predicate)
}

func parseArguments(source string, expectedParts int) ([]int, error) {
	arguments := strings.Split(source, ",")
	if len(arguments) != expectedParts {
		return nil, badArgsError
	}
	result := make([]int, len(arguments))
	for index, argument := range arguments {
		value, convertError := strconv.Atoi(argument)
		if convertError != nil {
			return nil, parseError
		}
		result[index] = value
	}
	return result, nil
}

func parsePredicate(operation int, argument any) (Predicate, error) {
	// TODO (std_string) : this is simple version for demo. in real version we must use type casting
	switch operation {
	case OperationEqual:
		return func(value any) bool {
			return value == argument
		}, nil
	case OperationNotEqual:
		return func(value any) bool {
			return value != argument
		}, nil
	default:
		return nil, unsupportedOperationError
	}
}

func createLogicalAnd(predicates ...PredicateWithError) PredicateWithError {
	return func(data *HttpData, manager IExecutionManager) (bool, error) {
		for _, predicate := range predicates {
			result, err := predicate(data, manager)
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil
			}
		}
		return true, nil
	}
}

func createLogicalOr(predicates ...PredicateWithError) PredicateWithError {
	return func(data *HttpData, manager IExecutionManager) (bool, error) {
		for _, predicate := range predicates {
			result, err := predicate(data, manager)
			if err != nil {
				return false, err
			}
			if result {
				return true, nil
			}
		}
		return false, nil
	}
}

func createLogicalNot(predicate PredicateWithError) PredicateWithError {
	return func(data *HttpData, manager IExecutionManager) (bool, error) {
		result, err := predicate(data, manager)
		if err != nil {
			return false, err
		}
		return !result, nil
	}
}

func createMatch(path DataPath, patternId uint) (PredicateWithError, error) {
	switch path.MainPath {
	case HttpDataKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchHttpData(patternId, data)
		}, nil
	case HttpDataHostKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchHttpDataHost(patternId, data)
		}, nil
	case HttpDataProtocolKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchHttpDataProtocol(patternId, data)
		}, nil
	case HttpDataPortKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchHttpDataPort(patternId, data)
		}, nil
	case HttpDataHttpVersionKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchHttpDataHttpVersion(patternId, data)
		}, nil
	case HttpDataTimestampKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchHttpDataTimestamp(patternId, data)
		}, nil
	case OptionsKey:
		return generateMatchOption(path, patternId)
	case ClientKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchClient(patternId, data)
		}, nil
	case ClientIdKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchClientId(patternId, data)
		}, nil
	case ClientIpKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchClientIp(patternId, data)
		}, nil
	case GeoIpKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchGeoIp(patternId, data)
		}, nil
	case GeoIpCountryKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchGeoIpCountry(patternId, data)
		}, nil
	case GeoIpCountryCodeKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchGeoIpCountryCode(patternId, data)
		}, nil
	case GeoIpCityKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchGeoIpCity(patternId, data)
		}, nil
	case GeoIpLatKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchGeoIpLat(patternId, data)
		}, nil
	case GeoIpLonKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchGeoIpLon(patternId, data)
		}, nil
	case GeoIpAccuracyRadiusKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchGeoIpAccuracyRadius(patternId, data)
		}, nil
	case OsKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchOs(patternId, data)
		}, nil
	case OsNameKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchOsName(patternId, data)
		}, nil
	case OsVersionKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchOsVersion(patternId, data)
		}, nil
	case BrowserKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchBrowser(patternId, data)
		}, nil
	case BrowserNameKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchBrowserName(patternId, data)
		}, nil
	case BrowserVersionKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchBrowserVersion(patternId, data)
		}, nil
	case BasicAuthKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchBasicAuth(patternId, data)
		}, nil
	case BasicAuthUsernameKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchBasicAuthUsername(patternId, data)
		}, nil
	case BasicAuthPasswordKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchBasicAuthPassword(patternId, data)
		}, nil
	case RequestKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchRequest(patternId, data)
		}, nil
	case RequestIdKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchRequestId(patternId, data)
		}, nil
	case RequestPathKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchRequestPath(patternId, data)
		}, nil
	case RequestPathsKey:
		return generateMatchRequestPaths(path, patternId)
	case RequestQueryKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchRequestQuery(patternId, data)
		}, nil
	case RequestMethodKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchRequestMethod(patternId, data)
		}, nil
	case RequestBodyKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchRequestBody(patternId, path.ContentPath, data)
		}, nil
	case RequestGetKey:
		return generateMatchRequestGet(path, patternId)
	case RequestPostKey:
		return generateMatchRequestPost(path, patternId)
	case RequestHeadersKey:
		return generateMatchRequestHeaders(path, patternId)
	case RequestTimeKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchRequestTime(patternId, data)
		}, nil
	case RequestCookiesKey:
		return generateMatchRequestCookies(path, patternId)
	case RequestLengthKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchRequestLength(patternId, data)
		}, nil
	case ResponseKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchResponse(patternId, data)
		}, nil
	case ResponseBodyKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchResponseBody(patternId, path.ContentPath, data)
		}, nil
	case ResponseCodeKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchResponseCode(patternId, data)
		}, nil
	case ResponseSourceKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchResponseSource(patternId, data)
		}, nil
	case ResponseHeadersKey:
		return generateMatchResponseHeaders(path, patternId)
	case ResponseLengthKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchResponseLength(patternId, data)
		}, nil
	default:
		return nil, unknownMainPathError
	}
}

func createExists(path DataPath) (PredicateWithError, error) {
	switch path.MainPath {
	case OptionsKey:
		if !path.ContentPath.IsSimple() {
			return nil, badContentPathError
		}
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckOptionExistence(path.ContentPath.Path, data)
		}, nil
	case RequestGetKey:
		if path.ContentPath.IsEmpty() {
			return nil, badContentPathError
		}
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckRequestGetValueExistence(path.ContentPath, data)
		}, nil
	case RequestPostKey:
		if path.ContentPath.IsEmpty() {
			return nil, badContentPathError
		}
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckRequestPostValueExistence(path.ContentPath, data)
		}, nil
	case RequestHeadersKey:
		if path.ContentPath.IsEmpty() {
			return nil, badContentPathError
		}
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckRequestHeaderValueExistence(path.ContentPath, data)
		}, nil
	case RequestCookiesKey:
		if path.ContentPath.IsEmpty() {
			return nil, badContentPathError
		}
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckRequestCookieValueExistence(path.ContentPath, data)
		}, nil
	case ResponseHeadersKey:
		if path.ContentPath.IsEmpty() {
			return nil, badContentPathError
		}
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckResponseHeaderValueExistence(path.ContentPath, data)
		}, nil
	default:
		return nil, unknownMainPathError
	}
}

func createCheck(path DataPath, predicate Predicate) (PredicateWithError, error) {
	switch path.MainPath {
	case HttpDataKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckHttpData(predicate, data)
		}, nil
	case HttpDataHostKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckHttpDataHost(predicate, data)
		}, nil
	case HttpDataProtocolKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckHttpDataProtocol(predicate, data)
		}, nil
	case HttpDataPortKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckHttpDataPort(predicate, data)
		}, nil
	case HttpDataHttpVersionKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckHttpDataHttpVersion(predicate, data)
		}, nil
	case HttpDataTimestampKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckHttpDataTimestamp(predicate, data)
		}, nil
	case OptionsKey:
		return generateCheckOption(path, predicate)
	case ClientKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckClient(predicate, data)
		}, nil
	case ClientIdKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckClientId(predicate, data)
		}, nil
	case ClientIpKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckClientIp(predicate, data)
		}, nil
	case GeoIpKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckGeoIp(predicate, data)
		}, nil
	case GeoIpCountryKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckGeoIpCountry(predicate, data)
		}, nil
	case GeoIpCountryCodeKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckGeoIpCountryCode(predicate, data)
		}, nil
	case GeoIpCityKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckGeoIpCity(predicate, data)
		}, nil
	case GeoIpLatKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckGeoIpLat(predicate, data)
		}, nil
	case GeoIpLonKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckGeoIpLon(predicate, data)
		}, nil
	case GeoIpAccuracyRadiusKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckGeoIpAccuracyRadius(predicate, data)
		}, nil
	case OsKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckOs(predicate, data)
		}, nil
	case OsNameKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckOsName(predicate, data)
		}, nil
	case OsVersionKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckOsVersion(predicate, data)
		}, nil
	case BrowserKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckBrowser(predicate, data)
		}, nil
	case BrowserNameKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckBrowserName(predicate, data)
		}, nil
	case BrowserVersionKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckBrowserVersion(predicate, data)
		}, nil
	case BasicAuthKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckBasicAuth(predicate, data)
		}, nil
	case BasicAuthUsernameKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckBasicAuthUsername(predicate, data)
		}, nil
	case BasicAuthPasswordKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckBasicAuthPassword(predicate, data)
		}, nil
	case RequestKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckRequest(predicate, data)
		}, nil
	case RequestIdKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckRequestId(predicate, data)
		}, nil
	case RequestPathKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckRequestPath(predicate, data)
		}, nil
	case RequestPathsKey:
		return generateCheckRequestPaths(path, predicate)
	case RequestQueryKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckRequestQuery(predicate, data)
		}, nil
	case RequestMethodKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckRequestMethod(predicate, data)
		}, nil
	case RequestBodyKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckRequestBody(predicate, path.ContentPath, data)
		}, nil
	case RequestGetKey:
		return generateCheckRequestGet(path, predicate)
	case RequestPostKey:
		return generateCheckRequestPost(path, predicate)
	case RequestHeadersKey:
		return generateCheckRequestHeaders(path, predicate)
	case RequestTimeKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckRequestTime(predicate, data)
		}, nil
	case RequestCookiesKey:
		return generateCheckRequestCookies(path, predicate)
	case RequestLengthKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckRequestLength(predicate, data)
		}, nil
	case ResponseKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckResponse(predicate, data)
		}, nil
	case ResponseBodyKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckResponseBody(predicate, path.ContentPath, data)
		}, nil
	case ResponseCodeKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckResponseCode(predicate, data)
		}, nil
	case ResponseSourceKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckResponseSource(predicate, data)
		}, nil
	case ResponseHeadersKey:
		return generateCheckResponseHeaders(path, predicate)
	case ResponseLengthKey:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckResponseLength(predicate, data)
		}, nil
	default:
		return nil, unknownMainPathError
	}
}

func generateMatchOption(path DataPath, patternId uint) (PredicateWithError, error) {
	switch {
	case path.ContentPath.IsEmpty():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchOptions(patternId, data)
		}, nil
	case path.ContentPath.IsSimple():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchOption(patternId, path.ContentPath.Path, data)
		}, nil
	default:
		return nil, badContentPathError
	}
}

func generateMatchRequestPaths(path DataPath, patternId uint) (PredicateWithError, error) {
	switch {
	case path.ContentPath.IsEmpty():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchRequestPaths(patternId, data)
		}, nil
	default:
		index, convertError := strconv.Atoi(path.ContentPath.Parts[0])
		if convertError != nil {
			return nil, badRequestPathIndexError
		}
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.MatchRequestPathsElement(patternId, index, path.ContentPath, data)
		}, nil
	}
}

func generateMatchRequestGet(path DataPath, patternId uint) (PredicateWithError, error) {
	switch {
	case path.ContentPath.IsEmpty():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchRequestGet(patternId, data)
		}, nil
	default:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchRequestGetValue(patternId, path.ContentPath, data)
		}, nil
	}
}

func generateMatchRequestPost(path DataPath, patternId uint) (PredicateWithError, error) {
	switch {
	case path.ContentPath.IsEmpty():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchRequestPost(patternId, data)
		}, nil
	default:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchRequestPostValue(patternId, path.ContentPath, data)
		}, nil
	}
}

func generateMatchRequestHeaders(path DataPath, patternId uint) (PredicateWithError, error) {
	switch {
	case path.ContentPath.IsEmpty():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchRequestHeaders(patternId, data)
		}, nil
	default:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchRequestHeaderValue(patternId, path.ContentPath, data)
		}, nil
	}
}

func generateMatchRequestCookies(path DataPath, patternId uint) (PredicateWithError, error) {
	switch {
	case path.ContentPath.IsEmpty():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchRequestCookies(patternId, data)
		}, nil
	default:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchRequestCookieValue(patternId, path.ContentPath, data)
		}, nil
	}
}

func generateMatchResponseHeaders(path DataPath, patternId uint) (PredicateWithError, error) {
	switch {
	case path.ContentPath.IsEmpty():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchResponseHeaders(patternId, data)
		}, nil
	default:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveMatchResponseHeaderValue(patternId, path.ContentPath, data)
		}, nil
	}
}

func generateCheckOption(path DataPath, predicate Predicate) (PredicateWithError, error) {
	switch {
	case path.ContentPath.IsEmpty():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckOptions(predicate, data)
		}, nil
	case path.ContentPath.IsSimple():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckOption(predicate, path.ContentPath.Path, data)
		}, nil
	default:
		return nil, badContentPathError
	}
}

func generateCheckRequestPaths(path DataPath, predicate Predicate) (PredicateWithError, error) {
	switch {
	case path.ContentPath.IsEmpty():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckRequestPaths(predicate, data)
		}, nil
	default:
		index, convertError := strconv.Atoi(path.ContentPath.Parts[0])
		if convertError != nil {
			return nil, badRequestPathIndexError
		}
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.CheckRequestPathsElement(predicate, index, path.ContentPath, data)
		}, nil
	}
}

func generateCheckRequestGet(path DataPath, predicate Predicate) (PredicateWithError, error) {
	switch {
	case path.ContentPath.IsEmpty():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckRequestGet(predicate, data)
		}, nil
	default:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckRequestGetValue(predicate, path.ContentPath, data)
		}, nil
	}
}

func generateCheckRequestPost(path DataPath, predicate Predicate) (PredicateWithError, error) {
	switch {
	case path.ContentPath.IsEmpty():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckRequestPost(predicate, data)
		}, nil
	default:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckRequestPostValue(predicate, path.ContentPath, data)
		}, nil
	}
}

func generateCheckRequestHeaders(path DataPath, predicate Predicate) (PredicateWithError, error) {
	switch {
	case path.ContentPath.IsEmpty():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckRequestHeaders(predicate, data)
		}, nil
	default:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckRequestHeaderValue(predicate, path.ContentPath, data)
		}, nil
	}
}

func generateCheckRequestCookies(path DataPath, predicate Predicate) (PredicateWithError, error) {
	switch {
	case path.ContentPath.IsEmpty():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckRequestCookies(predicate, data)
		}, nil
	default:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckRequestCookieValue(predicate, path.ContentPath, data)
		}, nil
	}
}

func generateCheckResponseHeaders(path DataPath, predicate Predicate) (PredicateWithError, error) {
	switch {
	case path.ContentPath.IsEmpty():
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckResponseHeaders(predicate, data)
		}, nil
	default:
		return func(data *HttpData, manager IExecutionManager) (bool, error) {
			return manager.RecursiveCheckResponseHeaderValue(predicate, path.ContentPath, data)
		}, nil
	}
}
