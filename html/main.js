let socket = null

window.onbeforeunload = () => {
    let json = {
        action: "left"
    }
    socket.send(JSON.stringify(json))
}

document.addEventListener("DOMContentLoaded", () => {
    let usernameInput = document.getElementById("username")
    let messageInput = document.getElementById("message")
    let usersListUl = document.getElementById("users_list")

    socket = new WebSocket("ws://localhost:8080/ws")

    socket.onopen = (event) => {
        console.log("successfully connected", event)
    }
    socket.onerror = (err) => {
        console.log("socket error", err)
    }
    socket.onclose = (closeEvent) => {
        console.log("socket close", closeEvent)
    }

    socket.onmessage = (event) => {
        let message = JSON.parse(event.data)
        console.log("OnMessage", message)
        switch (message.action) {
            case "userList":
                let userList = message.users_list
                // clear list
                while (usersListUl.firstChild) usersListUl.removeChild(usersListUl.firstChild)

                userList.forEach(element => {
                    let li = document.createElement("li")
                    li.textContent = element 
                    usersListUl.appendChild(li)
                });
                
                break;
        }
    }


    usernameInput.addEventListener("change", function() {
        let json = {
            action: "username",
            username: this.value
        }
        socket.send(JSON.stringify(json))
    })
    
})