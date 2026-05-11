"use strict";

const API = {
  summary: () => fetch("/api/summary").then(handle),
  countries: () => fetch("/api/countries").then(handle),
  years: () => fetch("/api/years").then(handle),
  emissions: (countries, from, to) => {
    const params = new URLSearchParams();
    params.set("countries", countries.join(","));
    if (from) params.set("from", from);
    if (to) params.set("to", to);
    return fetch("/api/emissions?" + params.toString()).then(handle);
  },
  top: (n) => fetch("/api/top?n=" + encodeURIComponent(n)).then(handle),
};

async function handle(res) {
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error || ("HTTP " + res.status));
  }
  return res.json();
}

const METRIC_LABELS = {
  co2_total_mt: "CO₂ total (Mt)",
  co2_per_capita_t: "CO₂ per capita (t)",
  renewable_share_pct: "Renováveis (%)",
};

const PALETTE = [
  "#34d399", "#60a5fa", "#f59e0b", "#f472b6",
  "#a78bfa", "#22d3ee", "#fb7185", "#facc15",
];

let compareChart = null;
let topChart = null;

document.addEventListener("DOMContentLoaded", () => {
  init().catch((err) => {
    console.error(err);
    setMessage("Falha ao iniciar: " + err.message);
  });
});

async function init() {
  const [summary, countries] = await Promise.all([
    API.summary(),
    API.countries(),
  ]);

  renderSummary(summary);
  fillCountrySelect(countries);
  fillYearSelects(summary.years_covered);

  document.getElementById("apply").addEventListener("click", refreshCompare);
  document.getElementById("top-n").addEventListener("change", refreshTop);

  preselectDefaultCountries(countries);
  await refreshCompare();
  await refreshTop();
}

function renderSummary(s) {
  set("card-year", s.latest_year || "—");
  set("card-countries", s.total_countries);
  set("card-global", fmt(s.global_co2_mt));
  set("card-renew", fmt(s.avg_renewable_share_pct) + " %");
}

function fillCountrySelect(countries) {
  const sel = document.getElementById("country-select");
  sel.innerHTML = "";
  for (const c of countries) {
    const opt = document.createElement("option");
    opt.value = c;
    opt.textContent = c;
    sel.appendChild(opt);
  }
}

function preselectDefaultCountries(countries) {
  const defaults = ["Brazil", "United States", "China", "Germany"];
  const sel = document.getElementById("country-select");
  for (const opt of sel.options) {
    if (defaults.includes(opt.value)) opt.selected = true;
  }
}

function fillYearSelects(years) {
  const from = document.getElementById("year-from");
  const to = document.getElementById("year-to");
  from.innerHTML = "";
  to.innerHTML = "";
  for (const y of years) {
    const a = document.createElement("option");
    a.value = y;
    a.textContent = y;
    const b = a.cloneNode(true);
    from.appendChild(a);
    to.appendChild(b);
  }
  if (years.length > 0) {
    from.value = years[0];
    to.value = years[years.length - 1];
  }
}

function selectedCountries() {
  const sel = document.getElementById("country-select");
  return Array.from(sel.selectedOptions).map((o) => o.value);
}

async function refreshCompare() {
  setMessage("");
  const countries = selectedCountries();
  if (countries.length === 0) {
    if (compareChart) { compareChart.destroy(); compareChart = null; }
    setMessage("Selecione ao menos um país.");
    return;
  }
  if (countries.length > 5) {
    setMessage("Selecione no máximo 5 países (você selecionou " + countries.length + ").");
    return;
  }
  const from = +document.getElementById("year-from").value || 0;
  const to = +document.getElementById("year-to").value || 0;
  if (from && to && from > to) {
    setMessage("Intervalo inválido: 'De' não pode ser maior que 'Até'.");
    return;
  }
  const metric = document.getElementById("metric").value;

  let series;
  try {
    series = await API.emissions(countries, from, to);
  } catch (err) {
    setMessage(err.message);
    return;
  }

  drawCompare(series, metric);
}

function drawCompare(series, metric) {
  const ctx = document.getElementById("chart-compare").getContext("2d");
  const yearsSet = new Set();
  for (const s of series) for (const p of s.points) yearsSet.add(p.year);
  const years = Array.from(yearsSet).sort((a, b) => a - b);

  const datasets = series.map((s, i) => {
    const byYear = new Map(s.points.map((p) => [p.year, p[metric]]));
    return {
      label: s.country,
      data: years.map((y) => (byYear.has(y) ? byYear.get(y) : null)),
      borderColor: PALETTE[i % PALETTE.length],
      backgroundColor: PALETTE[i % PALETTE.length] + "33",
      tension: 0.25,
      spanGaps: true,
    };
  });

  if (compareChart) compareChart.destroy();
  compareChart = new Chart(ctx, {
    type: "line",
    data: { labels: years, datasets },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { labels: { color: "#e6edf3" } },
        tooltip: { mode: "index", intersect: false },
        title: {
          display: true,
          text: METRIC_LABELS[metric] || metric,
          color: "#e6edf3",
          font: { size: 14, weight: "600" },
        },
      },
      scales: {
        x: { ticks: { color: "#94a3b8" }, grid: { color: "#243447" } },
        y: { ticks: { color: "#94a3b8" }, grid: { color: "#243447" }, beginAtZero: false },
      },
    },
  });
}

async function refreshTop() {
  const n = +document.getElementById("top-n").value || 10;
  let top;
  try {
    top = await API.top(n);
  } catch (err) {
    console.error(err);
    return;
  }
  document.getElementById("top-title").textContent =
    "Top " + top.length + " maiores emissores (" + (top[0]?.year || "—") + ")";
  drawTop(top);
}

function drawTop(top) {
  const ctx = document.getElementById("chart-top").getContext("2d");
  if (topChart) topChart.destroy();
  topChart = new Chart(ctx, {
    type: "bar",
    data: {
      labels: top.map((t) => t.country),
      datasets: [{
        label: "CO₂ total (Mt)",
        data: top.map((t) => t.co2_total_mt),
        backgroundColor: top.map((_, i) => PALETTE[i % PALETTE.length] + "cc"),
        borderColor: top.map((_, i) => PALETTE[i % PALETTE.length]),
        borderWidth: 1,
      }],
    },
    options: {
      indexAxis: "y",
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { display: false } },
      scales: {
        x: { ticks: { color: "#94a3b8" }, grid: { color: "#243447" }, beginAtZero: true },
        y: { ticks: { color: "#94a3b8" }, grid: { color: "#243447" } },
      },
    },
  });
}

function set(id, value) { document.getElementById(id).textContent = value; }

function setMessage(msg) { document.getElementById("compare-msg").textContent = msg; }

function fmt(n) {
  if (n == null) return "—";
  return new Intl.NumberFormat("pt-BR", { maximumFractionDigits: 2 }).format(n);
}
