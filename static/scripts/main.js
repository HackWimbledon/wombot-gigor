
console.log("Hello");

window.addEventListener("load", () => {
    fetch("/brains").then((response) => {
        response.json().then((brains) => {
            brainmap=brains.brains;
            var list = this.document.getElementById("brainrows");
            rows=""
            Object.entries(brainmap).forEach(([bk,bv]) => {
                rows=rows+brainrow(bv);
            });
            list.innerHTML=rows
        })
    })
});

function brainrow(brain) {
    return `<tr>
        <td id="brainstatus.${brain.Id}">
            <button class="pure-button">Start</button>
        </td>
        <td>
            ${brain.Name}
        </td>
    </tr>`
}