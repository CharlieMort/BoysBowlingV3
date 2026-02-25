let frameTable = document.getElementById("frameTable")
let addRow = document.getElementById("addRowButton")

addRow.addEventListener("click", () => {
    let userInp = document.createElement("input")
    userInp.type = "text"

    let scoreInp = document.createElement("input")
    scoreInp.type = "text"

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
