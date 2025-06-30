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
	KEYWORD        TokenType = 0
	ID             TokenType = 1
	NUMBER         TokenType = 2
	STRING_LITERAL TokenType = 3
	SYMBOL         TokenType = 4
	ERROR          TokenType = 5
)

func (tt TokenType) String() string {
	switch tt {
	case KEYWORD:
		return "KEYWORD"
	case ID:
		return "ID"
	case NUMBER:
		return "NUMBER"
	case STRING_LITERAL:
		return "STRING_LITERAL"
	case SYMBOL:
		return "SYMBOL"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

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
	Error          string  `json:"error,omitempty"`
}
type SymbolInfo struct {
	Type       string
	IsConstant bool
}

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
				nextCh, _, err := reader.ReadRune()
				if err != nil || (!unicode.IsLetter(nextCh) && !unicode.IsDigit(nextCh)) {
					if err == nil {
						reader.UnreadRune()
					}
					break
				}
				lexeme += string(nextCh)
			}
			if keywords[lexeme] {
				tokens = append(tokens, Token{Type: KEYWORD, Lexeme: lexeme})
			} else {
				tokens = append(tokens, Token{Type: ID, Lexeme: lexeme})
			}
		} else if unicode.IsDigit(ch) {
			lexeme := string(ch)
			for {
				nextCh, _, err := reader.ReadRune()
				if err != nil || !unicode.IsDigit(nextCh) {
					if err == nil {
						reader.UnreadRune()
					}
					break
				}
				lexeme += string(nextCh)
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
	uniqueErrors := make(map[string]bool)
	for i, token := range tokens {
		if token.Lexeme == "int" {
			if i+1 < len(tokens) && tokens[i+1].Type == ID {
				symbolTable[tokens[i+1].Lexeme] = true
			}
		}
		if token.Type == ID {
			if _, declared := symbolTable[token.Lexeme]; !declared {
				errorMsg := fmt.Sprintf("Error Semántico: La variable '%s' se usa pero no ha sido declarada.", token.Lexeme)
				if !uniqueErrors[errorMsg] {
					semanticErrors = append(semanticErrors, errorMsg)
					uniqueErrors[errorMsg] = true
				}
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
		if ch == '"' {
			lexeme := ""
			for {
				ch, _, err := reader.ReadRune()
				if err != nil || ch == '"' {
					break
				}
				lexeme += string(ch)
			}
			tokens = append(tokens, Token{Type: STRING_LITERAL, Lexeme: `"` + lexeme + `"`})
			continue
		}
		if unicode.IsLetter(ch) {
			lexeme := string(ch)
			for {
				nextCh, _, err := reader.ReadRune()
				if err != nil || (!unicode.IsLetter(nextCh) && !unicode.IsDigit(nextCh)) {
					if err == nil {
						reader.UnreadRune()
					}
					break
				}
				lexeme += string(nextCh)
			}
			if keywords[lexeme] {
				tokens = append(tokens, Token{Type: KEYWORD, Lexeme: lexeme})
			} else {
				tokens = append(tokens, Token{Type: ID, Lexeme: lexeme})
			}
			continue
		}
		if unicode.IsDigit(ch) || (ch == '-' && len(tokens) > 0 && tokens[len(tokens)-1].Type != NUMBER && tokens[len(tokens)-1].Type != ID) {
			lexeme := string(ch)
			for {
				nextCh, _, err := reader.ReadRune()
				if err != nil || (!unicode.IsDigit(nextCh) && nextCh != '.') {
					if err == nil {
						reader.UnreadRune()
					}
					break
				}
				lexeme += string(nextCh)
			}
			tokens = append(tokens, Token{Type: NUMBER, Lexeme: lexeme})
			continue
		}
		if strings.ContainsRune("=+-*/:(){}", ch) {
			tokens = append(tokens, Token{Type: SYMBOL, Lexeme: string(ch)})
			continue
		}
		tokens = append(tokens, Token{Type: ERROR, Lexeme: string(ch)})
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

func lexerForJava(code string) []Token {
	var tokens []Token
	keywords := map[string]bool{"public": true, "class": true, "static": true, "void": true, "main": true, "int": true, "String": true, "if": true, "System": true, "out": true, "println": true, "equals": true}
	reader := strings.NewReader(code)
	for {
		ch, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		if unicode.IsSpace(ch) {
			continue
		}
		if ch == '"' {
			lexeme := ""
			isClosed := false
			for {
				charInString, _, err := reader.ReadRune()
				if err != nil || charInString == '\n' {
					break
				}
				if charInString == '"' {
					isClosed = true
					break
				}
				lexeme += string(charInString)
			}
			if isClosed {
				tokens = append(tokens, Token{Type: STRING_LITERAL, Lexeme: `"` + lexeme + `"`})
			} else {
				tokens = append(tokens, Token{Type: ERROR, Lexeme: `String sin cerrar: "` + lexeme})
			}
			continue
		}
		if unicode.IsLetter(ch) {
			lexeme := string(ch)
			for {
				nextCh, _, err := reader.ReadRune()
				if err != nil || (!unicode.IsLetter(nextCh) && !unicode.IsDigit(nextCh) && nextCh != '_') {
					if err == nil {
						reader.UnreadRune()
					}
					break
				}
				lexeme += string(nextCh)
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
			for {
				nextCh, _, err := reader.ReadRune()
				if err != nil || !unicode.IsDigit(nextCh) {
					if err == nil {
						reader.UnreadRune()
					}
					break
				}
				lexeme += string(nextCh)
			}
			tokens = append(tokens, Token{Type: NUMBER, Lexeme: lexeme})
			continue
		}
		if strings.ContainsRune("(){}[]=;>.", ch) {
			tokens = append(tokens, Token{Type: SYMBOL, Lexeme: string(ch)})
			continue
		}
		tokens = append(tokens, Token{Type: ERROR, Lexeme: fmt.Sprintf("Caracter inesperado: %c", ch)})
	}
	return tokens
}
func analyzeJava(code string) AnalyzeResponse {
	tokens := lexerForJava(code)
	var syntaxErrors []string
	var semanticErrors []string
	symbolTable := make(map[string]SymbolInfo)
	for _, token := range tokens {
		if token.Type == ERROR {
			syntaxErrors = append(syntaxErrors, fmt.Sprintf("Error Léxico: %s", token.Lexeme))
		}
	}
	if len(syntaxErrors) > 0 {
		return AnalyzeResponse{LexicalTokens: tokens, SyntaxResult: strings.Join(syntaxErrors, "\n"), SemanticResult: "El análisis se detuvo debido a errores léxicos/sintácticos."}
	}
	for i := 0; i < len(tokens); {
		token := tokens[i]
		isStartOfStatement := (i == 0) || (tokens[i-1].Lexeme == "{" || tokens[i-1].Lexeme == ";")
		if isStartOfStatement && (token.Lexeme == "int" || token.Lexeme == "String") {
			if i+4 >= len(tokens) || tokens[i+1].Type != ID || tokens[i+2].Lexeme != "=" || (tokens[i+3].Type != NUMBER && tokens[i+3].Type != STRING_LITERAL && tokens[i+3].Type != ID) || tokens[i+4].Lexeme != ";" {
				syntaxErrors = append(syntaxErrors, fmt.Sprintf("Error de sintaxis en la declaración cerca de '%s'.", tokens[i+1].Lexeme))
				i++
				continue
			}
			varType := token.Lexeme
			varName := tokens[i+1].Lexeme
			valueToken := tokens[i+3]
			if varType == "int" && valueToken.Type != NUMBER {
				semanticErrors = append(semanticErrors, fmt.Sprintf("Error Semántico: No se puede asignar un valor de tipo '%s' a una variable de tipo 'int' (en la declaración de '%s').", valueToken.Type.String(), varName))
			} else if varType == "String" && valueToken.Type != STRING_LITERAL {
				semanticErrors = append(semanticErrors, fmt.Sprintf("Error Semántico: No se puede asignar un valor de tipo '%s' a una variable de tipo 'String' (en la declaración de '%s').", valueToken.Type.String(), varName))
			} else {
				if _, exists := symbolTable[varName]; exists {
					semanticErrors = append(semanticErrors, fmt.Sprintf("Error Semántico: La variable '%s' ya ha sido declarada.", varName))
				} else {
					symbolTable[varName] = SymbolInfo{Type: varType}
				}
			}
			i += 5
			continue
		}
		if token.Lexeme == "if" {
			if i+1 >= len(tokens) || tokens[i+1].Lexeme != "(" {
				syntaxErrors = append(syntaxErrors, "Error de sintaxis: se esperaba '(' después de 'if'.")
				i++
				continue
			}
			openParenIndex := i + 1
			closeParenIndex := -1
			for j := openParenIndex + 1; j < len(tokens); j++ {
				if tokens[j].Lexeme == ")" {
					closeParenIndex = j
					break
				}
			}
			if closeParenIndex == -1 {
				syntaxErrors = append(syntaxErrors, "Error de sintaxis: falta ')' para cerrar la condición del 'if'.")
			} else {
				if closeParenIndex == openParenIndex+4 && tokens[openParenIndex+1].Type == ID && tokens[openParenIndex+2].Lexeme == ">" {
					varName := tokens[openParenIndex+1].Lexeme
					if symbol, exists := symbolTable[varName]; !exists {
						semanticErrors = append(semanticErrors, fmt.Sprintf("Error Semántico: La variable '%s' no ha sido declarada.", varName))
					} else if symbol.Type != "int" {
						semanticErrors = append(semanticErrors, fmt.Sprintf("Error Semántico: El operador '>' solo aplica a 'int', no a '%s'.", symbol.Type))
					}
				}
			}
		}
		i++
	}
	syntaxResult := "El código es sintácticamente correcto."
	if len(syntaxErrors) > 0 {
		syntaxResult = strings.Join(syntaxErrors, "\n")
	}
	semanticResult := "El código es lógicamente correcto."
	if len(semanticErrors) > 0 {
		semanticResult = strings.Join(semanticErrors, "\n")
	}
	return AnalyzeResponse{LexicalTokens: tokens, SyntaxResult: syntaxResult, SemanticResult: semanticResult}
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decodificando", http.StatusBadRequest)
		return
	}
	var resp AnalyzeResponse
	switch req.Language {
	case "java":
		resp = analyzeJava(req.Code)
	case "swift":
		resp = analyzeSwift(req.Code)
	case "c_simple":
		resp = analyzeCSimple(req.Code)
	default:
		resp = analyzeCSimple(req.Code)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func main() {
	http.HandleFunc("/analyze", analyzeHandler)
	fmt.Println("Servidor de análisis MULTI-LENGUAJE v4 (FINAL) iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
