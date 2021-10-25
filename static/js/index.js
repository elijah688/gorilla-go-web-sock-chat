


const submit = (event, socket, name) => {
    event.preventDefault(event);
    if (name != "") {
        const a = { name, action: "broadcast", message: "hello to all " }
        const b = JSON.stringify(a)
        socket.send(b);
    }

}

const main = () => {

    let socket;
    const nameInputField = document.getElementById('name');

    nameInputField.addEventListener("change", () => {
        const name = nameInputField.value
        if (name != "") {
            if (socket) {
                console.log('closing....')
                socket.close()
            }
            socket = new WebSocket('ws://localhost:8080/ws');

            socket.addEventListener('open', () => {
                const a = { name, action: "connected", message: "hello server" }
                const b = JSON.stringify(a)
                socket.send(b);
            });

            socket.addEventListener('message', (event) => {
                console.log('Message from server ', event.data);
            });

            document.getElementById("chat-form")
                .addEventListener('submit', (event) => submit(event, socket, name));

        }
    })

}

main()





