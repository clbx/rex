import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { apiURL } from "../util/api";
import { Link } from "react-router-dom";



const SetId = (gameId) => {
    const [responseMessage, setResponseMessage] = useState("")
    const [id, setId] = useState("")


    const submitId = async (event) => {
        event.preventDefault();
        console.log(id)
        console.log(gameId)
        
        
        const res = await fetch(`${apiURL}/v1/games/setGameById?id=${gameId.gameId}&tgdbid=${id}`, {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            }
        });
    };

    return (
        <>
            <form onSubmit={submitId}>
                <input type="text" value={id} 
                onChange={(event) => {
                    setId(event.target.value);
                }}
                />
                <button type="submit">Submit</button>
            </form>
        </>
    );
};

export default SetId;