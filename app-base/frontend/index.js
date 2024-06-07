import generatePixelData from "./generatePixelData.js";
import fisrtDrawPixel from "./fisrtDrawPixel.js";
import drawPixel from "./drawPixel.js"
import onClickColor from "./events/onClickColor.js";
import onClickPixel from "./events/onClickPixel.js";


//PARTIE APP BASE


var ws;
const pixelSize = 10;
const canvasHeight = 600;
const canvasWidth = 600;
const canvasEl = document.querySelector('canvas');
const ctx = canvasEl.getContext('2d');
let selectedColor = null;
let selectedX = null;
let selectedY = null;

let pixelData = generatePixelData({
    pixelSize: pixelSize,
    width: canvasWidth,
    height: canvasHeight,
});

console.table(pixelData);

document.getElementById("connecter").onclick = function(evt) {
    if (ws) {
        return false;
    }

    var host = document.getElementById("host").value;
    var port = document.getElementById("port").value;

    addToLog("Tentative de connexion")
    addToLog("host = " + host + ", port = " + port)
    ws = new WebSocket("ws://"+host+":"+port+"/ws");

    ws.onopen = function(evt) {
        addToLog("Websocket ouverte");
        fisrtDrawPixel({canvasEl, ctx, pixelData, pixelSize, canvasWidth, canvasHeight});
        //drawPixel({ctx, pixelData, pixelSize, canvasWidth, canvasHeight});
        onClickColor({
            callback: (color) => {
                selectedColor = color;
                console.log("selectedColor", selectedColor);
            }
        });
        onClickPixel({ctx,selectedColor, canvasEl, pixelSize, pixelData, callback : (ctx, rowIndex, colIndex, pixelSize) => {
                if ( selectedColor?.length &&
                    rowIndex >= 0 &&
                    rowIndex < pixelData.length &&
                    colIndex >=0 &&
                    colIndex < pixelData[0].length) {
                        selectedX = rowIndex;
                        selectedY = colIndex;
                        console.log("selectedX", selectedX);
                        console.log("selectedY", selectedY);
                }
            }
        })
    }

    ws.onclose = function(evt) {
        addToLog("Websocket fermée");
        ws = null;
    }

    ws.onmessage = function(evt) {
        addToLog("Réception: " + evt.data);
        const parts = evt.data.split('/=');
        let x, y, R, G, B;

        parts.forEach(part => {
            if (part.startsWith('positionX=')) {
                x = parseInt(part.split('=')[1], 10);
            } else if (part.startsWith('positionY=')) {
                y = parseInt(part.split('=')[1], 10);
            } else if (part.startsWith('R=')) {
                R = parseInt(part.split('=')[1], 10);
            } else if (part.startsWith('G=')) {
                G = parseInt(part.split('=')[1], 10);
            } else if (part.startsWith('B=')) {
                B = parseInt(part.split('=')[1], 10);
            }
        });
        console.log("x:", x, "y:", y, "R:", R, "G:", G, "B:", B);

        if (x !== undefined && y !== undefined && R !== undefined && G !== undefined && B !== undefined) {
            const color = `rgb(${R}, ${G}, ${B})`;
            drawPixel({ ctx, color:color, rowIndex:x, colIndex:y, pixelData, pixelSize });
        } else {
            console.error("Une ou plusieurs valeurs sont indéfinies:", { x, y, R, G, B });
        }
    }

    ws.onerror = function(evt) {
        addToLog("Erreur: " + evt.data);
    }
    return false;
}

document.getElementById("fermer").onclick = function(evt) {
    if (!ws) {
        return false;
    }
    ws.close();
    return false;
}

document.getElementById("sauvegarder").onclick = function (evt){
    if (!ws) {
        return false;
    }
    let sndmsg = "sauvegarder"
    addToLog("Envoi: " + sndmsg);
    ws.send(sndmsg);
    return false;
}

document.getElementById("envoyer").onclick = function(evt) {
    if (!ws) {
        return false;
    }
    const button = evt.target;

    button.disabled = true;
    let R = 0;
    let G = 0;
    let B = 0;
    let posX = selectedX;
    let posY = selectedY;

    let regex = /^rgb\((\d+),\s*(\d+),\s*(\d+)\)$/;
    let result = selectedColor.match(regex);
    if (result) {
        R = parseInt(result[1]);
        G = parseInt(result[2]);
        B = parseInt(result[3]);

    } else {
        throw new Error("Le format de la chaîne n'est pas valide");
    }

    const sndmsg = "/=positionX="+posX.toString()+"/=positionY="+posY.toString()+"/=R="+R.toString()+"/=G="+G.toString()+"/=B="+B.toString();

    addToLog("Envoi: " + sndmsg);
    ws.send(sndmsg);

    setTimeout(function() {
        button.disabled = false;
    }, 10000);
    return false;
}

function addToLog(message) {
    var logs = document.getElementById("logs");
    var d = document.createElement("div");
    d.textContent = message;
    logs.appendChild(d);
    logs.scroll(0, logs.scrollHeight);
}








//PARTIE APP NET



var ws2;

document.getElementById("connecterNet").onclick = function(evt) {
    if (ws2) {
        return false;
    }

    var host = document.getElementById("hostNet").value;
    var port = document.getElementById("portNet").value;

    addToLogNet("Tentative de connexion")
    addToLogNet("host = " + host + ", port = " + port)
    ws2 = new WebSocket("ws://"+host+":"+port+"/ws");

    ws2.onopen = function(evt) {
        addToLogNet("Websocket ouverte");
    }

    ws2.onclose = function(evt) {
        addToLogNet("Websocket fermée");
        ws2 = null;
    }

    ws2.onmessage = function(evt) {
        addToLogNet(evt.data);
    }

    ws2.onerror = function(evt) {
        addToLog("Erreur: " + evt.data);
    }
    return false;
}

document.getElementById("fermerNet").onclick = function(evt) {
    if (!ws2) {
        return false;
    }
    ws2.close();
    return false;
}


function addToLogNet(message) {
    var logs = document.getElementById("logsNet");
    var d = document.createElement("div");
    d.textContent = message;
    logs.appendChild(d);
    logs.scroll(0, logs.scrollHeight);
}



// RECUPERATION AUTOMATIQUE DES PORTS
window.onload = function() {
    const urlParams = new URLSearchParams(window.location.search);
    const portB = urlParams.get('portB');
    const portN = urlParams.get('portN');

    const portBValue = parseInt(portB, 10);
    const portNValue = parseInt(portN, 10);


    if (!isNaN(portBValue)) {
        document.getElementById('port').value = portBValue;
    }

    if (!isNaN(portNValue)) {
        document.getElementById('portNet').value = portNValue;
    }

    console.log(portBValue)
    console.log(portNValue)
    setTimeout(function() {
        document.getElementById('connecter').click();
        document.getElementById('connecterNet').click();
    }, 5000);
};