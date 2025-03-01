async function calculate_1() {
    const inputs_1 = {
        U: parseFloat(document.getElementById("ui").value) || 1,
        Kz: parseFloat(document.getElementById("kz").value) || 1,
        Time: parseFloat(document.getElementById("time").value) || 1,
        Sm: parseFloat(document.getElementById("sm").value) || 1,
        Tm: parseFloat(document.getElementById("tm").value) || 1,
        input_type: 1,
    };

    const response = await fetch("/calculate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(inputs_1)
    });

    const data = await response.json();
    document.getElementById("results_1").innerHTML = 
    `<p>Броньований кабель: ${data.bronk.toFixed(2)}</p>
     <p>ААБ кабель: ${data.abbk.toFixed(2)}</p>`;
}

async function calculate_2() {
    const inputs_2 = {
        Kzu: parseFloat(document.getElementById("kzu").value) || 1,
        input_type: 2,
    };

    const response = await fetch("/calculate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(inputs_2)
    });
    
    const data = await response.json();

    document.getElementById("results_2").innerHTML = 
     `
        <p>Струм трифазного КЗ: ${data.strum3.toFixed(2)}</p>
    `;
}

async function calculate_3() {
    const inputs_3 = {
        Rh: parseFloat(document.getElementById("rh").value) || 1,
        Xh: parseFloat(document.getElementById("xh").value) || 1,
        Rm: parseFloat(document.getElementById("rm").value) || 1,
        Xm: parseFloat(document.getElementById("xm").value) || 1,
        input_type: 3,
    };

    const response = await fetch("/calculate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(inputs_3)
    });
    
    const data = await response.json();

    document.getElementById("results_3").innerHTML = 
     `
        <p>Струм трифазного КЗ.\nНормальний режим: ${data.ish_3.toFixed(2)}. Мінімальний режим: ${data.ish_3_min.toFixed(2)}</p>
        <p>Дійсний струм трифазного КЗ.\nНормальний режим: ${data.ish_2.toFixed(2)}. Мінімальний режим: ${data.ish_2_min.toFixed(2)}</p>
        <p>Дійсний струм трифазного КЗ.\nНормальний режим: ${data.dish_3.toFixed(2)}. Мінімальний режим: ${data.dish_3_min.toFixed(2)}</p>
        <p>Дійсний струм двофазного КЗ.\nНормальний режим: ${data.dish_2.toFixed(2)}. Мінімальний режим: ${data.dish_2_min.toFixed(2)}</p>
        <p>Аварійний режим на данній підстанції не передбачений.</p>
    `;
}
