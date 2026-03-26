async function loadGrid() {
  const res = await fetch("/api/species");
  const data = await res.json();

  const grid = document.getElementById("species-grid");
  const list = document.createElement("ul");

  list.innerHTML = data.map(s => `
    <li id="species-${s.slug}">
      <a href="/species/${s.slug}" onclick="loadSpecies('${s.slug}'); return false;">
        <img src="/images/${s.slug}.jpg" alt="${s.name}">
        <h2>${s.name}</h2>
        <p>Redlist: ${s.redlist}</p>
      </a>
    </li>
  `).join("");

  grid.innerHTML = "";
  grid.appendChild(list);
  document.title = "BETTA";
  document.getElementById("title").innerText = "BETTA";
}

async function loadSpecies(slug) {
  const res = await fetch(`/api/species/${slug}`);
  if (!res.ok) {
    alert("Species not found");
    return;
  }
  const s = await res.json();

  const grid = document.getElementById("species-grid");
  grid.innerHTML = `
    <div class="detail">
      <img src="/images/${s.slug}.jpg" alt="${s.name}">
      <h2>${s.name}</h2>
      <p><strong>Scientific:</strong> ${s.scientific}</p>
      <p><strong>Info:</strong> ${s.info}</p>
      <p><strong>Habitat:</strong> ${s.habitat}</p>
      <p><strong>Breeding:</strong> ${s.breeding}</p>
      <p><strong>Captive:</strong> ${s.captive}</p>
      <p><strong>Redlist:</strong> ${s.redlist}</p>
      ${s.photos.map(p => `<img src="${p.src}" alt="${p.caption}" style="max-width:200px;"><p>${p.caption}</p>`).join("")}
      <button onclick="loadGrid()">Back to list</button>
    </div>
  `;

  document.title = s.name;
  document.getElementById("title").innerText = s.name;
  history.pushState({}, "", `/species/${slug}`);
}

// Handle browser back/forward
window.addEventListener("popstate", () => {
  const path = window.location.pathname;
  if (path.startsWith("/species/")) {
    loadSpecies(path.split("/")[2]);
  } else {
    loadGrid();
  }
});

loadGrid();