import React, { useState, useCallback } from 'react';
import './App.css';

const initialCodeC = `int a = 0;
int b = 10;
do {
    a = 3 * b;
} 
while (x == 2);`;

const initialCodeSwift = `let playerName: String = "Kratos"
var playerLevel: Int = 1
// Error: Intentar cambiar una constante
playerName = "Fantasma de Esparta"`;

const initialCodeJava = `public class ejercicio {
  public static void main(String[] args) {
    int edad = 22;
    String escuela = "upchiapas";
    if (edad > 18) {
      System.out.println("Mayor de edad");
    }
  }
}`;

function App() {
  const [code, setCode] = useState(initialCodeJava);
  const [analysisResult, setAnalysisResult] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [language, setLanguage] = useState('java');

  const handleLanguageChange = (e) => {
    const newLang = e.target.value;
    setLanguage(newLang);
    if (newLang === 'swift') { setCode(initialCodeSwift); } 
    else if (newLang === 'c') { setCode(initialCodeC); }
    else { setCode(initialCodeJava); }
    setAnalysisResult(null); 
  };

  const handleAnalyze = useCallback(async () => {
    setIsLoading(true);
    setAnalysisResult(null);
    try {
      const response = await fetch('http://localhost:8080/analyze', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code, language }),
      });
      const data = await response.json();
      setAnalysisResult(data);
    } catch (error) {
      setAnalysisResult({ error: 'Error de conexión con el backend.' });
    } finally {
      setIsLoading(false);
    }
  }, [code, language]);

  const hasSyntaxError = analysisResult?.syntaxResult?.toLowerCase().startsWith('error');
  const hasSemanticError = analysisResult?.semanticResult?.toLowerCase().startsWith('error');

  return (
    <div className="App">
      <div className="ide-container">
        <h1>Analizador</h1>
        <div className="editor-container">
          <select value={language} onChange={handleLanguageChange} className="language-selector">
              <option value="c">C</option>
              <option value="swift">Swift</option>
              <option value="java">Java</option>
          </select>
          <textarea
            value={code}
            onChange={(e) => setCode(e.target.value)}
            className="code-editor"
            spellCheck="false"
          />
          <button onClick={handleAnalyze} disabled={isLoading} className="analyze-button">
            {isLoading ? 'Analizando...' : 'Analizar Código'}
          </button>
        </div>
      </div>

      {analysisResult && (
        <div className="results-container">
          {analysisResult.error ? (
            <div className="analysis-card error-card">
              <h3>Error de Conexión</h3>
              <p>{analysisResult.error}</p>
            </div>
          ) : (
            <>
              <div className="analysis-card">
                <h3>Análisis Léxico</h3>
                <table className="token-table detailed-table">
                  <thead>
                    <tr>
                      <th>Tokens</th>
                      <th>PR/Tipo</th>
                      <th>ID</th>
                      <th>Numero</th>
                      <th>String</th>
                      <th>Simbolo</th>
                      <th>Error</th>
                    </tr>
                  </thead>
                  <tbody>
                    {analysisResult.lexicalTokens.map((token, index) => (
                      <tr key={index}>
                        <td>{token.lexeme}</td>
                        <td>{token.type === 0 ? 'x' : ''}</td>
                        <td>{token.type === 1 ? 'x' : ''}</td>
                        <td>{token.type === 2 ? 'x' : ''}</td>
                        <td>{token.type === 3 ? 'x' : ''}</td>
                        <td>{token.type === 4 ? 'x' : ''}</td>
                        <td>{token.type === 5 ? 'x' : ''}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              <div className={`analysis-card ${hasSyntaxError ? 'has-error' : ''}`}>
                <h3 className={hasSyntaxError ? 'title-error' : ''}>Análisis Sintáctico</h3>
                <p>{analysisResult.syntaxResult}</p>
              </div>
              
              <div className={`analysis-card ${hasSemanticError ? 'has-error' : ''}`}>
                <h3 className={hasSemanticError ? 'title-error' : ''}>Análisis Semántico</h3>
                <p>{analysisResult.semanticResult}</p>
              </div>
            </>
          )}
        </div>
      )}
    </div>
  );
}

export default App;