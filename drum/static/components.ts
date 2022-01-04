type Tag = {
    el: HTMLElement;
    id: (identification: any) => any;
    onClick: (handler: any) => any;
    render: (content: string) => any;
}

function TagElement(tag: string, className: string | undefined): Tag {
    let el = document.createElement(tag)
    if (className) {
        el.setAttribute("class", className)
    }
    return {
        el : el,
        id: function(identification: string) {
            this.el.setAttribute("id", identification)
            return this
        },
        onClick: function(handler) {
            this.el.addEventListener('onClick', handler)
            return this
        },
        render: function(content: string | HTMLElement | HTMLElement[]): HTMLElement {
            if (Array.isArray(content)) {
                for(let child of content) {
                    if (~(child === undefined)) {
                        this.el.appendChild(child)
                    }
                }
            } else if  (typeof content == 'string') {
                el.textContent = content
            } else {
                el.appendChild(content)
            }
            return this.el
        }
    }
}