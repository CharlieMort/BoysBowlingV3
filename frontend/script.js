let rankings = document.getElementById("rankings")

function CreateRankingElement(scores, title) {
    let rankElement = document.createElement("div")
    rankElement.className = "rankings"
    let podiumElement = document.createElement("div")
    podiumElement.className = "podiums"
    for (let i = 0; i<3; i++) {
        let podWrapper = document.createElement("div")
        podWrapper.className = "podiumParent"
        let podName = document.createElement("p")
        podName.innerText = scores[i].name
        let pod = document.createElement("div")
        pod.classList.add("podium")
        switch(i) {
            case 0:
                pod.classList.add("p2")
                break;
            case 1:
                pod.classList.add("p1")
                break;
            case 2:
                pod.classList.add("p3")
                break;
        }
        let scoreEle = document.createElement("p")
        scoreEle.innerText = scores[i].score
        pod.appendChild(scoreEle)
        podWrapper.appendChild(podName)
        podWrapper.appendChild(pod)
        podiumElement.appendChild(podWrapper)
    }

    let titleEle = document.createElement("h1")
    titleEle.innerText = title

    rankElement.appendChild(podiumElement)
    rankElement.appendChild(titleEle)

    return rankElement
}

fetch("http://localhost:8888/stats")
    .then(resp => {
        resp.json()
            .then(jason => {
                console.log(jason)
                jason.map((stat) => {
                    rankings.appendChild(CreateRankingElement(stat.podium, stat.title))
                })
            })
    })
    .catch(err => {
        console.log(err)
    })