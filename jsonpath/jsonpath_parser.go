package jsonpath

import (
	"errors"
	"strings"
	"unicode"
)

func parsePathString(source string) ([]string, error) {
	state := initState()
	var err error
	for _, ch := range []rune(source) { //nolint
		state, err = state.handler(ch, state)
		if err != nil {
			return nil, err
		}
	}
	state, err = state.handler(eofRune, state)
	if err != nil {
		return nil, err
	}
	return state.path, nil
}

const eofRune rune = '\u0000'

const (
	parserNoneState            = 0
	parserNeutralState         = 100
	parserDotState             = 200
	parserBracketState         = 300
	parserSingleQuoteNameState = 400
	parserDoubleQuoteNameState = 500
)

type stateHandler func(ch rune, state *parseState) (*parseState, error)

type parseState struct {
	id      int
	handler stateHandler
	path    []string
	part    *strings.Builder
}

func initState() *parseState {
	return &parseState{id: parserNoneState, handler: processNoneState, path: make([]string, 0), part: &strings.Builder{}}
}

func (state *parseState) updateState(id int, handler stateHandler) *parseState {
	state.id = id
	state.handler = handler
	return state
}

func (state *parseState) appendPart() *parseState {
	state.path = append(state.path, state.part.String())
	state.part.Reset()
	return state
}

func processNoneState(char rune, state *parseState) (*parseState, error) {
	switch {
	case char == '$':
		return state.updateState(parserNeutralState, processNeutralState), nil
	default:
		return nil, errors.New("bad JSON path")
	}
}

func processNeutralState(char rune, state *parseState) (*parseState, error) {
	switch {
	case char == '.':
		return state.updateState(parserDotState, processDotState), nil
	case char == '[':
		state.part.WriteRune(char)
		return state.updateState(parserBracketState, processBracketState), nil
	case char == eofRune:
		return state, nil
	}
	return nil, errors.New("bad JSON path")
}

func processDotState(char rune, state *parseState) (*parseState, error) {
	switch {
	case unicode.IsSpace(char):
		return nil, errors.New("bad JSON path")
	case char == '.' && state.part.Len() == 0:
		return nil, errors.New("bad JSON path")
	case char == '.':
		return state.appendPart(), nil
	case char == '[' && state.part.Len() == 0:
		state.part.WriteRune(char)
		return state.updateState(parserBracketState, processBracketState), nil
	case char == '[':
		state.appendPart()
		state.part.WriteRune(char)
		return state.updateState(parserBracketState, processBracketState), nil
	case char == ']':
		return nil, errors.New("bad JSON path")
	case char == eofRune && state.part.Len() == 0:
		return nil, errors.New("bad JSON path")
	case char == eofRune:
		return state.appendPart(), nil
	default:
		// TODO (std_string) : probably add additional check on allowed symbols
		state.part.WriteRune(char)
		return state, nil
	}
}

func processBracketState(char rune, state *parseState) (*parseState, error) {
	switch {
	case unicode.IsSpace(char) && state.part.Len() == 1:
		return state, nil
	case char == ']' && state.part.Len() == 1:
		return nil, errors.New("bad JSON path")
	case char == ']':
		state.part.WriteRune(char)
		return state.appendPart().updateState(parserNeutralState, processNeutralState), nil
	case char == '\'' && state.part.Len() > 1:
		return nil, errors.New("bad JSON path")
	case char == '\'':
		state.part.Reset()
		return state.updateState(parserSingleQuoteNameState, processSingleQuoteNameState), nil
	case char == '"' && state.part.Len() > 1:
		return nil, errors.New("bad JSON path")
	case char == '"':
		state.part.Reset()
		return state.updateState(parserDoubleQuoteNameState, processDoubleQuoteNameState), nil
	case char == '[':
		return nil, errors.New("bad JSON path")
	case char == '.':
		return nil, errors.New("bad JSON path")
	default:
		// TODO (std_string) : probably add additional check on allowed symbols
		state.part.WriteRune(char)
		return state, nil
	}
}

func processSingleQuoteNameState(char rune, state *parseState) (*parseState, error) {
	const parserSingleQuoteNameEndState = parserSingleQuoteNameState + 1
	switch {
	case state.id == parserSingleQuoteNameState && char == '\'' && state.part.Len() == 0:
		return nil, errors.New("bad JSON path")
	case state.id == parserSingleQuoteNameState && char == '\'':
		return state.updateState(parserSingleQuoteNameEndState, processSingleQuoteNameState), nil
	case state.id == parserSingleQuoteNameEndState && unicode.IsSpace(char):
		return state, nil
	case state.id == parserSingleQuoteNameEndState && char == ']':
		return state.appendPart().updateState(parserNeutralState, processNeutralState), nil
	case state.id == parserSingleQuoteNameEndState && char != ']':
		return nil, errors.New("bad JSON path")
	default:
		// TODO (std_string) : probably add additional check on allowed symbols
		state.part.WriteRune(char)
		return state, nil
	}
}

func processDoubleQuoteNameState(char rune, state *parseState) (*parseState, error) {
	const parserDoubleQuoteNameEndState = parserDoubleQuoteNameState + 1
	switch {
	case state.id == parserDoubleQuoteNameState && char == '\'' && state.part.Len() == 0:
		return nil, errors.New("bad JSON path")
	case state.id == parserDoubleQuoteNameState && char == '"':
		return state.updateState(parserDoubleQuoteNameEndState, processDoubleQuoteNameState), nil
	case state.id == parserDoubleQuoteNameEndState && unicode.IsSpace(char):
		return state, nil
	case state.id == parserDoubleQuoteNameEndState && char == ']':
		return state.appendPart().updateState(parserNeutralState, processNeutralState), nil
	case state.id == parserDoubleQuoteNameEndState && char != ']':
		return nil, errors.New("bad JSON path")
	default:
		// TODO (std_string) : probably add additional check on allowed symbols
		state.part.WriteRune(char)
		return state, nil
	}
}
