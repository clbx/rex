import React, { useState, useEffect } from "react";
import { apiURL } from "../util/api";

import SearchInput from "../Search/SearchInput";
import FilterPlatform from "../FilterPlatform/FilterPlatform";

import { Link } from "react-router-dom";

import unknown from "../../assets/unknown.jpg"

const AllGames = () => {
  const [games, setGames] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  const getAllGames = async () => {
    try {
      const res = await fetch(`${apiURL}/v1/games`);

      if (!res.ok) throw new Error("Something went wrong!");

      const data = await res.json();

      console.log(data);

      setGames(data);

      setIsLoading(false);
    } catch (error) {
      setIsLoading(false);
      setError(error.message);
    }
  };

  const getGameByName = async (gameyName) => {
    try {
      const res = await fetch(`${apiURL}/v1/games`);
      // Filter Through Titles
      if (!res.ok) throw new Error("No games found");

      const data = await res.json();
      setGames(data);

      setIsLoading(false);
    } catch (error) {
      setIsLoading(false);
      setError(error.message);
    }
  };

  const getGameByPlatform = async (regionName) => {
    try {
      const res = await fetch(`${apiURL}/v1/games/byPlatform`);

      if (!res.ok) throw new Error("Failed..........");

      const data = await res.json();
      setGames(data);

      setIsLoading(false);
    } catch (error) {
      setIsLoading(false);
      setError(false);
    }
  };

  useEffect(() => {
    getAllGames();
  }, []);

  return (
    <div className="all__country__wrapper">
      <div className="country__top">
        <div className="search">
          <SearchInput onSearch={getGameByName} />
        </div>

        <div className="filter">
          <FilterPlatform onSelect={getGameByPlatform} />
        </div>
      </div>

      <div className="country__bottom">
        {isLoading && !error && <h4>Loading........</h4>}
        {error && !isLoading && <h4>{error}</h4>}


        {games?.map((game) => (
          <Link to={`/games/${game.ID}`}>
            <div className="country__card">
              <div className="country__img">

                {game.BoxartFrontPath === ""
                  ? <img src={unknown} alt="" />
                  : <img src={apiURL + game.BoxartFrontPath} alt = "" />
                }

                {console.log(game)}
              </div>

              <div className="country__data">
                <h3>{game.Name}</h3>
              </div>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
};

export default AllGames;