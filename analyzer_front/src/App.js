import React, { useState } from 'react';
import './App.css';

const initialCodeC = `int a = 0;
int b = 10;
do {
    a = 3 * b;
} 
while (x == 2);`;

const initialCodeSwift = `let playerName: String = "Kratos"
var playerLevel: Int = 1
playerLevel = 2

// Error: Intentar cambiar una constante (let)
playerName = "Fantasma de Esparta"`;


function App() {
  const [code, setCode] = useState(initialCodeC);
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);
  const [language, setLanguage] = useState('c_simple');

  const handleLanguageChange = (e) => {
    const newLang = e.target.value;
    setLanguage(newLang);
    if (newLang === 'swift') {
      setCode(initialCodeSwift);
    } else {
      setCode(initialCodeC);
    }
    setResults(null); 
  }

  const handleAnalyze = async () => {
    setLoading(true);
    setResults(null);
    try {
      const response = await fetch('http://localhost:8080/analyze', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ code, language }),
      });
      const data = await response.json();
      setResults(data);
    } catch (error) {
      console.error("Error al conectar con el backend:", error);
      alert("No se pudo conectar con el servidor de Go. ¿Está encendido?");
    } finally {
      setLoading(false);
    }
  };

  const renderTable = () => {
    if (!results || !results.lexicalTokens) return null;
    const totals = { PR: 0, ID: 0, Numeros: 0, Simbolos: 0, String: 0 };
    results.lexicalTokens.forEach(token => {
      switch (token.Type) {
        case 0: totals.PR++; break;
        case 1: totals.ID++; break;
        case 2: totals.Numeros++; break;
        case 3: totals.String++; break; 
        case 4: totals.Simbolos++; break; 
      }
    });

    return (
      <div className="results-table">
        <h3>Analizador Léxico</h3>
        <table>
          <thead>
            <tr>
              <th>Tokens</th>
              <th>PR/Tipo</th>
              <th>ID</th>
              <th>Numero</th>
              <th>String</th>
              <th>Simbolo</th>
            </tr>
          </thead>
          <tbody>
            {results.lexicalTokens.map((token, index) => (
              <tr key={index}>
                <td>{token.lexeme}</td>
                <td>{token.type === 0 ? 'x' : ''}</td>
                <td>{token.type === 1 ? 'x' : ''}</td>
                <td>{token.type === 2 ? 'x' : ''}</td>
                <td>{token.type === 3 ? 'x' : ''}</td>
                <td>{token.type === 4 ? 'x' : ''}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>Analizador Multi-Lenguaje</h1>
        <div className="container">
          <div className="editor">
            <h3>Selecciona el Lenguaje</h3>
            <select value={language} onChange={handleLanguageChange} className="language-selector">
                <option value="c_simple">Lenguaje Simple C</option>
                <option value="swift">Swift</option>
            </select>

            <h3>Código Fuente</h3>
            <textarea
              value={code}
              onChange={(e) => setCode(e.target.value)}
              rows="15"
              cols="50"
            />
            <button onClick={handleAnalyze} disabled={loading}>
              {loading ? 'Analizando...' : 'Analizar'}
            </button>
          </div>
          <div className="results">
            {renderTable()}
            {results && (
              <div className="analysis-results">
                <h3>Análisis Sintáctico y Semántico</h3>
                <p>{results.syntaxResult}</p>
                <p className="semantic-error">{results.semanticResult}</p>
              </div>
            )}
          </div>
        </div>
      </header>
    </div>
  );
}

export default App;