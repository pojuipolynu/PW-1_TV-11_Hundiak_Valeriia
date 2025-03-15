async function calculate_1() {

    const inputs_1 = {
        efficiency: parseFloat(document.getElementById("efficiency").value) || 0.92,
        powerFactor: parseFloat(document.getElementById("powerfactor").value) || 0.9,
        voltage: parseFloat(document.getElementById("voltage").value) || 0.38,
        quantity: parseFloat(document.getElementById("quantity").value) || 4,
        pH: parseFloat(document.getElementById("ph").value) || 20,
        kB: parseFloat(document.getElementById("kb").value) || 0.21,
        tg: parseFloat(document.getElementById("tg").value) || 1.55,
    };

    const response = await fetch("/calculate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(inputs_1)
    });

    // if (!response.ok) {
    //     throw new Error("HTTP error " + response.status);
    // }
    const data = await response.json();
    try {
        document.getElementById("results_1").innerHTML = 
        `<p>Розрахунковий струм: ${data.current.toFixed(2)}</p>
        <p>Груповий коефіцієнт використання: ${data.groupusage.toFixed(2)}</p>
        <p>Ефективна кількість ЕП: ${data.effectiveqty.toFixed(2)}</p>
        <p>Активне навантаження: ${data.activepower.toFixed(2)}</p>
        <p>Реактивне навантаження: ${data.reactivepower.toFixed(2)}</p>
        <p>Повна потужність: ${data.totalpower.toFixed(2)}</p>`;
    }
    catch(err) {
        document.getElementById("results_1").innerHTML = 
        `<p>Виникла помилка, спробуйте пізніше</p>`;
    }
}
