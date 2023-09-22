window.gear5thadserver = window.gear5thadserver || {}
window.gear5thadserver[document.currentScript.getAttribute("adslot")] = window.gear5thadserver[document.currentScript.getAttribute("adslot")] || {
    adSlotId: document.currentScript.getAttribute("adslot"),
    async main() {
        let ins = document.getElementById(this.adSlotId)

        let frame = document.createElement("iframe")
        frame.id = "gear5th-adframe-" + ins.dataset.adslotId
        frame.setAttribute("src", ins.dataset.adserverUrl)
        frame.setAttribute("width", ins.dataset.frameWidth)
        frame.setAttribute("height", ins.dataset.frameHeight)
        frame.setAttribute("loading", "lazy")
        frame.style.border = "none"
        frame.setAttribute("scrolling", "no")
        ins.appendChild(frame)


        let impressionCheckInterval = setInterval(() => {
            if (this.isPartiallyInViewport(frame)) {
                frame.contentWindow.postMessage("gear5th-ad-impression", "*");
                clearInterval(impressionCheckInterval)
            }
        }, 5000);

    },

    isPartiallyInViewport(elem) {

        let bounding = elem.getBoundingClientRect();
        let elemHeight = elem.offsetHeight;
        let elemWidth = elem.offsetWidth;

        return (bounding.top >= -elemHeight
            && bounding.left >= -elemWidth
            && bounding.right <= (window.innerWidth || document.documentElement.clientWidth) + elemWidth
            && bounding.bottom <= (window.innerHeight || document.documentElement.clientHeight) + elemHeight)
    },


    isFullyInParentViewport(elem) {
        let bounding = elem.getBoundingClientRect();
        return (bounding.top >= 0 &&
            bounding.left >= 0 &&
            bounding.right <= (window.parent.innerWidth || window.parent.document.documentElement.clientWidth) &&
            bounding.bottom <= (window.parent.innerHeight || window.parent.document.documentElement.clientHeight))
    }
}

window.gear5thadserver[document.currentScript.getAttribute("adslot")].main()
