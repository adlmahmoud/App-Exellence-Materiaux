<script>
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
      <div class="logo-box">EM</div>
      <div class="brand-text">
        <span class="brand-title">Exellence</span>
        <span class="brand-sub">Matériaux</span>
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
  .logo-box {
    width: 40px;
    height: 40px;
    background: #0284c7;
    color: white;
    font-weight: 800;
    font-size: 1.1rem;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
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
    padding: 12px 14px;
    background: transparent;
    border: none;
    border-radius: 8px;
    color: #94a3b8;
    font-size: 0.95rem;
    font-weight: 500;
    cursor: pointer;
    text-align: left;
    transition: all 0.15s ease;
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
  .icon {
    font-size: 1.15rem;
  }
  .sidebar-footer {
    padding: 16px 12px;
    border-top: 1px solid #1e293b;
  }
  .btn-apropos {
    width: 100%;
    padding: 10px 14px;
    background: transparent;
    border: 1px solid #334155;
    color: #94a3b8;
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.85rem;
    text-align: center;
  }
  .btn-apropos:hover {
    background: #1e293b;
    color: #ffffff;
  }
  .main-content {
    flex: 1;
    overflow-y: auto;
    background: #f8fafc;
  }
</style>
