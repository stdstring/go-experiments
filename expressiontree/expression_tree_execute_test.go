package expressiontree

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestParseExecuteExpression(t *testing.T) {
	someError := errors.New("some error")
	storage := &parseStorage{
		knownPath: []DataPath{
			CreateDataPathWithMainOnly(HttpDataKey),
			CreateDataPathWithSimpleContent(OptionsKey, "IDDQD"),
			CreateDataPathWithSimpleContent(RequestHeadersKey, "IDKFA"),
		},
		checkArguments: []any{
			"",
			"IDCLIP",
		},
	}
	httpData := &HttpData{}
	testCases := []struct {
		name           string
		source         string
		expectedCalls  func(*MockIExecutionManager) []*gomock.Call
		expectedResult bool
		expectedError  error
	}{
		{
			name:   "EXISTS(http.options.IDDQD)->true",
			source: "EXISTS(1)",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(true, nil),
				}
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:   "EXISTS(http.options.IDDQD)->false",
			source: "EXISTS(1)",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(false, nil),
				}
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:   "EXISTS(http.options.IDDQD)->SomeError",
			source: "EXISTS(1)",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(false, someError),
				}
			},
			expectedResult: false,
			expectedError:  someError,
		},
		{
			name:   "MATCH(http.options.IDDQD,666)->true",
			source: "MATCH(1,666)",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(true, nil),
				}
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:   "MATCH(http.options.IDDQD,666)->false",
			source: "MATCH(1,666)",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(false, nil),
				}
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:   "MATCH(http.options.IDDQD,666)->SomeError",
			source: "MATCH(1,666)",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(false, someError),
				}
			},
			expectedResult: false,
			expectedError:  someError,
		},
		{
			name:   `CHECK(http.options.IDDQD=="IDCLIP")->true`,
			source: "CHECK(1,0,1)",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOption(newPredicateMatcher(), "IDDQD", httpData).Return(true, nil),
				}
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:   `CHECK(http.options.IDDQD=="IDCLIP")->false`,
			source: "CHECK(1,0,1)",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOption(newPredicateMatcher(), "IDDQD", httpData).Return(false, nil),
				}
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:   `CHECK(http.options.IDDQD=="IDCLIP")->SomeError`,
			source: "CHECK(1,0,1)",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOption(newPredicateMatcher(), "IDDQD", httpData).Return(false, someError),
				}
			},
			expectedResult: false,
			expectedError:  someError,
		},
		{
			name:   "NOT(EXISTS(http.options.IDDQD)->true)",
			source: "NOT(EXISTS(1))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(true, nil),
				}
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:   "NOT(EXISTS(http.options.IDDQD)->false)",
			source: "NOT(EXISTS(1))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(false, nil),
				}
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:   "NOT(EXISTS(http.options.IDDQD)->SomeError)",
			source: "NOT(EXISTS(1))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(false, someError),
				}
			},
			expectedResult: false,
			expectedError:  someError,
		},
		{
			name:   "AND(EXISTS(http.options.IDDQD)->true,MATCH(http.options.IDDQD,666)->true)",
			source: "AND(EXISTS(1),MATCH(1,666))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(true, nil),
					mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(true, nil),
				}
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:   "AND(EXISTS(http.options.IDDQD)->true,MATCH(http.options.IDDQD,666)->false)",
			source: "AND(EXISTS(1),MATCH(1,666))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(true, nil),
					mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(false, nil),
				}
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:   "AND(EXISTS(http.options.IDDQD)->false,MATCH(http.options.IDDQD,666)->true)",
			source: "AND(EXISTS(1),MATCH(1,666))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(false, nil),
					//mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(true, nil),
				}
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:   "AND(EXISTS(http.options.IDDQD)->false,MATCH(http.options.IDDQD,666)->false)",
			source: "AND(EXISTS(1),MATCH(1,666))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(false, nil),
					//mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(false, nil),
				}
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:   "AND(EXISTS(http.options.IDDQD)->true,MATCH(http.options.IDDQD,666)->SomeError)",
			source: "AND(EXISTS(1),MATCH(1,666))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(true, nil),
					mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(false, someError),
				}
			},
			expectedResult: false,
			expectedError:  someError,
		},
		{
			name:   "AND(EXISTS(http.options.IDDQD)->SomeError,MATCH(http.options.IDDQD,666)->true)",
			source: "AND(EXISTS(1),MATCH(1,666))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(false, someError),
					//mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(true, nil),
				}
			},
			expectedResult: false,
			expectedError:  someError,
		},
		{
			name:   "AND(EXISTS(http.options.IDDQD)->false,MATCH(http.options.IDDQD,666)->SomeError)",
			source: "AND(EXISTS(1),MATCH(1,666))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(false, nil),
					//mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(false, someError),
				}
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:   "OR(EXISTS(http.options.IDDQD)->true,MATCH(http.options.IDDQD,666)->true)",
			source: "OR(EXISTS(1),MATCH(1,666))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(true, nil),
					//mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(true, nil),
				}
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:   "OR(EXISTS(http.options.IDDQD)->true,MATCH(http.options.IDDQD,666)->false)",
			source: "OR(EXISTS(1),MATCH(1,666))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(true, nil),
					//mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(false, nil),
				}
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:   "OR(EXISTS(http.options.IDDQD)->false,MATCH(http.options.IDDQD,666)->true)",
			source: "OR(EXISTS(1),MATCH(1,666))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(false, nil),
					mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(true, nil),
				}
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:   "OR(EXISTS(http.options.IDDQD)->false,MATCH(http.options.IDDQD,666)->false)",
			source: "OR(EXISTS(1),MATCH(1,666))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(false, nil),
					mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(false, nil),
				}
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:   "OR(EXISTS(http.options.IDDQD)->SomeError,MATCH(http.options.IDDQD,666)->true)",
			source: "OR(EXISTS(1),MATCH(1,666))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(false, someError),
					//mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(true, nil),
				}
			},
			expectedResult: false,
			expectedError:  someError,
		},
		{
			name:   "OR(EXISTS(http.options.IDDQD)->true,MATCH(http.options.IDDQD,666)->SomeError)",
			source: "OR(EXISTS(1),MATCH(1,666))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(true, nil),
					//mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(false, someError),
				}
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:   "OR(EXISTS(http.options.IDDQD)->false,MATCH(http.options.IDDQD,666)->SomeError)",
			source: "OR(EXISTS(1),MATCH(1,666))",
			expectedCalls: func(mock *MockIExecutionManager) []*gomock.Call {
				return []*gomock.Call{
					mock.EXPECT().CheckOptionExistence("IDDQD", httpData).Return(false, nil),
					mock.EXPECT().MatchOption(uint(666), "IDDQD", httpData).Return(false, someError),
				}
			},
			expectedResult: false,
			expectedError:  someError,
		},
	}
	for _, testCase := range testCases {
		currentTestCase := testCase
		t.Run(currentTestCase.name, func(t *testing.T) {
			source := testCase.source
			expression, expressionError := parseExpressionTree(source, storage)
			assert.NoError(t, expressionError)
			assert.NotNil(t, expression)
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			executionManager := NewMockIExecutionManager(mockController)
			gomock.InOrder(currentTestCase.expectedCalls(executionManager)...)
			actualResult, actualError := expression(httpData, executionManager)
			assert.Equal(t, currentTestCase.expectedResult, actualResult)
			assert.Equal(t, currentTestCase.expectedError, actualError)
		})
	}
}

type predicateMatcher struct {
}

func (m *predicateMatcher) Matches(x interface{}) bool {
	return x != nil
}

func (m *predicateMatcher) String() string {
	return "matcher for Predicate - for func(value any) bool"
}

func newPredicateMatcher() *predicateMatcher {
	return &predicateMatcher{}
}
