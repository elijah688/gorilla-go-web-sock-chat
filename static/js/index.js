


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

const main = () => {

    const socket = new WebSocket('ws://localhost:8080/ws');
    const nameInputField = document.getElementById('name');

    // socket.addEventListener('open', () => {
    // const a = { name: name, action: "connected", message: "hello from client" }
    // const b = JSON.stringify(a)
    // socket.send(b);
    // });

    socket.addEventListener('message', (event) => {
        console.log('Message from server ', event.data);
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





