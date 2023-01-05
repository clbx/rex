import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { apiURL } from "../util/api";
import { Link } from "react-router-dom";
import SetId from "../SetId/SetId"



const GameInfo = () => {
  const [game, setGame] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  const { gameId } = useParams();


  useEffect(() => {
    const getGameById = async () => {
      try {
        const res = await fetch(`${apiURL}/v1/games/byId?id=${gameId}`);
        if (!res.ok) throw new Error("Invalid Game ID");

        const data = await res.json();

        setGame(data);
        setIsLoading(false);
      } catch (error) {
        setIsLoading(false);
        setError(error.message);
      }
    };

    getGameById();
  }, [gameId]);

  return (
    <div className="country__info__wrapper">
      <button>
        <Link to="/">Back</Link>
      </button>

      {isLoading && !error && <h4>Loading........</h4>}
      {error && !isLoading && { error }}

      <div className="country__info__container">
          <div className="country__info-img">
            <img src={apiURL + game.BoxartFrontPath} alt="" width="500px"/>
          </div>

          <div className="country__info">
            <h3>{game.Name}</h3>

            <div className="country__info-left">
              <h5>
                Description: <span>{game.Overview}</span>
              </h5>
              <h5>
                Release Date: <span>{game.ReleaseDate}</span>
              </h5>
            </div>
            <div>
              <SetId 
                  gameId={game.ID}
              />
            </div>
          </div>
        </div>
    </div>
  );
};

export default GameInfo;