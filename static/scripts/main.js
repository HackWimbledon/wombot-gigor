var term;

window.addEventListener("load", () => {
  term = new Terminal();
  term.setOption("fontSize", 12);
  term.open(document.getElementById("terminal"));
  connectws();
});

var ws;

var print = function(message) {
  term.writeln(message);
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
          console.log(igorMsg);
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
  console.log("UPDATE BRAINS");
  brainmap = igormsg.resp;
  ids = [];
  var list = this.document.getElementById("brainrows");
  rows = "";
  Object.entries(brainmap).forEach(([bk, bv]) => {
    console.log(bk, bv);
    rows = rows + brainrow(bv);
    ids.push(bv.id);
  });
  list.innerHTML = rows;
  for (var id of ids) {
    button = this.document.getElementById("brainstatus." + id);
    button.onclick = (evt) => {
      sendIgorCmd("start", { brain: evt.target.getAttribute("data-arg1") });
    };
  }
}

function brainrow(brain) {
  return `<tr>
        <td>
            <button id="brainstatus.${brain.id}" data-arg1="${
    brain.id
  }" class="pure-button">Start</button>
        </td>
        <td>
            ${brain.name}
        </td>
    </tr>`;
}
