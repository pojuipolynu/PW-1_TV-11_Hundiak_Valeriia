async function calculate_1() {
    const inputs_1 = {
        omega: parseFloat(document.getElementById("omega").value) || 0.01,
        tsValue: parseFloat(document.getElementById("tsvalue").value) || 0.045,
        pmValue: parseFloat(document.getElementById("pmvalue").value) || 5120,
        tmValue: parseFloat(document.getElementById("tmvalue").value) || 6451,
        zavarValue: parseFloat(document.getElementById("zavarvalue").value) || 23.6,
        zplanValue: parseFloat(document.getElementById("zplanvalue").value) || 17.6,
        kpValue: parseFloat(document.getElementById("kpvalue").value) || 0.004,
        input_type: 1,
    };

    const response = await fetch("/calculate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(inputs_1)
    });

    if (!response.ok) {
        throw new Error("HTTP error " + response.status);
    }
    const data = await response.json();

    try {
        document.getElementById("results_1").innerHTML = 
        `<p>Очікувана відсутність енергопостачання в надзвичайних ситуаціях: ${data.mwneda.toFixed(2)}</p>
        <p>Очікуваний дефіцит енергії для запланованих: ${data.mwnedp.toFixed(2)}</p>
        <p>Загальна очікувана вартість перерв у роботі: ${data.mz.toFixed(2)}</p>`;
    }
    catch(err) {
        document.getElementById("results_1").innerHTML = 
        `<p>Виникла помилка, спробуйте пізніше</p>`;
    }
}

async function calculate_2() {
    let inputs = {
        input_type: 2,
        amounts: {}
    };

    document.querySelectorAll('.equipment-input').forEach(input => {
        inputs.amounts[input.id] = parseInt(input.value) || 0;
    });

    const response = await fetch('/calculate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(inputs)
    });

    if (!response.ok) {
        throw new Error("HTTP error " + response.status);
    }
    
    const result = await response.json();

    try {
        document.getElementById('results_2').innerHTML = `
        <p>Частота відмов одноколової системи: ${result.woc.toFixed(4)}</p>
        <p>Середня тривалість відновлення: ${result.tvoc.toFixed(4)}</p>
        <p>Коефіцієнт аварійного простою: ${result.kaoc.toFixed(4)}</p>
        <p>Коефіцієнт планового простою: ${result.kpoc.toFixed(4)}</p>
        <p>Частота відмов двох кіл: ${result.wdk.toFixed(4)}</p>
        <p>Частота відмов з урахуванням секційного вимикача: ${result.wds.toFixed(4)}</p>
    `;
    }
    catch(err) {
        document.getElementById("results_2").innerHTML = 
        `<p>Виникла помилка, спробуйте пізніше</p>`;
    }
}

