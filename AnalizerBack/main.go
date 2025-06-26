package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"
)

type TokenType int

const (
	KEYWORD TokenType = iota
	ID
	NUMBER
	STRING_LITERAL
	SYMBOL
	EOF
)

type Token struct {
	Type   TokenType `json:"type"`
	Lexeme string    `json:"lexeme"`
}
type AnalyzeRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}
type AnalyzeResponse struct {
	LexicalTokens  []Token `json:"lexicalTokens"`
	SyntaxResult   string  `json:"syntaxResult"`
	SemanticResult string  `json:"semanticResult"`
}

// ANALIZADOR PARA LENGUAJE C
func lexerForCSimple(input string) []Token {
	var tokens []Token
	keywords := map[string]bool{"int": true, "do": true, "while": true}
	reader := strings.NewReader(input)
	for {
		ch, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		if unicode.IsSpace(ch) {
			continue
		}
		if unicode.IsLetter(ch) {
			lexeme := string(ch)
			for {
				ch, _, err := reader.ReadRune()
				if err != nil || (!unicode.IsLetter(ch) && !unicode.IsDigit(ch)) {
					if err == nil {
						reader.UnreadRune()
					}
					break
				}
				lexeme += string(ch)
			}
			if keywords[lexeme] {
				tokens = append(tokens, Token{Type: KEYWORD, Lexeme: lexeme})
			} else {
				tokens = append(tokens, Token{Type: ID, Lexeme: lexeme})
			}
		} else if unicode.IsDigit(ch) {
			lexeme := string(ch)
			for {
				ch, _, err := reader.ReadRune()
				if err != nil || !unicode.IsDigit(ch) {
					if err == nil {
						reader.UnreadRune()
					}
					break
				}
				lexeme += string(ch)
			}
			tokens = append(tokens, Token{Type: NUMBER, Lexeme: lexeme})
		} else {
			lexeme := string(ch)
			if ch == '=' {
				if nextCh, _, _ := reader.ReadRune(); nextCh == '=' {
					lexeme += string(nextCh)
				} else {
					reader.UnreadRune()
				}
			}
			tokens = append(tokens, Token{Type: SYMBOL, Lexeme: lexeme})
		}
	}
	return tokens
}

func analyzeSemanticsForCSimple(tokens []Token) (string, string) {
	symbolTable := make(map[string]bool)
	var semanticErrors []string
	syntaxResult := "Análisis Sintáctico: La estructura del código es correcta."

	for i, token := range tokens {
		if token.Lexeme == "int" {
			if i+1 < len(tokens) && tokens[i+1].Type == ID {
				symbolTable[tokens[i+1].Lexeme] = true
			}
		}
		if token.Type == ID {
			if _, declared := symbolTable[token.Lexeme]; !declared {
				semanticErrors = append(semanticErrors, fmt.Sprintf("Error Semántico: La variable '%s' se usa pero no ha sido declarada.", token.Lexeme))
			}
		}
	}

	semanticResult := "Análisis Semántico: El código es lógicamente correcto."
	if len(semanticErrors) > 0 {
		semanticResult = strings.Join(semanticErrors, " ")
	}
	return syntaxResult, semanticResult
}

func analyzeCSimple(code string) AnalyzeResponse {
	tokens := lexerForCSimple(code)
	syntaxResult, semanticResult := analyzeSemanticsForCSimple(tokens)
	return AnalyzeResponse{LexicalTokens: tokens, SyntaxResult: syntaxResult, SemanticResult: semanticResult}
}

// ANALIZADOR PARA SWIFT
type SymbolInfo struct {
	Type       string
	IsConstant bool
}

