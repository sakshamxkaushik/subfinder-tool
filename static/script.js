document.getElementById("enumeration-form").addEventListener("submit", function (e) {
    e.preventDefault();
    
    const domain = document.getElementById("domain").value;
    const concurrency = document.getElementById("concurrency").value;

    fetch("/enumerate", {
        method: "POST",
        body: new URLSearchParams({ domain, concurrency }),
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
    })
    .then(response => response.json())
    .then(data => {
        // Handle and display the results in the "results" div
        const resultsList = document.getElementById("results");

        // Clear existing results
        resultsList.innerHTML = "";

        // Populate the results list
        for (const key in data) {
            if (data.hasOwnProperty(key)) {
                const ips = data[key];
                const listItem = document.createElement("li");
                listItem.textContent = key + ": " + ips.join(", ");
                resultsList.appendChild(listItem);
            }
        }
    })
    .catch(error => console.error(error));
});
