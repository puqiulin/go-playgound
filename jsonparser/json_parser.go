package jsonparser

import (
	"errors"
	"strconv"
	"unicode"
)

type JSON interface{}

type Parser struct {
	input string
	pos   int
}

func NewParser(input string) *Parser {
	return &Parser{input: input, pos: 0}
}

func (p *Parser) Parse() (JSON, error) {
	p.consumeWhitespace()
	if p.pos >= len(p.input) {
		return nil, errors.New("unexpected end of input")
	}

	switch p.input[p.pos] {
	case '{':
		return p.parseObject()
	case '[':
		return p.parseArray()
	case '"':
		return p.parseString()
	case 't', 'f':
		return p.parseBoolean()
	case 'n':
		return p.parseNull()
	default:
		if unicode.IsDigit(rune(p.input[p.pos])) || p.input[p.pos] == '-' {
			return p.parseNumber()
		}
		return nil, errors.New("unexpected character")
	}
}

func (p *Parser) parseObject() (map[string]JSON, error) {
	obj := make(map[string]JSON)
	p.pos++ // Skip opening brace

	for {
		p.consumeWhitespace()
		if p.pos >= len(p.input) {
			return nil, errors.New("unexpected end of input")
		}

		if p.input[p.pos] == '}' {
			p.pos++
			return obj, nil
		}

		key, err := p.parseString()
		if err != nil {
			return nil, err
		}

		p.consumeWhitespace()
		if p.pos >= len(p.input) || p.input[p.pos] != ':' {
			return nil, errors.New("expected colon")
		}
		p.pos++

		value, err := p.Parse()
		if err != nil {
			return nil, err
		}

		obj[key] = value

		p.consumeWhitespace()
		if p.pos >= len(p.input) {
			return nil, errors.New("unexpected end of input")
		}

		if p.input[p.pos] == ',' {
			p.pos++
		} else if p.input[p.pos] != '}' {
			return nil, errors.New("expected comma or closing brace")
		}
	}
}

func (p *Parser) parseArray() ([]JSON, error) {
	arr := make([]JSON, 0)
	p.pos++ // Skip opening bracket

	for {
		p.consumeWhitespace()
		if p.pos >= len(p.input) {
			return nil, errors.New("unexpected end of input")
		}

		if p.input[p.pos] == ']' {
			p.pos++
			return arr, nil
		}

		value, err := p.Parse()
		if err != nil {
			return nil, err
		}

		arr = append(arr, value)

		p.consumeWhitespace()
		if p.pos >= len(p.input) {
			return nil, errors.New("unexpected end of input")
		}

		if p.input[p.pos] == ',' {
			p.pos++
		} else if p.input[p.pos] != ']' {
			return nil, errors.New("expected comma or closing bracket")
		}
	}
}

func (p *Parser) parseString() (string, error) {
	p.pos++ // Skip opening quote
	start := p.pos

	for p.pos < len(p.input) && p.input[p.pos] != '"' {
		if p.input[p.pos] == '\\' {
			p.pos++
		}
		p.pos++
	}

	if p.pos >= len(p.input) {
		return "", errors.New("unterminated string")
	}

	result := p.input[start:p.pos]
	p.pos++ // Skip closing quote
	return result, nil
}

func (p *Parser) parseNumber() (float64, error) {
	start := p.pos

	for p.pos < len(p.input) && (unicode.IsDigit(rune(p.input[p.pos])) || p.input[p.pos] == '.' || p.input[p.pos] == 'e' || p.input[p.pos] == 'E' || p.input[p.pos] == '+' || p.input[p.pos] == '-') {
		p.pos++
	}

	numStr := p.input[start:p.pos]
	return strconv.ParseFloat(numStr, 64)
}

func (p *Parser) parseBoolean() (bool, error) {
	if len(p.input)-p.pos >= 4 && p.input[p.pos:p.pos+4] == "true" {
		p.pos += 4
		return true, nil
	}
	if len(p.input)-p.pos >= 5 && p.input[p.pos:p.pos+5] == "false" {
		p.pos += 5
		return false, nil
	}
	return false, errors.New("invalid boolean")
}

func (p *Parser) parseNull() (interface{}, error) {
	if len(p.input)-p.pos >= 4 && p.input[p.pos:p.pos+4] == "null" {
		p.pos += 4
		return nil, nil
	}
	return nil, errors.New("invalid null")
}

func (p *Parser) consumeWhitespace() {
	for p.pos < len(p.input) && unicode.IsSpace(rune(p.input[p.pos])) {
		p.pos++
	}
}
