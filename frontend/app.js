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
    topGridDiv.appendChild(currentRow)
    data.forEach(function (item) {
        if (index == 4) {
            //Append row, and create a new one
            topGridDiv.appendChild(currentRow)
            currentRow = document.createElement("div")
            currentRow.classList.add("row")
            index = 0
        }
        var currentCol = document.createElement("div")
        currentCol.classList.add("col")
        var elem = document.createElement("img")
        if(item.BoxartFrontPath != ""){
            elem.src = "http://localhost:8080" + item.BoxartFrontPath
        }
        else{
            elem.src = "assets/unknown.jpg"
        }
        elem.width="250"
        currentCol.appendChild(elem)
        currentCol.appendChild(document.createTextNode(item.Name))
        currentRow.appendChild(currentCol)
        index += 1
    })
    topGridDiv.appendChild(currentRow)
}