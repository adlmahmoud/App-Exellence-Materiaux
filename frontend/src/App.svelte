<script>
  import logo from './assets/images/gestion-stock.png';
  import Caisse from './pages/Caisse.svelte';
  import Produits from './pages/Produits.svelte';
  import Stock from './pages/Stock.svelte';
  import Factures from './pages/Factures.svelte';
  import Parametres from './pages/Parametres.svelte';
  import APropos from './pages/APropos.svelte';

  let pageActive = 'caisse';
  let modalAPropos = false;

  function changerPage(p) {
    pageActive = p;
  }
</script>

<div class="app-layout">
  <aside class="sidebar">
    <div class="brand">
      <img class="brand-logo" src={logo} alt="Logo Gestion de Stock" />
      <div class="brand-text">
        <span class="brand-title">Gestion</span>
        <span class="brand-sub">de Stock</span>
      </div>
    </div>

    <nav class="nav">
      <button
        class="nav-item {pageActive === 'caisse' ? 'active' : ''}"
        on:click={() => changerPage('caisse')}
      >
        <span class="icon">🛒</span>
        <span>Caisse</span>
      </button>

      <button
        class="nav-item {pageActive === 'produits' ? 'active' : ''}"
        on:click={() => changerPage('produits')}
      >
        <span class="icon">📦</span>
        <span>Produits</span>
      </button>

      <button
        class="nav-item {pageActive === 'stock' ? 'active' : ''}"
        on:click={() => changerPage('stock')}
      >
        <span class="icon">📊</span>
        <span>Stock</span>
      </button>

      <button
        class="nav-item {pageActive === 'factures' ? 'active' : ''}"
        on:click={() => changerPage('factures')}
      >
        <span class="icon">🧾</span>
        <span>Factures & Bons</span>
      </button>

      <button
        class="nav-item {pageActive === 'parametres' ? 'active' : ''}"
        on:click={() => changerPage('parametres')}
      >
        <span class="icon">⚙️</span>
        <span>Paramètres</span>
      </button>
    </nav>

    <div class="sidebar-footer">
      <button class="btn-apropos" on:click={() => modalAPropos = true}>
        ℹ️ À propos
      </button>
    </div>
  </aside>

  <main class="main-content">
    {#if pageActive === 'caisse'}
      <Caisse />
    {:else if pageActive === 'produits'}
      <Produits />
    {:else if pageActive === 'stock'}
      <Stock />
    {:else if pageActive === 'factures'}
      <Factures />
    {:else if pageActive === 'parametres'}
      <Parametres />
    {/if}
  </main>
</div>

<APropos visible={modalAPropos} onClose={() => modalAPropos = false} />

<style>
  :global(body) {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    background-color: #f8fafc;
    color: #1e293b;
    user-select: none;
    -webkit-user-select: none;
  }
  .app-layout {
    display: flex;
    height: 100vh;
    overflow: hidden;
  }
  .sidebar {
    width: 240px;
    background: #0f172a;
    color: white;
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
  }
  .brand {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 24px 20px;
    border-bottom: 1px solid #1e293b;
  }
  .brand-logo {
    width: 52px;
    height: 52px;
    flex-shrink: 0;
    object-fit: contain;
  }
  .brand-text {
    display: flex;
    flex-direction: column;
  }
  .brand-title {
    font-size: 1.1rem;
    font-weight: 700;
    color: #ffffff;
    line-height: 1.2;
  }
  .brand-sub {
    font-size: 0.8rem;
    color: #94a3b8;
  }
  .nav {
    flex: 1;
    padding: 16px 12px;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }
  .nav-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    border: none;
    background: transparent;
    color: #cbd5e1;
    font-size: 0.95rem;
    font-weight: 500;
    border-radius: 8px;
    cursor: pointer;
    text-align: left;
    transition: background 0.15s, color 0.15s;
  }
  .nav-item:hover {
    background: #1e293b;
    color: #ffffff;
  }
  .nav-item.active {
    background: #0284c7;
    color: #ffffff;
    font-weight: 600;
  }
  .nav-item .icon {
    font-size: 1.2rem;
  }
  .sidebar-footer {
    padding: 16px;
    border-top: 1px solid #1e293b;
  }
  .btn-apropos {
    width: 100%;
    padding: 10px 14px;
    background: #1e293b;
    border: 1px solid #334155;
    color: #cbd5e1;
    font-size: 0.85rem;
    border-radius: 6px;
    cursor: pointer;
    text-align: left;
    display: flex;
    align-items: center;
    gap: 8px;
    transition: background 0.15s, color 0.15s;
  }
  .btn-apropos:hover {
    background: #334155;
    color: #ffffff;
  }
  .main-content {
    flex: 1;
    overflow-y: auto;
    background-color: #f8fafc;
  }
</style>
