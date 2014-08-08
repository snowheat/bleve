//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package sorani_normalize

import (
	"bytes"
	"unicode"
	"unicode/utf8"

	"github.com/couchbaselabs/bleve/analysis"
)

const (
	YEH         = '\u064A'
	DOTLESS_YEH = '\u0649'
	FARSI_YEH   = '\u06CC'

	KAF   = '\u0643'
	KEHEH = '\u06A9'

	HEH             = '\u0647'
	AE              = '\u06D5'
	ZWNJ            = '\u200C'
	HEH_DOACHASHMEE = '\u06BE'
	TEH_MARBUTA     = '\u0629'

	REH        = '\u0631'
	RREH       = '\u0695'
	RREH_ABOVE = '\u0692'

	TATWEEL  = '\u0640'
	FATHATAN = '\u064B'
	DAMMATAN = '\u064C'
	KASRATAN = '\u064D'
	FATHA    = '\u064E'
	DAMMA    = '\u064F'
	KASRA    = '\u0650'
	SHADDA   = '\u0651'
	SUKUN    = '\u0652'
)

type SoraniNormalizeFilter struct {
}

func NewSoraniNormalizeFilter() *SoraniNormalizeFilter {
	return &SoraniNormalizeFilter{}
}

func (s *SoraniNormalizeFilter) Filter(input analysis.TokenStream) analysis.TokenStream {
	rv := make(analysis.TokenStream, 0)

	for _, token := range input {
		term := normalize(token.Term)
		token.Term = term
		rv = append(rv, token)
	}

	return rv
}

func normalize(input []byte) []byte {
	runes := bytes.Runes(input)
	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		case YEH, DOTLESS_YEH:
			runes[i] = FARSI_YEH
		case KAF:
			runes[i] = KEHEH
		case ZWNJ:
			if i > 0 && runes[i-1] == HEH {
				runes[i-1] = AE
			}
			runes = deleteRune(runes, i)
			i--
		case HEH:
			if i == len(runes)-1 {
				runes[i] = AE
			}
		case TEH_MARBUTA:
			runes[i] = AE
		case HEH_DOACHASHMEE:
			runes[i] = HEH
		case REH:
			if i == 0 {
				runes[i] = RREH
			}
		case RREH_ABOVE:
			runes[i] = RREH
		case TATWEEL, KASRATAN, DAMMATAN, FATHATAN, FATHA, DAMMA, KASRA, SHADDA, SUKUN:
			runes = deleteRune(runes, i)
			i--
		default:
			if unicode.In(runes[i], unicode.Cf) {
				runes = deleteRune(runes, i)
				i--
			}
		}
	}
	return buildTermFromRunes(runes)
}

func deleteRune(in []rune, pos int) []rune {
	if pos >= len(in) {
		return in
	}
	copy(in[pos:], in[pos+1:])
	return in[:len(in)-1]
}

func buildTermFromRunes(runes []rune) []byte {
	rv := make([]byte, 0, len(runes)*4)
	for _, r := range runes {
		runeBytes := make([]byte, utf8.RuneLen(r))
		utf8.EncodeRune(runeBytes, r)
		rv = append(rv, runeBytes...)
	}
	return rv
}
