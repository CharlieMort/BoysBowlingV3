let frameTable = document.getElementById("frameTable")
let addRow = document.getElementById("addRowButton")

addRow.addEventListener("click", () => {
    let userInp = document.createElement("input")
    userInp.type = "text"
    userInp.style = "width: 5rem;"
    userInp.name = "playerName"

    let scoreInp = document.createElement("input")
    scoreInp.type = "text"
    scoreInp.style = "width: 15rem;"
    scoreInp.name = "scorecard"

    let remBut = document.createElement("button")
    remBut.innerText = "X"
    remBut.id = "removeRow"
    remBut.onclick = RemoveRow

    let trUser = document.createElement("td")
    trUser.appendChild(userInp)

    let trScore = document.createElement("td")
    trScore.appendChild(scoreInp)

    let trRem = document.createElement("tr")
    trRem.appendChild(remBut)

    let tr = document.createElement("tr")
    tr.appendChild(trUser)
    tr.appendChild(trScore)
    tr.appendChild(trRem)

    frameTable.appendChild(tr)
})

function RemoveRow(e) {
    e.preventDefault()
    console.log(e.target)
    frameTable.removeChild(e.target.parentElement.parentElement)
}


let form = document.getElementById("addGameForm")
form.addEventListener("submit", (e) => {
    e.preventDefault()
    let formDat = new FormData(form)
    let playerNames = formDat.getAll("playerName")
    let scorecards = formDat.getAll("scorecard")
    let frames = []
    playerNames.map((player, idx) => {
        frames.push({
            name: player,
            scorecard: scorecards[idx]
        })
    })

    let body = {
        datePlayed: form.elements["datePlayed"].value,
        frames: frames
    }

    fetch("http://localhost:8888/game", {
        method: "POST",
        body: JSON.stringify(body)
    }).then((resp) => {
        resp.json(jason)
            .then(console.log(jason))
    })
})