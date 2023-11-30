const fb = document.getElementById('fb')
const username = document.getElementById('username')
const fbMsg = "Must choose an user name"
const players = document.getElementById('players-cards')

const ws = new WebSocket("ws://localhost:3000/ws" + location.pathname)

function makeWaitingPlayers(n) {
    for (let i = 0; i < n; i++) {
        players.innerHTML += `<div class="player"><p><span class="player-number">p${i+1}</span> Waiting...</p></div>`
    }
}

makeWaitingPlayers(2)

ws.onopen = _ => {
    console.log("ws open connection")
}

ws.onclose = _ => {
    console.log("ws close connection")
}

ws.onmessage = event => {
    let j = JSON.parse(event.data)

    console.log(j.event)
    
    switch (j.event) {
        case "addUser":
            let ps = [...players.querySelectorAll(".player")]

            for (let i = 0; i < ps.length; i++) {
                const p = ps[i];
                if (p.id == "") {
                    console.log(p)
                    p.outerHTML = j.content
                    break;
                }
            }

            let wsEvent = {
                "event": "userAdded",
                "content": j.userId
            }

            ws.send(JSON.stringify(wsEvent))
            
            break;
        case "start":
            // <- (maybe fuse with turn)
            // allow turn user to click
            break;
        case "click":
            // will not be here
            // ->
            // send card that was clicked
            break;
        case "turn":
            // <-
            // turn card with this id
            break;
        case "paired":
            // <-
            // remove cards from view
            break;
        case "notpaired":
            // <-
            // unturn cards
            break;
        case "allpaired":
            // <-
            // end game
            break;
        case "leave":
            // make a feedback and leave card visible
            document.getElementById(j.userId).remove()
            break;
    }
}

document.getElementById('name').addEventListener("submit", function(event) {
    event.preventDefault()

    if (username.value == "") {
        fb.innerHTML = fbMsg
        return
    }

    document.getElementById("modal").style = "display:none;"

    let wsEvent = {
        "event": "updateUserName",
        "content": document.getElementById("username").value
    }

    ws.send(JSON.stringify(wsEvent))

    memory.New()
})

username.addEventListener('keyup', function(){ fb.innerHTML = this.value == "" ? fbMsg : "" })

var cfg = {
  id: 'canvas',
  loading: {
    show: () => {
    //   document.getElementById('loading').style.display = 'block'
    },
    hide: () => {
    //   document.getElementById('loading').style.display = 'none'
    },
  },
  renderer : { antialias: true },
  cards: {
    geometry: '../assets/geometry/geometry.json',
    grid: [6, 6],
    // get from backend
    cardsIds: ['1', '2', '3'],
    cards: ["🛩️"] //shuffled
  },
  paired: {
    target: '.last-span'
  }
}

const memory = Memory(cfg)

memory.Cards.select(id => {
  console.log(id)
})