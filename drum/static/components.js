function newElement(obj) {
    let el
    let content 
    let attrs = new Map()
    for (let prop in obj) {
        if (prop.startsWith('div') || prop.startsWith('p') || prop.startsWith('span')) {
            let details = prop.split(".")
            el = document.createElement(details[0])
            if (details.length === 2) {
                el.setAttribute("class", details[1])
            }
            content = obj[prop]
        } else if (prop !== "click") {
            attrs.set(prop, obj[prop])
        }
    }
    if (el === undefined) {
        console.log("ops")
        return el
    }
    attrs.forEach((value, key) => el.setAttribute(key, value))
    if (content === undefined) {
        return el
    }
    if (Array.isArray(content)) {
        content.forEach(child => {
            let childEl = newElement(child)
            if (~(childEl === undefined)) { return el.appendChild(childEl) }
        })
    } else if  ((typeof content == 'string') || (content instanceof String)) {
        el.textContent = content
    } else {
        el.appendChild(newElement(content))
    }
    if (obj['click']) {
        el.addEventListener('onclik', obj['click'])
    }
    console.log(el)
    return el
}

function walletElement(wallet) {
    return newElement(
        {
            "p.walletdetails": 
            [
                {"span.wallet": '0' + wallet.token.substring(1,20)+'...'},
                {"span.walletbalance": wallet.balance.toLocaleString()}
            ] 
        }
    )
}

function postElement(post) {
    return newElement(
        {
            "div.post": [
                {"span.usercaption" : post.author},
                {"span.timestamp" : post.timestamp},
                {"p.postcontent" : post.content}
            ]
        }
    )
}

function newStage(author, stage, members) {
    return stage = {
        posts: [],
        stage: stage,
        author: author,
        members: members,
        append: function(post) {
            this.posts.push(post)
            let contentWindow=document.getElementById("contentWindow")
            contentWindow.appendChild(postElement(post))
        }
    }
}

let MyPosts = newStage("Ruben", "Aereum", [])

let MyWallets = {
    wallets: [],
    render: function() {
        this.wallets.sort((x,y) => x.balance < y.balance)
        let walletContainer = document.getElementById("wallets")
        walletContainer.innerHTML = ''
        this.wallets.forEach(wallet => walletContainer.appendChild(walletElement(wallet)))
    },
    remove: function(hash) {
        document.getElementById(hash).remove()
    },
    append: function(wallet) {
        this.wallets.push(wallet)
        this.render()
    },
    update: function(wallet) {
        for(let n=0; n<this.wallets.length;n++) {
            if (this.wallets[n].token == wallet.token) {
                this.wallets[n] = wallet
                break
            } 
        }
        this.render()
    }
}

let socket = new WebSocket("ws://localhost:7000/ws")
socket.onopen = () => {}
socket.onmessage = (msg) => {
    data = JSON.parse(msg.data)
    if (data.action === "NewWalletBalance") {
        MyWallets.append(data)
    } else if (data.action === "NewStagePost") {
        console.log(data)
        MyPosts.append(data)
    }
}
socket.onclose = (event) => {}
