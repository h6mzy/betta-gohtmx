async function loadGrid() {
  const res = await fetch("/api/species");
  const data = await res.json();

  const grid = document.getElementById("species-grid");

  grid.innerHTML = data.map(s => `
    <div class="card">
      <a href="/species/${s.id}" onclick="loadSpecies('${s.id}'); return false;">
        <img src="/images/${s.slug}.jpg">
        <h2>${s.name}</h2>
        <p>${s.redlist_status}</p>
      </a>
    </div>
  `).join("");
}

async function loadSpecies(id) {
  const res = await fetch(`/api/species/${id}`);
  const s = await res.json();

  document.getElementById("species-grid").innerHTML = `
    <div class="detail">
      <img src="/images/${s.slug}.jpg">
      <h2>${s.name}</h2>
      <p>Color: ${s.color}</p>
      <p>Size: ${s.size}</p>
      <p>Redlist: ${s.redlist_status}</p>
    </div>
  `;

  document.title = `BETTA ${s.name}`;
  document.getElementById("title").innerText = s.name;

  history.pushState({}, "", `/species/${id}`);
}

window.addEventListener("popstate", () => {
  loadGrid();
});

loadGrid();