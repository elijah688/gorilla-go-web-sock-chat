


const submit = (event, socket) => {
    event.preventDefault(event);

    const name = document.getElementById('name').value;
    const message = document.getElementById('message').value;

    console.log(name)
    console.log(message)
    const canSubmit = name != "" && message != "" && true || false
    if (canSubmit) {
        const a = { name, action: "broadcast", message }
        const b = JSON.stringify(a)
        socket.send(b);
    }

}

const openNav = () => document.getElementById("side-nav").style.width = "250px";

const closeNav = () => document.getElementById("side-nav").style.width = "0";

const handleConnectionEventMessages = (messages) => {
    var list = document.createElement('ul');
    messages.forEach(({ message, name, action }) => {

        if (action === "connected") {

            var li = document.createElement('li');
            li.classList = ["online"];


            var h = document.createElement('H4');
            let t = document.createTextNode(`${name}: `);
            h.appendChild(t);


            var p = document.createElement('p');
            t = document.createTextNode(message);
            p.appendChild(t);


            p.textContent = message;

            li.appendChild(h);
            li.appendChild(p);
            list.appendChild(li)
        }
    });
    document.querySelector('#online-items').appendChild(list);
    // document.querySelector('#online').innerHTML = list.innerHTML;
}

const main = () => {
    const socket = new WebSocket('ws://localhost:8080/ws');
    const nameInputField = document.getElementById('name');
    const messages = []


    socket.addEventListener('message', (event) => {
        const data = event.data
        messages.push(JSON.parse(data))
        handleConnectionEventMessages(messages)
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
        const a = { name, action: "disconnected", message: "goodby from client" }
        const b = JSON.stringify(a)
        socket.send(b);
        socket.close()
    }
}

main()


