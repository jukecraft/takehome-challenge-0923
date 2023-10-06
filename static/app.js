const Controller = {
    search: (ev) => {
        ev.preventDefault();
        const form = document.getElementById("form");
        const data = Object.fromEntries(new FormData(form));
        const response = fetch(`/search?q=${data.query}`).then((response) => {
            response.json().then((results) => {
                Controller.updateTable(results);
            });
        });
    },
    moreResults: (ev) => {
        ev.preventDefault();
        const form = document.getElementById("form");
        const data = Object.fromEntries(new FormData(form));
        const table = document.getElementById("table");
        const numberOfResultsSoFar = table.tBodies[0].rows.length
        fetch(`/search?q=${data.query}&existing=${numberOfResultsSoFar}`)
            .then((response) => {
                response.json().then((results) => {
                    Controller.updateTable(results);
                });
            });
    },
    updateTable: (results) => {
        const table = document.getElementById("table-body");
        const rows = [];
        for (let result of results) {
            rows.push(`<tr><td>${result}</td></tr>`);
        }
        table.innerHTML = rows;
    },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
const loadMoreButton = document.getElementById("load-more");
loadMoreButton.addEventListener("click", Controller.moreResults);
