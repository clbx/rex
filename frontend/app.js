'use strict';

const mainDiv = document.getElementById("main-div")
const urlParams = new URLSearchParams(window.location.search)



fetch("http://localhost:8080/v1/platforms")
    .then(response => response.json())
    .then((data) => {
        var linksDiv = document.getElementById("platform-links")
        data.forEach(function (item) {
            var platformLink = document.createElement('a')
            var platformText = document.createTextNode(item.Name)
            platformLink.title = item.Name + "-link"
            platformLink.href = window.location.href.split('?')[0] + "?platform=" + item.Platform
            platformLink.appendChild(platformText)
            linksDiv.append(platformLink)
            linksDiv.append(document.createTextNode(" "))
        })
    })

if (urlParams.has("platform")) {
    console.log("parameter search")
} else {
    fetch("http://localhost:8080/v1/games")
        .then(response => response.json())
        .then((data) => {
            makeGameGrid(data)
        })
}


function makeGameGrid(data) {
    var gamesDiv = document.getElementById("games")
    var topGridDiv = document.createElement("div")
    gamesDiv.appendChild(topGridDiv)
    topGridDiv.classList.add("container")
    topGridDiv.classList.add("text-center")
    var index = 0
    var currentRow = document.createElement("div")
    currentRow.classList.add("row")
    currentRow.classList.add("p-3")
    topGridDiv.appendChild(currentRow)
    data.forEach(function (game) {
        if (index == 4) {
            //Append row, and create a new one
            topGridDiv.appendChild(currentRow)
            currentRow = document.createElement("div")
            currentRow.classList.add("row")
            currentRow.classList.add("p-3")
            index = 0
        }
        currentRow.appendChild(makeGame(game))
        index += 1
    })
    topGridDiv.appendChild(currentRow)
}

/**
 * 
 * @param {*} game The game object returned from the server
 * @returns A div for the game
 */
function makeGame(game){
    var gameIcon = document.createElement("div")

    gameIcon.classList.add("container")
    gameIcon.classList.add("col")
    
    var elem = document.createElement("img")
    elem.width = "200"
    
    

    if(game.BoxartFrontPath != "") {
        elem.src = "http://localhost:8080" + game.BoxartFrontPath
        gameIcon.appendChild(elem)
    }
    else {
        elem.src = "assets/unknown.jpg"
        var text = document.createElement("div")
        text.classList.add("centered")
        text.appendChild(document.createTextNode(game.Name))
        gameIcon.appendChild(elem)
        gameIcon.appendChild(text)
    }
    return gameIcon
}