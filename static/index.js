

const openNav = () => document.getElementById("side-nav").style.width = "250px";

const closeNav = () => document.getElementById("side-nav").style.width = "0";

const showToast = () => {
    console.log('asd')
    var x = document.getElementById("snackbar");
  
    x.className = "show";
  
    setTimeout(function(){ x.className = x.className.replace("show", ""); }, 3000);
  }


const submit = (event, socket) => {
    event.preventDefault(event);

    const name = document.getElementById('name').value;
    const message = document.getElementById('message').value;


    const canSubmit = name != "" && message != "" && true || false
    if (canSubmit) {
        const a = { name, action: "broadcast", message }
        const b = JSON.stringify(a)
        socket.send(b);
    }

}


const handleConnectionUpdateNotification = (clients) => {
    const list = document.createElement('ul');
    clients.forEach(client => {
        const li = document.createElement('li');
        li.classList = ["online"];

        const h = document.createElement('H4');
        const t = document.createTextNode(`${client} connected...`);
        h.appendChild(t);
        li.appendChild(h);
        list.appendChild(li);
    });
    document.querySelector('#online-items').innerHTML = list.innerHTML


    openNav()
    const timeout = setTimeout(() => {
        closeNav()
        clearTimeout(timeout)
    }, 2000)

}
const handleConnectionEvents = ({ message, name, action, clients }) => {

    switch (action) {
        case "connected":
            if (message === `The name \"${name}\" is taken!`) {
                document.getElementById("snackbar").innerHTML = message;
                showToast()
            }
            handleConnectionUpdateNotification(clients)
            break;
        case "disconnected":
            handleConnectionUpdateNotification(clients)
        case "broadcast":
            const li = document.createElement('li');
            li.id = "list-item";
            const userName = document.getElementById('name').value;
            if (name === userName) {
                li.classList = ["home"];
            } else {
                li.classList = ["away"];
            }

            const h = document.createElement('H4');
            let t = document.createTextNode(`${name}: `);
            h.appendChild(t);


            const p = document.createElement('p');
            t = document.createTextNode(message);
            p.appendChild(t);

            li.appendChild(h);
            li.appendChild(p);

            document.querySelector('#messages').appendChild(li);
            break;
        default:
            break;
    }
}

const main = () => {
    const socket = new WebSocket('ws://localhost:8080/ws');
    const nameInputField = document.getElementById('name');


    socket.addEventListener('message', (event) => {
        handleConnectionEvents(JSON.parse(event.data))
    });

    nameInputField.addEventListener("change", () => {
        const name = nameInputField.value
        if (name != "") {
            const a = { name, action: "connected", message: "hello from client" }
            const b = JSON.stringify(a)
            socket.send(b);
        }
    })

    document.getElementById("chat-form")
        .addEventListener('submit', (event) => submit(event, socket));

    window.onbeforeunload = () => {
        const a = { name: nameInputField.value, action: "disconnected", message: "goodby from client" }
        const b = JSON.stringify(a)
        socket.send(b);
        socket.close()
    }
}

main()


