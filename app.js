function newColorButton(name) {
    let btn = document.createElement("button");
    btn.type = "button";
    btn.className = "list-group-item list-group-item-action";
    btn.innerText = name;
    btn.ariaValueText = name;
    return btn;
}

// Initialize settings modal with current app settings
function initSettingsModal(settings) {
    let colorsList = document.getElementById('colors');
    for (let i in settings.c) {
        colorsList.appendChild(newColorButton(i));
    }

    let i = 0;
    for (let e of colorsList.children) {
        (() => {
            e.addEventListener("click", btn => {
                displaySettingsColor(settings, btn.target.ariaValueText);
            });
            if (i == 0) {
                displaySettingsColor(settings, e.ariaValueText);
            }
            i++;
        })();
    }
}

// Reset settings modal to its initial state
function resetSettingsModal() {
    document.getElementById('colors').innerHTML = null;
}

// When a color button is clicked, setup the settings for this color
function displaySettingsColor(settings, color) {
    let colorsList = document.getElementById('colors');
    for (let c of colorsList.children) {
        (() => {
            if (c.ariaValueText == color) {
                c.classList.add("active");

                let colorInput = document.getElementById("colorInput");
                let colorInput2 = colorInput.cloneNode();
                colorInput2.value = color;
                //todo change color in config_settings with event listener
                colorInput.parentNode.replaceChild(colorInput2, colorInput);
                colorInput.remove();

                let populationInput = document.getElementById('populationInput');
                let populationInput2 = populationInput.cloneNode();
                populationInput2.value = settings.c[color];
                populationInput2.addEventListener("change", e => {
                    config_settings.c[color] = parseInt(e.target.value);
                });
                populationInput.parentNode.replaceChild(populationInput2, populationInput);
                populationInput.remove();
            } else {
                c.classList.remove("active");
            }
        })();
    }
}

function updateSettings() {
    setSettings(JSON.stringify(config_settings));
}

let config_settings = {}

document.addEventListener('DOMContentLoaded', function () {
    //* Initialize Bootstrap tooltips
    Array.from(document.querySelectorAll('[data-bs-toggle="tooltip"]')).forEach(t => {
        new bootstrap.Tooltip(t);
    });

    //* Initialize settings modal on show
    document.getElementById('settingsModal').addEventListener('show.bs.modal', e => {
        settings = JSON.parse(getSettings());
        config_settings = settings;
        console.log(settings);
        initSettingsModal(settings);
    })

    //* Clear settings modal on hide
    document.getElementById('settingsModal').addEventListener('hidden.bs.modal', e => {
        config_settings = {}
        resetSettingsModal();
    })

    //* fetch and start webasm app
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("app.wasm", { cache: "no-store" }), go.importObject).then((result) => {
        go.run(result.instance);
        startApp("{\"c\":{\"#FF0000\":20,\"#0000FF\":30},\"r\":[{\"c1\":\"#0000FF\",\"c2\":\"#0000FF\",\"f\":-0.1,\"r\":400},{\"c1\":\"#FF0000\",\"c2\":\"#FF0000\",\"f\":0.8,\"r\":200},{\"c1\":\"#FF0000\",\"c2\":\"#0000FF\",\"f\":0.5,\"r\":400},{\"c1\":\"#0000FF\",\"c2\":\"#FF0000\",\"f\":-0.8,\"r\":200}]}");
    })

}, false);