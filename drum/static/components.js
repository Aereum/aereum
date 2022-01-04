// placeholder implement a very simple container of updatable elements
// id must refer to a valid tag with respective id in the document, 
// render is a function that gets a value and returns the respective
// HTML element. Sort is an optional sort function for the values.
function placeholder(id, render, sort) {
    return {
        el: document.getElementById(id),
        values: [],
        has: function(token) {
           for(let value of this.values) {
               if (value.token === token) {
                   return value
               }
           }
           return false
        },
        update: function(value) {
            let found = -1
            let position = 0
            for(let n=0;n<this.values.length;n++) {
                // sort(a,b) return b <= a
                if (sort(this.values[n],value)) {
                    position = n
                } else if (found > 0) {
                    break
                }
                if (this.values[n].token === value.token) {
                    this.values[n] = value
                    found = n
                }
            }
            if (found < 0) {
                return
            }
            const old = document.getElementById(id+value.token)
            if (!old) {
                return
            }
            if (position === found) {
                this.el.replaceChild(render(value),old)
            } else if (postion < this.values.length - 1) {
                before = document.getElementById(id + this.values[position+1])
                el.removeChild(old)
                this.el.insertBefore(render(value),before)
            } else {
                el.removeChild(old)
                this.el.appendChild(render(value))
            }
        },
        append: function(value) {
            for(let n=0;n<this.values.length;n++) {
                // sort(a,b) return b <= a
                if (sort(this.values[n],value)) {
                    this.el.insertBefore(render(value), document.getElementById(id+this.values[n].token))
                    this.values.splice(n,0,value)
                    return
                }
            }
            this.el.appendChild(render(value))
            this.values.push(value)
        },
        remove: function(value) {
            for(let n=0;n<this.values.length;n++) {
                if (this.values[n].token === value.token) {
                    this.values.splice(n,1)
                }
                const element = document.getElementById(id+value.token)
                this.el.removeChild(element)
                return
            }            
        }
    }
}

// returns a new tag element with className set if any is provided. it is rendered
// with content that can be a text, another element or an array of other elements. 
// id and onclick can be set with the helper functions that returns the instance 
// of the tag element with those properties set. 
function newTagElement(tag, className) {
    let el = document.createElement(tag)
    if (className) {
        el.setAttribute("class", className)
    }
    return {
        el : el,
        id: function(identification) {
            this.el.setAttribute("id", identification)
            return this
        },
        onClick: function(handler) {
            this.el.addEventListener('onClick', handler)
            return this
        },
        render: function(content) {
            if (Array.isArray(content)) {
                for(let child of content) {
                    if (~(child === undefined)) {
                        this.el.appendChild(child)
                    }
                }
            } else if  ((typeof content == 'string') || (content instanceof String)) {
                el.textContent = content
            } else {
                el.appendChild(content)
            }
            return this.el
        }
    }
}

const div = (className) => newTagElement("div", className)
const span = (className) => newTagElement("span", className)
const p = (className) => newTagElement("p", className)

const walletBalance = (wallet) => p("walletdetails")
    .id("balance"+wallet.token).render([
        span("wallet").render('0' + wallet.token.substring(1,20)+'...'),
        span("walletbalance").render(wallet.balance.toLocaleString())
])

const post = (post) => {
    subComponents = [
        span("usercaption").render(post.author),
        span("timestamp").render(post.timestamp),
        p("postcontent").render(post.content)
    ]
    if (post.moderation) {
        subComponents.push(
            p("aprove")
            .onClick(() => aprovetoken(post.token))
            .render("Aprove")
        )
    }
    return div("post "+post.status).render(subComponents)
}

const renderStageHeader = (stage) => {
    const header = p().render([
        span("stageAuthor").render(stage.author + ":"),
        span("stageName").render(stage.Name)
    ])
    document.getElementById("midHeader").innerHTML = header
}

const wellcome = div("wellcome").render([
    p().render("Welcome to aereum network")
])

const MyAttorneys = placeholder("attorneys", attorney)
const MyWallets = placeholder("wallets", wallet, (x,y) => y.balance >= x.balance)
const MyStages = placeholder("mystages", stage, (x,y) => x.Caption >= y.Caption)
const MyModerations = placeholder("mymoderation", moderation, (x,y) => x.Caption >= y.Caption)
const MyEngagement = placeholder("myengagements", engagement, (x,y) => x.Name >= y.Name)

let socket = new WebSocket("ws://localhost:7000/ws")
socket.onopen = () => {}
socket.onmessage = (msg) => {
    data = JSON.parse(msg.data)
    if (data.action === "NewWalletBalance") {
        MyWallets.append(data)
    } else if (data.action === "NewStagePost") {    
        const moderated = MyModerations.has(data.stage.token)
        if (moderated) {
            moderated.count += 1
            MyModerations.update(moderated)
        }
        if (MainView.token === data.stage.token) {
            MainView.insert(data.post)
        }
        const engaged = MyEngagement(data.stage.token)
        if (engaged) {
            engaged.count += 1
            MyEngagement.update(engaged)
        }
    } else if (data.action === "NewStage") {
        if (data.moderation) {
            ModerationStages.add(data.token)
        }
        ActiveStages.add(data.token)
    }
}
socket.onclose = (event) => {}
