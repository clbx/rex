import "./app.css";
import { Routes, Route } from "react-router-dom";
import AllGames from "./components/AllCountries/AllCountries";
import GameInfo from "./components/CountryInfo/GameInfo";

function App() {
  return (
    <>
      <div className="header">
        <div className="container">
          <h5>Rex Game Library</h5>
        </div>
      </div>
      <div className="container">
        <Routes>
          <Route path="/" element={<AllGames />} />
          <Route path="/games/:gameId" element={<GameInfo />} />
        </Routes>
      </div>
    </>
  );
}

export default App;