func lexerForSwift(input string) []Token {
	var tokens []Token
	keywords := map[string]bool{"let": true, "var": true, "if": true, "else": true, "func": true, "return": true, "Int": true, "String": true, "Double": true, "Bool": true, "print": true}
	reader := strings.NewReader(input)
	for {
		ch, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		if unicode.IsSpace(ch) {
			continue
		}
		if unicode.IsLetter(ch) {
			lexeme := string(ch)
			for {
				ch, _, err := reader.ReadRune()
				if err != nil || (!unicode.IsLetter(ch) && !unicode.IsDigit(ch)) {
					if err == nil {
						reader.UnreadRune()
					}
					break
				}
				lexeme += string(ch)
			}
			if keywords[lexeme] {
				tokens = append(tokens, Token{Type: KEYWORD, Lexeme: lexeme})
			} else {
				tokens = append(tokens, Token{Type: ID, Lexeme: lexeme})
			}
			continue
		}
		if unicode.IsDigit(ch) {
			lexeme := string(ch)
			hasDot := false
			for {
				ch, _, err := reader.ReadRune()
				if err != nil || (!unicode.IsDigit(ch) && ch != '.') {
					if err == nil {
						reader.UnreadRune()
					}
					break
				}
				if ch == '.' {
					if hasDot {
						reader.UnreadRune()
						break
					}
					hasDot = true
				}
				lexeme += string(ch)
			}
			tokens = append(tokens, Token{Type: NUMBER, Lexeme: lexeme})
			continue
		}
		if ch == '"' {
			lexeme := ""
			for {
				ch, _, err := reader.ReadRune()
				if err != nil || ch == '"' {
					break
				}
				lexeme += string(ch)
			}
			tokens = append(tokens, Token{Type: STRING_LITERAL, Lexeme: lexeme})
			continue
		}
		tokens = append(tokens, Token{Type: SYMBOL, Lexeme: string(ch)})
	}
	return tokens
}

func analyzeSemanticsForSwift(tokens []Token) (string, string) {
	symbolTable := make(map[string]SymbolInfo)
	var semanticErrors []string
	syntaxResult := "Análisis Sintáctico: La estructura del código es válida para el subconjunto analizado."
	i := 0
	for i < len(tokens) {
		token := tokens[i]
		if token.Lexeme == "let" || token.Lexeme == "var" {
			isConst := token.Lexeme == "let"
			if i+3 >= len(tokens) || tokens[i+1].Type != ID || tokens[i+2].Lexeme != ":" || tokens[i+3].Type != KEYWORD {
				syntaxResult = "Error Sintáctico: Se esperaba la estructura 'let/var nombre: Tipo'."
				break
			}
			name := tokens[i+1].Lexeme
			typeName := tokens[i+3].Lexeme
			if _, exists := symbolTable[name]; exists {
				semanticErrors = append(semanticErrors, fmt.Sprintf("Error Semántico: Redefinición inválida de la variable '%s'.", name))
			} else {
				symbolTable[name] = SymbolInfo{Type: typeName, IsConstant: isConst}
			}
			i += 4
			continue
		}
		if token.Type == ID && i+1 < len(tokens) && tokens[i+1].Lexeme == "=" {
			name := token.Lexeme
			info, exists := symbolTable[name]
			if !exists {
				semanticErrors = append(semanticErrors, fmt.Sprintf("Error Semántico: La variable '%s' se usa pero no ha sido declarada.", name))
			} else if info.IsConstant {
				semanticErrors = append(semanticErrors, fmt.Sprintf("Error Semántico: No se puede asignar a '%s' porque es una constante (let).", name))
			}
			i += 2
			continue
		}
		i++
	}
	semanticResult := "Análisis Semántico: El código es lógicamente correcto."
	if len(semanticErrors) > 0 {
		semanticResult = strings.Join(semanticErrors, " ")
	}
	return syntaxResult, semanticResult
}

func analyzeSwift(code string) AnalyzeResponse {
	tokens := lexerForSwift(code)
	syntaxResult, semanticResult := analyzeSemanticsForSwift(tokens)
	return AnalyzeResponse{LexicalTokens: tokens, SyntaxResult: syntaxResult, SemanticResult: semanticResult}
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp AnalyzeResponse

	// ROUTER: Decide qué analizador usar
	switch req.Language {
	case "swift":
		resp = analyzeSwift(req.Code)
	case "c_simple":
		resp = analyzeCSimple(req.Code)
	default:
		// Por defecto, usa el de C-Simple si no se especifica.
		resp = analyzeCSimple(req.Code)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/analyze", analyzeHandler)
	fmt.Println("Servidor de análisis iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
