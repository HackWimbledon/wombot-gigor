window.addEventListener("load", () => {
  connectws();
  document.getElementById("send").onclick = send;
});

var ws;

var print = function(message) {
  var d = document.createElement("div");
  d.innerHTML = message;
  output.appendChild(d);
};

function sendIgorCmd(cmd, params) {
  igorCmd = {};
  igorCmd.cmd = cmd;
  igorCmd.args = params;
  ws.send(JSON.stringify(igorCmd));
}

function requestBrains() {
  sendIgorCmd("request", { for: "brains" });
}

function send(evt) {
  if (!ws) {
    return false;
  }
  input = document.getElementById("input");
  print("SEND: " + input.value);
  ws.send(input.value);
  return false;
}

function connectws() {
  return fetch("/config").then(response => {
    response.json().then(config => {
      ws = new WebSocket(config.websocket);
      ws.onopen = function(evt) {
        print("OPEN");
        requestBrains();
      };
      ws.onclose = function(evt) {
        print("CLOSE");
        ws = null;
      };
      ws.onmessage = function(evt) {
        igorMsg = JSON.parse(evt.data);
        if (igorMsg.cmd == "brains") {
          // Response contains current brains list
          updateBrains(igorMsg);
        } else {
          print("MESSAGE " + evt.data);
        }
      };
      ws.onerror = function(evt) {
        print("ERROR: " + evt.data);
      };
      return ws;
    });
  });
}

function updateBrains(igormsg) {
    console.log("UpdateBrains");
    brainmap = igormsg.resp;
    ids=[];
    var list = this.document.getElementById("brainrows");
    rows = "";
    Object.entries(brainmap).forEach(([bk, bv]) => {
      console.log(bk, bv);
      rows = rows + brainrow(bv);
      ids.push(bv.id)
    });
    list.innerHTML = rows;
    for (var id of ids) {
        button=this.document.getElementById("brainstatus."+id);
        button.onclick=() => { sendIgorCmd("start", { "brain": id}) };
    }
  }

function brainrow(brain) {
  return `<tr>
        <td>
            <button id="brainstatus.${
              brain.id
            }" class="pure-button">Start</button>
        </td>
        <td>
            ${brain.name}
        </td>
    </tr>`;
}
